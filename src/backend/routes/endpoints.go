/*
  routes
    endpoints
      declare all API routes paths and access configurations
*/

package routes

import (
	"net/http"

	"database/sql"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// GetEndpoints ...
// group all API endpoints
func GetEndpoints(db *sql.DB) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: "/system/initialized",
			HandlerFunc:  GetSystemInitialized(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/system/version",
			HandlerFunc:  HTTPuseMiddleware(GetVersion, HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/system/flatName",
			HandlerFunc:  HTTPuseMiddleware(GetSettingsFlatName(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/settings/flatName",
			HandlerFunc:  HTTPuseMiddleware(SetSettingsFlatName(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/settings/notes",
			HandlerFunc:  HTTPuseMiddleware(GetSettingsShoppingListNotes(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/settings/shoppingListNotes",
			HandlerFunc:  HTTPuseMiddleware(PutSettingsShoppingList(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/admin/register",
			HandlerFunc:  HTTPuseMiddleware(PostAdminRegister(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(PostUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/admin/useraccountconfirms",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirms(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/admin/useraccountconfirms/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirm(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/auth",
			HandlerFunc:  UserAuthValidate(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/auth",
			HandlerFunc:  UserAuth(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/auth/reset",
			HandlerFunc:  UserAuthReset(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/confirm/{id}",
			HandlerFunc:  GetUserConfirmValid(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/confirm/{id}",
			HandlerFunc:  PostUserConfirm(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(GetProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(PutProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(PatchProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/groups",
			HandlerFunc:  HTTPuseMiddleware(GetAllGroups(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/groups/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetGroup(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/user/can-i/group/{name}",
			HandlerFunc:  HTTPuseMiddleware(UserCanIgroup(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingLists(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/completed",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListCompleted(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(PostShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItems(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(PostItemToShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{id}/obtained",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItemObtained(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItemTags(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/lists/{listId}/tags/{tagName}",
			HandlerFunc:  HTTPuseMiddleware(UpdateShoppingListItemTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  HTTPuseMiddleware(PostShoppingTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags",
			HandlerFunc:  HTTPuseMiddleware(GetAllShoppingTags(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  HTTPuseMiddleware(UpdateShoppingTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: "/apps/shoppinglist/tags/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodDelete,
		},
	}
}
