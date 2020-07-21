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
	username = common.GetDBusername()
	password = common.GetDBpassword()
	hostname = common.GetDBhost()
	database = common.GetDBdatabase()
	sslmode = common.GetDBsslMode()
)

// DB ...
// given database credentials, return a database connection
func DB(username string, password string, hostname string, database string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=%v", username, password, hostname, database, sslmode)
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
