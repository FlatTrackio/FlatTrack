package httpserver

import (
	"bufio"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"slices"
	"time"

	"gitlab.com/flattrack/flattrack/internal/common"
)

type statusRecorder struct {
	http.ResponseWriter
	Status            int
	statusCodeWritten bool
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
	r.statusCodeWritten = true
}

func (w *statusRecorder) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	if !w.statusCodeWritten {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	return n, err
}

func (w *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijack not supported")
}

// scrubHeaders to remove sensitive data logged
func scrubHeaders(in http.Header) (o http.Header) {
	headers := []string{"Authorization", "Cookie"}
	o = http.Header{}
	maps.Copy(o, in)
	for _, h := range headers {
		if az := o.Get(h); az != "" {
			o.Set(h, "[REDACTED]")
		}
	}
	return o
}

// logging ...
// log the HTTP requests
func logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestIP := GetRequestIP(r)
		recorder := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		scrubbedHeaders := scrubHeaders(r.Header)
		slog.Info(
			"HTTP Request",
			"status", recorder.Status,
			"method", r.Method,
			"url", r.URL.String(),
			"proto", r.Proto,
			"requestIP", requestIP,
			"remoteAddr", r.RemoteAddr,
			"headers", scrubbedHeaders,
			"duration", time.Since(start).String(),
		)
	})
}

// HTTPHeaderBackendAllowTypes headers to check for content type
type HTTPHeaderBackendAllowTypes string

const (
	// HTTPHeaderBackendAllowTypesContentType use the Content-Type http header
	HTTPHeaderBackendAllowTypesContentType HTTPHeaderBackendAllowTypes = "Content-Type"
	// HTTPHeaderBackendAllowTypesAccept use the Accept http header
	HTTPHeaderBackendAllowTypesAccept HTTPHeaderBackendAllowTypes = "Accept"
)

func (h *HTTPServer) RewriteToDomain(next http.Handler) http.Handler {
	noRedirectDomains := common.GetInstanceURLNoRedirectDomains()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.instanceURL != nil &&
			r.Host != h.instanceURL.Host &&
			!slices.Contains(noRedirectDomains, r.Host) &&
			r.URL.Path != "/_healthz" {
			sourceHost := r.Host
			r.URL.Host = h.instanceURL.Host
			r.URL.Scheme = h.instanceURL.Scheme
			slog.Info("Redirecting domain", "source", sourceHost, "destination", h.instanceURL.Host, "url", r.URL.String())
			http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
