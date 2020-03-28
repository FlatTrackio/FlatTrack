package settings

import (
	"database/sql"
)

// GetFlatName
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

// SetFlatName
// given a flatName, set the name of the flat
func SetFlatName(db *sql.DB, flatName string) (err error) {
	sqlStatement := `update settings set value = $1 where name = $2;`
	_, err = db.Query(sqlStatement, flatName, "flatName")
	return err
}
