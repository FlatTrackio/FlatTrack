/*
  shoppinglist
    manage shopping lists
*/

package shoppinglist

import (
	"database/sql"
	"errors"

	"github.com/imdario/mergo"

	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

func ValidateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec) (valid bool, err error) {
	if len(shoppingList.Name) == 0 || len(shoppingList.Name) > 30 || shoppingList.Name == "" {
		return valid, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if shoppingList.Notes != "" && len(shoppingList.Notes) > 100 {
		return valid, errors.New("Unable to save shopping list notes, as they are too long")
	}
	if len(shoppingList.Author) == 0 || shoppingList.Author == "" {
		return valid, errors.New("No shopping list author has been provided")
	}
	user, err := users.GetUserById(db, shoppingList.Author, false)
	if err != nil || user.Id == "" {
		return valid, errors.New("Unable to find author for shopping list")
	}
	if shoppingList.TemplateId != "" {
		list, err := GetShoppingList(db, shoppingList.TemplateId)
		if err != nil || list.Id == "" {
			return valid, errors.New("Unable to find list to use as template from provided id")
		}
	}
	return true, err
}

func ValidateShoppingListItem(db *sql.DB, item types.ShoppingItemSpec) (valid bool, err error) {
	if len(item.Name) == 0 || len(item.Name) > 30 || item.Name == "" {
		return valid, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if item.Notes != "" && len(item.Notes) > 40 {
		return valid, errors.New("Unable to save shopping list notes, as they are too long")
	}
	if item.Quantity == 0 {
		return valid, errors.New("Item quanity must be greater than zero")
	}
	user, err := users.GetUserById(db, item.Author, false)
	if err != nil || user.Id == "" {
		return valid, errors.New("Unable to find author for shopping list")
	}
	return true, err
}

// GetShoppingLists
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
		shoppingList, err := ShoppingListObjectFromRows(rows)
		if err != nil {
			return shoppingLists, err
		}
		shoppingList.Count, err = GetListItemCount(db, shoppingList.Id)
		if err != nil {
			return shoppingLists, err
		}
		shoppingLists = append(shoppingLists, shoppingList)
	}
	return shoppingLists, err
}

