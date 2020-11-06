/*
  groups
    manage the groups and privileges of user accounts
*/

package groups

import (
	"database/sql"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// AddUserToGroup ...
// given a userID and a groupID, adds a user to a group
func AddUserToGroup(db *sql.DB, userID string, groupID string) (err error) {
	sqlStatement := `insert into user_to_groups (userid, groupid) values ($1, $2)`
	rows, err := db.Query(sqlStatement, userID, groupID)
	defer rows.Close()
	return err
}

// RemoveUserFromGroup ...
// given a userID and a groupID, removes a user from a group
func RemoveUserFromGroup(db *sql.DB, userID string, groupID string) (err error) {
	sqlStatement := `delete from user_to_groups where userid = $1 and groupid = $2`
	rows, err := db.Query(sqlStatement, userID, groupID)
	defer rows.Close()
	return err
}

// GroupObjectFromRows ...
// constructs a group object from database rows
func GroupObjectFromRows(rows *sql.Rows) (group types.GroupSpec, err error) {
	rows.Scan(&group.ID, &group.Name, &group.DefaultGroup, &group.Description, &group.CreationTimestamp, &group.ModificationTimestamp, &group.DeletionTimestamp)
	err = rows.Err()
	return group, err
}

// GetAllGroups ...
// returns a list of all groups
func GetAllGroups(db *sql.DB) (groups []types.GroupSpec, err error) {
	sqlStatement := `select * from groups where deletionTimestamp = 0`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group types.GroupSpec
		group, err = GroupObjectFromRows(rows)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

// GetGroupByName ...
// given a group name, return the group
func GetGroupByName(db *sql.DB, name string) (group types.GroupSpec, err error) {
	sqlStatement := `select * from groups where name = $1`
	rows, err := db.Query(sqlStatement, name)
	if err != nil {
		return group, err
	}
	defer rows.Close()
	rows.Next()
	group, err = GroupObjectFromRows(rows)
	return group, err
}

// GetGroupByID ...
// given a group id, return the group
func GetGroupByID(db *sql.DB, id string) (group types.GroupSpec, err error) {
	sqlStatement := `select * from groups where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return group, err
	}
	defer rows.Close()
	rows.Next()
	group, err = GroupObjectFromRows(rows)
	return group, err
}

// GetGroupsOfUserByID ...
// given a userID, return the group which the user account belongs to
func GetGroupsOfUserByID(db *sql.DB, userID string) (groups []types.GroupSpec, err error) {
	var groupIDs []string
	sqlStatement := `select groupid from user_to_groups where userid = $1`
	rows, err := db.Query(sqlStatement, userID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var groupID string
		rows.Scan(&groupID)
		groupIDs = append(groupIDs, groupID)
	}
	for _, groupID := range groupIDs {
		group, err := GetGroupByID(db, groupID)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

// GetGroupNamesOfUserByID ...
// given a userID, return the group names which the user account belongs to
func GetGroupNamesOfUserByID(db *sql.DB, userID string) (groups []string, err error) {
	groupsFull, err := GetGroupsOfUserByID(db, userID)
	if err != nil {
		return groups, err
	}
	for _, groupItem := range groupsFull {
		groups = append(groups, groupItem.Name)
	}
	return groups, err
}

// CheckUserInGroup ...
// return bool if user is in a group
func CheckUserInGroup(db *sql.DB, userID string, group string) (found bool, err error) {
	groups, err := GetGroupNamesOfUserByID(db, userID)
	if err != nil {
		return found, err
	}
	for _, groupItem := range groups {
		if groupItem == group {
			return true, err
		}
	}
	return found, err
}

// GetDefaultGroups ...
// return a list of default groups
func GetDefaultGroups(db *sql.DB) (groups []types.GroupSpec, err error) {
	sqlStatement := `select * from groups where defaultGroup = true`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group types.GroupSpec
		group, err = GroupObjectFromRows(rows)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

// UpdateUserGroups ...
// manages a user account's groups according to what's provided
func UpdateUserGroups(db *sql.DB, userID string, groups []string) (complete bool, err error) {
	allGroups, err := GetAllGroups(db)
	if err != nil {
		return false, err
	}
	for _, group := range allGroups {
		inGroup, err := CheckUserInGroup(db, userID, group.Name)
		if err != nil {
			return false, err
		}
		shouldBeInGroup := common.StringInStringSlice(group.Name, groups)
		groupFull, err := GetGroupByName(db, group.Name)
		if inGroup == true && shouldBeInGroup == false {
			err = RemoveUserFromGroup(db, userID, groupFull.ID)
			if err != nil {
				return false, err
			}
		} else if inGroup == false && shouldBeInGroup == true {
			err = AddUserToGroup(db, userID, groupFull.ID)
			if err != nil {
				return false, err
			}
		} else {
			continue
		}
	}
	return true, err
}
