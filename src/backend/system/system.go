/*
  system
    manage system level data
*/

package system

import (
	"database/sql"
)

// GetHasInitialized
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

// GetJWTsecret
// return the JWT secret, used in authentication
func GetJWTsecret(db *sql.DB) (jwtSecret string, err error) {
	sqlStatement := `select value from system where name = 'jetSecret'`
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
