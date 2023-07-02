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
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/files"
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
	if _, err := w.Write(response); err != nil {
		log.Printf("error: failed to write response; %v\n", err)
	}
}

// GetHTTPresponseBodyContents ...
// convert the body of a HTTP response into a JSONMessageResponse
func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(responseData, &output); err != nil {
		log.Printf("error: failed to unmarshal response body contents; %v", err)
	}
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
	router := mux.NewRouter().StrictSlash(true)
	r := router.Handle("/_healthz", Healthz(db))
	server := &http.Server{
		Handler:           r.GetHandler(),
		Addr:              common.GetAppMetricsPort(),
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("Health listening on %v", port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
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
		} else if len(requestPath) >= 1 && requestPath[1] == "files" {
			pathSection = "files   "
		}
		requestIP := GetRequestIP(r)
		log.Printf("[%v] %v %v %v %v %v", pathSection, r.Method, r.URL, r.Proto, r.Response, requestIP)
		next.ServeHTTP(w, r)
	})
}

// HTTPHeaderBackendAllowTypes heads to check for content type
type HTTPHeaderBackendAllowTypes string

const (
	// HTTPHeaderBackendAllowTypesContentType use the Content-Type http header
	HTTPHeaderBackendAllowTypesContentType HTTPHeaderBackendAllowTypes = "Content-Type"
	// HTTPHeaderBackendAllowTypesAccept use the Accept http header
	HTTPHeaderBackendAllowTypesAccept HTTPHeaderBackendAllowTypes = "Accept"
)

// RequireContentType ...
// 404s requests if content-type isn't what is expected
func RequireContentType(expectedContentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			foundInRequiredTypes := false
			for _, t := range []HTTPHeaderBackendAllowTypes{HTTPHeaderBackendAllowTypesContentType, HTTPHeaderBackendAllowTypesAccept} {
				if r.Header.Get(string(t)) == expectedContentType {
					foundInRequiredTypes = true
					break
				}
			}
			if foundInRequiredTypes {
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

// Router http router
type Router struct {
	DB         *sql.DB
	FileAccess files.FileAccess
}

// Handle ...
// manage the launching of the API's webserver
func (r Router) Handle() {
	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(true)
	apiEndpointPrefix := "/api"
	passthrough := FrontendOptions{
		SetupMessage: common.GetAppSetupMessage(),
		LoginMessage: common.GetAppLoginMessage(),
	}

	apiRouters := router.PathPrefix(apiEndpointPrefix).Subrouter()
	apiRouters.Use(RequireContentType("application/json"))
	apiRouters.HandleFunc("", Root)
	for _, endpoint := range GetEndpoints(r.DB) {
		apiRouters.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HTTPMethod, http.MethodOptions)
	}

	apiRouters.HandleFunc(apiEndpointPrefix+"/{.*}", UnknownEndpoint)
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/robots.txt")
	})

	router.HandleFunc("/_healthz", Healthz(r.DB)).Methods(http.MethodGet)
	fileRouter := router.PathPrefix("/files").Subrouter()
	fileRouter.HandleFunc("/{.*}", r.GetServeFilestoreObjects("/files")).Methods(http.MethodGet)
	// TODO add files post

	router.PathPrefix("/").Handler(FrontendHandler(common.GetAppDistFolder(), passthrough)).Methods(http.MethodGet)

	router.Use(Logging)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "User-Agent", "Accept-Encoding"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := gziphandler.GzipHandler(c.Handler(router))

	srv := &http.Server{
		Handler:      handler,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("HTTP listening on", port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
