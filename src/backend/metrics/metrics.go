package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"gitlab.com/flattrack/flattrack/src/backend/common"
)

// Handle ...
// HTTP handler for metrics
func Handle() {
	port := common.GetAppMetricsPort()
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Metrics listening on %v\n", port)
	http.ListenAndServe(port, nil)
}
