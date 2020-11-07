/*
  system
    manage system level data
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

package system

import (
	"database/sql"
)

// GetHasInitialized ...
// return if the FlatTrack instance has initialized
func GetHasInitialized(db *sql.DB) (initialized string, err error) {
	sqlStatement := `select value from system where name = 'initialized'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return initialized, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&initialized)
	}
	return initialized, err
}

// SetHasInitialized ...
// set if the FlatTrack instance has been initialized
func SetHasInitialized(db *sql.DB) (err error) {
	sqlStatement := `update system set value = 'true' where name = 'initialized'`
	rows, err := db.Query(sqlStatement)
	defer rows.Close()
	return err
}

// GetJWTsecret ...
// return the JWT secret, used in authentication
func GetJWTsecret(db *sql.DB) (jwtSecret string, err error) {
	sqlStatement := `select value from system where name = 'jwtSecret'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return jwtSecret, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&jwtSecret)
	}
	return jwtSecret, err
}
