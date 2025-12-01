package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/database"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

const (
	//nolint:gosec
	FlatTrackSchedulerSecretHeader = "X-FlatTrack-Scheduler-Secret"
)

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
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID

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
		selectors.NotID = jwtUserID
	}

	users, err := h.users.List(false, selectors)
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
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

	user, err := h.users.GetByID(id, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user",
			},
			Spec: types.UserSpec{},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user account",
		},
		Spec: user,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostUser ...
// create a user
func (h *HTTPServer) PostUser(w http.ResponseWriter, r *http.Request) {
	var context string

	var user types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	userAccount, err := h.users.Create(user, user.Password == "")
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create user account",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PutUser ...
// updates a user account by their id
func (h *HTTPServer) PutUser(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	if _, err := h.users.GetByID(userID, false); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}

	// TODO disallow admins to remove their own admin group access
	userAccountUpdated, err := h.users.UpdateAsAdmin(userID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated user account",
		},
		Spec: userAccountUpdated,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchUser ...
// patches a user account by their id
func (h *HTTPServer) PatchUser(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	if _, err := h.users.GetByID(userID, false); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}

	// TODO disallow admins to remove their own admin group access
	userAccountPatched, err := h.users.PatchAsAdmin(userID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched user account",
		},
		Spec: userAccountPatched,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchUserDisabled ...
// patches a user account's disabled field by their id
func (h *HTTPServer) PatchUserDisabled(w http.ResponseWriter, r *http.Request) {
	var context string

	var userAccount types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	if _, err := h.users.GetByID(userID, false); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	userAccountPatched, err := h.users.PatchDisabledAsAdmin(userID, userAccount.Disabled)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "disabled user account",
		},
		Spec: userAccountPatched,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteUser ...
// delete a user
func (h *HTTPServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	vars := mux.Vars(r)
	userID := vars["id"]

	userInDB, err := h.users.GetByID(userID, false)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account by id: " + userID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if jwtUserID == userID {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to delete user account of invoker",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusForbidden, JSONresp)
		return
	}

	if err := h.users.DeleteByID(userInDB.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user account id from token",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted user account",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched profile",
		},
		Spec: user,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutProfile ...
// Update a user account their id from their JWT
func (h *HTTPServer) PutProfile(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var userAccount types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	userAccountUpdated, err := h.users.Update(jwtUserID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update user account",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated user account",
		},
		Spec: userAccountUpdated,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchProfile ...
// patches a user account their id from their JWT
func (h *HTTPServer) PatchProfile(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var userAccount types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	userAccountPatched, err := h.users.Patch(jwtUserID, userAccount)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch user account",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UserAuth ...
// authenticate a user
func (h *HTTPServer) UserAuth(w http.ResponseWriter, r *http.Request) {
	var user types.UserSpec
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	userInDB, err := h.users.GetByEmail(user.Email, false)
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
		slog.Error("error checking password", "error", err)
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
		slog.Error("error checking password", "error", err)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusUnauthorized, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "auth token is valid",
		},
		Data: valid,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UserAuthReset ...
// invalidates all JWTs
func (h *HTTPServer) UserAuthReset(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	if err := h.users.GenerateNewAuthNonce(jwtUserID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to find user account with id: " + jwtUserID,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "reset all authentication tokens",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UserCanIgroup ...
// respond whether the current user account is in a group
func (h *HTTPServer) UserCanIgroup(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	vars := mux.Vars(r)
	groupName := vars["name"]

	group, err := h.groups.GetByName(groupName)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group by name: " + groupName,
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	userIsInGroup, err := h.groups.CheckUserInGroup(jwtUserID, group.Name)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to check whether user is in group",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user is in group",
		},
		Data: userIsInGroup,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if flatName == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "flat name is not set",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched flat name",
		},
		Spec: flatName,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// SetSettingsFlatName ...
// update the flat's name
func (h *HTTPServer) SetSettingsFlatName(w http.ResponseWriter, r *http.Request) {
	var context string

	var flatName types.FlatName
	if err := json.NewDecoder(r.Body).Decode(&flatName); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set flat name",
		},
		Spec: true,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if initialized {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "system is initialised",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusOK, JSONresp)
		return
	}

	var registrationForm types.Registration
	if err := json.NewDecoder(r.Body).Decode(&registrationForm); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}
	registrationForm.Secret = r.FormValue("secret")

	registered, jwt, err := h.registration.Register(registrationForm)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to register instance",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// GetShoppingList ...
// responds with list of shopping lists
func (h *HTTPServer) GetShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	shoppingList, err := h.shoppinglist.ShoppingList().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping list",
		},
		Spec: shoppingList,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingLists ...
