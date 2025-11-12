/*
  system
    manage system level data
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

package system

import (
	"database/sql"
	"encoding/json"
	"log/slog"

	"gitlab.com/flattrack/flattrack/pkg/types"
)

// Manager for system configuration
type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) getValue(name string) (output string, err error) {
	sqlStatement := `select value from system where name = $1`
	rows, err := m.db.Query(sqlStatement, name)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&output); err != nil {
			return "", err
		}
	}
	return output, nil
}

func (m *Manager) setValue(name, value string) (err error) {
	sqlStatement := `update system set value = $2 where name = $1`
	rows, err := m.db.Query(sqlStatement, name, value)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()
	return nil
}

// GetHasInitialized ...
// return if the FlatTrack instance has initialized
func (m *Manager) GetHasInitialized() (bool, error) {
	val, err := m.getValue("initialized")
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

// SetHasInitialized ...
// set if the FlatTrack instance has been initialized
func (m *Manager) SetHasInitialized() (err error) {
	return m.setValue("initialized", "true")
}

// GetJWTsecret ...
// return the JWT secret, used in authentication
func (m *Manager) GetJWTsecret() (string, error) {
	return m.getValue("jwtSecret")
}

// GetInstanceUUID returns the instance UUID
func (m *Manager) GetInstanceUUID() (string, error) {
	return m.getValue("instanceUUID")
}

// GetSchedulerLastRun ...
// returns the date of the last run of the scheduler
func (m *Manager) GetSchedulerLastRun() (types.SchedulerLastRun, error) {
	val, err := m.getValue("schedulerLastRun")
	if err != nil {
		return types.SchedulerLastRun{}, err
	}
	var lastRun types.SchedulerLastRun
	if err := json.Unmarshal([]byte(val), &lastRun); err != nil {
		return types.SchedulerLastRun{}, err
	}
	return lastRun, nil
}

// SetSchedulerLastRun ...
// set if the FlatTrack instance has been initialized
func (m *Manager) SetSchedulerLastRun(lastRun types.SchedulerLastRun) (err error) {
	b, err := json.Marshal(lastRun)
	if err != nil {
		return err
	}
	return m.setValue("schedulerLastRun", string(b))
}
