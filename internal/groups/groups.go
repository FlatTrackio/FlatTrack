/*
  groups
    manage the groups and privileges of user accounts
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

package groups

import (
	"database/sql"
	"log"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	GroupFlatmember = "flatmember"
	GroupAdmin      = "admin"
)

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

// AddUserToGroup ...
// given a userID and a groupID, adds a user to a group
func (m *Manager) AddUserToGroup(userID string, groupID string) (err error) {
	sqlStatement := `insert into user_to_groups (userid, groupid) values ($1, $2)`
	rows, err := m.db.Query(sqlStatement, userID, groupID)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// RemoveUserFromGroup ...
// given a userID and a groupID, removes a user from a group
func (m *Manager) RemoveUserFromGroup(userID string, groupID string) (err error) {
	sqlStatement := `delete from user_to_groups where userid = $1 and groupid = $2`
	rows, err := m.db.Query(sqlStatement, userID, groupID)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}

// groupObjectFromRows ...
// constructs a group object from database rows
func groupObjectFromRows(rows *sql.Rows) (group types.GroupSpec, err error) {
	err = rows.Scan(&group.ID, &group.Name, &group.DefaultGroup, &group.Description, &group.CreationTimestamp, &group.ModificationTimestamp, &group.DeletionTimestamp)
	if err != nil {
		return types.GroupSpec{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.GroupSpec{}, err
	}
	return group, nil
}

// List ...
// returns a list of all groups
func (m *Manager) List() (groups []types.GroupSpec, err error) {
	sqlStatement := `select * from groups where deletionTimestamp = 0`
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return []types.GroupSpec{}, err
	}
	if err != nil {
		return []types.GroupSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		var group types.GroupSpec
		group, err = groupObjectFromRows(rows)
		if err != nil {
			return []types.GroupSpec{}, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// GetByName ...
// given a group name, return the group
func (m *Manager) GetByName(name string) (group types.GroupSpec, err error) {
	sqlStatement := `select * from groups where name = $1`
	rows, err := m.db.Query(sqlStatement, name)
	if err != nil {
		return types.GroupSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	group, err = groupObjectFromRows(rows)
	if err != nil {
		return types.GroupSpec{}, err
	}
	return group, nil
}

// GetByID ...
// given a group id, return the group
func (m *Manager) GetByID(id string) (group types.GroupSpec, err error) {
	sqlStatement := `select * from groups where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.GroupSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	group, err = groupObjectFromRows(rows)
	if err != nil {
		return types.GroupSpec{}, err
	}
	return group, nil
}

// GetGroupsOfUserByID ...
// given a userID, return the group which the user account belongs to
func (m *Manager) GetGroupsOfUserByID(userID string) (groups []types.GroupSpec, err error) {
	var groupIDs []string
	sqlStatement := `select groupid from user_to_groups where userid = $1`
	rows, err := m.db.Query(sqlStatement, userID)
	if err != nil {
		return []types.GroupSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return []types.GroupSpec{}, err
		}
		groupIDs = append(groupIDs, groupID)
	}
	for _, groupID := range groupIDs {
		group, err := m.GetByID(groupID)
		if err != nil {
			return []types.GroupSpec{}, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// GetGroupNamesOfUserByID ...
// given a userID, return the group names which the user account belongs to
func (m *Manager) GetGroupNamesOfUserByID(userID string) (groups []string, err error) {
	groupsFull, err := m.GetGroupsOfUserByID(userID)
	if err != nil {
		return []string{}, err
	}
	for _, groupItem := range groupsFull {
		groups = append(groups, groupItem.Name)
	}
	return groups, nil
}

// CheckUserInGroup ...
// return bool if user is in a group
func (m *Manager) CheckUserInGroup(userID string, group string) (found bool, err error) {
	groups, err := m.GetGroupNamesOfUserByID(userID)
	if err != nil {
		return false, err
	}
	for _, groupItem := range groups {
		if groupItem == group {
			return true, nil
		}
	}
	return false, nil
}

// GetDefault ...
// return a list of default groups
func (m *Manager) GetDefault() (groups []types.GroupSpec, err error) {
	sqlStatement := `select * from groups where defaultGroup = true`
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return []types.GroupSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		var group types.GroupSpec
		group, err = groupObjectFromRows(rows)
		if err != nil {
			return []types.GroupSpec{}, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

// UpdateUserGroups ...
// manages a user account's groups according to what's provided
func (m *Manager) UpdateUserGroups(userID string, groups []string) (err error) {
	allGroups, err := m.List()
	if err != nil {
		return err
	}
	for _, group := range allGroups {
		inGroup, err := m.CheckUserInGroup(userID, group.Name)
		if err != nil {
			return err
		}
		shouldBeInGroup := common.StringInStringSlice(group.Name, groups)
		groupFull, err := m.GetByName(group.Name)
		if err != nil {
			return err
		}
		if inGroup && !shouldBeInGroup {
			if err := m.RemoveUserFromGroup(userID, groupFull.ID); err != nil {
				return err
			}
		} else if !inGroup && shouldBeInGroup {
			if err := m.AddUserToGroup(userID, groupFull.ID); err != nil {
				return err
			}
		} else {
			continue
		}
	}
	return nil
}
