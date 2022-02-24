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

	"gitlab.com/flattrack/flattrack/pkg/common"
)

type Database struct {
	Protocol     string
	DatabaseType string
	Username     string
	Password     string
	Hostname     string
	Port         string
	Database     string
	SSLMode      string

	ConnectionString string
	DB               *sql.DB
}

func NewDatabase() *Database {
	return &Database{
		Protocol:     common.GetDBprotocol(),
		DatabaseType: common.GetDBdatabaseType(),
		Username:     common.GetDBusername(),
		Password:     common.GetDBpassword(),
		Hostname:     common.GetDBhost(),
		Port:         common.GetDBport(),
		Database:     common.GetDBdatabase(),
		SSLMode:      common.GetDBsslMode(),
	}
}

func (d *Database) GetConnectionString() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v", d.Protocol, d.Username, d.Password, d.Hostname, d.Port, d.Database, d.SSLMode)
}

// Open ...
// given database credentials, return a database connection
func (d *Database) Open() (*sql.DB, error) {
	connStr := d.GetConnectionString()
	db, err := sql.Open(d.Protocol, connStr)
	if err != nil {
		return &sql.DB{}, err
	}
	d.ConnectionString = connStr
	return db, nil
}

// Close ...
// close the connection to the database
func (d *Database) Close() (err error) {
	return d.DB.Close()
}

// Ping ...
// ping the database
func (d *Database) Ping() (err error) {
	var one int
	rows, err := d.DB.Query(`SELECT 1`)
	if err != nil {
		log.Println("Error querying database", err.Error())
		return err
	}
	rows.Scan(&one)
	if one != 0 {
		return fmt.Errorf("Wild, this error should never occur.")
	}
	return rows.Err()
}
