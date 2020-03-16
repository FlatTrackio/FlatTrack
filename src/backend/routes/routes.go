/*
	route related
*/

package routes

import (
	"encoding/json"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

// PostUser
// create a user
func PostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to create user account"
		vars := mux.Vars(r)

		accountName := mux.Vars["name"]
		accountEmail := mux.Vars["email"]
		accountGroups := mux.Var["groups"]
		accountPassword := mux.Var["password"]
		accountPhoneNumber := mux.Var["phoneNumber"]
		accountCreationTimestamp := time.Now()

		userAccount, err := users.CreateUser(db, user)
		if err == nil {
			code = 200
			response = "Created user account"
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccount,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	// root of API
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hey! you're talking to the Flattrack API",
		},
	}
	JSONResponse(r, w, 200, JSONresp)
}

func UnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	JSONResponse(r, w, 404, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}
