/*
	initialise the API
*/

package main

import (
	"log"
	"net/http"
	"time"
	// "os"
	// "strings"
	// "fmt"

	"database/sql"
	"github.com/ddo/go-vue-handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/migrations"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
)

// bring up the API
func handleWebserver(db *sql.DB) {
	port := common.GetAppPort()
	router := mux.NewRouter().StrictSlash(true)
	apiEndpointPrefix := "/api"

	for _, endpoint := range routes.GetEndpoints(apiEndpointPrefix, db) {
		router.HandleFunc(endpoint.EndpointPath, endpoint.HandlerFunc).Methods(endpoint.HttpMethod, http.MethodOptions)
	}

	router.HandleFunc(apiEndpointPrefix+"/{.*}", routes.UnknownEndpoint)
	router.HandleFunc(apiEndpointPrefix, routes.Root)
	router.PathPrefix("/").Handler(vue.Handler(common.GetAppDistFolder())).Methods("GET")

	router.Use(common.Logging)

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

// initialise the app
func main() {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())
	_ = godotenv.Load(".env")
	dbUsername := common.GetDBusername()
	dbPassword := common.GetDBpassword()
	dbHostname := common.GetDBhost()
	dbDatabase := common.GetDBdatabase()
	db, err := database.DB(dbUsername, dbPassword, dbHostname, dbDatabase)
	if err != nil {
		log.Println(err)
		return
	}
	err = migrations.Migrate(db)
	if err != nil {
		log.Println(err)
		return
	}

	handleWebserver(db)
}
