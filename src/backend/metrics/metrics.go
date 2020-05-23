package metrics

import (
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handle() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Metrics listening on :2112")
	http.ListenAndServe(":2112", nil)
}
