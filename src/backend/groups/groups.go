/*
  groups
    manage the groups and privileges of user accounts
*/

package groups

import (
	"database/sql"

	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// AddUserToGroup
// given a userId and a groupId, adds a user to a group
func AddUserToGroup(db *sql.DB, userId string, groupId string) (err error) {
	sqlStatement := `insert into user_to_groups (userid, groupid) values ($1, $2)`
	_, err = db.Query(sqlStatement, userId, groupId)
	return err
}

// RemoveUserToGroup
// given a userId and a groupId, removes a user from a group
func RemoveUserFromGroup(db *sql.DB, userId string, groupId string) (err error) {
	sqlStatement := `delete from user_to_groups where userid = $1 and groupid = $2`
	_, err = db.Query(sqlStatement, userId, groupId)
	return err
}

// GetAllGroups
// returns a list of all groups
func GetAllGroups(db *sql.DB) (groups []types.GroupSpec, err error) {
	sqlStatement := `select * from groups`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group types.GroupSpec
		rows.Scan(&group.Id, &group.Name, &group.DefaultGroup)
		groups = append(groups, group)
	}
	return groups, err
}

// GetGroupByName
// given a group name, return the group
func GetGroupByName(db *sql.DB, name string) (group types.GroupSpec, err error) {
	sqlStatement := `select id, name, defaultGroup from groups where name = $1`
	rows, err := db.Query(sqlStatement, name)
	if err != nil {
		return group, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&group.Id, &group.Name, &group.DefaultGroup)
	}
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
	for rows.Next() {
		rows.Scan(&group.Id, &group.Name, &group.DefaultGroup)
	}
	return group, err
}

// GetGroupsOfUserById
// given a userId, return the group which the user belongs to
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
