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

const (
	//nolint:gosec
	FlatTrackSchedulerSecretHeader = "X-FlatTrack-Scheduler-Secret"
)

// HTTPuseMiddleware ...
// append functions to run before the endpoint handler
func httpUseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

// TODO handle io.ReadAll err

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
					Response: "failed to get user account id from token",
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
				Response: "failed to get a list of all users",
			},
		})
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user accounts",
		},
		List: users,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUser ...
// get a user by id or email (whatever is provided in the given respective order)
func (h *HTTPServer) GetUser(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.users.GetUserByID(id, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user",
			},
			Spec: types.UserSpec{},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if user.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user",
			},
			Spec: types.UserSpec{},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user account",
		},
		Spec: user,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostUser ...
// create a user
func (h *HTTPServer) PostUser(w http.ResponseWriter, r *http.Request) {
	var context string

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
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create user account",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	context = fmt.Sprintf("'%v'", userAccount.ID)
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "created user account",
		},
		Spec: userAccount,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PutUser ...
// updates a user account by their id
func (h *HTTPServer) PutUser(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
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

	_, err = h.users.GetUserByID(userID, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}

	// TODO disallow admins to remove their own admin group access
	userAccountUpdated, err := h.users.UpdateProfileAdmin(userID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated user account",
		},
		Spec: userAccountUpdated,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchUser ...
// patches a user account by their id
func (h *HTTPServer) PatchUser(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
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

	_, err = h.users.GetUserByID(userID, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}

	// TODO disallow admins to remove their own admin group access
	userAccountPatched, err := h.users.PatchProfileAdmin(userID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched user account",
		},
		Spec: userAccountPatched,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchUserDisabled ...
// patches a user account's disabled field by their id
func (h *HTTPServer) PatchUserDisabled(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
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

	_, err = h.users.GetUserByID(userID, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	jwtID, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	isAdmin, err := h.groups.CheckUserInGroup(jwtID, "admin")
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to check user in group",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	if isAdmin && userID == jwtID {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to disable user account of invoker",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusForbidden, JSONresp)
		return
	}

	userAccountPatched, err := h.users.PatchUserDisabledAdmin(userID, userAccount.Disabled)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "disabled user account",
		},
		Spec: userAccountPatched,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteUser ...
// delete a user
func (h *HTTPServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	userID := vars["id"]

	userInDB, err := h.users.GetUserByID(userID, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	myUserID, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	if myUserID == userID {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to delete user account of invoker",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusForbidden, JSONresp)
		return
	}

	if err := h.users.DeleteUserByID(userInDB.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted user account",
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetProfile ...
// returns the authenticated user's profile
func (h *HTTPServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	var context string

	user, err := h.users.GetProfile(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if user.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user",
			},
			Spec: types.UserSpec{},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched profile",
		},
		Spec: user,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutProfile ...
// Update a user account their id from their JWT
func (h *HTTPServer) PutProfile(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
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
	if err := json.Unmarshal(body, &userAccount); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	userAccountUpdated, err := h.users.UpdateProfile(id, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update user account",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated user account",
		},
		Spec: userAccountUpdated,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchProfile ...
// patches a user account their id from their JWT
func (h *HTTPServer) PatchProfile(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
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
	if err := json.Unmarshal(body, &userAccount); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	userAccountPatched, err := h.users.PatchProfile(id, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched user account",
		},
		Spec: userAccountPatched,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetSystemInitialized ...
// check if the server has been initialized
func (h *HTTPServer) GetSystemInitialized(w http.ResponseWriter, r *http.Request) {
	var context string

	initialised, err := h.system.GetHasInitialized()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get system initialise status",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	response := "not initialised"
	if initialised {
		response = "initialised"
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: response,
		},
		Data: initialised,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
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

	valid, claims, err := h.users.ValidateJWTauthToken(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to validate auth token",
			},
			Data: false,
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusUnauthorized, JSONresp)
		return
	}
	context = fmt.Sprintf("for user with ID '%v'", claims.ID)
	if !valid {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "auth token is not valid",
			},
			Data: valid,
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusUnauthorized, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "auth token is valid",
		},
		Data: valid,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UserAuthReset ...
// invalidates all JWTs
func (h *HTTPServer) UserAuthReset(w http.ResponseWriter, r *http.Request) {
	var context string

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user account with id: " + id,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if err := h.users.GenerateNewAuthNonce(id); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user account with id: " + id,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "reset all authentication tokens",
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UserCanIgroup ...
// respond whether the current user account is in a group
func (h *HTTPServer) UserCanIgroup(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	groupName := vars["name"]

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	group, err := h.groups.GetGroupByName(groupName)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group by name: " + groupName,
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	userIsInGroup, err := h.groups.CheckUserInGroup(id, group.Name)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to check whether user is in group",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user is in group",
		},
		Data: userIsInGroup,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetSettingsFlatName ...
// responds with the name of the flat
func (h *HTTPServer) GetSettingsFlatName(w http.ResponseWriter, r *http.Request) {
	var context string
	flatName, err := h.settings.GetFlatName()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get flat name setting",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if flatName == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "flat name is not set",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched flat name",
		},
		Spec: flatName,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// SetSettingsFlatName ...
// update the flat's name
func (h *HTTPServer) SetSettingsFlatName(w http.ResponseWriter, r *http.Request) {
	var context string

	var flatName types.FlatName
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
	if err := json.Unmarshal(body, &flatName); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	if err := h.settings.SetFlatName(flatName.FlatName); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to set flat name setting",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set flat name",
		},
		Spec: true,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostAdminRegister ...
// register the instance of FlatTrack
func (h *HTTPServer) PostAdminRegister(w http.ResponseWriter, r *http.Request) {
	var context string

	initialized, err := h.system.GetHasInitialized()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get system initialised status",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if initialized {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "system is initialised",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusOK, JSONresp)
		return
	}

	var registrationForm types.Registration
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
	if err := json.Unmarshal(body, &registrationForm); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}
	registrationForm.InstanceIDConfirm = r.FormValue("instanceIDConfirm")

	registered, jwt, err := h.registration.Register(registrationForm)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to register instance",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "registered",
		},
		Spec: registered,
		Data: jwt,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// GetShoppingList ...
// responds with list of shopping lists
func (h *HTTPServer) GetShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	shoppingList, err := h.shoppinglist.ShoppingList().GetShoppingList(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping list",
		},
		Spec: shoppingList,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingLists ...
// responds with shopping list by id
func (h *HTTPServer) GetShoppingLists(w http.ResponseWriter, r *http.Request) {
	var context string
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
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping lists",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping lists",
		},
		List: shoppingLists,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostShoppingList ...
// creates a new shopping list to add items to
func (h *HTTPServer) PostShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingList types.ShoppingListSpec
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
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	shoppingList.Author = id
	shoppingListInserted, err := h.shoppinglist.ShoppingList().CreateShoppingList(shoppingList, options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "created shopping list",
		},
		Spec: shoppingListInserted,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PatchShoppingList ...
// patches an existing shopping list
func (h *HTTPServer) PatchShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingList types.ShoppingListSpec
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
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	shoppingList.AuthorLast = id
	shoppingListPatched, err := h.shoppinglist.ShoppingList().PatchShoppingList(list.ID, shoppingList)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched shopping list",
		},
		Spec: shoppingListPatched,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutShoppingList ...
// updates an existing shopping list
func (h *HTTPServer) PutShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingList types.ShoppingListSpec
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	shoppingList.AuthorLast = id
	shoppingListUpdated, err := h.shoppinglist.ShoppingList().UpdateShoppingList(list.ID, shoppingList)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping list",
		},
		Spec: shoppingListUpdated,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteShoppingList ...
// delete a new shopping list by it's id
func (h *HTTPServer) DeleteShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingList().DeleteShoppingList(list.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to delete shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted shopping list",
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingListItems ...
// responds with shopping items by list id
func (h *HTTPServer) GetShoppingListItems(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	options := types.ShoppingItemOptions{
		SortBy: r.FormValue("sortBy"),
		Selector: types.ShoppingItemSelector{
			Obtained: r.FormValue("obtained"),
		},
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	// TODO add item selectors for this endpoint
	shoppingListItems, err := h.shoppinglist.ShoppingItem().GetShoppingListItems(list.ID, options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list items",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping list items",
		},
		List: shoppingListItems,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingListItem ...
// responds with list of shopping lists
func (h *HTTPServer) GetShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	itemID := vars["itemId"]
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	shoppingListItem, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if shoppingListItem.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetch shopping list item",
		},
		Spec: shoppingListItem,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostItemToShoppingList ...
// adds an item to a shopping list
func (h *HTTPServer) PostItemToShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingItem types.ShoppingItemSpec
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.Author = id
	shoppingItemInserted, err := h.shoppinglist.ShoppingItem().AddItemToList(list.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to add item to shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "added item to shopping list",
		},
		Spec: shoppingItemInserted,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PatchShoppingListCompleted ...
// adds an item to a shopping list
func (h *HTTPServer) PatchShoppingListCompleted(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingList types.ShoppingListSpec
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

	userID, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	patchedList, err := h.shoppinglist.ShoppingList().SetListCompleted(list.ID, shoppingList.Completed, userID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to set shopping list as completed",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "shopping list set as completed",
		},
		Spec: patchedList,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchShoppingListItem ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingItem types.ShoppingItemSpec
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.AuthorLast = id
	patchedItem, err := h.shoppinglist.ShoppingItem().PatchItem(listID, item.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched shopping list item",
		},
		Spec: patchedItem,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutShoppingListItem ...
// updates an item in a shopping list
func (h *HTTPServer) PutShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingItem types.ShoppingItemSpec
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.AuthorLast = id
	updatedItem, err := h.shoppinglist.ShoppingItem().UpdateItem(listID, item.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping list item",
		},
		Spec: updatedItem,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchShoppingListItemObtained ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItemObtained(w http.ResponseWriter, r *http.Request) {
	var context string

	var shoppingItem types.ShoppingItemSpec
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

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	patchedItem, err := h.shoppinglist.ShoppingItem().SetItemObtained(listID, item.ID, shoppingItem.Obtained, id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set shopping list item as obtained",
		},
		Spec: patchedItem,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteShoppingListItem ...
// delete a shopping list item by it's id
func (h *HTTPServer) DeleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	itemID := vars["itemId"]
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().GetShoppingListItem(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingItem().RemoveItemFromList(item.ID, list.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to remove item from shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "removed item from shopping list",
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingListItemTags ...
// responds with tags used in shopping list items from a list
func (h *HTTPServer) GetShoppingListItemTags(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	tags, err := h.shoppinglist.ShoppingTag().GetShoppingListTags(list.ID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get tags from shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched tags from shopping list",
		},
		List: tags,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UpdateShoppingListItemTag ...
// updates then tag name used in shopping list items from a list
func (h *HTTPServer) UpdateShoppingListItemTag(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	listID := vars["listId"]
	tag := vars["tagName"]

	list, err := h.shoppinglist.ShoppingList().GetShoppingList(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	var tagUpdate types.ShoppingTag
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
	if err := json.Unmarshal(body, &tagUpdate); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	tagUpdated, err := h.shoppinglist.ShoppingTag().UpdateShoppingListTag(list.ID, tag, tagUpdate.Name)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping list tag",
		},
		Spec: tagUpdated,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetAllShoppingTags ...
// responds with all tags used in shopping list items
func (h *HTTPServer) GetAllShoppingTags(w http.ResponseWriter, r *http.Request) {
	var context string
	options := types.ShoppingTagOptions{
		SortBy: r.FormValue("sortBy"),
	}

	tags, err := h.shoppinglist.ShoppingTag().GetAllShoppingTags(options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list tags",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping list tags",
		},
		List: tags,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostShoppingTag ...
// creates a tag name
func (h *HTTPServer) PostShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string

	var tag types.ShoppingTag
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
	if err := json.Unmarshal(body, &tag); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	id, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	tag.Author = id
	tagCreated, err := h.shoppinglist.ShoppingTag().CreateShoppingTag(tag)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "created shopping tag",
		},
		Spec: tagCreated,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// GetShoppingTag ...
// gets a shopping tag by id
func (h *HTTPServer) GetShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	tag, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping tag",
		},
		Spec: tag,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UpdateShoppingTag ...
// updates a tag name
func (h *HTTPServer) UpdateShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := h.users.GetIDFromJWT(r)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account from id",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if userID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account from id",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	tag, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	var tagUpdate types.ShoppingTag
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
	if err := json.Unmarshal(body, &tagUpdate); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}
	tagUpdate.AuthorLast = userID
	tagUpdated, err := h.shoppinglist.ShoppingTag().UpdateShoppingTag(tag.ID, tagUpdate)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping tag",
		},
		Spec: tagUpdated,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteShoppingTag ...
// deletes a shopping tag by id
func (h *HTTPServer) DeleteShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	id := vars["id"]

	tag, err := h.shoppinglist.ShoppingTag().GetShoppingTag(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingTag().DeleteShoppingTag(tag.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to delete shopping tag",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted shopping tag",
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetSettingsShoppingListNotes ...
// responds with the notes for shopping lists
func (h *HTTPServer) GetSettingsShoppingListNotes(w http.ResponseWriter, r *http.Request) {
	var context string
	notes, err := h.settings.GetShoppingListNotes()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping notes",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping notes",
		},
		Spec: notes,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutSettingsShoppingList ...
// update the notes for shopping lists
func (h *HTTPServer) PutSettingsShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var notes types.ShoppingListNotes
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
	if err := json.Unmarshal(body, &notes); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	if err := h.settings.SetShoppingListNotes(notes.Notes); err != nil {
		context = err.Error()
		code := http.StatusInternalServerError
		if err.Error() == "Unable to set shopping list notes as it is either invalid, too short, or too long" {
			code = http.StatusBadRequest
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping notes",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, code, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set shopping notes",
		},
		Spec: notes.Notes,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetSettingsFlatNotes ...
// responds with the notes for flat
func (h *HTTPServer) GetSettingsFlatNotes(w http.ResponseWriter, r *http.Request) {
	var context string
	notes, err := h.settings.GetFlatNotes()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get flat notes",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched flat notes",
		},
		Spec: types.FlatNotes{
			Notes: notes,
		},
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutSettingsFlatNotes ...
// update the notes for flat
func (h *HTTPServer) PutSettingsFlatNotes(w http.ResponseWriter, r *http.Request) {
	var context string

	var notes types.FlatNotes
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
	if err := json.Unmarshal(body, &notes); err != nil {
		log.Printf("error: failed to unmarshal; %v\n", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	if err := h.settings.SetFlatNotes(notes.Notes); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get flat notes",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set flat notes",
		},
		Spec: notes.Notes,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetAllGroups ...
// returns a list of all groups
func (h *HTTPServer) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var context string
	groups, err := h.groups.GetAllGroups()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get groups",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched groups",
		},
		List: groups,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetGroup ...
// returns a group by id
func (h *HTTPServer) GetGroup(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	id := vars["id"]

	group, err := h.groups.GetGroupByID(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if group.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONResponse(r, w, http.StatusOK, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched group",
		},
		Spec: group,
	})
}

// GetUserConfirms ...
// returns a list of account confirms
// TODO should this exist?
func (h *HTTPServer) GetUserConfirms(w http.ResponseWriter, r *http.Request) {
	var context string
	userIDSelector := r.FormValue("userId")
	userCreationSecretSelector := types.UserCreationSecretSelector{
		UserID: userIDSelector,
	}

	creationSecrets, err := h.users.GetAllUserCreationSecrets(userCreationSecretSelector)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secrets",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user creation secrets",
		},
		List: creationSecrets,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUserConfirm ...
// returns an account confirm by id
func (h *HTTPServer) GetUserConfirm(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.GetUserCreationSecret(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if creationSecret.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	log.Printf("%+v\n", creationSecret)
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user creation secret",
		},
		Spec: creationSecret,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUserConfirmValid ...
// returns if an account confirm is valid by id
// TODO should this exist?
func (h *HTTPServer) GetUserConfirmValid(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.GetUserCreationSecret(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if creationSecret.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return

	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user creation secret valid",
		},
		Data: creationSecret.ID != "" && creationSecret.Secret != "" && creationSecret.Valid,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostUserConfirm ...
// confirm a user account
func (h *HTTPServer) PostUserConfirm(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	id := vars["id"]

	secret := r.FormValue("secret")

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
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	tokenString, err := h.users.ConfirmUserAccount(id, secret, user)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to confirm user account",
			},
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "confirmed user account",
		},
		Data: tokenString,
	}
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// GetVersion ...
// returns version information about the instance
func (h *HTTPServer) GetVersion(w http.ResponseWriter, r *http.Request) {
	var context string

	version := common.GetAppBuildVersion()
	commitHash := common.GetAppBuildHash()
	mode := common.GetAppBuildMode()
	date := common.GetAppBuildDate()
	golangVersion := runtime.Version()
	osType := runtime.GOOS
	osArch := runtime.GOARCH

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched version information",
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
	log.Println(JSONresp.Metadata.Response, context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

func (h *HTTPServer) PostSchedulerRun(w http.ResponseWriter, r *http.Request) {
	if !h.scheduling.GetEndpointEnabled() {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`Not found`)); err != nil {
			log.Printf("failed to write response: %v\n", err)
		}
		return
	}
	secret, expectedSecret := r.Header.Get(FlatTrackSchedulerSecretHeader), h.scheduling.GetEndpointSecret()
	if secret != expectedSecret {
		JSONResponse(r, w, http.StatusUnauthorized, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unexpected secret",
			},
		})
		return
	}
	// NOTE make this async?
	if err := h.scheduling.PerformWork(); err != nil {
		JSONResponse(r, w, http.StatusInternalServerError, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to run work",
			},
		})
		return
	}
	JSONResponse(r, w, http.StatusOK, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "completed work",
		},
	})
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

// Healthz ...
// HTTP handler for health checks
func (h *HTTPServer) Healthz(w http.ResponseWriter, r *http.Request) {
	var context string

	if err := h.health.Healthy(); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "not healthy",
			},
			Data: false,
		}
		log.Println(JSONresp.Metadata.Response, context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "healthy",
		},
		Data: true,
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
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
			EndpointPath: "/system/schedule",
			HandlerFunc:  h.PostSchedulerRun,
			HTTPMethod:   http.MethodPost,
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
