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

// GetEndpoints
// group all API endpoints
func GetEndpoints(endpointPrefix string, db *sql.DB) types.Endpoints {
	return types.Endpoints{
		{
			EndpointPath: endpointPrefix + "/system/initialized",
			HandlerFunc:  GetSystemInitialized(db),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/system/flatName",
			HandlerFunc:  HTTPuseMiddleware(GetSettingsFlatName(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/settings/flatName",
			HandlerFunc:  HTTPuseMiddleware(SetSettingsFlatName(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/register",
			HandlerFunc:  HTTPuseMiddleware(PostAdminRegister(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(PostUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteUser(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/admin/useraccountconfirms",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirms(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/useraccountconfirms/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUserConfirm(db), HTTPvalidateJWT(db), HTTPcheckGroupsFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/auth",
			HandlerFunc:  UserAuthValidate(db),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/auth",
			HandlerFunc:  UserAuth(db),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/user/confirm/{id}",
			HandlerFunc:  GetUserConfirmValid(db),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/confirm/{id}",
			HandlerFunc:  PostUserConfirm(db),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/user/can-i/group/{name}",
			HandlerFunc:  HTTPuseMiddleware(UserCanIgroup(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(GetProfile(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(PatchProfile(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/groups",
			HandlerFunc:  HTTPuseMiddleware(GetAllGroups(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/groups/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetGroup(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingLists(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/completed",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListCompleted(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists",
			HandlerFunc:  HTTPuseMiddleware(PostShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItems(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{.*}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItem(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{id}/items",
			HandlerFunc:  HTTPuseMiddleware(PostItemToShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{.*}/items/{id}",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItem(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{.*}/items/{id}/obtained",
			HandlerFunc:  HTTPuseMiddleware(PatchShoppingListItemObtained(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPatch,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/lists/{listId}/items/{itemId}",
			HandlerFunc:  HTTPuseMiddleware(DeleteShoppingListItem(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodDelete,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/tags",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingListItemTags(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
	}
}