// responds with shopping list by id
func (h *HTTPServer) GetShoppingLists(w http.ResponseWriter, r *http.Request) {
	var context string
	modificationTimestampAfterString := r.FormValue("modificationTimestampAfter")
	creationTimestampAfterString := r.FormValue("creationTimestampAfter")
	limitString := r.FormValue("limit")
	pageString := r.FormValue("page")
	modificationTimestampAfter, err := strconv.ParseInt(modificationTimestampAfterString, 10, 64)
	if err != nil && modificationTimestampAfterString != "" {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to parse value for limiting request for shopping lists",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	creationTimestampAfter, err := strconv.ParseInt(creationTimestampAfterString, 10, 64)
	if err != nil && creationTimestampAfterString != "" {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to parse value for limiting request for shopping lists",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil && pageString != "" {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to parse value for limiting request for shopping lists",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	page, err := strconv.Atoi(pageString)
	if err != nil && pageString != "" {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "unable to parse value for limiting request for shopping lists",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	options := types.ShoppingListOptions{
		SortBy: r.FormValue("sortBy"),
		Limit:  limit,
		Page:   page,
		Selector: types.ShoppingListSelector{
			Completed:                  r.FormValue("completed"),
			ModificationTimestampAfter: modificationTimestampAfter,
			CreationTimestampAfter:     creationTimestampAfter,
		},
	}

	shoppingLists, err := h.shoppinglist.ShoppingList().List(options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping lists",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping lists",
		},
		List: shoppingLists,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostShoppingList ...
// creates a new shopping list to add items to
func (h *HTTPServer) PostShoppingList(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingList types.ShoppingListSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingList); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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

	shoppingList.Author = jwtUserID
	shoppingListInserted, err := h.shoppinglist.ShoppingList().Create(shoppingList, options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "created shopping list",
		},
		Spec: shoppingListInserted,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PatchShoppingList ...
// patches an existing shopping list
func (h *HTTPServer) PatchShoppingList(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingList types.ShoppingListSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingList); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	shoppingList.AuthorLast = jwtUserID
	shoppingListPatched, err := h.shoppinglist.ShoppingList().Patch(list.ID, shoppingList)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched shopping list",
		},
		Spec: shoppingListPatched,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutShoppingList ...
// updates an existing shopping list
func (h *HTTPServer) PutShoppingList(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingList types.ShoppingListSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingList); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	shoppingList.AuthorLast = jwtUserID
	shoppingListUpdated, err := h.shoppinglist.ShoppingList().Update(list.ID, shoppingList)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingList().Delete(list.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to delete shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted shopping list",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	list, err := h.shoppinglist.ShoppingList().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	// TODO add item selectors for this endpoint
	shoppingListItems, err := h.shoppinglist.ShoppingItem().List(list.ID, options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list items",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	shoppingListItem, err := h.shoppinglist.ShoppingItem().Get(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if shoppingListItem.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetch shopping list item",
		},
		Spec: shoppingListItem,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostItemToShoppingList ...
// adds an item to a shopping list
func (h *HTTPServer) PostItemToShoppingList(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingItem types.ShoppingItemSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.Author = jwtUserID
	shoppingItemInserted, err := h.shoppinglist.ShoppingItem().AddItemToList(list.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to add item to shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusBadRequest, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "added item to shopping list",
		},
		Spec: shoppingItemInserted,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// PatchShoppingListCompleted ...
// adds an item to a shopping list
func (h *HTTPServer) PatchShoppingListCompleted(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingList types.ShoppingListSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingList); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	listID := vars["id"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	patchedList, err := h.shoppinglist.ShoppingList().SetListCompleted(list.ID, shoppingList.Completed, jwtUserID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to set shopping list as completed",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "shopping list set as completed",
		},
		Spec: patchedList,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchShoppingListItem ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItem(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingItem types.ShoppingItemSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().Get(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.AuthorLast = jwtUserID
	patchedItem, err := h.shoppinglist.ShoppingItem().Patch(listID, item.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to patch shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "patched shopping list item",
		},
		Spec: patchedItem,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutShoppingListItem ...
// updates an item in a shopping list
func (h *HTTPServer) PutShoppingListItem(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingItem types.ShoppingItemSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().Get(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	shoppingItem.AuthorLast = jwtUserID
	updatedItem, err := h.shoppinglist.ShoppingItem().Update(listID, item.ID, shoppingItem)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping list item",
		},
		Spec: updatedItem,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PatchShoppingListItemObtained ...
// patches an item in a shopping list
func (h *HTTPServer) PatchShoppingListItemObtained(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var shoppingItem types.ShoppingItemSpec
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().Get(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	patchedItem, err := h.shoppinglist.ShoppingItem().SetItemObtained(listID, item.ID, shoppingItem.Obtained, jwtUserID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set shopping list item as obtained",
		},
		Spec: patchedItem,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteShoppingListItem ...
// delete a shopping list item by it's id
func (h *HTTPServer) DeleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	vars := mux.Vars(r)
	itemID := vars["itemId"]
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	item, err := h.shoppinglist.ShoppingItem().Get(list.ID, itemID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if item.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list item",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingItem().Delete(item.ID, list.ID, jwtUserID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to remove item from shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "removed item from shopping list",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// DeleteShoppingListTagItems ...
// delete a shopping list items by matching a tag
func (h *HTTPServer) DeleteShoppingListTagItems(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string
	var item types.ShoppingItemSpec

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	vars := mux.Vars(r)
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingItem().DeleteTagItems(list.ID, item.Tag, jwtUserID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to remove items from shopping list by tag name",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "removed items from shopping list by tag name",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetShoppingListItemTags ...
// responds with tags used in shopping list items from a list
func (h *HTTPServer) GetShoppingListItemTags(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	listID := vars["listId"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	tags, err := h.shoppinglist.ShoppingTag().ListTagsInList(list.ID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get tags from shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched tags from shopping list",
		},
		List: tags,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UpdateShoppingListItemTag ...
// updates then tag name used in shopping list items from a list
func (h *HTTPServer) UpdateShoppingListItemTag(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	listID := vars["listId"]
	tag := vars["tagName"]

	list, err := h.shoppinglist.ShoppingList().Get(listID)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if list.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	var tagUpdate types.ShoppingTag
	if err := json.NewDecoder(r.Body).Decode(&tagUpdate); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	tagUpdated, err := h.shoppinglist.ShoppingTag().UpdateInList(list.ID, tag, tagUpdate.Name)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping list tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "updated shopping list tag",
		},
		Spec: tagUpdated,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetAllShoppingTags ...
// responds with all tags used in shopping list items
func (h *HTTPServer) GetAllShoppingTags(w http.ResponseWriter, r *http.Request) {
	var context string
	options := types.ShoppingTagOptions{
		SortBy: r.FormValue("sortBy"),
	}

	tags, err := h.shoppinglist.ShoppingTag().List(options)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping list tags",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping list tags",
		},
		List: tags,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PostShoppingTag ...
// creates a tag name
func (h *HTTPServer) PostShoppingTag(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string

	var tag types.ShoppingTag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	tag.Author = jwtUserID
	tagCreated, err := h.shoppinglist.ShoppingTag().Create(tag)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to create shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "created shopping tag",
		},
		Spec: tagCreated,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusCreated, JSONresp)
}

// GetShoppingTag ...
// gets a shopping tag by id
func (h *HTTPServer) GetShoppingTag(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	tag, err := h.shoppinglist.ShoppingTag().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping tag",
		},
		Spec: tag,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// UpdateShoppingTag ...
// updates a tag name
func (h *HTTPServer) UpdateShoppingTag(w http.ResponseWriter, r *http.Request) {
	reqClaims := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
	jwtUserID := reqClaims.ID
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	tag, err := h.shoppinglist.ShoppingTag().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	var tagUpdate types.ShoppingTag
	if err := json.NewDecoder(r.Body).Decode(&tagUpdate); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}
	tagUpdate.AuthorLast = jwtUserID
	tagUpdated, err := h.shoppinglist.ShoppingTag().Update(tag.ID, tagUpdate)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to update shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	tag, err := h.shoppinglist.ShoppingTag().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if tag.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}

	if err := h.shoppinglist.ShoppingTag().Delete(tag.ID); err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to delete shopping tag",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "deleted shopping tag",
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping notes",
		},
		Spec: notes,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutSettingsShoppingList ...
// update the notes for shopping lists
func (h *HTTPServer) PutSettingsShoppingList(w http.ResponseWriter, r *http.Request) {
	var context string

	var notes types.ShoppingListNotes
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, code, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set shopping notes",
		},
		Spec: notes.Notes,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutSettingsFlatNotes ...
// update the notes for flat
func (h *HTTPServer) PutSettingsFlatNotes(w http.ResponseWriter, r *http.Request) {
	var context string

	var notes types.FlatNotes
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		slog.Error("failed to unmarshal", "error", err)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set flat notes",
		},
		Spec: notes.Notes,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetSettingsShoppingListKeepPolicy ...
// responds with the keepPolicy for shopping lists
func (h *HTTPServer) GetSettingsShoppingListKeepPolicy(w http.ResponseWriter, r *http.Request) {
	var context string
	keepPolicy, err := h.settings.GetShoppingListKeepPolicy()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping keep policy",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched shopping keep policy",
		},
		Spec: keepPolicy,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// PutSettingsShoppingListKeepPolicy ...
// update the keep policy for shopping lists
func (h *HTTPServer) PutSettingsShoppingListKeepPolicy(w http.ResponseWriter, r *http.Request) {
	var context string

	var spec types.ShoppingListKeepPolicySpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		slog.Error("failed to unmarshal", "error", err)
		JSONResponse(r, w, http.StatusBadRequest, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to read request body",
			},
		})
		return
	}

	if err := h.settings.SetShoppingListKeepPolicy(spec.KeepPolicy); err != nil {
		context = err.Error()
		code := http.StatusInternalServerError
		if err.Error() == "Unable to set shopping list keep policy as it is invalid" {
			code = http.StatusBadRequest
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get shopping keep policy",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, code, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "set shopping keep policy",
		},
		Spec: spec.KeepPolicy,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetAllGroups ...
// returns a list of all groups
func (h *HTTPServer) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var context string
	groups, err := h.groups.List()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get groups",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	group, err := h.groups.GetByID(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if group.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get group",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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

	creationSecrets, err := h.users.UserCreationSecrets().List(userCreationSecretSelector)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secrets",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user creation secrets",
		},
		List: creationSecrets,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUserConfirm ...
// returns an account confirm by id
func (h *HTTPServer) GetUserConfirm(w http.ResponseWriter, r *http.Request) {
	var context string

	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.UserCreationSecrets().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if creationSecret.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched user creation secret",
		},
		Spec: creationSecret,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

// GetUserConfirmValid ...
// returns if an account confirm is valid by id
// TODO should this exist?
func (h *HTTPServer) GetUserConfirmValid(w http.ResponseWriter, r *http.Request) {
	var context string
	vars := mux.Vars(r)
	id := vars["id"]

	creationSecret, err := h.users.UserCreationSecrets().Get(id)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusNotFound, JSONresp)
		return
	}
	if creationSecret.ID == "" {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get user creation secret",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "confirmed user account",
		},
		Data: tokenString,
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
	postgresVersion, err := database.GetVersion(h.db)
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get postgres version",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	osType := runtime.GOOS
	osArch := runtime.GOARCH
	schedulerLastRun, err := h.system.GetSchedulerLastRun()
	if err != nil {
		context = err.Error()
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "failed to get scheduler last run info",
			},
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "fetched version information",
		},
		Data: types.SystemVersion{
			Version:          version,
			CommitHash:       commitHash,
			Mode:             mode,
			Date:             date,
			GolangVersion:    golangVersion,
			PostgresVersion:  postgresVersion,
			OSType:           osType,
			OSArch:           osArch,
			SchedulerLastRun: schedulerLastRun,
		},
	}
	slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
	JSONResponse(r, w, http.StatusOK, JSONresp)
}

func (h *HTTPServer) PostSchedulerRun(w http.ResponseWriter, r *http.Request) {
	if !h.scheduling.GetEndpointEnabled() {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`Not found`)); err != nil {
			slog.Error("failed to write response", "error", err)
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
// 				slog.Error("failed to write response", "error", err)
// 			}
// 			return
// 		}
// 		log.Println(objectInfo.Key, objectInfo.Size, objectInfo.ContentType)
// 		w.Header().Set("content-length", fmt.Sprintf("%d", objectInfo.Size))
// 		w.Header().Set("content-type", objectInfo.ContentType)
// 		w.Header().Set("accept-ranges", "bytes")
// 		w.WriteHeader(http.StatusOK)
// 		if _, err := w.Write(object); err != nil {
// 			slog.Error("failed to write response", "error", err)
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
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
		JSONResponse(r, w, http.StatusInternalServerError, JSONresp)
		return
	}
	if h.maintenanceMode {
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "instance in maintenance mode",
			},
			Data: false,
		}
		slog.Info("request log", "response", JSONresp.Metadata.Response, "context", context)
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
func (h *HTTPServer) HTTPvalidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contextMsg string
		valid, claims, err := h.users.ValidateJWTauthToken(r)
		if err != nil {
			slog.Error("Failed to validate JWT", "error", err)
			JSONResponse(r, w, http.StatusUnauthorized, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unauthorized",
				},
			})
			return
		}
		if claims.ID != "" {
			contextMsg = fmt.Sprintf("with user id '%v'", claims.ID)
		}
		if !valid {
			slog.Info("Unauthorized request with token", "context", contextMsg)
			JSONResponse(r, w, http.StatusUnauthorized, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unauthorized",
				},
			})
			return
		}
		ctx := context.WithValue(r.Context(), types.RequestContextKeyClaimAuth, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// HTTPcheckGroupsFromID ...
// middleware for checking if a route can be accessed given a ID and groupID
func (h *HTTPServer) HTTPcheckGroupsFromID(next http.HandlerFunc, groupsAllowed ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqClaims, ok := r.Context().Value(types.RequestContextKeyClaimAuth).(*types.JWTclaim)
		if !ok {
			JSONResponse(r, w, http.StatusInternalServerError, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Unable to find claims",
				},
			})
			return
		}
		jwtUserID := reqClaims.ID
		found := 0
		for _, group := range groupsAllowed {
			if userInGroup, err := h.groups.CheckUserInGroup(jwtUserID, group); userInGroup && err == nil {
				found++
			}
		}
		if found == len(groupsAllowed) {
			next.ServeHTTP(w, r)
			return
		}
		slog.Info("User tried to access route that is protected by group access", "uid", jwtUserID, "groupsAllowed", groupsAllowed)
		JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Forbidden",
			},
		})
	}
}