// GetShoppingList
// returns a given shopping list, by it's Id
func GetShoppingList(db *sql.DB, listId string) (shoppingList types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list where id = $1 and deletionTimestamp = 0`
	rows, err := db.Query(sqlStatement, listId)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	rows.Next()
	shoppingList, err = ShoppingListObjectFromRows(rows)
	if err != nil {
		return shoppingList, err
	}
	shoppingList.Count, err = GetListItemCount(db, shoppingList.Id)
	return shoppingList, err
}

// GetShoppingListItems
// returns a list of items on a shopping list
func GetShoppingListItems(db *sql.DB, listId string, itemSelector types.ShoppingItemSelector) (items []types.ShoppingItemSpec, err error) {
	sqlStatement := `select * from shopping_item where listId = $1 order by tag, name`
	rows, err := db.Query(sqlStatement, listId)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := ShoppingItemObjectFromRows(rows)
		if err != nil {
			return items, err
		}
		if itemSelector.NotObtained == true {
			if item.Obtained == true {
				continue
			}
		}
		items = append(items, item)
	}
	return items, err
}

// GetShoppingListItem
// given an item id, return it's properties
func GetShoppingListItem(db *sql.DB, itemId string) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `select * from shopping_item where id = $1`
	rows, err := db.Query(sqlStatement, itemId)
	if err != nil {
		return item, err
	}
	defer rows.Close()

	rows.Next()
	return ShoppingItemObjectFromRows(rows)
}

// CreateShoppingList
// creates a shopping list for adding items to
func CreateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec, itemSelector types.ShoppingItemSelector) (shoppingListInserted types.ShoppingListSpec, err error) {
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
	rows.Next()
	shoppingListInserted, err = ShoppingListObjectFromRows(rows)
	if err != nil || shoppingListInserted.Id == "" {
		return shoppingListInserted, errors.New("Failed to create shopping list")
	}
	if shoppingList.TemplateId == "" {
		return shoppingListInserted, err
	}

	// if using other list as a template
	shoppingListItems, err := GetShoppingListItems(db, shoppingList.TemplateId, itemSelector)
	if err != nil {
		return shoppingListInserted, errors.New("Failed to fetch items from shopping list")
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
		newItem, err := AddItemToList(db, shoppingListInserted.Id, newItem)
		if err != nil || newItem.Id == "" {
			return shoppingListInserted, errors.New("Failed to add new item to new shopping list from template")
		}
	}

	return shoppingListInserted, err
}

// PatchShoppingList
// patches a shopping list
func PatchShoppingList(db *sql.DB, listId string, shoppingList types.ShoppingListSpec) (shoppingListPatched types.ShoppingListSpec, err error) {
	existingList, err := GetShoppingList(db, listId)
	if err != nil || existingList.Id == "" {
		return shoppingListPatched, errors.New("Failed to fetch existing shopping list")
	}
	err = mergo.Merge(&shoppingList, existingList)
	if err != nil {
		return shoppingListPatched, errors.New("Failed to update fields in the item")
	}
	valid, err := ValidateShoppingList(db, existingList)
	if !valid || err != nil {
		return shoppingListPatched, err
	}

	sqlStatement := `update shopping_list set name = $1, notes = $2, authorLast = $3, completed = $4, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $5
                         returning *`
	rows, err := db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.AuthorLast, shoppingList.Completed, listId)
	if err != nil {
		return shoppingListPatched, err
	}
	rows.Next()
	shoppingListPatched, err = ShoppingListObjectFromRows(rows)
	if err != nil || shoppingListPatched.Id == "" {
		return shoppingListPatched, errors.New("Failed to create shopping list")
	}
	return shoppingListPatched, err
}

// SetListCompleted
// updates the list's completed field
func SetListCompleted(db *sql.DB, listId string, completed bool) (list types.ShoppingListSpec, err error) {
	sqlStatement := `update shopping_list set completed = $1 where id = $2 returning *`
	rows, err := db.Query(sqlStatement, completed, listId)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		list, err = ShoppingListObjectFromRows(rows)
		if err != nil {
			return list, err
		}
	}
	return list, err
}

// ShoppingListObjectFromRows
// returns a shopping list object from rows
func ShoppingListObjectFromRows(rows *sql.Rows) (list types.ShoppingListSpec, err error) {
	rows.Scan(&list.Id, &list.Name, &list.Notes, &list.Author, &list.AuthorLast, &list.Completed, &list.CreationTimestamp, &list.ModificationTimestamp, &list.DeletionTimestamp)
	err = rows.Err()
	return list, err
}

// DeleteShoppingList
// deletes a shopping list, given a shopping list Id
func DeleteShoppingList(db *sql.DB, listId string) (err error) {
	err = RemoveAllItemsFromList(db, listId)
	if err != nil {
		return errors.New("Failed to remove all items from list")
	}

	sqlStatement := `delete from shopping_list where id = $1`
	rows, err := db.Query(sqlStatement, listId)
	defer rows.Close()
	return err
}

// AddItemToList
// adds a new item
func AddItemToList(db *sql.DB, listId string, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	valid, err := ValidateShoppingListItem(db, item)
	if !valid || err != nil {
		return itemInserted, err
	}

	item.AuthorLast = item.Author

	sqlStatement := `insert into shopping_item (listId, name, price, quantity, notes, author, authorLast, tag)
                         values ($1, $2, $3, $4, $5, $6, $7, $8)
                         returning *`
	rows, err := db.Query(sqlStatement, listId, item.Name, item.Price, item.Quantity, item.Notes, item.Author, item.AuthorLast, item.Tag)
	if err != nil {
		return itemInserted, err
	}
	defer rows.Close()
	for rows.Next() {
		itemInserted, err = ShoppingItemObjectFromRows(rows)
		if err != nil {
			return itemInserted, err
		}
	}
	return itemInserted, err
}

// PatchItem
// patches a shopping item
func PatchItem(db *sql.DB, itemId string, item types.ShoppingItemSpec) (itemPatched types.ShoppingItemSpec, err error) {
	existingItem, err := GetShoppingListItem(db, itemId)
	if err != nil || existingItem.Id == "" {
		return itemPatched, errors.New("Failed to fetch existing shopping list")
	}
	err = mergo.Merge(&item, existingItem)
	if err != nil {
		return itemPatched, errors.New("Failed to update fields in the item")
	}

	valid, err := ValidateShoppingListItem(db, existingItem)
	if !valid || err != nil {
		return itemPatched, err
	}

	sqlStatement := `update shopping_item set name = $2, price = $3, quantity = $4, notes = $5, authorLast = $6, tag = $7, obtained = $8, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := db.Query(sqlStatement, itemId, item.Name, item.Price, item.Quantity, item.Notes, item.AuthorLast, item.Tag, item.Obtained)
	if err != nil {
		return itemPatched, err
	}
	defer rows.Close()
	for rows.Next() {
		itemPatched, err = ShoppingItemObjectFromRows(rows)
		if err != nil {
			return itemPatched, err
		}
	}
	return itemPatched, err
}

// SetItemObtained
// updates the item's obtained field
func SetItemObtained(db *sql.DB, itemId string, obtained bool) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `update shopping_item set obtained = $1 where id = $2 returning *`
	rows, err := db.Query(sqlStatement, obtained, itemId)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err = ShoppingItemObjectFromRows(rows)
		if err != nil {
			return item, err
		}
	}
	return item, err
}

// ShoppingItemObjectFromRows
// returns an item object from rows
func ShoppingItemObjectFromRows(rows *sql.Rows) (item types.ShoppingItemSpec, err error) {
	rows.Scan(&item.Id, &item.ListId, &item.Name, &item.Price, &item.Quantity, &item.Notes, &item.Obtained, &item.Tag, &item.Author, &item.AuthorLast, &item.CreationTimestamp, &item.ModificationTimestamp, &item.DeletionTimestamp)
	err = rows.Err()
	return item, err
}

// RemoveItemFromList
// given an item id, remove it
func RemoveItemFromList(db *sql.DB, itemId string, listId string) (err error) {
	sqlStatement := `delete from shopping_item where id = $1 and listId = $2`
	rows, err := db.Query(sqlStatement, itemId, listId)
	defer rows.Close()
	return err
}

// RemoveAllItemsFromList
// given an item id, remove all items
func RemoveAllItemsFromList(db *sql.DB, listId string) (err error) {
	sqlStatement := `delete from shopping_item where listId = $1`
	rows, err := db.Query(sqlStatement, listId)
	defer rows.Close()
	return err
}

// GetListCount
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

// GetListItemCount
// returns a count of the items in a list
func GetListItemCount(db *sql.DB, listId string) (count int, err error) {
	sqlStatement := `select count(*) from shopping_item where listId = $1`
	rows, err := db.Query(sqlStatement, listId)
	if err != nil {
		return count, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)

	return count, err
}

// GetShoppingListTags
// returns a list of tags used in items across lists
func GetShoppingListTags(db *sql.DB) (tags []string, err error) {
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
