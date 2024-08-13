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
	"time"
)

// Manager for system configuration
type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

type systemSetting struct {
	ID                    string
	Name                  string
	Value                 string
	CreationTimestamp     int64
	ModificationTimestamp int64
	DeletionTimestamp     int64
}

func (m *Manager) getValue(name string) (output string, err error) {
	sqlStatement := `select value from system where name = $1`
	rows, err := m.db.Query(sqlStatement, name)
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

func (m *Manager) setValue(name, value string) (err error) {
	sqlStatement := `
        update system set
          value = $2,
          modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int
        where name = $1`
	rows, err := m.db.Query(sqlStatement, name, value)
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

func (m *Manager) getSystemSetting(name string) (output systemSetting, err error) {
	sqlStatement := `select * from system where name = $1`
	rows, err := m.db.Query(sqlStatement, name)
	if err != nil {
		return systemSetting{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(
			&output.ID,
			&output.Name,
			&output.Value,
			&output.CreationTimestamp,
			&output.ModificationTimestamp,
			&output.DeletionTimestamp,
		); err != nil {
			return systemSetting{}, err
		}
	}
	return output, nil
}

// GetHasInitialized ...
// return if the FlatTrack instance has initialized
func (m *Manager) GetHasInitialized() (bool, error) {
	val, err := m.getValue("initialized")
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

// SetHasInitialized ...
// set if the FlatTrack instance has been initialized
func (m *Manager) SetHasInitialized() (err error) {
	return m.setValue("initialized", "true")
}

// GetJWTsecret ...
// return the JWT secret, used in authentication
func (m *Manager) GetJWTsecret() (string, error) {
	return m.getValue("jwtSecret")
}

func (m *Manager) generateJWTSecret() (err error) {
	sqlStatement := `
        update system set
          value = md5(random()::text || clock_timestamp()::text)::uuid,
          modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int
          where name = 'jwtSecret'`
	rows, err := m.db.Query(sqlStatement)
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

// GetInstanceUUID returns the instance UUID
func (m *Manager) GetInstanceUUID() (string, error) {
	return m.getValue("instanceUUID")
}

func (m *Manager) RefreshJWTSecret() error {
	// TODO factor last signed JWT and usage of it
	const oneMonth = 32 * 24 * time.Hour
	setting, err := m.getSystemSetting("jwtSecret")
	if err != nil {
		return err
	}
	mod := time.Unix(setting.ModificationTimestamp, 0)
	age := time.Now().Sub(mod)
	if age < oneMonth {
		return nil
	}
	if err := m.generateJWTSecret(); err != nil {
		return err
	}
	log.Println("[system]: updated jwt secret")
	return nil
}
