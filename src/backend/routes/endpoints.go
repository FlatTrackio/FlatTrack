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
func GetEndpoints(endpointPrefix string, db *sql.DB) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: endpointPrefix + "/system/initialized",
			HandlerFunc:  GetSystemInitialized(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/system/version",
			HandlerFunc:  HTTPuseMiddleware(GetVersion, HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/system/flatName",
			HandlerFunc:  HTTPuseMiddleware(GetSettingsFlatName(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/settings/flatName",
			HandlerFunc:  HTTPuseMiddleware(SetSettingsFlatName(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/register",
			HandlerFunc:  HTTPuseMiddleware(PostAdminRegister(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(PostUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/admin/useraccountconfirms",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirms(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/useraccountconfirms/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirm(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromID(db, "admin")),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/auth",
			HandlerFunc:  UserAuthValidate(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/auth",
			HandlerFunc:  UserAuth(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/user/auth/reset",
			HandlerFunc:  UserAuthReset(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/user/confirm/{id}",
			HandlerFunc:  GetUserConfirmValid(db),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/confirm/{id}",
			HandlerFunc:  PostUserConfirm(db),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(GetProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(PutProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(PatchProfile(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/groups",
			HandlerFunc:  HTTPuseMiddleware(GetAllGroups(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/groups/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetGroup(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/can-i/group/{name}",
			HandlerFunc:  HTTPuseMiddleware(UserCanIgroup(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingLists(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/completed",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListCompleted(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(PostShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItems(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(PostItemToShoppingList(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(PutShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{id}/obtained",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItemObtained(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingListItem(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/tags",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItemTags(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/tags/{tagName}",
			HandlerFunc:  HTTPuseMiddleware(UpdateShoppingListItemTag(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodPut,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/tags",
			HandlerFunc:  HTTPuseMiddleware(GetAllShoppingListItemTags(db), HTTPvalidateJWT(db)),
			HTTPMethod:   http.MethodGet,
		},
	}
}
