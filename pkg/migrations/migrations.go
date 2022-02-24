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
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	// allow file-based migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
)

type Migration struct {
	DBConfig      *database.Database
	MigrationPath string
}

func NewMigration(dbConfig *database.Database) *Migration {
	return &Migration{
		DBConfig:      dbConfig,
		MigrationPath: common.GetMigrationsPath(),
	}
}

func (mi *Migration) migraterPostgres() (m *migrate.Migrate, err error) {
	driver, err := postgres.WithInstance(mi.DBConfig.DB, &postgres.Config{})
	if err != nil {
		return &migrate.Migrate{}, err
	}
	m, err = migrate.NewWithDatabaseInstance(fmt.Sprintf("file:///%v", mi.MigrationPath), mi.DBConfig.DatabaseType, driver)
	if err != nil {
		return &migrate.Migrate{}, err
	}
	return m, nil
}

func (mi Migration) migraterCockroach() (m *migrate.Migrate, err error) {
	c := &cockroachdb.CockroachDb{}
	d, err := c.Open(mi.DBConfig.ConnectionString)
	if err != nil {
		return &migrate.Migrate{}, err
	}
	m, err = migrate.NewWithDatabaseInstance("file://./migrations", "migrate", d)
	if err != nil {
		return &migrate.Migrate{}, err
	}
	return m, nil
}

// Migrate ...
// creates all the tables via the migration sql files
func (mi *Migration) Migrate() (err error) {
	var m *migrate.Migrate
	if mi.DBConfig.DatabaseType == "cockroachdb" {
		m, err = mi.migraterCockroach()
		if err != nil {
			return err
		}
	} else {
		m, err = mi.migraterPostgres()
		if err != nil {
			return err
		}
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
	return nil
}

// Reset ...
// removes all tables
func (mi *Migration) Reset() (err error) {
	m, err := mi.migraterPostgres()
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
	return nil
}
