package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.com/flattrack/flattrack/src/backend/common"
)

var (
	DB_USERNAME = common.GetDBusername()
	DB_PASSWORD = common.GetDBpassword()
	DB_HOSTNAME = common.GetDBhost()
	DB_DATABASE = common.GetDBdatabase()
)

func DB(username string, password string, hostname string, database string) (*sql.DB, error) {
	username = common.SetFirstOrSecond(username, DB_USERNAME)
	password = common.SetFirstOrSecond(password, DB_PASSWORD)
	hostname = common.SetFirstOrSecond(hostname, DB_HOSTNAME)
	database = common.SetFirstOrSecond(database, DB_DATABASE)
	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", username, password, hostname, database)
	return sql.Open("postgres", connStr)
}
