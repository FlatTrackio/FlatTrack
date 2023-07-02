/*
  shoppinglist
    list
      manage shopping list tags
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

package shoppinglist

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/flattrack/flattrack/pkg/types"
)

// CreateShoppingTag ...
// adds a new tag to be used in lists
func CreateShoppingTag(db *sql.DB, newTag types.ShoppingTag) (tag types.ShoppingTag, err error) {
	// validate
	valid, err := ValidateShoppingTag(db, newTag)
	if !valid || err != nil {
		return types.ShoppingTag{}, err
	}
	newTag.AuthorLast = newTag.Author
	// create
	sqlStatement := `insert into shopping_list_tag (name, author, authorLast)
                         values ($1, $2, $3)
                         returning *`
	rows, err := db.Query(sqlStatement, newTag.Name, newTag.Author, newTag.AuthorLast)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()

	// return
	tag, err = GetTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
}

// GetShoppingListTags ...
// returns a list of tags used in items in a list
func GetShoppingListTags(db *sql.DB, listID string) (tags []string, err error) {
	sqlStatement := `select distinct tag from shopping_item where listId = $1 order by tag`
	rows, err := db.Query(sqlStatement, listID)
	if err != nil {
		return []string{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return []string{}, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// GetShoppingListTag ...
// returns a list of tags used in items in a list
func GetShoppingListTag(db *sql.DB, listID string, tag string) (tagInDB string, err error) {
	sqlStatement := `select tag from shopping_item where listId = $1 and tag = $2`
	rows, err := db.Query(sqlStatement, listID, tag)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	if err := rows.Scan(&tagInDB); err != nil {
		return "", err
	}
	return tagInDB, nil
}

// UpdateShoppingListTag ...
// updates a tag's name in a list
func UpdateShoppingListTag(db *sql.DB, listID string, tag string, tagUpdate string) (tagNew string, err error) {
	tagInDB, err := GetShoppingListTag(db, listID, tag)
	if tagInDB == "" || err != nil {
		return "", fmt.Errorf("Unable to find tag to update")
	}
	if tagUpdate != "" && len(tagUpdate) == 0 || len(tagUpdate) > 30 {
		return "", fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	sqlStatement := `update shopping_item set tag = $3 where listId = $1 and tag = $2 returning tag`
	rows, err := db.Query(sqlStatement, listID, tag, tagUpdate)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	if err := rows.Scan(&tagNew); err != nil {
		return "", err
	}
	return tagNew, nil
}

// GetShoppingTag ...
// returns a tag, given an id
func GetShoppingTag(db *sql.DB, id string) (tag types.ShoppingTag, err error) {
	sqlStatement := `select * from shopping_list_tag where id = $1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	tag, err = GetTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
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
		return []types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		tag, err := GetTagObjectFromRows(rows)
		if err != nil {
			return []types.ShoppingTag{}, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// UpdateShoppingTag ...
// updates a tag's name
func UpdateShoppingTag(db *sql.DB, id string, tag types.ShoppingTag) (tagUpdated types.ShoppingTag, err error) {
	tagInDB, err := GetShoppingTag(db, id)
	if tagInDB.ID == "" || err != nil {
		return types.ShoppingTag{}, fmt.Errorf("Unable to find tag to update")
	}
	if tag.Name != "" && len(tag.Name) == 0 || len(tag.Name) > 30 {
		return types.ShoppingTag{}, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	sqlStatement := `update shopping_list_tag set name = $2, authorLast = $3, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := db.Query(sqlStatement, id, tag.Name, tag.AuthorLast)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	tagUpdated, err = GetTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tagUpdated, nil
}

// DeleteShoppingTag ...
// deletes a shopping tag
func DeleteShoppingTag(db *sql.DB, id string) (err error) {
	sqlStatement := `delete from shopping_list_tag where id = $1`
	rows, err := db.Query(sqlStatement, id)
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

// GetTagObjectFromRows ...
// returns a shopping tag object from rows
func GetTagObjectFromRows(rows *sql.Rows) (tag types.ShoppingTag, err error) {
	if err := rows.Scan(&tag.ID, &tag.Name, &tag.Author, &tag.AuthorLast, &tag.CreationTimestamp, &tag.ModificationTimestamp, &tag.DeletionTimestamp); err != nil {
		return types.ShoppingTag{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
}
