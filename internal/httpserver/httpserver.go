package httpserver

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
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
	server *http.Server
	db     *sql.DB

	users           *users.Manager
	shoppinglist    *shoppinglist.Manager
	emails          *emails.Manager
	groups          *groups.Manager
	health          *health.Manager
	migrations      *migrations.Manager
	registration    *registration.Manager
	settings        *settings.Manager
	system          *system.Manager
	scheduling      *scheduling.Manager
	maintenanceMode bool
	instanceURL     *url.URL
}

func NewHTTPServer(
	db *sql.DB,
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
	maintenanceMode bool,
) (h *HTTPServer) {
	var err error
	h = &HTTPServer{}
	h.db = db
	h.users = users
	h.shoppinglist = shoppinglist
	h.emails = emails
	h.groups = groups
	h.health = health
	h.migrations = migrations
	h.registration = registration
	h.settings = settings
	h.system = system
	h.scheduling = scheduling
	h.maintenanceMode = maintenanceMode
	h.instanceURL, err = common.GetInstanceURL()
	if err != nil {
		slog.Info("Failed to get instance URL")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/_healthz", h.Healthz)
	apiRouter := router.
		PathPrefix("/api").
		Headers("Accept", "application/json").
		Subrouter()
	apiRouter.NotFoundHandler = h.HTTP404()
	apiRouter.MethodNotAllowedHandler = h.HTTPMethodNotAllowed()
	h.registerAPIHandlers(apiRouter)

	passthrough := &frontendOptions{
		LoginMessage:           common.GetAppLoginMessage(),
		MaintenanceModeMessage: common.GetAppMaintenanceModeMessage(),
		SetupMessage:           common.GetAppSetupMessage(),
	}
	router.PathPrefix("/").Handler(frontendHandler(common.GetAppWebFolder(), passthrough)).Methods(http.MethodGet)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "User-Agent", "Accept-Encoding"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})
	router.Use(logging)
	router.Use(h.RewriteToDomain)
	router.Use(c.Handler)
	router.Use(gziphandler.GzipHandler)
	router.NotFoundHandler = h.HTTP404()
	router.MethodNotAllowedHandler = h.HTTPMethodNotAllowed()

	h.server = &http.Server{
		Handler:      router,
		Addr:         common.GetAppPort(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return h
}

func (h *HTTPServer) Listen() {
	slog.Info("HTTP listening on " + h.server.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to listen on HTTP port", "error", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		slog.Error("Server didn't exit gracefully", "error", err)
	}
}
