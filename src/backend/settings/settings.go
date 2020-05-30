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
