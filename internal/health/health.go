package health

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/database"
)

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

// Healthy ...
// returns if the instance is healthy
func (m *Manager) Healthy() (err error) {
	if err := database.Ping(m.db); err != nil {
		return err
	}
	return nil
}

// Listen ...
func (m *Manager) Listen() {
	if common.GetAppHealthEnabled() != "true" {
		return
	}

	port := common.GetAppHealthPort()
	router := mux.NewRouter().StrictSlash(true)
	r := router.HandleFunc("/_healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := m.Healthy(); err != nil {
			slog.Error("error app unhealthy", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("not healthy")); err != nil {
				slog.Error("Failed to write response", "error", err)
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("healthy")); err != nil {
			slog.Error("Failed to write response", "error", err)
			return
		}
	})
	server := &http.Server{
		Handler:           r.GetHandler(),
		Addr:              common.GetAppHealthPort(),
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	slog.Info("Health listening on" + port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to listen on HTTP health port", "error", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server didn't exit gracefully", "error", err)
	}
}
