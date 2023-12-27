package health

import (
	"context"
	"database/sql"
	"log"
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
			log.Printf("error app unhealthy: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("not healthy")); err != nil {
				log.Fatal(err)
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("healthy")); err != nil {
			log.Fatal(err)
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
	log.Printf("Health listening on %v", port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
