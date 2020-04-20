/*
  groups
    manage the groups and privileges of user accounts
*/

package groups

import (
	"database/sql"

	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// AddUserToGroup
// given a userId and a groupId, adds a user to a group
func AddUserToGroup(db *sql.DB, userId string, groupId string) (err error) {
	sqlStatement := `insert into user_to_groups (userid, groupid) values ($1, $2)`
	rows, err := db.Query(sqlStatement, userId, groupId)
	defer rows.Close()
	return err
}

// RemoveUserToGroup
// given a userId and a groupId, removes a user from a group
func RemoveUserFromGroup(db *sql.DB, userId string, groupId string) (err error) {
	sqlStatement := `delete from user_to_groups where userid = $1 and groupid = $2`
	rows, err := db.Query(sqlStatement, userId, groupId)
	defer rows.Close()
	return err
}

// GroupObjectFromRows
// constructs a group object from database rows
func GroupObjectFromRows(rows *sql.Rows) (group types.GroupSpec, err error) {
	rows.Scan(&group.Id, &group.Name, &group.DefaultGroup, &group.Description, &group.CreationTimestamp, &group.ModificationTimestamp, &group.DeletionTimestamp)
	err = rows.Err()
	return group, err
}

// GetAllGroups
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

// GetGroupByName
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

// GetGroupById
// given a group id, return the group
func GetGroupById(db *sql.DB, id string) (group types.GroupSpec, err error) {
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

// GetGroupsOfUserById
// given a userId, return the group which the user account belongs to
func GetGroupsOfUserById(db *sql.DB, userId string) (groups []types.GroupSpec, err error) {
	var groupIds []string
	sqlStatement := `select groupid from user_to_groups where userid = $1`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var groupId string
		rows.Scan(&groupId)
		groupIds = append(groupIds, groupId)
	}
	for _, groupId := range groupIds {
		group, err := GetGroupById(db, groupId)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

// GetGroupNamesOfUserById
// given a userId, return the group names which the user account belongs to
func GetGroupNamesOfUserById(db *sql.DB, userId string) (groups []string, err error) {
	groupsFull, err := GetGroupsOfUserById(db, userId)
	if err != nil {
		return groups, err
	}
	for _, groupItem := range groupsFull {
		groups = append(groups, groupItem.Name)
	}
	return groups, err
}

// CheckUserInGroup
// return bool if user is in a group
func CheckUserInGroup(db *sql.DB, userId string, group string) (found bool, err error) {
	groups, err := GetGroupNamesOfUserById(db, userId)
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

// GetDefaultGroups
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

// UpdateUserGroups
// manages a user account's groups according to what's provided
func UpdateUserGroups(db *sql.DB, userId string, groups []string) (complete bool, err error) {
	allGroups, err := GetAllGroups(db)
	if err != nil {
		return false, err
	}
	for _, group := range allGroups {
		inGroup, err := CheckUserInGroup(db, userId, group.Name)
		if err != nil {
			return false, err
		}
		shouldBeInGroup := common.StringInStringSlice(group.Name, groups)
		groupFull, err := GetGroupByName(db, group.Name)
		if inGroup == true && shouldBeInGroup == false {
			err = RemoveUserFromGroup(db, userId, groupFull.Id)
			if err != nil {
				return false, err
			}
		} else if inGroup == false && shouldBeInGroup == true {
			err = AddUserToGroup(db, userId, groupFull.Id)
			if err != nil {
				return false, err
			}
		} else {
			continue
		}
	}
	return true, err
}
