/*
  shoppinglist
    list
      manage shopping lists
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

	"github.com/imdario/mergo"

	"gitlab.com/flattrack/flattrack/pkg/types"
)

// ValidateShoppingList ...
// given a shopping list, return it's validity
func ValidateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec) (valid bool, err error) {
	if len(shoppingList.Name) == 0 || len(shoppingList.Name) >= 30 || shoppingList.Name == "" {
		return valid, fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if shoppingList.Notes != "" && len(shoppingList.Notes) > 100 {
		return valid, fmt.Errorf("Unable to save shopping list notes, as they are too long")
	}
	if shoppingList.TemplateID != "" {
		list, err := GetShoppingList(db, shoppingList.TemplateID)
		if err != nil || list.ID == "" {
			return valid, fmt.Errorf("Unable to find list to use as template from provided id")
		}
	}
	return true, err
}

// ValidateShoppingTag ...
// given a shopping tag, return it's validity
func ValidateShoppingTag(db *sql.DB, tag types.ShoppingTag) (valid bool, err error) {
	if tag.Name != "" && len(tag.Name) == 0 || len(tag.Name) >= 30 {
		return valid, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	return true, err
}

// GetShoppingLists ...
// returns a list of all shopping lists (name, notes, author, etc...)
func GetShoppingLists(db *sql.DB, options types.ShoppingListOptions) (shoppingLists []types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list where deletionTimestamp = 0 `
	fields := []interface{}{}

	if options.Selector.ModificationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and modificationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.ModificationTimestampAfter)
	}
	if options.Selector.CreationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and creationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.CreationTimestampAfter)
	}

	if options.SortBy == types.ShoppingListSortByRecentlyUpdated {
		sqlStatement += `order by modificationTimestamp desc `
	} else if options.SortBy == types.ShoppingListSortByLastUpdated {
		sqlStatement += `order by modificationTimestamp asc `
	} else if options.SortBy == types.ShoppingListSortByRecentlyAdded {
		sqlStatement += `order by creationTimestamp asc `
	} else if options.SortBy == types.ShoppingListSortByLastAdded {
		sqlStatement += `order by creationTimestamp asc `
	} else if options.SortBy == types.ShoppingListSortByAlphabeticalDescending {
		sqlStatement += `order by name asc `
	} else if options.SortBy == types.ShoppingListSortByAlphabeticalAscending {
		sqlStatement += `order by name desc `
	} else {
		sqlStatement += `order by creationTimestamp desc `
	}

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(`limit $%v `, len(fields)+1)
		fields = append(fields, options.Limit)
	}

	rows, err := db.Query(sqlStatement, fields...)
	if err != nil {
		return shoppingLists, err
	}
	defer rows.Close()
	for rows.Next() {
		shoppingList, err := GetListObjectFromRows(rows)
		if err != nil {
			return shoppingLists, err
		}
		shoppingList.Count, err = GetListItemCount(db, shoppingList.ID)
		if err != nil {
			return shoppingLists, err
		}

		if options.Selector.Completed == "true" && shoppingList.Completed != true {
			continue
		} else if options.Selector.Completed == "false" && shoppingList.Completed != false {
			continue
		}
		shoppingLists = append(shoppingLists, shoppingList)
	}
	return shoppingLists, err
}

// GetShoppingList ...
// returns a given shopping list, by it's ID
func GetShoppingList(db *sql.DB, listID string) (shoppingList types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list where id = $1 and deletionTimestamp = 0`
	rows, err := db.Query(sqlStatement, listID)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()
	rows.Next()
	shoppingList, err = GetListObjectFromRows(rows)
	if err != nil {
		return shoppingList, err
	}
	shoppingList.Count, err = GetListItemCount(db, shoppingList.ID)
	return shoppingList, err
}

// CreateShoppingList ...
// creates a shopping list for adding items to
func CreateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec, options types.ShoppingItemOptions) (shoppingListInserted types.ShoppingListSpec, err error) {
	valid, err := ValidateShoppingList(db, shoppingList)
	if !valid || err != nil {
		return shoppingListInserted, err
	}

	shoppingList.AuthorLast = shoppingList.Author
	shoppingList.Completed = false

	sqlStatement := `insert into shopping_list (name, notes, author, authorLast, completed, templateId)
                         values ($1, $2, $3, $4, $5, $6)
                         returning *`
	rows, err := db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.Author, shoppingList.AuthorLast, shoppingList.Completed, shoppingList.TemplateID)
	if err != nil {
		return shoppingListInserted, err
	}
	defer rows.Close()
	rows.Next()
	shoppingListInserted, err = GetListObjectFromRows(rows)
	if err != nil || shoppingListInserted.ID == "" {
		return shoppingListInserted, fmt.Errorf("Failed to create shopping list")
	}
	if shoppingList.TemplateID == "" {
		return shoppingListInserted, err
	}

	// if using other list as a template
	shoppingListItems, err := GetShoppingListItems(db, shoppingList.TemplateID, options)
	if err != nil {
		return shoppingListInserted, fmt.Errorf("Failed to fetch items from shopping list")
	}

	for _, item := range shoppingListItems {
		newItem := types.ShoppingItemSpec{
			Name:       item.Name,
			Notes:      item.Notes,
			Price:      item.Price,
			Quantity:   item.Quantity,
			Tag:        item.Tag,
			Author:     shoppingList.Author,
			AuthorLast: shoppingList.Author,
			TemplateID: shoppingList.TemplateID,
		}
		newItem, err := AddItemToList(db, shoppingListInserted.ID, newItem)
		if err != nil || newItem.ID == "" {
			return shoppingListInserted, fmt.Errorf("Failed to add new item to new shopping list from template")
		}
	}

	return shoppingListInserted, err
}

// PatchShoppingList ...
// patches a shopping list
func PatchShoppingList(db *sql.DB, listID string, shoppingList types.ShoppingListSpec) (shoppingListPatched types.ShoppingListSpec, err error) {
	existingList, err := GetShoppingList(db, listID)
	if err != nil || existingList.ID == "" {
		return shoppingListPatched, fmt.Errorf("Failed to fetch existing shopping list")
	}
	err = mergo.Merge(&shoppingList, existingList)
	if err != nil {
		return shoppingListPatched, fmt.Errorf("Failed to update fields in the item")
	}
	valid, err := ValidateShoppingList(db, existingList)
	if !valid || err != nil {
		return shoppingListPatched, err
	}

	sqlStatement := `update shopping_list set name = $1, notes = $2, authorLast = $3, completed = $4, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $5
                         returning *`
	rows, err := db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.AuthorLast, shoppingList.Completed, listID)
	if err != nil {
		return shoppingListPatched, err
	}
	defer rows.Close()
	rows.Next()
	shoppingListPatched, err = GetListObjectFromRows(rows)
	if err != nil || shoppingListPatched.ID == "" {
		return shoppingListPatched, fmt.Errorf("Failed to patch shopping list")
	}
	return shoppingListPatched, err
}

// UpdateShoppingList ...
// updates a shopping list
func UpdateShoppingList(db *sql.DB, listID string, shoppingList types.ShoppingListSpec) (shoppingListUpdated types.ShoppingListSpec, err error) {
	valid, err := ValidateShoppingList(db, shoppingList)
	if !valid || err != nil {
		return shoppingListUpdated, err
	}

	sqlStatement := `update shopping_list set name = $1, notes = $2, authorLast = $3, completed = $4, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $5
                         returning *`
	rows, err := db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.AuthorLast, shoppingList.Completed, listID)
	if err != nil {
		return shoppingListUpdated, err
	}
	defer rows.Close()
	rows.Next()
	shoppingListUpdated, err = GetListObjectFromRows(rows)
	if err != nil || shoppingListUpdated.ID == "" {
		return shoppingListUpdated, fmt.Errorf("Failed to create shopping list")
	}
	return shoppingListUpdated, err
}

// SetListCompleted ...
// updates the list's completed field
func SetListCompleted(db *sql.DB, listID string, completed bool, userID string) (list types.ShoppingListSpec, err error) {
	sqlStatement := `update shopping_list set completed = $1 where id = $2 returning *`
	rows, err := db.Query(sqlStatement, completed, listID)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		list, err = GetListObjectFromRows(rows)
		if err != nil {
			return list, err
		}
	}

	return list, err
}

// GetListObjectFromRows ...
// returns a shopping list object from rows
func GetListObjectFromRows(rows *sql.Rows) (list types.ShoppingListSpec, err error) {
	rows.Scan(&list.ID, &list.Name, &list.Notes, &list.Author, &list.AuthorLast, &list.Completed, &list.CreationTimestamp, &list.ModificationTimestamp, &list.DeletionTimestamp, &list.TemplateID)
	err = rows.Err()
	return list, err
}

// DeleteShoppingList ...
// deletes a shopping list, given a shopping list Id
func DeleteShoppingList(db *sql.DB, listID string) (err error) {
	err = RemoveAllItemsFromList(db, listID)
	if err != nil {
		return fmt.Errorf("Failed to remove all items from list")
	}

	sqlStatement := `delete from shopping_list where id = $1`
	rows, err := db.Query(sqlStatement, listID)
	defer rows.Close()
	return err
}

// GetListCount ...
// returns a count lists
func GetListCount(db *sql.DB) (count int, err error) {
	sqlStatement := `select count(*) from shopping_list`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return count, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)

	return count, err
}
