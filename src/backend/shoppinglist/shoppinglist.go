/*
  shoppinglist
    manage shopping lists
*/

package shoppinglist

import (
	"database/sql"
	"fmt"

	"github.com/imdario/mergo"

	"gitlab.com/flattrack/flattrack/src/backend/types"
)

// ValidateShoppingList ...
// given a shopping list, return it's validitiy
func ValidateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec) (valid bool, err error) {
	if len(shoppingList.Name) == 0 || len(shoppingList.Name) > 30 || shoppingList.Name == "" {
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

// ValidateShoppingListItem ...
// given a shopping list item, return it's validitiy
func ValidateShoppingListItem(db *sql.DB, item types.ShoppingItemSpec) (valid bool, err error) {
	if len(item.Name) == 0 || len(item.Name) > 30 || item.Name == "" {
		return valid, fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if item.Notes != "" && len(item.Notes) > 40 {
		return valid, fmt.Errorf("Unable to save shopping list notes, as they are too long")
	}
	if item.Tag != "" && len(item.Tag) == 0 || len(item.Tag) > 30 {
		return valid, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	if item.Quantity < 1 {
		return valid, fmt.Errorf("Item quantity must be at least one")
	}
	return true, err
}

// GetShoppingLists ...
// returns a list of all shopping lists (name, notes, author, etc...)
func GetShoppingLists(db *sql.DB) (shoppingLists []types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list
                         where deletionTimestamp = 0
	                 order by creationTimestamp desc`

	rows, err := db.Query(sqlStatement)
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

// GetShoppingListItems ...
// returns a list of items on a shopping list
func GetShoppingListItems(db *sql.DB, listID string, options types.ShoppingItemOptions) (items []types.ShoppingItemSpec, err error) {
	sqlStatement := `select * from shopping_item where listId = $1 order by tag asc, name asc`
	if options.SortBy == types.ShoppingItemSortByHighestPrice {
		sqlStatement = `select * from shopping_item where listId = $1 order by price desc, name asc`
	} else if options.SortBy == types.ShoppingItemSortByHighestQuantity {
		sqlStatement = `select * from shopping_item where listId = $1 order by quantity desc, name desc`
	} else if options.SortBy == types.ShoppingItemSortByLowestPrice {
		sqlStatement = `select * from shopping_item where listId = $1 order by price asc, name asc`
	} else if options.SortBy == types.ShoppingItemSortByLowestQuantity {
		sqlStatement = `select * from shopping_item where listId = $1 order by quantity asc, name desc`
	} else if options.SortBy == types.ShoppingItemSortByRecentlyAdded {
		sqlStatement = `select * from shopping_item where listId = $1 order by creationTimestamp desc`
	} else if options.SortBy == types.ShoppingItemSortByRecentlyUpdated {
		sqlStatement = `select * from shopping_item where listId = $1 order by modificationTimestamp desc`
	} else if options.SortBy == types.ShoppingItemSortByLastAdded {
		sqlStatement = `select * from shopping_item where listId = $1 order by creationTimestamp asc`
	} else if options.SortBy == types.ShoppingItemSortByLastUpdated {
		sqlStatement = `select * from shopping_item where listId = $1 order by modificationTimestamp asc`
	} else if options.SortBy == types.ShoppingItemSortByAlphabeticalDescending {
		sqlStatement = `select * from shopping_item where listId = $1 order by name asc`
	} else if options.SortBy == types.ShoppingItemSortByAlphabeticalAscending {
		sqlStatement = `select * from shopping_item where listId = $1 order by name desc`
	}
	rows, err := db.Query(sqlStatement, listID)
	if err != nil {
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err := GetItemObjectFromRows(rows)
		if err != nil {
			return items, err
		}
		if options.Selector.TemplateListItemSelector == "obtained" {
			if item.Obtained != true {
				continue
			}
		} else if options.Selector.TemplateListItemSelector == "unobtained" {
			if item.Obtained != false {
				continue
			}
		}
		items = append(items, item)
	}
	return items, err
}

// GetShoppingListItem ...
// given an item id, return it's properties
func GetShoppingListItem(db *sql.DB, listid, itemID string) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `select * from shopping_item where listid = $1 and id = $2`
	rows, err := db.Query(sqlStatement, listid, itemID)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	rows.Next()
	return GetItemObjectFromRows(rows)
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

	sqlStatement := `insert into shopping_list (name, notes, author, authorLast, completed)
                         values ($1, $2, $3, $4, $5)
                         returning *`
	rows, err := db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.Author, shoppingList.AuthorLast, shoppingList.Completed)
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
	rows.Scan(&list.ID, &list.Name, &list.Notes, &list.Author, &list.AuthorLast, &list.Completed, &list.CreationTimestamp, &list.ModificationTimestamp, &list.DeletionTimestamp)
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

// AddItemToList ...
// adds a new item
func AddItemToList(db *sql.DB, listID string, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	valid, err := ValidateShoppingListItem(db, item)
	if !valid || err != nil {
		return itemInserted, err
	}

	item.AuthorLast = item.Author

	sqlStatement := `insert into shopping_item (listId, name, price, quantity, notes, author, authorLast, tag)
                         values ($1, $2, $3, $4, $5, $6, $7, $8)
                         returning *`
	rows, err := db.Query(sqlStatement, listID, item.Name, item.Price, item.Quantity, item.Notes, item.Author, item.AuthorLast, item.Tag)
	if err != nil {
		return itemInserted, err
	}
	defer rows.Close()
	for rows.Next() {
		itemInserted, err = GetItemObjectFromRows(rows)
		if err != nil {
			return itemInserted, err
		}
	}
	return itemInserted, err
}

// PatchItem ...
// patches a shopping item
func PatchItem(db *sql.DB, listid string, itemID string, item types.ShoppingItemSpec) (itemPatched types.ShoppingItemSpec, err error) {
	existingItem, err := GetShoppingListItem(db, listid, itemID)
	if err != nil || existingItem.ID == "" {
		return itemPatched, fmt.Errorf("Failed to fetch existing shopping list")
	}
	err = mergo.Merge(&item, existingItem)
	if err != nil {
		return itemPatched, fmt.Errorf("Failed to update fields in the item")
	}

	valid, err := ValidateShoppingListItem(db, existingItem)
	if !valid || err != nil {
		return itemPatched, err
	}

	sqlStatement := `update shopping_item set name = $2, price = $3, quantity = $4, notes = $5, authorLast = $6, tag = $7, obtained = $8, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := db.Query(sqlStatement, itemID, item.Name, item.Price, item.Quantity, item.Notes, item.AuthorLast, item.Tag, item.Obtained)
	if err != nil {
		return itemPatched, err
	}
	defer rows.Close()
	for rows.Next() {
		itemPatched, err = GetItemObjectFromRows(rows)
		if err != nil {
			return itemPatched, err
		}
	}

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: item.AuthorLast,
	}
	_, err = PatchShoppingList(db, itemPatched.ListID, shoppingListPatch)
	return itemPatched, err
}

// UpdateItem ...
// patches a shopping item
func UpdateItem(db *sql.DB, listID string, itemID string, item types.ShoppingItemSpec) (itemUpdated types.ShoppingItemSpec, err error) {
	valid, err := ValidateShoppingListItem(db, item)
	if !valid || err != nil {
		return itemUpdated, err
	}

	sqlStatement := `update shopping_item set name = $3, price = $4, quantity = $5, notes = $6, authorLast = $7, tag = $8, obtained = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where listId = $1 and id = $2 returning *`
	rows, err := db.Query(sqlStatement, listID, itemID, item.Name, item.Price, item.Quantity, item.Notes, item.AuthorLast, item.Tag, item.Obtained)
	if err != nil {
		return itemUpdated, err
	}
	defer rows.Close()
	for rows.Next() {
		itemUpdated, err = GetItemObjectFromRows(rows)
		if err != nil {
			return itemUpdated, err
		}
	}

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: item.AuthorLast,
	}
	_, err = PatchShoppingList(db, itemUpdated.ListID, shoppingListPatch)
	return itemUpdated, err
}

// SetItemObtained ...
// updates the item's obtained field
func SetItemObtained(db *sql.DB, listID string, itemID string, obtained bool, authorLast string) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `update shopping_item set obtained = $3 where listId = $1 and id = $2 returning *`
	rows, err := db.Query(sqlStatement, listID, itemID, obtained)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err = GetItemObjectFromRows(rows)
		if err != nil {
			return item, err
		}
	}

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: authorLast,
	}
	_, err = PatchShoppingList(db, item.ListID, shoppingListPatch)
	return item, err
}

