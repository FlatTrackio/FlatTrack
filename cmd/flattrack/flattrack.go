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
	"log"

	"github.com/joho/godotenv"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
	"gitlab.com/flattrack/flattrack/pkg/files"
	"gitlab.com/flattrack/flattrack/pkg/metrics"
	"gitlab.com/flattrack/flattrack/pkg/migrations"
	"gitlab.com/flattrack/flattrack/pkg/routes"
	"gitlab.com/flattrack/flattrack/pkg/system"
)

// Start ...
// initialise the app
func Start() {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)

	log.Println(common.GetMigrationsPath(), common.GetAppDistFolder())

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
	fileAccess, err := files.Open(minioHost, minioAccessKey, minioSecretKey, minioBucket, minioUseSSL == "true")
	if err != nil {
		log.Println("Minio error:", err)
		return
	}

	systemManager := system.SystemManager{DB: db}
	systemUUID, err := systemManager.GetInstanceUUID()
	if err != nil {
		log.Println("Error getting system UUID:", err)
	}
	fileAccess.Prefix = systemUUID

	router := routes.Router{
		DB:         db,
		FileAccess: fileAccess,
	}

	go func() {
		if router.FileAccess.Client == nil {
			log.Println("Error: no Minio client available, will not serve files")
			return
		}
		if router.FileAccess.BucketName == "" {
			log.Println("Error: no Minio bucket name was provided")
			return
		}
		err = router.FileAccess.Init()
		if err != nil {
			log.Println("Error initialising Minio bucket:", err)
		}
	}()
	go metrics.Handle()
	go routes.HealthHandler(db)
	router.Handle()
}
