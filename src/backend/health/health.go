package health

import (
	"database/sql"
	"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"log"
	"net/http"
)

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

// Handle
func Handle(db *sql.DB) {
	http.Handle("/_healthz", healthz(db))
	log.Println("Health listening on :8081")
	http.ListenAndServe(":8081", nil)
}
