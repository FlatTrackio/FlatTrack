package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// Handle ...
// HTTP handler for metrics
func Handle() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Metrics listening on :2112")
	http.ListenAndServe(":2112", nil)
}
