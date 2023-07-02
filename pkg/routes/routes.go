/*
  routes
    routes
      declare all API routes
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

package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"gitlab.com/flattrack/flattrack/pkg/common"
	// TODO integrate
	_ "gitlab.com/flattrack/flattrack/pkg/files"
	"gitlab.com/flattrack/flattrack/pkg/groups"
	"gitlab.com/flattrack/flattrack/pkg/health"
	"gitlab.com/flattrack/flattrack/pkg/registration"
	"gitlab.com/flattrack/flattrack/pkg/settings"
	"gitlab.com/flattrack/flattrack/pkg/shoppinglist"
	"gitlab.com/flattrack/flattrack/pkg/system"
	"gitlab.com/flattrack/flattrack/pkg/types"
	"gitlab.com/flattrack/flattrack/pkg/users"
)

// RouteHandler handle routes
type RouteHandler struct {
	// db         *sql.DB
	// fileAccess files.FileAccess
}

// TODO rejig to use NewRoutes(db, fileAccess)
// TODO restructure to not use a single JSON response in each handler
//      and instead respond on error afterwards

// GetAllUsers ...
// swagger:route GET /users users getAllUsers
// responses:
//
//	200: []userSpec
//
// get a list of all users
func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		userSelectorID := r.FormValue("id")
		userSelectorNotID := r.FormValue("notId")
		userSelectorGroup := r.FormValue("group")
		var errJWT error = nil
		if r.FormValue("notSelf") == "true" {
			userSelectorNotID, errJWT = users.GetIDFromJWT(db, r)
		}

		selectors := types.UserSelector{
			ID:    userSelectorID,
			NotID: userSelectorNotID,
			Group: userSelectorGroup,
		}

		users, err := users.GetAllUsers(db, false, selectors)
		if err == nil && errJWT == nil {
			code = http.StatusOK
			response = "Fetched user accounts"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: users,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetUser ...
// get a user by id or email (whatever is provided in the given respective order)
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string
		vars := mux.Vars(r)
		id := vars["id"]

		user, err := users.GetUserByID(db, id, false)
		if err != nil {
			context = err.Error()
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}
		if user.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}
		if err == nil {
			code = http.StatusOK
			response = "Fetched user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: user,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostUser ...
// create a user
func PostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var user types.UserSpec
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}
		if err := json.Unmarshal(body, &user); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		userAccount, err := users.CreateUser(db, user, user.Password == "")
		if err == nil && userAccount.ID != "" {
			code = http.StatusCreated
			response = "Created user account"
			context = fmt.Sprintf("'%v'", userAccount.ID)
		} else {
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccount,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutUser ...
// updates a user account by their id
func PutUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var userAccount types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &userAccount); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		userID := vars["id"]

		user, err := users.GetUserByID(db, userID, false)
		if err != nil || user.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		// TODO disallow admins to remove their own admin group access
		userAccountUpdated, err := users.UpdateProfileAdmin(db, userID, userAccount)
		if err == nil && userAccountUpdated.ID != "" {
			code = http.StatusOK
			response = "Successfully updated the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccountUpdated,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchUser ...
// patches a user account by their id
func PatchUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var userAccount types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &userAccount); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		userID := vars["id"]

		user, err := users.GetUserByID(db, userID, false)
		if err != nil || user.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		// TODO disallow admins to remove their own admin group access
		userAccountPatched, err := users.PatchProfileAdmin(db, userID, userAccount)
		if err == nil && userAccountPatched.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccountPatched,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchUserDisabled ...
// patches a user account's disabled field by their id
func PatchUserDisabled(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var userAccount types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &userAccount); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		userID := vars["id"]

		user, err := users.GetUserByID(db, userID, false)
		if err != nil || user.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		jwtID, errID := users.GetIDFromJWT(db, r)
		if errID != nil {
			code = http.StatusBadRequest
			response = "Failed get user ID from token"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		isAdmin, errGroup := groups.CheckUserInGroup(db, jwtID, "admin")
		if errGroup != nil {
			code = http.StatusInternalServerError
			response = "Failed to check for user in group"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		if isAdmin && userID == jwtID {
			code = http.StatusForbidden
			response = "Unable to disabled own user account"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.UserSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		userAccountPatched, err := users.PatchUserDisabledAdmin(db, userID, userAccount.Disabled)
		if err == nil {
			code = http.StatusOK
			response = "Successfully patched the user account disabled field"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccountPatched,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// DeleteUser ...
// delete a user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		vars := mux.Vars(r)
		userID := vars["id"]

		userInDB, err := users.GetUserByID(db, userID, false)
		if err != nil || userInDB.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find user"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: false,
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		myUserID, errID := users.GetIDFromJWT(db, r)
		if errID != nil {
			code = http.StatusBadRequest
			response = "Failed to get user account ID from token"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: false,
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		if myUserID == userID {
			code = http.StatusForbidden
			response = "Deleting own user account is disallowed"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: nil,
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		err = users.DeleteUserByID(db, userInDB.ID)
		if err == nil {
			code = http.StatusOK
			response = "Deleted user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: err == nil,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetProfile ...
// returns the authenticated user's profile
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		user, err := users.GetProfile(db, r)
		if err == nil && user.ID != "" {
			code = http.StatusOK
			response = "Fetched user account"
		} else {
			code = http.StatusNotFound
			response = "Failed to find your profile"
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: user,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutProfile ...
// Update a user account their id from their JWT
func PutProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var userAccount types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &userAccount); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		userAccountUpdated, err := users.UpdateProfile(db, id, userAccount)
		if err == nil && errID == nil && userAccountUpdated.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccountUpdated,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchProfile ...
// patches a user account their id from their JWT
func PatchProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var userAccount types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &userAccount); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		userAccountPatched, err := users.PatchProfile(db, id, userAccount)
		if err == nil && errID == nil && userAccountPatched.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccountPatched,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSystemInitialized ...
// check if the server has been initialized
func GetSystemInitialized(db *sql.DB) http.HandlerFunc {
	systemManager := &system.Manager{DB: db}
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string
		initialized, err := systemManager.GetHasInitialized()
		if err == nil {
			code = http.StatusOK
		}
		if err == nil && initialized == "true" {
			response = "This FlatTrack instance has initialized"
		} else if err == nil && initialized == "false" {
			response = "This FlatTrack instance has not initialized"
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: initialized == "true",
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserAuth ...
// authenticate a user
func UserAuth(db *sql.DB) http.HandlerFunc {
	userManager := users.UserManager{DB: db}
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int
		jwtToken := ""

		var user types.UserSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &user); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		userInDB, err := users.GetUserByEmail(db, user.Email, false)
		if err != nil {
			log.Println(response, context)
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Data: jwtToken,
			}
			JSONResponse(r, w, code, JSONresp)
			return
		}
		if userInDB.ID != "" && !userInDB.Registered {
			response = "Account not yet registered"
			code = http.StatusForbidden
		} else if userInDB.Disabled {
			response = "User account has been disabled"
			code = http.StatusForbidden
		}
		// Check password locally, fall back to remote if incorrect
		matches, err := users.CheckUserPassword(db, userInDB.Email, user.Password)
		if err == nil && matches && code == http.StatusUnauthorized {
			jwtToken, _ = userManager.GenerateJWTauthToken(userInDB.ID, userInDB.AuthNonce, 0)
			response = "Successfully authenticated user"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: jwtToken,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserAuthValidate ...
// validate an auth token
func UserAuthValidate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		valid, claims, err := users.ValidateJWTauthToken(db, r)
		if valid && err == nil {
			response = "Authentication token is valid"
			code = http.StatusOK
		} else {
			response = err.Error()
			context = fmt.Sprintf("for user with ID '%v'", claims.ID)
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: valid,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserAuthReset ...
// invalidates all JWTs
func UserAuthReset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		id, err := users.GetIDFromJWT(db, r)
		if err != nil {
			response = "Failed to find account"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}
		err = users.GenerateNewAuthNonce(db, id)
		if err == nil {
			response = "Successfully signed out all devices logged in"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UserCanIgroup ...
// respond whether the current user account is in a group
func UserCanIgroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		groupName := vars["name"]

		group, err := groups.GetGroupByName(db, groupName)
		if err != nil || group.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find group"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.GroupSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		userIsInGroup, err := groups.CheckUserInGroup(db, id, groupName)
		if err == nil && errID == nil {
			response = "Determined if user account can perform tasks of group"
			code = http.StatusOK
		}

		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: userIsInGroup,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSettingsFlatName ...
// responds with the name of the flat
func GetSettingsFlatName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		flatName, err := settings.GetFlatName(db)
		if flatName == "" {
			response = "Flat name is not set"
			code = http.StatusOK
		} else if err == nil {
			response = "Fetched the flat name"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: flatName,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// SetSettingsFlatName ...
// update the flat's name
func SetSettingsFlatName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var flatName types.FlatName
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &flatName); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		err := settings.SetFlatName(db, flatName.FlatName)
		if err == nil {
			code = http.StatusOK
			response = "Successfully set flat name"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: err == nil,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostAdminRegister ...
// register the instance of FlatTrack
func PostAdminRegister(db *sql.DB) http.HandlerFunc {
	systemManager := &system.Manager{DB: db}
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		initialized, err := systemManager.GetHasInitialized()
		if err == nil && initialized == "true" {
			response = "This instance is already registered"
			code = http.StatusOK
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: true,
				Data: "",
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		var registrationForm types.Registration
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &registrationForm); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		registered, jwt, err := registration.Register(db, registrationForm)
		if err == nil {
			code = http.StatusCreated
			response = "Successfully registered the FlatTrack instance"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
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

// GetShoppingList ...
// responds with list of shopping lists
func GetShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		shoppingList, err := shoppinglist.GetShoppingList(db, id)
		if err == nil && shoppingList.ID != "" {
			response = "Fetched the shopping lists"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingList,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetShoppingLists ...
// responds with shopping list by id
func GetShoppingLists(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		modificationTimestampAfter, _ := strconv.Atoi(r.FormValue("modificationTimestampAfter"))
		creationTimestampAfter, _ := strconv.Atoi(r.FormValue("creationTimestampAfter"))
		limit, _ := strconv.Atoi(r.FormValue("limit"))

		options := types.ShoppingListOptions{
			SortBy: r.FormValue("sortBy"),
			Limit:  limit,
			Selector: types.ShoppingListSelector{
				Completed:                  r.FormValue("completed"),
				ModificationTimestampAfter: modificationTimestampAfter,
				CreationTimestampAfter:     creationTimestampAfter,
			},
		}

		shoppingLists, err := shoppinglist.GetShoppingLists(db, options)
		if err == nil {
			response = "Fetched the shopping lists"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: shoppingLists,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostShoppingList ...
// creates a new shopping list to add items to
func PostShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingList types.ShoppingListSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingList); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		options := types.ShoppingItemOptions{
			Selector: types.ShoppingItemSelector{
				TemplateListItemSelector: r.FormValue("templateListItemSelector"),
			},
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingList.Author = id
		shoppingListInserted, err := shoppinglist.CreateShoppingList(db, shoppingList, options)
		if err == nil && errID == nil {
			code = http.StatusCreated
			response = "Successfully created the shopping list"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingListInserted,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchShoppingList ...
// patches an existing shopping list
func PatchShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingList types.ShoppingListSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingList); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		listID := vars["id"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingList.AuthorLast = id
		shoppingListPatched, err := shoppinglist.PatchShoppingList(db, listID, shoppingList)
		if err == nil && errID == nil && shoppingListPatched.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the shopping list"
		} else {
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingListPatched,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutShoppingList ...
// updates an existing shopping list
func PutShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingList types.ShoppingListSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingList); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		listID := vars["id"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingList.AuthorLast = id
		shoppingListUpdated, err := shoppinglist.UpdateShoppingList(db, listID, shoppingList)
		if err == nil && errID == nil && shoppingListUpdated.ID != "" {
			code = http.StatusOK
			response = "Successfully updated the shopping list"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingListUpdated,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// DeleteShoppingList ...
// delete a new shopping list by it's id
func DeleteShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		vars := mux.Vars(r)
		listID := vars["id"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		err = shoppinglist.DeleteShoppingList(db, listID)
		if err == nil {
			code = http.StatusOK
			response = "Successfully deleted the shopping list"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetShoppingListItems ...
// responds with shopping items by list id
func GetShoppingListItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		options := types.ShoppingItemOptions{
			SortBy: r.FormValue("sortBy"),
			Selector: types.ShoppingItemSelector{
				Obtained: r.FormValue("obtained"),
			},
		}

		list, err := shoppinglist.GetShoppingList(db, id)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		// TODO add item selectors for this endpoint
		shoppingListItems, err := shoppinglist.GetShoppingListItems(db, id, options)
		if err == nil {
			response = "Fetched the shopping list items"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: shoppingListItems,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetShoppingListItem ...
// responds with list of shopping lists
func GetShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		itemID := vars["itemId"]
		listID := vars["listId"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		shoppingListItem, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err == nil && shoppingListItem.ID != "" {
			response = "Fetched the shopping list item"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingListItem,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostItemToShoppingList ...
// adds an item to a shopping list
func PostItemToShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingItem types.ShoppingItemSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingItem); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		listID := vars["id"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingItem.Author = id
		shoppingItemInserted, err := shoppinglist.AddItemToList(db, listID, shoppingItem)
		if err == nil && errID == nil {
			code = http.StatusCreated
			response = "Successfully created the shopping list item"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: shoppingItemInserted,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchShoppingListCompleted ...
// adds an item to a shopping list
func PatchShoppingListCompleted(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingList types.ShoppingListSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingList); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		listID := vars["id"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		patchedList, err := shoppinglist.SetListCompleted(db, listID, shoppingList.Completed, id)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully patched the shopping list completed field"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: patchedList,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchShoppingListItem ...
// patches an item in a shopping list
func PatchShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingItem types.ShoppingItemSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingItem); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		itemID := vars["id"]
		listID := vars["listId"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		item, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err != nil || item.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list item"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingItemSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingItem.AuthorLast = id
		patchedItem, err := shoppinglist.PatchItem(db, listID, itemID, shoppingItem)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully patched the shopping list item"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: patchedItem,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutShoppingListItem ...
// updates an item in a shopping list
func PutShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingItem types.ShoppingItemSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingItem); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		itemID := vars["id"]
		listID := vars["listId"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		item, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err != nil || item.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list item"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingItemSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingItem.AuthorLast = id
		updatedItem, err := shoppinglist.UpdateItem(db, listID, itemID, shoppingItem)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully updated the shopping list item"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: updatedItem,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PatchShoppingListItemObtained ...
// patches an item in a shopping list
func PatchShoppingListItemObtained(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var shoppingItem types.ShoppingItemSpec
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &shoppingItem); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		vars := mux.Vars(r)
		itemID := vars["id"]
		listID := vars["listId"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		item, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err != nil || item.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list item"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingItemSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		patchedItem, err := shoppinglist.SetItemObtained(db, listID, itemID, shoppingItem.Obtained, id)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully patched the shopping list item obtained field"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: patchedItem,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// DeleteShoppingListItem ...
// delete a shopping list item by it's id
func DeleteShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		vars := mux.Vars(r)
		itemID := vars["itemId"]
		listID := vars["listId"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		item, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err != nil || item.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list item"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingItemSpec{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		err = shoppinglist.RemoveItemFromList(db, itemID, listID)
		if err == nil {
			code = http.StatusOK
			response = "Successfully deleted the shopping list item"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetShoppingListItemTags ...
// responds with tags used in shopping list items from a list
func GetShoppingListItemTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		listID := vars["listId"]

		tags, err := shoppinglist.GetShoppingListTags(db, listID)
		if err == nil {
			response = "Fetched the shopping list item tags"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: tags,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UpdateShoppingListItemTag ...
// updates then tag name used in shopping list items from a list
func UpdateShoppingListItemTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		listID := vars["listId"]
		tag := vars["tagName"]

		list, err := shoppinglist.GetShoppingList(db, listID)
		if err != nil || list.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping list"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingListSpec{},
			}
			JSONResponse(r, w, code, JSONresp)
			log.Println(response, context)
			return
		}

		var tagUpdate types.ShoppingTag
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &tagUpdate); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		tag, err = shoppinglist.UpdateShoppingListTag(db, listID, tag, tagUpdate.Name)
		if err == nil {
			response = "Updated the shopping list item tag name"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: tag,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetAllShoppingTags ...
// responds with all tags used in shopping list items
func GetAllShoppingTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		options := types.ShoppingTagOptions{
			SortBy: r.FormValue("sortBy"),
		}

		tags, err := shoppinglist.GetAllShoppingTags(db, options)
		if err == nil {
			response = "Fetched the shopping list item tags"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: tags,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostShoppingTag ...
// creates a tag name
func PostShoppingTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		var tag types.ShoppingTag
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &tag); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		tag.Author = id
		tag, err := shoppinglist.CreateShoppingTag(db, tag)
		if err == nil && errID == nil {
			response = "Updated a shopping tag"
			code = http.StatusCreated
		} else {
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: tag,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetShoppingTag ...
// gets a shopping tag by id
func GetShoppingTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		tag, err := shoppinglist.GetShoppingTag(db, id)
		if err == nil {
			response = "Fetched a shopping tag"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: tag,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// UpdateShoppingTag ...
// updates a tag name
func UpdateShoppingTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		tagInDB, err := shoppinglist.GetShoppingTag(db, id)
		if err != nil || tagInDB.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping tag"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingTag{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		var tagUpdate types.ShoppingTag
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &tagUpdate); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		userID, errID := users.GetIDFromJWT(db, r)
		tagUpdate.AuthorLast = userID
		tag, err := shoppinglist.UpdateShoppingTag(db, id, tagUpdate)
		if err == nil && errID == nil {
			response = "Updated a shopping tag"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: tag,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// DeleteShoppingTag ...
// deletes a shopping tag by id
func DeleteShoppingTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		vars := mux.Vars(r)
		id := vars["id"]

		tagInDB, err := shoppinglist.GetShoppingTag(db, id)
		if err != nil || tagInDB.ID == "" {
			code = http.StatusNotFound
			response = "Failed to find shopping tag"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.ShoppingTag{},
			}
			log.Println(response, context)
			JSONResponse(r, w, code, JSONresp)
			return
		}

		err = shoppinglist.DeleteShoppingTag(db, id)
		if err == nil {
			code = http.StatusOK
			response = "Successfully deleted the shopping tag"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSettingsShoppingListNotes ...
// responds with the notes for shopping lists
func GetSettingsShoppingListNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		notes, err := settings.GetShoppingListNotes(db)
		if notes == "" {
			response = "There are no notes set for the shopping list"
			code = http.StatusOK
		} else if err == nil {
			response = "Fetched the shopping list notes"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: notes,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutSettingsShoppingList ...
// update the notes for shopping lists
func PutSettingsShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var notes types.ShoppingListNotes
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &notes); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		err := settings.SetShoppingListNotes(db, notes.Notes)
		if err == nil {
			code = http.StatusOK
			response = "Successfully set shopping list notes"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
			notes.Notes = ""
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: notes.Notes,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSettingsFlatNotes ...
// responds with the notes for flat
func GetSettingsFlatNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		notes, err := settings.GetFlatNotes(db)
		if notes == "" {
			response = "There are no notes set for the flat"
			code = http.StatusOK
		} else if err == nil {
			response = "Fetched the flat notes"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: types.FlatNotes{
				Notes: notes,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PutSettingsFlatNotes ...
// update the notes for flat
func PutSettingsFlatNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var code int
		var response string

		var notes types.FlatNotes
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &notes); err != nil {
			log.Printf("error: failed to unmarshal; %v\n", err)
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		err := settings.SetFlatNotes(db, notes.Notes)
		if err == nil {
			code = http.StatusOK
			response = "Successfully set flat notes"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
			notes.Notes = ""
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: notes.Notes,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetAllGroups ...
// returns a list of all groups
func GetAllGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		groups, err := groups.GetAllGroups(db)
		if err == nil {
			response = "Fetched all groups"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: groups,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetGroup ...
// returns a group by id
func GetGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		group, err := groups.GetGroupByID(db, id)
		if err == nil && group.ID != "" {
			response = "Fetched the groups"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: group,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetUserConfirms ...
// returns a list of account confirms
func GetUserConfirms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		userIDSelector := r.FormValue("userId")
		userCreationSecretSelector := types.UserCreationSecretSelector{
			UserID: userIDSelector,
		}

		creationSecrets, err := users.GetAllUserCreationSecrets(db, userCreationSecretSelector)
		if err == nil {
			response = "Fetched the user account creation secrets"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			List: creationSecrets,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetUserConfirm ...
// returns an account confirm by id
func GetUserConfirm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.ID != "" {
			response = "Fetched the user account creation secret"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: creationSecret,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetUserConfirmValid ...
// returns if an account confirm is valid by id
func GetUserConfirmValid(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.ID != "" {
			response = "Fetched the user account creation secret"
			code = http.StatusOK
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: creationSecret.ID != "",
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostUserConfirm ...
// confirm a user account
func PostUserConfirm(db *sql.DB) http.HandlerFunc {
	userManager := users.UserManager{DB: db}
	return func(w http.ResponseWriter, r *http.Request) {
		var context string
		var response string
		var code int

		vars := mux.Vars(r)
		id := vars["id"]

		secret := r.FormValue("secret")

		var user types.UserSpec
		body, errUnmarshal := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &user); err != nil {
			JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "failed to read request body",
				},
			})
			return
		}

		tokenString, err := userManager.ConfirmUserAccount(id, secret, user)
		if err == nil && errUnmarshal == nil {
			response = "Your user account has been confirmed"
			code = http.StatusOK
		} else {
			response = err.Error()
		}
		log.Println(response, context)
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: tokenString,
		}
		JSONResponse(r, w, code, JSONresp)

	}
}

// GetVersion ...
// returns version information about the instance
func GetVersion(w http.ResponseWriter, r *http.Request) {
	var response string
	version := common.GetAppBuildVersion()
	commitHash := common.GetAppBuildHash()
	mode := common.GetAppBuildMode()
	date := common.GetAppBuildDate()
	golangVersion := runtime.Version()
	osType := runtime.GOOS
	osArch := runtime.GOARCH

	log.Println(response)
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: response,
		},
		Data: types.SystemVersion{
			Version:       version,
			CommitHash:    commitHash,
			Mode:          mode,
			Date:          date,
			GolangVersion: golangVersion,
			OSType:        osType,
			OSArch:        osArch,
		},
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetServeFilestoreObjects ...
// serves files or 404
func (router Router) GetServeFilestoreObjects(prefix string) http.HandlerFunc {
	if router.FileAccess.Client == nil {
		return HTTP404()
	}
	return func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		defer close(doneCh)

		path := strings.TrimPrefix(r.URL.Path, prefix+"/")
		object, objectInfo, err := router.FileAccess.Get(path)
		if objectInfo.Size == 0 {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte("File not found")); err != nil {
				log.Printf("error: failed to write response; %v", err)
			}
			return
		}
		if err != nil {
			log.Printf("%#v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("An error occurred with retrieving the requested object")); err != nil {
				log.Printf("error: failed to write response; %v\n", err)
			}
			return
		}
		log.Println(objectInfo.Key, objectInfo.Size, objectInfo.ContentType)
		w.Header().Set("content-length", fmt.Sprintf("%d", objectInfo.Size))
		w.Header().Set("content-type", objectInfo.ContentType)
		w.Header().Set("accept-ranges", "bytes")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(object); err != nil {
			log.Printf("error: failed to write response; %v\n", err)
		}
	}
}

// Root ...
// /api endpoint
func Root(w http.ResponseWriter, r *http.Request) {
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hey! you're talking to the Flattrack API",
		},
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UnknownEndpoint ...
// response for hitting an unknown endpoint
func UnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}

// Healthz ...
// HTTP handler for health checks
func Healthz(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response string
		var code int

		err := health.Healthy(db)
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
		JSONResponse(r, w, code, JSONresp)
	}
}

// HTTPvalidateJWT ...
// middleware for checking JWT auth token validity
func HTTPvalidateJWT(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var context string
			completed, claims, err := users.ValidateJWTauthToken(db, r)
			if completed && err == nil {
				h.ServeHTTP(w, r)
				return
			}
			if claims.ID != "" {
				context = fmt.Sprintf("with user id '%v'", claims.ID)
			}
			log.Printf("Unauthorized request with token %v\n", context)
			JSONResponse(r, w, http.StatusUnauthorized, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unauthorized",
				},
			})
		}
	}
}

// HTTPcheckGroupsFromID ...
// middleware for checking if a route can be accessed given a ID and groupID
func HTTPcheckGroupsFromID(db *sql.DB, groupsAllowed ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id, errID := users.GetIDFromJWT(db, r)
			for _, group := range groupsAllowed {
				if userInGroup, err := groups.CheckUserInGroup(db, id, group); userInGroup && err == nil && err == errID {
					h.ServeHTTP(w, r)
					return
				}
			}
			log.Printf("User '%v' tried to access route that is protected by '%v' group access", id, groupsAllowed)
			JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Forbidden",
				},
			})
		}
	}
}

// HTTP404 ...
// responds with 404
func HTTP404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`404 not found`)); err != nil {
			log.Printf("error: failed to write repsonse; %v\n", err)
		}
	}
}
