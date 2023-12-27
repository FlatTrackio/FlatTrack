package httpserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// HTTPuseMiddleware ...
// append functions to run before the endpoint handler
func httpUseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

// FrontendOptions ...
// options to send to the frontend index.html for templating
type frontendOptions struct {
	SetupMessage string
	LoginMessage string
	EmbeddedHTML template.HTML
}

// FrontendHandler ...
// handles rewriting and API root setting
func frontendHandler(publicDir string, passthrough *frontendOptions) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.Replace(req.URL.Path, "/", "", 1)
		if len(req.URL.Path) > 0 && req.URL.Path[len(req.URL.Path)-1:] != "/" {
			req.URL.Path = path.Join("/", req.URL.Path)
		}

		// TODO redirect to subPath + /unknown-page if _path does not include subPath at the front

		// static files
		if strings.Contains(req.URL.Path, ".") {
			handler.ServeHTTP(w, req)
			return
		}

		// frontend views
		indexPath := path.Join(publicDir, "/index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, passthrough); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (h *HTTPServer) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	selectors := types.UserSelector{}
	if userSelectorID := r.FormValue("id"); userSelectorID != "" {
		selectors.ID = userSelectorID
	}
	if userSelectorNotID := r.FormValue("notId"); userSelectorNotID != "" {
		selectors.NotID = userSelectorNotID
	}
	if userSelectorGroup := r.FormValue("group"); userSelectorGroup != "" {
		selectors.Group = userSelectorGroup
	}
	if r.FormValue("notSelf") == "true" {
		id, err := h.users.GetIDFromJWT(r)
		if err != nil {
			log.Printf("error getting id from jwt: %v\n", err)
			JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unable to get user id from token",
				},
			})
			return
		}
		selectors.NotID = id
	}

	users, err := h.users.GetAllUsers(false, selectors)
	if err != nil {
		log.Printf("error getting all users: %v\n", err)
		JSONResponse(r, w, http.StatusInternalServerError, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Error getting a list of all users",
			},
		})
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Fetched user accounts",
		},
		List: users,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUser ...
// get a user by id or email (whatever is provided in the given respective order)
func (h *HTTPServer) GetUser(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.users.GetUserByID(id, false)
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

// PostUser ...
// create a user
func (h *HTTPServer) PostUser(w http.ResponseWriter, r *http.Request) {
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

	userAccount, err := h.users.CreateUser(user, user.Password == "")
	if err == nil && userAccount.ID != "" {
		code = http.StatusCreated
		response = "Created user account"
		context = fmt.Sprintf("'%v'", userAccount.ID)
	} else {
		code = http.StatusBadRequest
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

// PutUser ...
// updates a user account by their id
func (h *HTTPServer) PutUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.users.GetUserByID(userID, false)
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
	userAccountUpdated, err := h.users.UpdateProfileAdmin(userID, userAccount)
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

// PatchUser ...
// patches a user account by their id
func (h *HTTPServer) PatchUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.users.GetUserByID(userID, false)
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
	userAccountPatched, err := h.users.PatchProfileAdmin(userID, userAccount)
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

// PatchUserDisabled ...
// patches a user account's disabled field by their id
func (h *HTTPServer) PatchUserDisabled(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.users.GetUserByID(userID, false)
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

	jwtID, errID := h.users.GetIDFromJWT(r)
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

	isAdmin, errGroup := h.groups.CheckUserInGroup(jwtID, "admin")
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

	userAccountPatched, err := h.users.PatchUserDisabledAdmin(userID, userAccount.Disabled)
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

// DeleteUser ...
// delete a user
func (h *HTTPServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	vars := mux.Vars(r)
	userID := vars["id"]

	userInDB, err := h.users.GetUserByID(userID, false)
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

	myUserID, errID := h.users.GetIDFromJWT(r)
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

	err = h.users.DeleteUserByID(userInDB.ID)
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

// GetProfile ...
// returns the authenticated user's profile
func (h *HTTPServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	user, err := h.users.GetProfile(r)
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

// PutProfile ...
// Update a user account their id from their JWT
func (h *HTTPServer) PutProfile(w http.ResponseWriter, r *http.Request) {
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

	id, errID := h.users.GetIDFromJWT(r)
	userAccountUpdated, err := h.users.UpdateProfile(id, userAccount)
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

// PatchProfile ...
// patches a user account their id from their JWT
func (h *HTTPServer) PatchProfile(w http.ResponseWriter, r *http.Request) {
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

	id, errID := h.users.GetIDFromJWT(r)
	userAccountPatched, err := h.users.PatchProfile(id, userAccount)
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

// GetSystemInitialized ...
// check if the server has been initialized
func (h *HTTPServer) GetSystemInitialized(w http.ResponseWriter, r *http.Request) {
	var context string
	code := http.StatusInternalServerError
	response := "Failed to fetch if this FlatTrack instance has initialized"

	initialized, err := h.system.GetHasInitialized()
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

// UserAuth ...
// authenticate a user
func (h *HTTPServer) UserAuth(w http.ResponseWriter, r *http.Request) {
	var user types.UserSpec
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: faled to read request body; %v\n", err)
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

	userInDB, err := h.users.GetUserByEmail(user.Email, false)
	if err != nil {
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Failed to get user account",
			},
		})
		return
	}
	if userInDB.ID != "" && !userInDB.Registered {
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "User account is not yet registered",
			},
		})
		return
	}
	if userInDB.Disabled {
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "User account has been disabled",
			},
		})
		return
	}
	// Check password locally, fall back to remote if incorrect
	matches, err := h.users.CheckUserPassword(userInDB.Email, user.Password)
	if err != nil {
		log.Printf("error checking password: %v\n", err)
		JSONResponse(r, w, http.StatusInternalServerError, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Failed to check user account password",
			},
		})
		return
	}
	if !matches {
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Unable to authenticate",
			},
		})
		return
	}
	jwtToken, err := h.users.GenerateJWTauthToken(userInDB.ID, userInDB.AuthNonce, 0)
	if err != nil {
		log.Printf("error checking password: %v\n", err)
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Successfully authenticated user",
			},
			Data: jwtToken,
		})
		return
	}
	JSONResponse(r, w, http.StatusOK, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Successfully authenticated user",
		},
		Data: jwtToken,
	})
}

