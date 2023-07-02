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
)

// Open ...
// given database credentials, return a database connection
func Open(username string, password string, hostname string, port string, database string, sslMode string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", username, password, hostname, port, database, sslMode)
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
	var zero int
	rows, err := db.Query(`SELECT 0`)
	if err != nil {
		log.Println("Error querying database", err.Error())
		return err
	}
	if err := rows.Scan(&zero); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if zero != 0 {
		return fmt.Errorf("Wild, this error should never occur.")
	}
	return nil
}
