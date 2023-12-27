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

	// include Pg
	_ "github.com/lib/pq"
	"gitlab.com/flattrack/flattrack/internal/common"
)

var (
	connectionString = ""
)

type databaseConnection struct {
	username string
	password string
	host     string
	port     string
	database string
	sslMode  string
}

func panicOnEmptyRequiredDatabaseField(name string, input string) {
	if input == "" {
		log.Panicf("error: database connection field '%v' must be not empty", name)
	}
}

func GetConnectionString() (output string) {
	if connectionString != "" {
		return connectionString
	}
	output = common.GetDBConnectionString()
	if output != "" {
		return output
	}
	conn := databaseConnection{
		username: common.GetDBusername(),
		password: common.GetDBpassword(),
		host:     common.GetDBhost(),
		port:     common.GetDBport(),
		database: common.GetDBdatabase(),
		sslMode:  common.GetDBsslMode(),
	}
	panicOnEmptyRequiredDatabaseField("username", conn.username)
	panicOnEmptyRequiredDatabaseField("password", conn.password)
	panicOnEmptyRequiredDatabaseField("host", conn.host)
	panicOnEmptyRequiredDatabaseField("port", conn.port)
	panicOnEmptyRequiredDatabaseField("database", conn.database)
	panicOnEmptyRequiredDatabaseField("sslMode", conn.sslMode)

	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		conn.username, conn.password, conn.host, conn.port, conn.database, conn.sslMode)
}

// Open ...
// given database credentials, return a database connection
func Open() (*sql.DB, error) {
	connectionString = GetConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close ...
// close the connection to the database
func Close(db *sql.DB) (err error) {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

// Ping ...
// ping the database
func Ping(db *sql.DB) (err error) {
	var zero int
	rows, err := db.Query(`SELECT 0`)
	if err != nil {
		log.Println("Error querying database", err.Error())
		return err
	}
	rows.Next()
	if err := rows.Scan(&zero); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if zero != 0 {
		return fmt.Errorf("wild, this error should never occur")
	}
	return nil
}
