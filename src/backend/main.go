/*
  main
    initialize the app
*/

// Package classification for FlatTrack API.
//
//     Schemes: http
//     Host: localhost
//     BasePath: /api
//     Version: 0.0.1-alpha6
//     License: AGPL-3.0 http://www.gnu.org/licenses/agpl-3.0.html
//     Contact: Caleb Woodbine <calebwoodbine.public@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"github.com/joho/godotenv"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/files"
	"gitlab.com/flattrack/flattrack/src/backend/metrics"
	"gitlab.com/flattrack/flattrack/src/backend/migrations"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
	"log"
	"strconv"
)

// main
// initialise the app
func main() {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)

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
		log.Println("migrations:", err)
		return
	}

	minioHost := common.GetAppMinioHost()
	minioAccessKey := common.GetAppMinioAccessKey()
	minioSecretKey := common.GetAppMinioSecretKey()
	minioUseSSL := common.GetAppMinioUseSSL()
	minioUseSSLBool, err := strconv.ParseBool(minioUseSSL)
	if err != nil {
		log.Println(err)
	}
	minioClient, err := files.Open(minioHost, minioAccessKey, minioSecretKey, minioUseSSLBool)
	if err != nil {
		log.Println(err)
		return
	}

	go metrics.Handle()
	go routes.HealthHandler(db)
	routes.Handle(db, minioClient)
}
