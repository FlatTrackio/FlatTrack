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
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab.com/flattrack/flattrack/internal/common"
)

type Manager struct {
	enabled bool

	server *http.Server
}

func NewManager() *Manager {
	router := mux.NewRouter().StrictSlash(true)
	r := router.Handle("/metrics", promhttp.Handler())
	return &Manager{
		enabled: common.GetAppMetricsEnabled(),
		server: &http.Server{
			Handler:           r.GetHandler(),
			Addr:              common.GetAppMetricsPort(),
			WriteTimeout:      15 * time.Second,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// Listen ...
// HTTP handler for metrics
func (m *Manager) Listen() {
	if !m.enabled {
		return
	}
	log.Printf("Metrics listening on %v\n", m.server.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
