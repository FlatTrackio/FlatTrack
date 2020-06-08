package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"log"
	"net/http"
)

// Handle ...
// HTTP handler for metrics
func Handle() {
	if common.GetAppMetricsEnabled() != "true" {
		return
	}

	port := common.GetAppMetricsPort()
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Metrics listening on %v\n", port)
	http.ListenAndServe(port, nil)
}