// UserAuthValidate ...
// validate an auth token
func (h *HTTPServer) UserAuthValidate(w http.ResponseWriter, r *http.Request) {
	var context string
	var response string
	code := http.StatusUnauthorized

	valid, claims, err := h.users.ValidateJWTauthToken(r)
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

// UserAuthReset ...
// invalidates all JWTs
func (h *HTTPServer) UserAuthReset(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to sign out all devices logged in"
	code := http.StatusInternalServerError

	id, err := h.users.GetIDFromJWT(r)
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
	err = h.users.GenerateNewAuthNonce(id)
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

// UserCanIgroup ...
// respond whether the current user account is in a group
func (h *HTTPServer) UserCanIgroup(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to determine group privileges"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	groupName := vars["name"]

	group, err := h.groups.GetGroupByName(groupName)
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

	id, errID := h.users.GetIDFromJWT(r)
	userIsInGroup, err := h.groups.CheckUserInGroup(id, groupName)
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

// GetSettingsFlatName ...
// responds with the name of the flat
func (h *HTTPServer) GetSettingsFlatName(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch the flat name"
	code := http.StatusInternalServerError

	flatName, err := h.settings.GetFlatName()
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

// SetSettingsFlatName ...
// update the flat's name
func (h *HTTPServer) SetSettingsFlatName(w http.ResponseWriter, r *http.Request) {
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

	err := h.settings.SetFlatName(flatName.FlatName)
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

// PostAdminRegister ...
// register the instance of FlatTrack
func (h *HTTPServer) PostAdminRegister(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	initialized, err := h.system.GetHasInitialized()
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

	registered, jwt, err := h.registration.Register(registrationForm)
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

// GetShoppingList ...
// responds with list of shopping lists
func (h *HTTPServer) GetShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping lists"
	code := http.StatusNotFound

	vars := mux.Vars(r)
	id := vars["id"]

	shoppingList, err := h.shoppinglist.ShoppingList().GetShoppingList(id)
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

// GetShoppingLists ...
// responds with shopping list by id
func (h *HTTPServer) GetShoppingLists(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping lists"
	code := http.StatusInternalServerError

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

	shoppingLists, err := h.shoppinglist.ShoppingList().GetShoppingLists(options)
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

// PostShoppingList ...
// creates a new shopping list to add items to
func (h *HTTPServer) PostShoppingList(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		log.Printf("error: %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get userid from token",
			},
		})
		return
	}
	shoppingList.Author = id
	shoppingListInserted, err := h.shoppinglist.ShoppingList().CreateShoppingList(shoppingList, options)
	if err == nil {
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

// PatchShoppingList ...
// patches an existing shopping list
func (h *HTTPServer) PatchShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string
	code := http.StatusBadRequest
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	id, errID := h.users.GetIDFromJWT(r)
	shoppingList.AuthorLast = id
	shoppingListPatched, err := h.shoppinglist.ShoppingList().PatchShoppingList(listID, shoppingList)
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

// PutShoppingList ...
// updates an existing shopping list
func (h *HTTPServer) PutShoppingList(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	id, errID := h.users.GetIDFromJWT(r)
	shoppingList.AuthorLast = id
	shoppingListUpdated, err := h.shoppinglist.ShoppingList().UpdateShoppingList(listID, shoppingList)
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

// DeleteShoppingList ...
// delete a new shopping list by it's id
func (h *HTTPServer) DeleteShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	err = h.shoppinglist.ShoppingList().DeleteShoppingList(listID)
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

// GetShoppingListItems ...
// responds with shopping items by list id
func (h *HTTPServer) GetShoppingListItems(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping list items"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	id := vars["id"]

	options := types.ShoppingItemOptions{
		SortBy: r.FormValue("sortBy"),
		Selector: types.ShoppingItemSelector{
			Obtained: r.FormValue("obtained"),
		},
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(id)
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
	shoppingListItems, err := h.shoppinglist.ShoppingItem().GetShoppingListItems(id, options)
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

// GetShoppingListItem ...
// responds with list of shopping lists
func (h *HTTPServer) GetShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping list item"
	code := http.StatusNotFound

	vars := mux.Vars(r)
	itemID := vars["itemId"]
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	shoppingListItem, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(listID, itemID)
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

// PostItemToShoppingList ...
// adds an item to a shopping list
func (h *HTTPServer) PostItemToShoppingList(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	id, errID := h.users.GetIDFromJWT(r)
	shoppingItem.Author = id
	shoppingItemInserted, err := h.shoppinglist.ShoppingItem().AddItemToList(listID, shoppingItem)
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

// PatchShoppingListCompleted ...
// adds an item to a shopping list
func (h *HTTPServer) PatchShoppingListCompleted(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	id, errID := h.users.GetIDFromJWT(r)
	patchedList, err := h.shoppinglist.ShoppingList().SetListCompleted(listID, shoppingList.Completed, id)
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

// PatchShoppingListItem ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItem(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(listID, itemID)
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

	id, errID := h.users.GetIDFromJWT(r)
	shoppingItem.AuthorLast = id
	patchedItem, err := h.shoppinglist.ShoppingItem().PatchItem(listID, itemID, shoppingItem)
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

// PutShoppingListItem ...
// updates an item in a shopping list
func (h *HTTPServer) PutShoppingListItem(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(listID, itemID)
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

	id, errID := h.users.GetIDFromJWT(r)
	shoppingItem.AuthorLast = id
	updatedItem, err := h.shoppinglist.ShoppingItem().UpdateItem(listID, itemID, shoppingItem)
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

// PatchShoppingListItemObtained ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItemObtained(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(listID, itemID)
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

	id, errID := h.users.GetIDFromJWT(r)
	patchedItem, err := h.shoppinglist.ShoppingItem().SetItemObtained(listID, itemID, shoppingItem.Obtained, id)
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

// DeleteShoppingListItem ...
// delete a shopping list item by it's id
func (h *HTTPServer) DeleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	vars := mux.Vars(r)
	itemID := vars["itemId"]
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(listID, itemID)
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

	err = h.shoppinglist.ShoppingItem().RemoveItemFromList(itemID, listID)
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

// GetShoppingListItemTags ...
// responds with tags used in shopping list items from a list
func (h *HTTPServer) GetShoppingListItemTags(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping list item tags"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	listID := vars["listId"]

	tags, err := h.shoppinglist.ShoppingTag().GetShoppingListTags(listID)
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

// UpdateShoppingListItemTag ...
// updates then tag name used in shopping list items from a list
func (h *HTTPServer) UpdateShoppingListItemTag(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to update shopping list item tag name"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	listID := vars["listId"]
	tag := vars["tagName"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
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

	tag, err = h.shoppinglist.ShoppingTag().UpdateShoppingListTag(listID, tag, tagUpdate.Name)
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

// GetAllShoppingTags ...
// responds with all tags used in shopping list items
func (h *HTTPServer) GetAllShoppingTags(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch shopping list item tags"
	code := http.StatusInternalServerError

	options := types.ShoppingTagOptions{
		SortBy: r.FormValue("sortBy"),
	}

	tags, err := h.shoppinglist.ShoppingTag().GetAllShoppingTags(options)
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

// PostShoppingTag ...
// creates a tag name
func (h *HTTPServer) PostShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	var response string
	code := http.StatusInternalServerError

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

	id, errID := h.users.GetIDFromJWT(r)
	tag.Author = id
	tag, err := h.shoppinglist.ShoppingTag().CreateShoppingTag(tag)
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

// GetShoppingTag ...
// gets a shopping tag by id
func (h *HTTPServer) GetShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch a shopping tag"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	id := vars["id"]

	tag, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
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

// UpdateShoppingTag ...
// updates a tag name
func (h *HTTPServer) UpdateShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to update a shopping tag"
	code := http.StatusInternalServerError

	vars := mux.Vars(r)
	id := vars["id"]

	tagInDB, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
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

	userID, errID := h.users.GetIDFromJWT(r)
	tagUpdate.AuthorLast = userID
	tag, err := h.shoppinglist.ShoppingTag().UpdateShoppingTag(id, tagUpdate)
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

// DeleteShoppingTag ...
// deletes a shopping tag by id
func (h *HTTPServer) DeleteShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	var code int
	var response string

	vars := mux.Vars(r)
	id := vars["id"]

	tagInDB, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
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

	err = h.shoppinglist.ShoppingTag().DeleteShoppingTag(id)
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

// GetSettingsShoppingListNotes ...
// responds with the notes for shopping lists
func (h *HTTPServer) GetSettingsShoppingListNotes(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch the shopping list notes"
	code := http.StatusInternalServerError

	notes, err := h.settings.GetShoppingListNotes()
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

// PutSettingsShoppingList ...
// update the notes for shopping lists
func (h *HTTPServer) PutSettingsShoppingList(w http.ResponseWriter, r *http.Request) {
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

	err := h.settings.SetShoppingListNotes(notes.Notes)
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

// GetSettingsFlatNotes ...
// responds with the notes for flat
func (h *HTTPServer) GetSettingsFlatNotes(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch the flat notes"
	code := http.StatusInternalServerError

	notes, err := h.settings.GetFlatNotes()
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

// PutSettingsFlatNotes ...
// update the notes for flat
func (h *HTTPServer) PutSettingsFlatNotes(w http.ResponseWriter, r *http.Request) {
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

	err := h.settings.SetFlatNotes(notes.Notes)
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

// GetAllGroups ...
// returns a list of all groups
func (h *HTTPServer) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch groups"
	code := http.StatusInternalServerError

	groups, err := h.groups.GetAllGroups()
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

// GetGroup ...
// returns a group by id
func (h *HTTPServer) GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	group, err := h.groups.GetGroupByID(id)
	if err != nil {
		log.Printf("error: failed to fetch group (%v); %v\n", id, err)
		JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Failed to fetch a group",
			},
		})
		return
	}
	if group.ID == "" {
		log.Printf("error: failed to fetch group (%v); %v\n", id, err)
		JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Group not found",
			},
		})
		return
	}
	JSONResponse(r, w, http.StatusOK, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Fetched the group",
		},
		Spec: group,
	})
}

// GetUserConfirms ...
// returns a list of account confirms
func (h *HTTPServer) GetUserConfirms(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch user account creation secrets"
	code := http.StatusInternalServerError

	userIDSelector := r.FormValue("userId")
	userCreationSecretSelector := types.UserCreationSecretSelector{
		UserID: userIDSelector,
	}

	creationSecrets, err := h.users.GetAllUserCreationSecrets(userCreationSecretSelector)
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

// GetUserConfirm ...
// returns an account confirm by id
func (h *HTTPServer) GetUserConfirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.GetUserCreationSecret(id)
	if err != nil {
		log.Printf("error getting user creation secret: %v\n", err)
		JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Failed to get user creation secret",
			},
		})
		return
	}
	if creationSecret.ID == "" {
		JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Failed to get user creation secret",
			},
		})
		return
	}
	JSONResponse(r, w, http.StatusOK, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Fetched the user account creation secret",
		},
		Spec: creationSecret,
	})
}

// GetUserConfirmValid ...
// returns if an account confirm is valid by id
func (h *HTTPServer) GetUserConfirmValid(w http.ResponseWriter, r *http.Request) {
	var context string
	response := "Failed to fetch user account creation secret"
	code := http.StatusNotFound

	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.GetUserCreationSecret(id)
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

// PostUserConfirm ...
// confirm a user account
func (h *HTTPServer) PostUserConfirm(w http.ResponseWriter, r *http.Request) {
	var context string
	var response string
	code := http.StatusInternalServerError

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

	tokenString, err := h.users.ConfirmUserAccount(id, secret, user)
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

// GetVersion ...
// returns version information about the instance
func (h *HTTPServer) GetVersion(w http.ResponseWriter, r *http.Request) {
	response := "Fetched version information"
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
// func (router Router) GetServeFilestoreObjects(prefix string) http.HandlerFunc {
// 	if router.FileAccess.Client == nil {
// 		return HTTP404()
// 	}
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		doneCh := make(chan struct{})
// 		defer close(doneCh)

// 		path := strings.TrimPrefix(r.URL.Path, prefix+"/")
// 		object, objectInfo, err := router.FileAccess.Get(path)
// 		if objectInfo.Size == 0 {
// 			w.WriteHeader(http.StatusNotFound)
// 			if _, err := w.Write([]byte("File not found")); err != nil {
// 				log.Printf("error: failed to write response; %v", err)
// 			}
// 			return
// 		}
// 		if err != nil {
// 			log.Printf("%#v\n", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			if _, err := w.Write([]byte("An error occurred with retrieving the requested object")); err != nil {
// 				log.Printf("error: failed to write response; %v\n", err)
// 			}
// 			return
// 		}
// 		log.Println(objectInfo.Key, objectInfo.Size, objectInfo.ContentType)
// 		w.Header().Set("content-length", fmt.Sprintf("%d", objectInfo.Size))
// 		w.Header().Set("content-type", objectInfo.ContentType)
// 		w.Header().Set("accept-ranges", "bytes")
// 		w.WriteHeader(http.StatusOK)
// 		if _, err := w.Write(object); err != nil {
// 			log.Printf("error: failed to write response; %v\n", err)
// 		}
// 	}
// }

// Root ...
// /api endpoint
func (h *HTTPServer) Root(w http.ResponseWriter, r *http.Request) {
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hey! you're talking to the Flattrack API",
		},
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UnknownEndpoint ...
// response for hitting an unknown endpoint
func (h *HTTPServer) UnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	JSONResponse(r, w, http.StatusNotFound, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}

// Healthz ...
// HTTP handler for health checks
func (h *HTTPServer) Healthz(w http.ResponseWriter, r *http.Request) {
	response := "App unhealthy"
	code := http.StatusInternalServerError

	err := h.health.Healthy()
	if err != nil {
		log.Printf("error app unhealth: %v\n", err)
	}
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

// HTTPvalidateJWT ...
// middleware for checking JWT auth token validity
func (h *HTTPServer) HTTPvalidateJWT() func(http.HandlerFunc) http.HandlerFunc {
	return func(ha http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var context string
			completed, claims, err := h.users.ValidateJWTauthToken(r)
			if completed && err == nil {
				ha.ServeHTTP(w, r)
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
func (h *HTTPServer) HTTPcheckGroupsFromID(groupsAllowed ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(ha http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id, errID := h.users.GetIDFromJWT(r)
			for _, group := range groupsAllowed {
				if userInGroup, err := h.groups.CheckUserInGroup(id, group); userInGroup && err == nil && err == errID {
					ha.ServeHTTP(w, r)
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
func (h *HTTPServer) HTTP404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`404 not found`)); err != nil {
			log.Printf("error: failed to write repsonse; %v\n", err)
		}
	}
}

func (h *HTTPServer) registerAPIHandlers(router *mux.Router) {
	routes := []struct {
		EndpointPath string
		HandlerFunc  http.HandlerFunc
		HTTPMethod   string
	}{
		{
			EndpointPath: "",
			HandlerFunc:  h.Root,
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/system/initialized",
			HandlerFunc:  h.GetSystemInitialized,
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/system/version",
			HandlerFunc:  httpUseMiddleware(h.GetVersion, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/system/flatName",
			HandlerFunc:  httpUseMiddleware(h.GetSettingsFlatName, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/settings/flatName",
			HandlerFunc:  httpUseMiddleware(h.SetSettingsFlatName, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/admin/settings/shoppingListNotes",
			HandlerFunc:  httpUseMiddleware(h.PutSettingsShoppingList, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/admin/settings/flatNotes",
			HandlerFunc:  httpUseMiddleware(h.GetSettingsFlatNotes, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/settings/flatNotes",
			HandlerFunc:  httpUseMiddleware(h.PutSettingsFlatNotes, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/admin/register",
			HandlerFunc:  httpUseMiddleware(h.PostAdminRegister),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/admin/users",
			HandlerFunc:  httpUseMiddleware(h.GetAllUsers, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetUser, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/users",
			HandlerFunc:  httpUseMiddleware(h.PostUser, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  httpUseMiddleware(h.PatchUser, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/admin/users/{id}/disabled",
			HandlerFunc:  httpUseMiddleware(h.PatchUserDisabled, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  httpUseMiddleware(h.PutUser, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  httpUseMiddleware(h.DeleteUser, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/admin/useraccountconfirms",
			HandlerFunc:  httpUseMiddleware(h.GetUserConfirms, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/useraccountconfirms/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetUserConfirm, h.HTTPvalidateJWT(), h.HTTPcheckGroupsFromID("admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/auth",
			HandlerFunc:  h.UserAuthValidate,
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/auth",
			HandlerFunc:  h.UserAuth,
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/auth/reset",
			HandlerFunc:  h.UserAuthReset,
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/confirm/{id}",
			HandlerFunc:  h.GetUserConfirmValid,
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/confirm/{id}",
			HandlerFunc:  h.PostUserConfirm,
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  httpUseMiddleware(h.GetProfile, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  httpUseMiddleware(h.PutProfile, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  httpUseMiddleware(h.PatchProfile, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/users",
			HandlerFunc:  httpUseMiddleware(h.GetAllUsers, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/users/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetUser, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/groups",
			HandlerFunc:  httpUseMiddleware(h.GetAllGroups, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/groups/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetGroup, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/can-i/group/{name}",
			HandlerFunc:  httpUseMiddleware(h.UserCanIgroup, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/settings/notes",
			HandlerFunc:  httpUseMiddleware(h.GetSettingsShoppingListNotes, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingLists, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  httpUseMiddleware(h.PatchShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  httpUseMiddleware(h.PutShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/completed",
			HandlerFunc:  httpUseMiddleware(h.PatchShoppingListCompleted, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  httpUseMiddleware(h.DeleteShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  httpUseMiddleware(h.PostShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingListItems, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingListItem, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  httpUseMiddleware(h.PostItemToShoppingList, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  httpUseMiddleware(h.PatchShoppingListItem, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  httpUseMiddleware(h.PutShoppingListItem, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}/obtained",
			HandlerFunc:  httpUseMiddleware(h.PatchShoppingListItemObtained, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  httpUseMiddleware(h.DeleteShoppingListItem, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingListItemTags, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags/{tagName}",
			HandlerFunc:  httpUseMiddleware(h.UpdateShoppingListItemTag, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  httpUseMiddleware(h.PostShoppingTag, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  httpUseMiddleware(h.GetAllShoppingTags, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  httpUseMiddleware(h.GetShoppingTag, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  httpUseMiddleware(h.UpdateShoppingTag, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  httpUseMiddleware(h.DeleteShoppingTag, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/flat/info",
			HandlerFunc:  httpUseMiddleware(h.GetSettingsFlatNotes, h.HTTPvalidateJWT()),
			HTTPMethod:   http.MethodGet,
		},
	}
	for _, r := range routes {
		router.HandleFunc(r.EndpointPath, r.HandlerFunc).Methods(r.HTTPMethod, http.MethodOptions)
	}
}
