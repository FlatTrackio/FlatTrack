/*
	route related
*/

package routes

import (
	"database/sql"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

// GetAllUsers
// get a list of all users
func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to fetch user accounts"

		users, err := users.GetAllUsers(db)
		if err == nil {
			code = 200
			response = "Fetched user accounts"
		} else {
			code = 400
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: users,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetUser
// get a user by id or email (whatever is provided in the given respective order)
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to fetch user account"
		vars := mux.Vars(r)
		id := vars["id"]

		user := types.UserSpec{
			Id:    id,
		}
		fmt.Println(user)

		user, err := users.GetUserById(db, id)
		if err == nil {
			code = 200
			response = "Fetched user account"
		} else {
			code = 400
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: user,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostUser
// create a user
func PostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
		// - handle existing user
		// - handle empty fields

		code := 500
		response := "Failed to create user account"

		var user types.UserSpec
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
		json.Unmarshal(body, &user)

		userAccount, err := users.CreateUser(db, user)
		if err == nil {
			code = 200
			response = "Created user account"
		} else {
			code = 400
			response = err.Error()
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

// Root
// usually /api endpoint
func Root(w http.ResponseWriter, r *http.Request) {
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hey! you're talking to the Flattrack API",
		},
	}
	JSONResponse(r, w, 200, JSONresp)
}

// UnknownEndpoint
// response for hitting an unknown endpoint
func UnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	JSONResponse(r, w, 404, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}
