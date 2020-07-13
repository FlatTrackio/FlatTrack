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
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/groups"
	"gitlab.com/flattrack/flattrack/src/backend/registration"
	"gitlab.com/flattrack/flattrack/src/backend/settings"
	"gitlab.com/flattrack/flattrack/src/backend/shoppinglist"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
	"gitlab.com/flattrack/flattrack/src/backend/health"
)

// GetAllUsers ...
// swagger:route GET /users users getAllUsers
// responses:
//  200: []userSpec
// get a list of all users
func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to fetch user accounts"

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
		code := http.StatusInternalServerError
		response := "Failed to fetch user account"
		vars := mux.Vars(r)
		id := vars["id"]

		user := types.UserSpec{
			ID: id,
		}

		user, err := users.GetUserByID(db, id, false)
		if err != nil || user.ID == "" {
			code = http.StatusNotFound
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
			code = http.StatusOK
			response = "Fetched user account"
		} else {
			code = http.StatusBadRequest
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

// PostUser ...
// create a user
func PostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusBadRequest
		response := "Failed to create user account"

		var user types.UserSpec
		body, err := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		userAccount, err := users.CreateUser(db, user, user.Password == "")
		if err == nil && userAccount.ID != "" {
			code = http.StatusOK
			response = "Created user account"
		} else {
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

// PutUser ...
// updates a user account by their id
func PutUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to updat the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

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
		code := http.StatusInternalServerError
		response := "Failed to patch the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

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
		code := http.StatusInternalServerError
		response := "Failed to delete user account"

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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		// TODO make sure that if any errors occur the user account still won't be deleted
		// TODO make sure that admins can't be deleted

		err = users.DeleteUserByID(db, userInDB.ID)
		if err == nil {
			code = http.StatusOK
			response = "Deleted user account"
		} else {
			code = http.StatusBadRequest
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

// GetProfile ...
// returns the authenticated user's profile
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to fetch user account"

		user, err := users.GetProfile(db, r)
		if err == nil && user.ID != "" {
			code = http.StatusOK
			response = "Fetched user account"
		} else {
			code = http.StatusNotFound
			response = "Failed to find your profile"
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

// PutProfile ...
// Update a user account their id from their JWT
func PutProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to patch the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

		id, errID := users.GetIDFromJWT(db, r)
		userAccountUpdated, err := users.UpdateProfile(db, id, userAccount)
		if err == nil && errID == nil && userAccountUpdated.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
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
		code := http.StatusInternalServerError
		response := "Failed to patch the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

		id, errID := users.GetIDFromJWT(db, r)
		userAccountPatched, err := users.PatchProfile(db, id, userAccount)
		if err == nil && errID == nil && userAccountPatched.ID != "" {
			code = http.StatusOK
			response = "Successfully patched the user account"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
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
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to fetch if this FlatTrack instance has initialized"
		initialized, err := system.GetHasInitialized(db)
		if err == nil {
			code = http.StatusOK
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

// UserAuth ...
// authenticate a user
func UserAuth(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to authenticate user, incorrect email or password"
		code := http.StatusUnauthorized
		jwtToken := ""

		var user types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		userInDB, err := users.GetUserByEmail(db, user.Email, false)
		if err != nil || userInDB.ID == "" {
			response = "Unable to find the account"
			code = http.StatusNotFound
		}
		if userInDB.ID != "" && userInDB.Registered == false {
			response = "Account not yet registered"
			code = http.StatusForbidden
		}
		// Check password locally, fall back to remote if incorrect
		matches, err := users.CheckUserPassword(db, userInDB.Email, user.Password)
		if err == nil && matches == true && code == http.StatusUnauthorized {
			jwtToken, _ = users.GenerateJWTauthToken(db, userInDB.ID, userInDB.AuthNonce, 0)
			response = "Successfully authenticated user"
			code = http.StatusOK
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

// UserAuthValidate ...
// validate an auth token
func UserAuthValidate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to validate authentication token"
		code := http.StatusUnauthorized

		valid, err := users.ValidateJWTauthToken(db, r)
		if valid == true && err == nil {
			response = "Authentication token is valid"
			code = http.StatusOK
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

// UserAuthReset ...
// invalidates all JWTs
func UserAuthReset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to sign out all devices logged in"
		code := http.StatusInternalServerError

		id, err := users.GetIDFromJWT(db, r)
		if err != nil {
			response = "Failed to find account"
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
			}
			JSONResponse(r, w, code, JSONresp)
			return
		}
		err = users.GenerateNewAuthNonce(db, id)
		if err == nil {
			response = "Successfully signed out all devices logged in"
			code = http.StatusOK
		}
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
		response := "Failed to determine group privileges"
		code := http.StatusInternalServerError

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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		userIsInGroup, err := groups.CheckUserInGroup(db, id, groupName)
		if err == nil && errID == nil {
			response = "Determined if user account can perform tasks of group"
			code = http.StatusOK
		}

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
		response := "Failed to fetch the flat name"
		code := http.StatusInternalServerError

		flatName, err := settings.GetFlatName(db)
		if flatName == "" {
			response = "Flat name is not set"
			code = http.StatusOK
		} else if err == nil {
			response = "Fetched the flat name"
			code = http.StatusOK
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

// SetSettingsFlatName ...
// update the flat's name
func SetSettingsFlatName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to update the flat name"

		var flatName types.FlatName
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &flatName)

		err := settings.SetFlatName(db, flatName.FlatName)
		if err == nil {
			code = http.StatusOK
			response = "Successfully set flat name"
		} else {
			code = http.StatusBadRequest
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

// PostAdminRegister ...
// register the instance of FlatTrack
func PostAdminRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		response := "Failed to register the FlatTrack instance"

		initialized, err := system.GetHasInitialized(db)
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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		var registrationForm types.Registration
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &registrationForm)

		registered, jwt, err := registration.Register(db, registrationForm)
		if err == nil {
			code = http.StatusOK
			response = "Successfully registered the FlatTrack instance"
		} else {
			code = http.StatusBadRequest
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

// GetShoppingList ...
// responds with list of shopping lists
func GetShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping lists"
		code := http.StatusNotFound

		vars := mux.Vars(r)
		id := vars["id"]

		shoppingList, err := shoppinglist.GetShoppingList(db, id)
		if err == nil && shoppingList.ID != "" {
			response = "Fetched the shopping lists"
			code = http.StatusOK
		}
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
		response := "Failed to fetch shopping lists"
		code := http.StatusInternalServerError

		options := types.ShoppingListOptions{
			SortBy: r.FormValue("sortBy"),
			Selector: types.ShoppingListSelector{
				Completed: r.FormValue("completed"),
			},
		}

		shoppingLists, err := shoppinglist.GetShoppingLists(db, options)
		if err == nil {
			response = "Fetched the shopping lists"
			code = http.StatusOK
		}
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
		code := http.StatusInternalServerError
		response := "Failed to create the shopping list"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

		options := types.ShoppingItemOptions{
			Selector: types.ShoppingItemSelector{
				TemplateListItemSelector: r.FormValue("templateListItemSelector"),
			},
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingList.Author = id
		shoppingListInserted, err := shoppinglist.CreateShoppingList(db, shoppingList, options)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully created the shopping list"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
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
		code := http.StatusBadRequest
		response := "Failed to patch the shopping list"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

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
		code := http.StatusInternalServerError
		response := "Failed to update the shopping list"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

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
		code := http.StatusInternalServerError
		response := "Failed to delete the shopping list"

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
		response := "Failed to fetch shopping list items"
		code := http.StatusInternalServerError

		vars := mux.Vars(r)
		id := vars["id"]

		options := types.ShoppingItemOptions{
			SortBy: r.FormValue("sortBy"),
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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		// TODO add item selectors for this endpoint
		shoppingListItems, err := shoppinglist.GetShoppingListItems(db, id, options)
		if err == nil {
			response = "Fetched the shopping list items"
			code = http.StatusOK
		}
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
		response := "Failed to fetch shopping list item"
		code := http.StatusNotFound

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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		shoppingListItem, err := shoppinglist.GetShoppingListItem(db, listID, itemID)
		if err == nil && shoppingListItem.ID != "" {
			response = "Fetched the shopping list item"
			code = http.StatusOK
		}
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
		code := http.StatusInternalServerError
		response := "Failed to create the shopping list item"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		id, errID := users.GetIDFromJWT(db, r)
		shoppingItem.Author = id
		shoppingItemInserted, err := shoppinglist.AddItemToList(db, listID, shoppingItem)
		if err == nil && errID == nil {
			code = http.StatusOK
			response = "Successfully created the shopping list item"
		} else {
			code = http.StatusBadRequest
			response = err.Error()
		}
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
		code := http.StatusInternalServerError
		response := "Failed to patch the shopping list completed field"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

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
		code := http.StatusInternalServerError
		response := "Failed to patch the shopping list item"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

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
		code := http.StatusInternalServerError
		response := "Failed to update the shopping list item"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

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
		code := http.StatusInternalServerError
		response := "Failed to patch the shopping list item obtained field"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

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
		code := http.StatusInternalServerError
		response := "Failed to delete the shopping list item"

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
		response := "Failed to fetch shopping list item tags"
		code := http.StatusInternalServerError

		vars := mux.Vars(r)
		listID := vars["listId"]

		tags, err := shoppinglist.GetShoppingListTags(db, listID)
		if err == nil {
			response = "Fetched the shopping list item tags"
			code = http.StatusOK
		}
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
		response := "Failed to update shopping list item tag name"
		code := http.StatusInternalServerError

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
			return
		}

		var tagUpdate types.ShoppingTag
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &tagUpdate)

		tag, err = shoppinglist.UpdateShoppingListTag(db, listID, tag, tagUpdate.Name)
		if err == nil {
			response = "Updated the shopping list item tag name"
			code = http.StatusOK
		}
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
		response := "Failed to fetch shopping list item tags"
		code := http.StatusInternalServerError

		options := types.ShoppingTagOptions{
			SortBy: r.FormValue("sortBy"),
		}

		tags, err := shoppinglist.GetAllShoppingTags(db, options)
		if err == nil {
			response = "Fetched the shopping list item tags"
			code = http.StatusOK
		}
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
		response := "Failed to create a shopping tag"
		code := http.StatusInternalServerError

		var tag types.ShoppingTag
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &tag)

		id, errID := users.GetIDFromJWT(db, r)
		tag.Author = id
		tag, err := shoppinglist.CreateShoppingTag(db, tag)
		if err == nil && errID == nil {
			response = "Updated a shopping tag"
			code = http.StatusOK
		} else {
			response = err.Error()
		}
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
		response := "Failed to fetch a shopping tag"
		code := http.StatusInternalServerError

		vars := mux.Vars(r)
		id := vars["id"]

		tag, err := shoppinglist.GetShoppingTag(db, id)
		if err == nil {
			response = "Fetched a shopping tag"
			code = http.StatusOK
		}
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
		response := "Failed to update a shopping tag"
		code := http.StatusInternalServerError

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
			JSONResponse(r, w, code, JSONresp)
			return
		}

		var tagUpdate types.ShoppingTag
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &tagUpdate)

		tag, err := shoppinglist.UpdateShoppingTag(db, id, tagUpdate)
		if err == nil {
			response = "Updated a shopping tag"
			code = http.StatusOK
		}
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
		code := http.StatusInternalServerError
		response := "Failed to delete the shopping tag"

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
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetAllGroups ...
// returns a list of all groups
func GetAllGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch groups"
		code := http.StatusInternalServerError

		groups, err := groups.GetAllGroups(db)
		if err == nil {
			response = "Fetched all groups"
			code = http.StatusOK
		}
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
		response := "Failed to fetch groups"
		code := http.StatusNotFound

		vars := mux.Vars(r)
		id := vars["id"]

		group, err := groups.GetGroupByID(db, id)
		if err == nil && group.ID != "" {
			response = "Fetched the groups"
			code = http.StatusOK
		}
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
		response := "Failed to fetch user account creation secrets"
		code := http.StatusInternalServerError

		userIDSelector := r.FormValue("userId")
		userCreationSecretSelector := types.UserCreationSecretSelector{
			UserID: userIDSelector,
		}

		creationSecrets, err := users.GetAllUserCreationSecrets(db, userCreationSecretSelector)
		if err == nil {
			response = "Fetched the user account creation secrets"
			code = http.StatusOK
		}
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
		response := "Failed to fetch user account creation secret"
		code := http.StatusNotFound

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.ID != "" {
			response = "Fetched the user account creation secret"
			code = http.StatusOK
		}
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
		response := "Failed to fetch user account creation secret"
		code := http.StatusNotFound

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.ID != "" {
			response = "Fetched the user account creation secret"
			code = http.StatusOK
		}
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
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to confirm your user account"
		code := http.StatusInternalServerError

		vars := mux.Vars(r)
		id := vars["id"]

		secret := r.FormValue("secret")

		var user types.UserSpec
		body, errUnmarshal := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		tokenString, err := users.ConfirmUserAccount(db, id, secret, user)
		if err == nil && errUnmarshal == nil {
			response = "Your user account has been confirmed"
			code = http.StatusOK
		} else {
			response = err.Error()
		}
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
	version := common.GetAppBuildVersion()
	commitHash := common.GetAppBuildHash()
	mode := common.GetAppBuildMode()
	date := common.GetAppBuildDate()

	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Fetched version information",
		},
		Data: types.SystemVersion{
			Version:    version,
			CommitHash: commitHash,
			Mode:       mode,
			Date:       date,
		},
	}
	JSONResponse(r, w, http.StatusOK, JSONresp)
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
		response := "App unhealthy"
		code := http.StatusInternalServerError

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
			if completed, err := users.ValidateJWTauthToken(db, r); completed == true && err == nil {
				h.ServeHTTP(w, r)
				return
			}
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
				if userInGroup, err := groups.CheckUserInGroup(db, id, group); userInGroup == true && err == nil && err == errID {
					h.ServeHTTP(w, r)
					return
				}
			}
			JSONResponse(r, w, http.StatusForbidden, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Forbidden",
				},
			})
		}
	}
}
