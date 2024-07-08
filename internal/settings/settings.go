/*
  settings
    admin settings for instances
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

package settings

import (
	"database/sql"
	"fmt"
	"log"
)

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) get(key string) (output string, err error) {
	sqlStatement := `select value from settings where name = $1`
	rows, err := m.db.Query(sqlStatement, key)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		err = rows.Scan(&output)
		if err != nil {
			return "", err
		}
	}
	return output, nil
}

func (m *Manager) set(key string, value string, validation func() (err error)) (err error) {
	if err := validation(); err != nil {
		return err
	}
	sqlStatement := `update settings set value = $1 where name = $2;`
	rows, err := m.db.Query(sqlStatement, value, key)
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

// GetFlatName ...
// returns the name of the flat
func (m *Manager) GetFlatName() (output string, err error) {
	output, err = m.get("flatName")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetFlatName ...
// given a flatName, set the name of the flat
func (m *Manager) SetFlatName(value string) (err error) {
	if err := m.set("flatName", value, func() error {
		if value == "" || len(value) == 0 || len(value) > 60 {
			return fmt.Errorf("Unable to set the flat name as it is either invalid, too short, or too long")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetTimezone ...
// returns the timezone
func (m *Manager) GetTimezone() (output string, err error) {
	output, err = m.get("timezone")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetTimezone ...
// given a timezone, set the timezone of the FlatTrack instance
func (m *Manager) SetTimezone(value string) (err error) {
	if err := m.set("timezone", value, func() error { return nil }); err != nil {
		return err
	}
	return nil
}

// GetLanguage ...
// returns the language
func (m *Manager) GetLanguage() (output string, err error) {
	output, err = m.get("language")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetLanguage ...
// given a language, set the language of the FlatTrack instance
func (m *Manager) SetLanguage(value string) (err error) {
	if err := m.set("language", value, func() error { return nil }); err != nil {
		return err
	}
	return nil
}

// GetShoppingListNotes ...
// returns the shopping list notes
func (m *Manager) GetShoppingListNotes() (output string, err error) {
	output, err = m.get("shoppingListNotes")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetShoppingListNotes ...
// sets shoppingListNotes
func (m *Manager) SetShoppingListNotes(value string) (err error) {
	if err := m.set("shoppingListNotes", value, func() error {
		if len(value) > 250 {
			return fmt.Errorf("Unable to set shopping list notes as it is either invalid: too short or too long")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetCostNotes ...
// returns the cost notes
func (m *Manager) GetCostNotes() (output string, err error) {
	output, err = m.get("costsNotes")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetCostNotes ...
// sets shoppingListNotes
func (m *Manager) SetCostNotes(value string) (err error) {
	if err := m.set("costsNotes", value, func() error {
		if len(value) > 250 {
			return fmt.Errorf("Unable to set cost notes as it is either invalid: too short or too long")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetCostNotes ...
// returns the cost notes
func (m *Manager) GetCostsWriteRequireGroupAdmin() (output bool, err error) {
	val, err := m.get("costsWriteRequireGroupAdmin")
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

// SetCostNotes ...
// sets shoppingListNotes
func (m *Manager) SetCostsWriteRequireGroupAdmin(input bool) (err error) {
	value := "false"
	if input {
		value = "true"
	}
	if err := m.set("costsWriteRequireGroupAdmin", value, func() error {
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetFlatNotes ...
// returns the flat notes
func (m *Manager) GetFlatNotes() (output string, err error) {
	output, err = m.get("flatNotes")
	if err != nil {
		return "", err
	}
	return output, nil
}

// SetFlatNotes ...
// sets flat notes
func (m *Manager) SetFlatNotes(value string) (err error) {
	if err := m.set("flatNotes", value, func() error {
		if len(value) > 500 {
			return fmt.Errorf("Unable to set flat notes as it is either invalid, too short, or too long")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
