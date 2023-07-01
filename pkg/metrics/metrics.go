/*
  metrics
    prometheus metrics
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/flattrack/flattrack/pkg/common"
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
