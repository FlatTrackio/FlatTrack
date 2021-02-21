/*
  database
    handle connections to the database
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

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	// include Pg
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/flattrack/flattrack/pkg/common"
)

// Database connection fields
var (
	username = common.GetDBusername()
	password = common.GetDBpassword()
	hostname = common.GetDBhost()
	database = common.GetDBdatabase()
	sslmode  = common.GetDBsslMode()
)

// DB ...
// given database credentials, return a database connection
func DB(dbType string, username string, password string, hostname string, database string) (*sql.DB, error) {
	var connectionString string
	switch (dbType) {
	case "postgres":
		connectionString = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=%v", username, password, hostname, database, sslmode)
	case "sqlite3":
		pwd, _ := os.Getwd()
		connectionString = path.Join(pwd, "flattrack-dev.db")
	default:
		log.Panicf("Invalid database type '%v'", dbType)
	}
	return sql.Open(dbType, connectionString)
}

// Close ...
// close the connection to the database
func Close(db *sql.DB) (err error) {
	return db.Close()
}

// Ping ...
// ping the database
func Ping(db *sql.DB) (err error) {
	return db.Ping()
}
