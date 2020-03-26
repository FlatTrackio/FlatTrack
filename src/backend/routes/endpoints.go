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
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  GetAllUsers(db),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users/{id}",
			HandlerFunc:  GetUser(db),
			HttpMethod:   http.MethodGet,
		},
		{
			EndpointPath: endpointPrefix + "/admin/users",
			HandlerFunc:  PostUser(db),
			HttpMethod:   http.MethodPost,
		},
	}
}