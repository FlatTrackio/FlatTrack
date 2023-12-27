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

// Logging ...
// log the HTTP requests
func logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIP := GetRequestIP(r)
		recorder := &statusRecorder{
			ResponseWriter: w,
		}
		log.Printf("%v %v %v %v %v %v %#v", recorder.Status, r.Method, r.URL, r.Proto, requestIP, r.RemoteAddr, r.Header)
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
