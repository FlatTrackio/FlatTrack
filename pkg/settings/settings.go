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
)

// GetFlatName ...
// returns the name of the flat
func GetFlatName(db *sql.DB) (flatName string, err error) {
	sqlStatement := `select value from settings where name = 'flatName'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return flatName, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&flatName)
	}
	return flatName, err
}

// SetFlatName ...
// given a flatName, set the name of the flat
func SetFlatName(db *sql.DB, flatName string) (err error) {
	if flatName == "" || len(flatName) == 0 || len(flatName) > 60 {
		return fmt.Errorf("Unable to set the flat name as it is either invalid, too short, or too long")
	}
	sqlStatement := `update settings set value = $1 where name = $2;`
	rows, err := db.Query(sqlStatement, flatName, "flatName")
	defer rows.Close()
	return err
}

// GetTimezone ...
// returns the timezone
func GetTimezone(db *sql.DB) (timezone string, err error) {
	sqlStatement := `select value from settings where name = 'timezone'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return timezone, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&timezone)
	}
	return timezone, err
}

// SetTimezone ...
// given a timezone, set the timezone of the FlatTrack instance
func SetTimezone(db *sql.DB, timezone string) (err error) {
	sqlStatement := `update settings set value = $1 where name = 'timezone';`
	rows, err := db.Query(sqlStatement, timezone)
	defer rows.Close()
	return err
}

// GetLanguage ...
// returns the language
func GetLanguage(db *sql.DB) (language string, err error) {
	sqlStatement := `select value from settings where name = 'language'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return language, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&language)
	}
	return language, err
}

// SetLanguage ...
// given a language, set the language of the FlatTrack instance
func SetLanguage(db *sql.DB, language string) (err error) {
	sqlStatement := `update settings set value = $1 where name = 'language';`
	rows, err := db.Query(sqlStatement, language)
	defer rows.Close()
	return err
}

// GetShoppingListNotes ...
// returns the shopping list notes
func GetShoppingListNotes(db *sql.DB) (notes string, err error) {
	sqlStatement := `select value from settings where name = 'shoppingListNotes'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return notes, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&notes)
	}
	return notes, err
}

// SetShoppingListNotes ...
// sets shoppingListNotes
func SetShoppingListNotes(db *sql.DB, notes string) (err error) {
	if len(notes) > 250 {
		return fmt.Errorf("Unable to set shopping list notes as it is either invalid, too short, or too long")
	}
	sqlStatement := `update settings set value = $1 where name = 'shoppingListNotes';`
	rows, err := db.Query(sqlStatement, notes)
	defer rows.Close()
	return err
}
