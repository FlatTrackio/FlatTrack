/*
  migrations
    handle database migrations
*/

package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gitlab.com/flattrack/flattrack/src/backend/common"
)

func Migrate(db *sql.DB) (err error) {
	migrationPath := common.GetMigrationsPath()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", migrationPath), "postgres", driver)
	if err != nil {
		return err
	}
	log.Println("migrating database")
	err = m.Up()
	if err != nil && err.Error() == "no change" {
		log.Println("database is up to date")
		err = nil
	} else if err != nil && err.Error() != "no change" {
		return err
	} else if err == nil {
		log.Println("database migrated successfully")
	}
	return err
}
