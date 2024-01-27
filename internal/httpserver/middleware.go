package httpserver

import (
	"log"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// scrubHeaders to remove sensitive data logged
func scrubHeaders(in http.Header) (o http.Header) {
	o = http.Header{}
	for k, v := range in {
		o[k] = v
	}
	if az := o.Get("Authorization"); az != "" {
		o.Set("Authorization", "bearer [REDACTED]")
	}
	return o
}

// logging ...
// log the HTTP requests
func logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIP := GetRequestIP(r)
		recorder := &statusRecorder{
			ResponseWriter: w,
		}
		scrubbedHeaders := scrubHeaders(r.Header)
		log.Printf("%v %v %v %v %v %v %#v", recorder.Status, r.Method, r.URL, r.Proto, requestIP, r.RemoteAddr, scrubbedHeaders)
		next.ServeHTTP(recorder, r)
	})
}

// HTTPHeaderBackendAllowTypes heads to check for content type
type HTTPHeaderBackendAllowTypes string

const (
	// HTTPHeaderBackendAllowTypesContentType use the Content-Type http header
	HTTPHeaderBackendAllowTypesContentType HTTPHeaderBackendAllowTypes = "Content-Type"
	// HTTPHeaderBackendAllowTypesAccept use the Accept http header
	HTTPHeaderBackendAllowTypesAccept HTTPHeaderBackendAllowTypes = "Accept"
)
