/*
  shoppinglist
    list
      manage shopping list tags
*/

package shoppinglist

import (
	"database/sql"
	"fmt"

	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// CreateShoppingTag ...
// adds a new tag to be used in lists
func CreateShoppingTag(db *sql.DB, newTag types.ShoppingTag) (tag types.ShoppingTag, err error) {
	// validate
	valid, err := ValidateShoppingTag(db, newTag)
	if !valid || err != nil {
		return tag, err
	}

	newTag.AuthorLast = newTag.Author

	// create
	sqlStatement := `insert into shopping_list_tag (name, author, authorLast)
                         values ($1, $2, $3)
                         returning *`
	rows, err := db.Query(sqlStatement, newTag.Name, newTag.Author, newTag.AuthorLast)
	if err != nil {
		return newTag, err
	}
	defer rows.Close()
	rows.Next()

	// return
	return GetTagObjectFromRows(rows)
}

// GetShoppingListTags ...
// returns a list of tags used in items in a list
func GetShoppingListTags(db *sql.DB, listID string) (tags []string, err error) {
	sqlStatement := `select distinct tag from shopping_item where listId = $1 order by tag`
	rows, err := db.Query(sqlStatement, listID)
	if err != nil {
		return tags, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags, err
}

// GetShoppingListTag ...
// returns a list of tags used in items in a list
func GetShoppingListTag(db *sql.DB, listID string, tag string) (tagInDB string, err error) {
	sqlStatement := `select tag from shopping_item where listId = $1 and tag = $2`
	rows, err := db.Query(sqlStatement, listID, tag)
	if err != nil {
		return tagInDB, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&tagInDB)
	return tagInDB, err
}

// UpdateShoppingListTag ...
// updates a tag's name in a list
func UpdateShoppingListTag(db *sql.DB, listID string, tag string, tagUpdate string) (tagNew string, err error) {
	tagInDB, err := GetShoppingListTag(db, listID, tag)
	if tagInDB == "" || err != nil {
		return tagNew, fmt.Errorf("Unable to find tag to update")
	}
	if tagUpdate != "" && len(tagUpdate) == 0 || len(tagUpdate) > 30 {
		return tagNew, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	sqlStatement := `update shopping_item set tag = $3 where listId = $1 and tag = $2 returning tag`
	rows, err := db.Query(sqlStatement, listID, tag, tagUpdate)
	if err != nil {
		return tagNew, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&tagNew)
	return tagNew, err
}

// GetShoppingTag ...
// returns a tag, given an id
func GetShoppingTag(db *sql.DB, id string) (tag types.ShoppingTag, err error) {
	sqlStatement := `select * from shopping_list_tag where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return tag, err
	}
	defer rows.Close()
	rows.Next()
	return GetTagObjectFromRows(rows)
}

// GetAllShoppingTags ...
// returns a list of all tags used in items across lists
func GetAllShoppingTags(db *sql.DB, options types.ShoppingTagOptions) (tags []types.ShoppingTag, err error) {
	sqlStatement := `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name asc`
	if options.SortBy == types.ShoppingTagSortByRecentlyUpdated {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by modificationTimestamp desc`
	} else if options.SortBy == types.ShoppingTagSortByLastUpdated {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by modificationTimestamp asc`
	} else if options.SortBy == types.ShoppingTagSortByRecentlyUpdated {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
                         order by creationTimestamp desc`
	} else if options.SortBy == types.ShoppingTagSortByLastAdded {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by creationTimestamp asc`
	} else if options.SortBy == types.ShoppingTagSortByAlphabeticalDescending {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name asc`
	} else if options.SortBy == types.ShoppingTagSortByAlphabeticalAscending {
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name desc`
	}
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return tags, err
	}
	defer rows.Close()
	for rows.Next() {
		tag, err := GetTagObjectFromRows(rows)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	return tags, err
}

// UpdateShoppingTag ...
// updates a tag's name
func UpdateShoppingTag(db *sql.DB, id string, tag types.ShoppingTag) (tagUpdated types.ShoppingTag, err error) {
	tagInDB, err := GetShoppingTag(db, id)
	if tagInDB.ID == "" || err != nil {
		return tagUpdated, fmt.Errorf("Unable to find tag to update")
	}
	if tag.Name != "" && len(tag.Name) == 0 || len(tag.Name) > 30 {
		return tagUpdated, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	sqlStatement := `update shopping_list_tag set name = $2, authorLast = $3, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := db.Query(sqlStatement, id, tag.Name, tag.AuthorLast)
	if err != nil {
		return tagUpdated, err
	}
	defer rows.Close()
	rows.Next()
	return GetTagObjectFromRows(rows)
}

// DeleteShoppingTag ...
// deletes a shopping tag
func DeleteShoppingTag(db *sql.DB, id string) (err error) {
	sqlStatement := `delete from shopping_list_tag where id = $1`
	rows, err := db.Query(sqlStatement, id)
	defer rows.Close()
	return err
}

// GetTagObjectFromRows ...
// returns a shopping tag object from rows
func GetTagObjectFromRows(rows *sql.Rows) (tag types.ShoppingTag, err error) {
	rows.Scan(&tag.ID, &tag.Name, &tag.Author, &tag.AuthorLast, &tag.CreationTimestamp, &tag.ModificationTimestamp, &tag.DeletionTimestamp)
	err = rows.Err()
	return tag, err
}