func (h *HTTPServer) HTTPMaintenanceMode(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(r, w, http.StatusServiceUnavailable, types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: "Instance in maintenance mode",
			},
		})
	}
}

// HTTP404 ...
// responds with 404
func (h *HTTPServer) HTTP404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`Not found`)); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	}
}

// HTTP404 ...
// responds with 404
func (h *HTTPServer) HTTPMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(`Method not allowed`)); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	}
}

func (h *HTTPServer) registerAPIHandlers(router *mux.Router) {
	routes := []struct {
		EndpointPath     string
		HandlerFunc      http.HandlerFunc
		HTTPMethod       string
		RequireAuth      bool
		RequireAllGroups []string
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
			HandlerFunc:  h.GetVersion,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/system/flatName",
			HandlerFunc:  h.GetSettingsFlatName,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath:     "/admin/settings/flatName",
			HandlerFunc:      h.SetSettingsFlatName,
			HTTPMethod:       http.MethodPost,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/settings/shoppingListNotes",
			HandlerFunc:      h.PutSettingsShoppingList,
			HTTPMethod:       http.MethodPut,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/settings/flatNotes",
			HandlerFunc:      h.GetSettingsFlatNotes,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/settings/flatNotes",
			HandlerFunc:      h.PutSettingsFlatNotes,
			HTTPMethod:       http.MethodPut,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/settings/shoppingListKeepPolicy",
			HandlerFunc:      h.GetSettingsShoppingListKeepPolicy,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/settings/shoppingListKeepPolicy",
			HandlerFunc:      h.PutSettingsShoppingListKeepPolicy,
			HTTPMethod:       http.MethodPut,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath: "/admin/register",
			HandlerFunc:  h.PostAdminRegister,
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath:     "/admin/users",
			HandlerFunc:      h.GetAllUsers,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users/{id}",
			HandlerFunc:      h.GetUser,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users",
			HandlerFunc:      h.PostUser,
			HTTPMethod:       http.MethodPost,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users/{id}",
			HandlerFunc:      h.PatchUser,
			HTTPMethod:       http.MethodPatch,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users/{id}/disabled",
			HandlerFunc:      h.PatchUserDisabled,
			HTTPMethod:       http.MethodPatch,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users/{id}",
			HandlerFunc:      h.PutUser,
			HTTPMethod:       http.MethodPut,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/users/{id}",
			HandlerFunc:      h.DeleteUser,
			HTTPMethod:       http.MethodDelete,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/useraccountconfirms",
			HandlerFunc:      h.GetUserConfirms,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
		},
		{
			EndpointPath:     "/admin/useraccountconfirms/{id}",
			HandlerFunc:      h.GetUserConfirm,
			HTTPMethod:       http.MethodGet,
			RequireAuth:      true,
			RequireAllGroups: []string{"admin"},
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
			RequireAuth:  true,
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
			HandlerFunc:  h.GetProfile,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  h.PutProfile,
			HTTPMethod:   http.MethodPut,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  h.PatchProfile,
			HTTPMethod:   http.MethodPatch,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/users",
			HandlerFunc:  h.GetAllUsers,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/users/{id}",
			HandlerFunc:  h.GetUser,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/groups",
			HandlerFunc:  h.GetAllGroups,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/groups/{id}",
			HandlerFunc:  h.GetGroup,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/user/can-i/group/{name}",
			HandlerFunc:  h.UserCanIgroup,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/settings/notes",
			HandlerFunc:  h.GetSettingsShoppingListNotes,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  h.GetShoppingLists,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  h.GetShoppingList,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  h.PatchShoppingList,
			HTTPMethod:   http.MethodPatch,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  h.PutShoppingList,
			HTTPMethod:   http.MethodPut,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/completed",
			HandlerFunc:  h.PatchShoppingListCompleted,
			HTTPMethod:   http.MethodPatch,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  h.DeleteShoppingList,
			HTTPMethod:   http.MethodDelete,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  h.PostShoppingList,
			HTTPMethod:   http.MethodPost,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  h.GetShoppingListItems,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  h.GetShoppingListItem,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  h.PostItemToShoppingList,
			HTTPMethod:   http.MethodPost,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  h.PatchShoppingListItem,
			HTTPMethod:   http.MethodPatch,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  h.PutShoppingListItem,
			HTTPMethod:   http.MethodPut,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}/obtained",
			HandlerFunc:  h.PatchShoppingListItemObtained,
			HTTPMethod:   http.MethodPatch,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  h.DeleteShoppingListItem,
			HTTPMethod:   http.MethodDelete,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tag",
			HandlerFunc:  h.DeleteShoppingListTagItems,
			HTTPMethod:   http.MethodDelete,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags",
			HandlerFunc:  h.GetShoppingListItemTags,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags/{tagName}",
			HandlerFunc:  h.UpdateShoppingListItemTag,
			HTTPMethod:   http.MethodPut,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  h.PostShoppingTag,
			HTTPMethod:   http.MethodPost,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  h.GetAllShoppingTags,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  h.GetShoppingTag,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  h.UpdateShoppingTag,
			HTTPMethod:   http.MethodPut,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  h.DeleteShoppingTag,
			HTTPMethod:   http.MethodDelete,
			RequireAuth:  true,
		},
		{
			EndpointPath: "/flat/info",
			HandlerFunc:  h.GetSettingsFlatNotes,
			HTTPMethod:   http.MethodGet,
			RequireAuth:  true,
		},
	}
	for _, r := range routes {
		handler := r.HandlerFunc
		if h.maintenanceMode {
			handler = h.HTTPMaintenanceMode(handler)
		}
		for _, g := range r.RequireAllGroups {
			handler = h.HTTPcheckGroupsFromID(handler, g)
		}
		if r.RequireAuth {
			handler = h.HTTPvalidateJWT(handler)
		}
		// NOTE handlers go in reverse order of dependency
		router.HandleFunc(r.EndpointPath, handler).Methods(r.HTTPMethod)
	}
}
