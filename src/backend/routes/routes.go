/*
  routes
    routes
      declare all API routes
*/

package routes

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/src/backend/registration"
	"gitlab.com/flattrack/flattrack/src/backend/settings"
	"gitlab.com/flattrack/flattrack/src/backend/system"
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
			Id: id,
		}

		user, err := users.GetUserById(db, id)
		if err != nil || user.Id == "" {
			code = 404
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			JSONResponse(r, w, code, JSONresp)
			return
		}
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
		code := 500
		response := "Failed to create user account"

		var user types.UserSpec
		body, err := ioutil.ReadAll(r.Body)
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

// DeleteUser
// delete a user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to delete user account"

		vars := mux.Vars(r)
		userId := vars["id"]

		userInDB, err := users.GetUserById(db, userId)
		if err != nil || userInDB.Id == "" {
			code = 404
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: false,
			}
			JSONResponse(r, w, code, JSONresp)
			return
		}

		err = users.DeleteUserById(db, userInDB.Id)
		if err == nil {
			code = 200
			response = "Deleted user account"
		} else {
			code = 400
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: err == nil,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// Initialized
// check if the server has been initialized
func GetSystemInitialized(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to fetch if this FlatTrack instance has initialized"
		initialized, err := system.GetHasInitialized(db)
		if err == nil {
			code = 200
		}
		if err == nil && initialized == "true" {
			response = "This FlatTrack instance has initialized"
		} else if err == nil && initialized == "false" {
			response = "This FlatTrack instance has not initialized"
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: initialized == "true",
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserAuth
// authenticate a user
func UserAuth(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to authenticate user, incorrect email or password"
		code := 403
		jwtToken := ""

		var user types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		userInDB, err := users.GetUserByEmail(db, user.Email)
		if err != nil {
			response = "Failed to find user"
		}
		// Check password locally, fall back to remote if incorrect
		matches, err := users.CheckUserPassword(db, userInDB.Email, user.Password)
		if err == nil && matches == true {
			jwtToken, _ = users.GenerateJWTauthToken(db, user.Id)
			response = "Successfully authenticated user"
			code = 200
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: jwtToken,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserAuthValidate
// validate an auth token
func UserAuthValidate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to validate authentication token"
		code := 400

		valid, err := users.ValidateJWTauthToken(db, r)
		if valid == true && err == nil {
			response = "Authentication token is valid"
			code = 200
		} else {
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: valid,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSettingsFlatName
// response with the name of the flat
func GetSettingsFlatName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch the flat name"
		code := 500

		flatName, err := settings.GetFlatName(db)
		if flatName == "" {
			response = "Flat name is not set"
			code = 200
		} else if err == nil {
			response = "Fetched the flat name"
			code = 200
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: flatName,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// SetSettingsFlatName
// update the flat's name
func SetSettingsFlatName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to update the flat name"

		var flatName types.FlatName
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &flatName)

		err := settings.SetFlatName(db, flatName.FlatName)
		if err == nil {
			code = 200
			response = "Successfully set flat name"
		} else {
			code = 400
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: err == nil,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostAdminRegister
// register the instance of FlatTrack
func PostAdminRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to register the FlatTrack instance"

		initialized, err := system.GetHasInitialized(db)
		if err == nil && initialized == "true" {
			response = "This instance is already registered"
			code = 200
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: true,
				Data: "",
			}
			JSONResponse(r, w, code, JSONresp)
		}

		var registrationForm types.Registration
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &registrationForm)

		registered, jwt, err := registration.Register(db, registrationForm)
		if err == nil {
			code = 200
			response = "Successfully registered the FlatTrack instance"
		} else {
			code = 400
			response = err.Error()
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: registered,
			Data: jwt,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// Root
// /api endpoint
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

// HTTPvalidateJWT
// middleware for checking JWT auth token validity
func HTTPvalidateJWT(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if completed, err := users.ValidateJWTauthToken(db, r); completed == true && err == nil {
				h.ServeHTTP(w, r)
				return
			}
			JSONResponse(r, w, 401, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unauthorized",
				},
			})
		}
	}
}
