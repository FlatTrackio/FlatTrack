package httpserver

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/rs/cors"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/emails"
	"gitlab.com/flattrack/flattrack/internal/groups"
	"gitlab.com/flattrack/flattrack/internal/health"
	"gitlab.com/flattrack/flattrack/internal/migrations"
	"gitlab.com/flattrack/flattrack/internal/registration"
	"gitlab.com/flattrack/flattrack/internal/scheduling"
	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/internal/shoppinglist"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/internal/users"
)

type HTTPServer struct {
	server   *http.Server
	listener *pq.Listener
	db       *sql.DB

	users        *users.Manager
	shoppinglist *shoppinglist.Manager
	emails       *emails.Manager
	groups       *groups.Manager
	health       *health.Manager
	migrations   *migrations.Manager
	registration *registration.Manager
	settings     *settings.Manager
	system       *system.Manager
	scheduling   *scheduling.Manager
}

func NewHTTPServer(
	db *sql.DB,
	listener *pq.Listener,
	users *users.Manager,
	shoppinglist *shoppinglist.Manager,
	emails *emails.Manager,
	groups *groups.Manager,
	health *health.Manager,
	migrations *migrations.Manager,
	registration *registration.Manager,
	settings *settings.Manager,
	system *system.Manager,
	scheduling *scheduling.Manager,
) (h *HTTPServer) {
	h = &HTTPServer{
		db:           db,
		listener:     listener,
		users:        users,
		shoppinglist: shoppinglist,
		emails:       emails,
		groups:       groups,
		health:       health,
		migrations:   migrations,
		registration: registration,
		settings:     settings,
		system:       system,
		scheduling:   scheduling,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/_healthz", h.Healthz)
	apiRouter := router.
		PathPrefix("/api").
		Subrouter()
	apiRouter.NotFoundHandler = h.HTTP404()
	h.registerAPIHandlers(apiRouter)

	passthrough := &frontendOptions{
		SetupMessage: common.GetAppSetupMessage(),
		LoginMessage: common.GetAppLoginMessage(),
	}
	router.PathPrefix("/").Handler(frontendHandler(common.GetAppDistFolder(), passthrough)).Methods(http.MethodGet)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "User-Agent", "Accept-Encoding"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})
	router.Use(logging)
	router.Use(c.Handler)
	router.Use(gziphandler.GzipHandler)

	h.server = &http.Server{
		Handler:      router,
		Addr:         common.GetAppPort(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return h
}

func (h *HTTPServer) Listen() {
	log.Println("HTTP listening on", h.server.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
