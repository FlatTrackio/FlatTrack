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
	"log"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	// allow file-based migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"gitlab.com/flattrack/flattrack/pkg/common"
)

type Migrations struct {
	DB *sql.DB
	Folder embed.FS
}

// Migrate ...
// creates all the tables via the migration sql files
func (mi Migrations) Migrate() (err error) {
	var m *migrate.Migrate
	driver, err := postgres.WithInstance(mi.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	_, err = mi.Folder.ReadDir("migrations")
	if err != nil {
		return err
	}
	d, err := iofs.New(mi.Folder, "migrations")
	if err != nil {
		return err
	}
	m, err = migrate.NewWithInstance("file:///migrations", d, "postgres", driver)
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

// Reset ...
// removes all tables
func Reset(db *sql.DB) (err error) {
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
	err = m.Down()
	if err != nil && err.Error() == "no change" {
		err = nil
	} else if err != nil && err.Error() != "no change" {
		return err
	}
	return err
}
