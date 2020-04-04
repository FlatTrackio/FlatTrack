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
			HandlerFunc:  HTTPuseMiddleware(SetSettingsFlatName(db), HTTPvalidateJWT(db), HTTPcheckGroupFromId(db, "admin")),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/register",
			HandlerFunc:  HTTPuseMiddleware(PostAdminRegister(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db), HTTPcheckGroupFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db), HTTPcheckGroupFromId(db, "admin")),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(PostUser(db), HTTPvalidateJWT(db), HTTPcheckGroupFromId(db, "admin")),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteUser(db), HTTPvalidateJWT(db), HTTPcheckGroupFromId(db, "admin")),
			HttpMethod:   http.MethodDelete,
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
			EndpointPath: endpointPrefix + "/user/profile",
			HandlerFunc:  HTTPuseMiddleware(GetProfile(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
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
			EndpointPath: endpointPrefix + "/apps/shoppinglist",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingLists(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/apps/shoppinglist",
			HandlerFunc:  HTTPuseMiddleware(PostShoppingList(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPost,
		},
	}
}
