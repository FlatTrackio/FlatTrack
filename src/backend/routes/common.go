/*
  routes
    common
      handle generic request related actions
*/

package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"database/sql"
	"github.com/ddo/go-vue-handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// JSONResponse
// form generic JSON responses
func JSONResponse(r *http.Request, w http.ResponseWriter, code int, output types.JSONMessageResponse) {
	// simpilify sending a JSON response
	output.Metadata.URL = r.RequestURI
	output.Metadata.Timestamp = time.Now().Unix()
	output.Metadata.Version = common.GetAppBuildVersion()
	response, _ := json.Marshal(output)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetHTTPresponseBodyContents
// convert the body of a HTTP response into a JSONMessageResponse
func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &output)
	return output
}

// HTTPuseMiddleware
// append functions to run before the endpoint handler
func HTTPuseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

// Logging
// log the HTTP requests
func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pathSection string
		requestPath := strings.Split(r.URL.Path, "/")
		if len(requestPath) >= 1 && requestPath[1] == "api" {
			pathSection = "backend "
		} else if len(requestPath) >= 1 && requestPath[1] == "metrics" {
			pathSection = "metrics "
		} else {
			pathSection = "frontend"
		}
		log.Printf("[%v] %v %v %v %v %v", pathSection, r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// HandleWebserver
// manage the launching of the API's webserver
func HandleWebserver(db *sql.DB) {
	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(true)
	apiEndpointPrefix := "/api"

	router.HandleFunc(apiEndpointPrefix, Root)
	for _, endpoint := range GetEndpoints(apiEndpointPrefix, db) {
		router.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HttpMethod, http.MethodOptions)
	}

	router.HandleFunc(apiEndpointPrefix+"/{.*}", UnknownEndpoint)
	// TODO implement /metrics for prometheus
	// TODO implement /healthz for healthiness checks
	// TODO implement /readyz for readiness checks
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/robots.txt")
	})
	router.PathPrefix("/").Handler(vue.Handler(common.GetAppDistFolder())).Methods("GET")

	router.Use(Logging)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("listening on", port)
	log.Fatal(srv.ListenAndServe())
}
