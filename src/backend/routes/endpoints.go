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
			EndpointPath: endpointPrefix + "/admin/settings/flatName",
			HandlerFunc:  HTTPuseMiddleware(GetSettingsFlatName(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/settings/flatName",
			HandlerFunc:  HTTPuseMiddleware(SetSettingsFlatName(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(GetAllUsers(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(GetUser(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  HTTPuseMiddleware(PostUser(db), HTTPvalidateJWT(db)),
			HttpMethod:   http.MethodPost,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  HTTPuseMiddleware(DeleteUser(db), HTTPvalidateJWT(db)),
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
	}
}
