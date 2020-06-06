package health

import (
	"database/sql"
	"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"log"
	"net/http"
)

// healthz ...
// HTTP handler for health checks
func healthz(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "App unhealthy"
		code := http.StatusInternalServerError

		err := database.Ping(db)
		if err == nil {
			response = "App healthy"
			code = http.StatusOK
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: err == nil,
		}
		routes.JSONResponse(r, w, code, JSONresp)
	}
}

// Handle ...
func Handle(db *sql.DB) {
	if common.GetAppHealthEnabled() != "true" {
		return
	}

	port := common.GetAppHealthPort()
	http.Handle("/_healthz", healthz(db))
	log.Printf("Health listening on %v", port)
	http.ListenAndServe(port, nil)
}
