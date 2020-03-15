package routes

import (
	// "net/http"

	"gitlab.com/flattrack/flattrack/src/backend/types"
)

func GetEndpoints(endpointPrefix string) types.Endpoints {
	return types.Endpoints{
		// {
		// 	EndpointPath: endpointPrefix + "/hello",
		// 	HandlerFunc:  APIgetHello,
		// 	HttpMethod:   http.MethodGet,
		// },
	}
}
