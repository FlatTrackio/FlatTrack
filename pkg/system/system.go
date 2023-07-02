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
	"log"
)

type SystemManager struct {
	DB *sql.DB
}

func (s SystemManager) getValue(name string) (output string, err error) {
	sqlStatement := `select value from system where name = $1`
	rows, err := s.DB.Query(sqlStatement, name)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&output); err != nil {
			return "", err
		}
	}
	return output, nil
}

func (s SystemManager) setValue(name, value string) (err error) {
	sqlStatement := `update system set value = $2 where name = $1`
	rows, err := s.DB.Query(sqlStatement, name, value)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// GetHasInitialized ...
// return if the FlatTrack instance has initialized
func (s SystemManager) GetHasInitialized() (string, error) {
	return s.getValue("initialized")
}

// SetHasInitialized ...
// set if the FlatTrack instance has been initialized
func (s SystemManager) SetHasInitialized() (err error) {
	return s.setValue("initialized", "true")
}

// GetJWTsecret ...
// return the JWT secret, used in authentication
func (s SystemManager) GetJWTsecret() (string, error) {
	return s.getValue("jwtSecret")
}

// GeInstanceUUID ...
// returns the instance UUID
func (s SystemManager) GetInstanceUUID() (string, error) {
	return s.getValue("instanceUUID")
}
