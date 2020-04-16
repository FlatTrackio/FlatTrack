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
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/flattrack/flattrack/src/backend/groups"
	"gitlab.com/flattrack/flattrack/src/backend/registration"
	"gitlab.com/flattrack/flattrack/src/backend/settings"
	"gitlab.com/flattrack/flattrack/src/backend/shoppinglist"
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

		userSelectorId := r.FormValue("id")
		userSelectorNotId := r.FormValue("notId")
		userSelectorGroup := r.FormValue("group")
		var errJWT error = nil
		if r.FormValue("notSelf") == "true" {
			userSelectorNotId, errJWT = users.GetIdFromJWT(db, r)
		}

		selectors := types.UserSelector{
			Id:    userSelectorId,
			NotId: userSelectorNotId,
			Group: userSelectorGroup,
		}

		users, err := users.GetAllUsers(db, false, selectors)
		if err == nil && errJWT == nil {
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

		user, err := users.GetUserById(db, id, false)
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

		userAccount, err := users.CreateUser(db, user, user.Password == "")
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

// PatchUser
// patches a user account by their id
func PatchUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

		vars := mux.Vars(r)
		userId := vars["id"]

		// TODO disallow admins to remove their own admin group access
		userAccountPatched, err := users.PatchProfile(db, userId, userAccount)
		if err == nil && userAccountPatched.Id != "" {
			code = 200
			response = "Successfully patched the user account"
		} else {
			code = 400
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

// DeleteUser
// delete a user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to delete user account"

		vars := mux.Vars(r)
		userId := vars["id"]

		userInDB, err := users.GetUserById(db, userId, false)
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

		// TODO make sure that if any errors occur the user account still won't be deleted
		// TODO make sure that admins can't be deleted

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

// GetProfile
// returns the authenticated user's profile
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to fetch user account"

		user, err := users.GetProfile(db, r)
		if err == nil && user.Id != "" {
			code = 200
			response = "Fetched user account"
		} else {
			code = 404
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

// PatchProfile
// patches a user account their id from their JWT
func PatchProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the user account"

		var userAccount types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &userAccount)

		id, errId := users.GetIdFromJWT(db, r)
		userAccountPatched, err := users.PatchProfile(db, id, userAccount)
		if err == nil && errId == nil && userAccountPatched.Id != "" {
			code = 200
			response = "Successfully patched the user account"
		} else {
			code = 400
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

// GetSystemInitialized
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
		code := 401
		jwtToken := ""

		var user types.UserSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		userInDB, err := users.GetUserByEmail(db, user.Email, false)
		if err != nil || userInDB.Id == "" {
			response = "Unable to find the account"
			code = 404
		}
		if userInDB.Id != "" && userInDB.Registered == false {
			response = "Account not yet registered"
			code = 403
		}
		// Check password locally, fall back to remote if incorrect
		matches, err := users.CheckUserPassword(db, userInDB.Email, user.Password)
		if err == nil && matches == true && code == 401 {
			jwtToken, _ = users.GenerateJWTauthToken(db, userInDB.Id, userInDB.AuthNonce)
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
		code := 401

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

// UserCanIgroup
// respond whether the current user account is in a group
func UserCanIgroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to determine group privileges"
		code := 500

		vars := mux.Vars(r)
		groupName := vars["name"]

		id, errId := users.GetIdFromJWT(db, r)
		userIsInGroup, err := groups.CheckUserInGroup(db, id, groupName)
		if err == nil && errId == nil {
			response = "Determined if user account can perform tasks of group"
			code = 200
		}

		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userIsInGroup,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// GetSettingsFlatName
// responds with the name of the flat
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
			return
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

// GetShoppingList
// responds with list of shopping lists
func GetShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping lists"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		shoppingList, err := shoppinglist.GetShoppingList(db, id)
		if err == nil && shoppingList.Id != "" {
			response = "Fetched the shopping lists"
			code = 200
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

// GetShoppingLists
// responds with shopping list by id
func GetShoppingLists(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping lists"
		code := 500

		shoppingLists, err := shoppinglist.GetShoppingLists(db)
		if err == nil {
			response = "Fetched the shopping lists"
			code = 200
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

// PostShoppingList
// creates a new shopping list to add items to
func PostShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to create the shopping list"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

		selectors := types.ShoppingItemSelector{
			NotObtained: r.FormValue("notObtained") == "true",
		}

		id, errId := users.GetIdFromJWT(db, r)
		shoppingList.Author = id
		shoppingListInserted, err := shoppinglist.CreateShoppingList(db, shoppingList, selectors)
		if err == nil && errId == nil {
			code = 200
			response = "Successfully created the shopping list"
		} else {
			code = 400
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

// PatchShoppingList
// patches an existing shopping list
func PatchShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the shopping list"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

		vars := mux.Vars(r)
		listId := vars["id"]

		id, errId := users.GetIdFromJWT(db, r)
		shoppingList.AuthorLast = id
		shoppingListPatched, err := shoppinglist.PatchShoppingList(db, listId, shoppingList)
		if err == nil && errId == nil && shoppingListPatched.Id != "" {
			code = 200
			response = "Successfully patched the shopping list"
		} else {
			code = 400
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

// DeleteShoppingList
// delete a new shopping list by it's id
func DeleteShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to delete the shopping list"

		vars := mux.Vars(r)
		listId := vars["id"]

		err := shoppinglist.DeleteShoppingList(db, listId)
		if err == nil {
			code = 200
			response = "Successfully deleted the shopping list"
		} else {
			code = 400
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

// GetShoppingListItems
// responds with shopping items by list id
func GetShoppingListItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping list items"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		// TODO add item selectors for this endpoint
		shoppingListItems, err := shoppinglist.GetShoppingListItems(db, id, types.ShoppingItemSelector{})
		if err == nil {
			response = "Fetched the shopping list items"
			code = 200
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

// GetShoppingListItem
// responds with list of shopping lists
func GetShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping list item"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		shoppingListItem, err := shoppinglist.GetShoppingListItem(db, id)
		if err == nil && shoppingListItem.Id != "" {
			response = "Fetched the shopping list item"
			code = 200
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

// PostItemToShoppingList
// adds an item to a shopping list
func PostItemToShoppingList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to create the shopping list item"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

		vars := mux.Vars(r)
		listId := vars["id"]

		id, errId := users.GetIdFromJWT(db, r)
		shoppingItem.Author = id
		shoppingItemInserted, err := shoppinglist.AddItemToList(db, listId, shoppingItem)
		if err == nil && errId == nil {
			code = 200
			response = "Successfully created the shopping list item"
		} else {
			code = 400
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

// PatchShoppingListCompleted
// adds an item to a shopping list
func PatchShoppingListCompleted(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the shopping list completed field"

		var shoppingList types.ShoppingListSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingList)

		vars := mux.Vars(r)
		itemId := vars["id"]

		id, errId := users.GetIdFromJWT(db, r)
		shoppingList.AuthorLast = id
		patchedList, err := shoppinglist.SetListCompleted(db, itemId, shoppingList.Completed)
		if err == nil && errId == nil {
			code = 200
			response = "Successfully patched the shopping list completed field"
		} else {
			code = 400
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

// PatchShoppingListItem
// patches an item in a shopping list
func PatchShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the shopping list item"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

		vars := mux.Vars(r)
		itemId := vars["id"]

		id, errId := users.GetIdFromJWT(db, r)
		shoppingItem.AuthorLast = id
		patchedItem, err := shoppinglist.PatchItem(db, itemId, shoppingItem)
		if err == nil && errId == nil {
			code = 200
			response = "Successfully patched the shopping list item"
		} else {
			code = 400
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

// PatchShoppingListItemObtained
// patches an item in a shopping list
func PatchShoppingListItemObtained(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to patch the shopping list item obtained field"

		var shoppingItem types.ShoppingItemSpec
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &shoppingItem)

		vars := mux.Vars(r)
		itemId := vars["id"]

		id, errId := users.GetIdFromJWT(db, r)
		shoppingItem.AuthorLast = id
		patchedItem, err := shoppinglist.SetItemObtained(db, itemId, shoppingItem.Obtained)
		if err == nil && errId == nil {
			code = 200
			response = "Successfully patched the shopping list item obtained field"
		} else {
			code = 400
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

// DeleteShoppingListItem
// delete a shopping list item by it's id
func DeleteShoppingListItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := 500
		response := "Failed to delete the shopping list item"

		vars := mux.Vars(r)
		itemId := vars["itemId"]
		listId := vars["listId"]

		err := shoppinglist.RemoveItemFromList(db, itemId, listId)
		if err == nil {
			code = 200
			response = "Successfully deleted the shopping list item"
		} else {
			code = 400
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

// GetShoppingListItemTags
// responds with tags used in shopping list items
func GetShoppingListItemTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch shopping list item tags"
		code := 500

		tags, err := shoppinglist.GetShoppingListTags(db)
		if err == nil {
			response = "Fetched the shopping list item tags"
			code = 200
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

// GetAllGroups
// returns a list of all groups
func GetAllGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch groups"
		code := 500

		groups, err := groups.GetAllGroups(db)
		if err == nil {
			response = "Fetched all groups"
			code = 200
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

// GetGroup
// returns a group by id
func GetGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch groups"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		group, err := groups.GetGroupById(db, id)
		if err == nil && group.Id != "" {
			response = "Fetched the groups"
			code = 200
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

// GetUserConfirms
// returns a list of account confirms
func GetUserConfirms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch user account creation secrets"
		code := 500

		userIdSelector := r.FormValue("userId")
		userCreationSecretSelector := types.UserCreationSecretSelector{
			UserId: userIdSelector,
		}

		creationSecrets, err := users.GetAllUserCreationSecrets(db, userCreationSecretSelector)
		if err == nil {
			response = "Fetched the user account creation secrets"
			code = 200
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

// GetUserConfirm
// returns an account confirm by id
func GetUserConfirm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch user account creation secret"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.Id != "" {
			response = "Fetched the user account creation secret"
			code = 200
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

// GetUserConfirmValid
// returns if an account confirm is valid by id
func GetUserConfirmValid(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to fetch user account creation secret"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		creationSecret, err := users.GetUserCreationSecret(db, id)
		if err == nil && creationSecret.Id != "" {
			response = "Fetched the user account creation secret"
			code = 200
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Data: creationSecret.Id != "",
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

// PostUserConfirm
// confirm a user account
func PostUserConfirm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Failed to confirm account"
		code := 500

		vars := mux.Vars(r)
		id := vars["id"]

		secret := r.FormValue("secret")

		var user types.UserSpec
		body, errUnmarshal := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		tokenString, err := users.ConfirmUserAccount(db, id, secret, user)
		log.Println(err)
		if err == nil && errUnmarshal == nil {
			response = "Confirmed the account"
			code = 200
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

// HTTPcheckGroupsFromId
// middleware for checking if a route can be accessed given a Id and groupId
func HTTPcheckGroupsFromId(db *sql.DB, groupsAllowed ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id, errId := users.GetIdFromJWT(db, r)
			for _, group := range groupsAllowed {
				if userInGroup, err := groups.CheckUserInGroup(db, id, group); userInGroup == true && err == nil && err == errId {
					h.ServeHTTP(w, r)
					return
				}
			}
			JSONResponse(r, w, 403, types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: "Forbidden",
				},
			})
		}
	}
}
