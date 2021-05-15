/*
  routes
    common
      handle generic request related actions
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package routes

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// JSONResponse ...
// form generic JSON responses
func JSONResponse(r *http.Request, w http.ResponseWriter, code int, output types.JSONMessageResponse) {
	// simpilify sending a JSON response
	output.Metadata.URL = r.RequestURI
	output.Metadata.Timestamp = time.Now().Unix()
	output.Metadata.Version = common.GetAppBuildVersion()
	response, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetHTTPresponseBodyContents ...
// convert the body of a HTTP response into a JSONMessageResponse
func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &output)
	return output
}

// HTTPuseMiddleware ...
// append functions to run before the endpoint handler
func HTTPuseMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

// HealthHandler ...
func HealthHandler(db *sql.DB) {
	if common.GetAppHealthEnabled() != "true" {
		return
	}

	port := common.GetAppHealthPort()
	http.Handle("/_healthz", Healthz(db))
	log.Printf("Health listening on %v", port)
	http.ListenAndServe(port, nil)
}

// GetRequestIP ...
// returns r.RemoteAddr unless RealIPHeader is set
func GetRequestIP(r *http.Request) (requestIP string) {
	realIPHeader := common.GetAppRealIPHeader()
	headerValue := r.Header.Get(realIPHeader)
	if realIPHeader == "" || headerValue == "" {
		return r.RemoteAddr
	}
	return headerValue
}

// Logging ...
// log the HTTP requests
func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathSection := "frontend"
		requestPath := strings.Split(r.URL.Path, "/")
		if len(requestPath) >= 1 && requestPath[1] == "api" {
			pathSection = "backend "
		} else if len(requestPath) >= 1 && requestPath[1] == "metrics" {
			pathSection = "metrics "
		}
		requestIP := GetRequestIP(r)
		log.Printf("[%v] %v %v %v %v %v", pathSection, r.Method, r.URL, r.Proto, r.Response, requestIP)
		next.ServeHTTP(w, r)
	})
}

// RequireContentType ...
// 404s requests if content-type isn't what is expected
func RequireContentType(expectedContentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if (len(r.Header["Content-Type"]) > 0 && r.Header["Content-Type"][0] == expectedContentType) ||
				(len(r.Header["Accept"]) > 0 && r.Header["Accept"][0] == expectedContentType) {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/unknown-page", http.StatusMovedPermanently)
		})
	}
}

// FrontendOptions ...
// options to send to the frontend index.html for templating
type FrontendOptions struct {
	SetupMessage string
	LoginMessage string
	EmbeddedHTML template.HTML
}

// FrontendHandler ...
// handles rewriting and API root setting
func FrontendHandler(publicDir string, passthrough FrontendOptions) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.Replace(req.URL.Path, "/", "", 1)
		if len(req.URL.Path) > 0 && req.URL.Path[len(req.URL.Path)-1:] != "/" {
			req.URL.Path = path.Join("/", req.URL.Path)
		}

		// TODO redirect to subPath + /unknown-page if _path does not include subPath at the front

		// static files
		if strings.Contains(req.URL.Path, ".") {
			handler.ServeHTTP(w, req)
			return
		}

		// frontend views
		indexPath := path.Join(publicDir, "/index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, passthrough); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// Handle ...
// manage the launching of the API's webserver
func Handle(db *sql.DB) {
	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(true)
	apiEndpointPrefix := "/api"
	passthrough := FrontendOptions{
		SetupMessage: common.GetAppSetupMessage(),
		LoginMessage: common.GetAppLoginMessage(),
		EmbeddedHTML: template.HTML(common.GetAppEmbeddedHTML()),
	}

	apiRouters := router.PathPrefix(apiEndpointPrefix).Subrouter()
	apiRouters.Use(RequireContentType("application/json"))
	apiRouters.HandleFunc("", Root)
	for _, endpoint := range GetEndpoints(db) {
		apiRouters.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HTTPMethod, http.MethodOptions)
	}

	apiRouters.HandleFunc(apiEndpointPrefix+"/{.*}", UnknownEndpoint)
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/robots.txt")
	})

	router.HandleFunc("/_healthz", Healthz(db)).Methods("GET")
	router.PathPrefix("/").Handler(FrontendHandler(common.GetAppDistFolder(), passthrough)).Methods("GET")

	router.Use(Logging)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "User-Agent", "Accept-Encoding"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("HTTP listening on", port)
	log.Fatal(srv.ListenAndServe())
}
