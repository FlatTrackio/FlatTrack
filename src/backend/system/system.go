package system

import (
	"database/sql"
)

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
