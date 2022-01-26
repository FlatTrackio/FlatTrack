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

// Package classification for FlatTrack API.
//
//     Schemes: http
//     Host: localhost
//     BasePath: /api
//     Version: 0.16.1
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

// Package flattrack ...
// backend cmd
package flattrack

import (
	"github.com/joho/godotenv"
	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
	"gitlab.com/flattrack/flattrack/pkg/metrics"
	"gitlab.com/flattrack/flattrack/pkg/migrations"
	"gitlab.com/flattrack/flattrack/pkg/routes"
	"gitlab.com/flattrack/flattrack/pkg/files"
	"log"
)

// Start ...
// initialise the app
func Start() {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)

	dbUsername := common.GetDBusername()
	dbPassword := common.GetDBpassword()
	dbHostname := common.GetDBhost()
	dbPort := common.GetDBport()
	dbDatabase := common.GetDBdatabase()
	dbSSLmode := common.GetDBsslMode()
	db, err := database.Open(dbUsername, dbPassword, dbHostname, dbPort, dbDatabase, dbSSLmode)
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
	minioBucket := common.GetAppMinioBucket()
	mc, err := files.Open(minioHost, minioAccessKey, minioSecretKey, minioUseSSL == "true")
	if err != nil {
		log.Println("Minio error:", err)
		return
	}

	go func(){
		err = files.Init(mc, minioBucket)
		if err != nil {
			log.Println("Minio error initialising bucket:", err)
		}
	}()
	go metrics.Handle()
	go routes.HealthHandler(db)
	routes.Handle(db, mc)
}
