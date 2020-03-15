package migrations

import (
	"database/sql"
	"fmt"

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
	if err := m.Up(); err != nil {
		return err
	}
	return err
}
