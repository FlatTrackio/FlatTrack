/*
  migrations
    handle database migrations
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

package migrations

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// allow file-based migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gitlab.com/flattrack/flattrack/internal/common"
)

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

// Migrate ...
// creates all the tables via the migration sql files
func (m *Manager) Migrate() (err error) {
	migrationPath := common.GetMigrationsPath()
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return err
	}
	mi, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file:///%v", migrationPath), "postgres", driver)
	if err != nil {
		return err
	}
	slog.Info("migrating database")
	err = mi.Up()
	if err != nil && err.Error() == "no change" {
		slog.Info("database is up to date")
		err = nil
	} else if err != nil && err.Error() != "no change" {
		return err
	} else if err == nil {
		slog.Info("database migrated successfully")
	}
	return err
}

// Reset ...
// removes all tables
func (m *Manager) Reset() (err error) {
	migrationPath := common.GetMigrationsPath()
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return err
	}
	mi, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", migrationPath), "postgres", driver)
	if err != nil {
		return err
	}
	slog.Info("migrating database")
	err = mi.Down()
	if err != nil && err.Error() == "no change" {
		err = nil
	} else if err != nil && err.Error() != "no change" {
		return err
	}
	return err
}
