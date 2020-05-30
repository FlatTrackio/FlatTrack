/*
  database
    handle connections to the database
*/

package database

import (
	"database/sql"
	"fmt"

	// include Pg
	_ "github.com/lib/pq"
	"gitlab.com/flattrack/flattrack/src/backend/common"
)

// Database connection fields
var (
	DbUsername = common.GetDBusername()
	DbPassword = common.GetDBpassword()
	DbHostname = common.GetDBhost()
	DbDatabase = common.GetDBdatabase()
)

// DB ...
// given database credentials, return a database connection
func DB(username string, password string, hostname string, database string) (*sql.DB, error) {
	username = common.SetFirstOrSecond(username, DbUsername)
	password = common.SetFirstOrSecond(password, DbPassword)
	hostname = common.SetFirstOrSecond(hostname, DbHostname)
	database = common.SetFirstOrSecond(database, DbDatabase)
	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", username, password, hostname, database)
	return sql.Open("postgres", connStr)
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