// GetItemObjectFromRows ...
// returns an item object from rows
func GetItemObjectFromRows(rows *sql.Rows) (item types.ShoppingItemSpec, err error) {
	rows.Scan(&item.ID, &item.ListID, &item.Name, &item.Price, &item.Quantity, &item.Notes, &item.Obtained, &item.Tag, &item.Author, &item.AuthorLast, &item.CreationTimestamp, &item.ModificationTimestamp, &item.DeletionTimestamp)
	err = rows.Err()
	return item, err
}

// RemoveItemFromList ...
// given an item id, remove it
func RemoveItemFromList(db *sql.DB, itemID string, listID string) (err error) {
	sqlStatement := `delete from shopping_item where id = $1 and listId = $2`
	rows, err := db.Query(sqlStatement, itemID, listID)
	defer rows.Close()
	return err
}

// RemoveAllItemsFromList ...
// given an item id, remove all items
func RemoveAllItemsFromList(db *sql.DB, listID string) (err error) {
	sqlStatement := `delete from shopping_item where listId = $1`
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

// GetListItemCount ...
// returns a count of the items in a list
func GetListItemCount(db *sql.DB, listID string) (count int, err error) {
	sqlStatement := `select count(*) from shopping_item where listId = $1`
	rows, err := db.Query(sqlStatement, listID)
	if err != nil {
		return count, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)

	return count, err
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

// GetAllShoppingListTags ...
// returns a list of all tags used in items across lists
func GetAllShoppingListTags(db *sql.DB) (tags []string, err error) {
	sqlStatement := `select distinct tag from shopping_item order by tag`
	rows, err := db.Query(sqlStatement)
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
