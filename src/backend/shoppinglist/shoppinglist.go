/*
  shoppinglist
    manage shopping lists
*/

package shoppinglist

import (
	"database/sql"
	"errors"

	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

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
	sqlStatement := `select * from shopping_list where id = $1`
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
	sqlStatement := `select * from shopping_item`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		found := true
		item, err := ShoppingItemObjectFromRows(rows)
		if err != nil {
			return items, err
		}
		if itemSelector.Regular == true {
			if item.Regular == false {
				found = false
			}
		}
		if found == false {
			continue
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
func CreateShoppingList(db *sql.DB, shoppingList types.ShoppingListSpec) (shoppingListInserted types.ShoppingListSpec, err error) {
	if len(shoppingList.Name) == 0 || len(shoppingList.Name) > 30 || shoppingList.Name == "" {
		return shoppingListInserted, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if shoppingList.Notes != "" && len(shoppingList.Notes) > 100 {
		return shoppingListInserted, errors.New("Unable to save shopping list notes, as they are too long")
	}
	if len(shoppingList.Author) == 0 || shoppingList.Author == "" {
		return shoppingListInserted, errors.New("No shopping list author has been provided")
	}
	user, err := users.GetUserById(db, shoppingList.Author, false)
	if err != nil || user.Id == "" {
		return shoppingListInserted, errors.New("Unable to find author for shopping list")
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
	err = AddRegularItemsToList(db, shoppingListInserted.Id)

	return shoppingListInserted, err
}

// AddRegularItemsToList
// finds and adds the items marked as regular to a list via new entries
func AddRegularItemsToList(db *sql.DB, listId string) (err error) {
	itemSelector := types.ShoppingItemSelector{Regular: true}
	regularItems, err := GetShoppingListItems(db, listId, itemSelector)
	if err != nil {
		return err
	}
	for _, regularItem := range regularItems {
		// ensure there aren't duplicate regular items
		newRegularItem := types.ShoppingItemSpec{
			Name: regularItem.Name,
			Price: regularItem.Price,
			Regular: false,
			Notes: regularItem.Notes,
		}
		itemInserted, err := AddItemToList(db, listId, newRegularItem)
		if err != nil || itemInserted.Id == "" {
			return err
		}
	}
	return err
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
	sqlStatement := `update shopping_list set deletionTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1`
	_, err = db.Query(sqlStatement, listId)
	return err
}

// AddItem
// adds a new item
func AddItem(db *sql.DB, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	if len(item.Name) == 0 || len(item.Name) > 30 || item.Name == "" {
		return itemInserted, errors.New("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if item.Notes != "" && len(item.Notes) > 40 {
		return itemInserted, errors.New("Unable to save shopping list notes, as they are too long")
	}
	user, err := users.GetUserById(db, item.Author, false)
	if err != nil || user.Id == "" {
		return itemInserted, errors.New("Unable to find author for shopping list")
	}

	item.AuthorLast = item.Author

	sqlStatement := `insert into shopping_item (name, price, regular, notes, author, authorLast)
                         values ($1, $2, $3, $4, $5, $6)
                         returning *`
	rows, err := db.Query(sqlStatement, item.Name, item.Price, item.Regular, item.Notes, item.Author, item.AuthorLast)
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

// ShoppingItemObjectFromRows
// returns an item object from rows
func ShoppingItemObjectFromRows(rows *sql.Rows) (item types.ShoppingItemSpec, err error) {
	defer rows.Close()
	rows.Scan(&item.Id, &item.Name, &item.Price, &item.Regular, &item.Notes, &item.Obtained, &item.Author, &item.AuthorLast, &item.CreationTimestamp, &item.ModificationTimestamp, &item.DeletionTimestamp)
	err = rows.Err()
	return item, err
}

// RemoveItem
// given an item id, remove it
func RemoveItem(db *sql.DB, itemId string) (err error) {
	sqlStatement := `delete from shopping_item where id = $1`
	_, err = db.Query(sqlStatement, itemId)
	return err
}

// AssignItemToLinkList
// links an item to a list
func AssignItemToLinkList(db *sql.DB, itemId string, listId string) (err error) {
	sqlStatement := `insert into shopping_item_to_list (itemId, listId) values ($1, $2)`
	_, err = db.Query(sqlStatement, itemId, listId)
	return err
}

// RemoveItemFromLinkList
// unlinks an item from a list
func RemoveItemFromLinkList(db *sql.DB, itemId string, listId string) (err error) {
	sqlStatement := `delete from shopping_item_to_list where itemId = $1 and listId = $2`
	_, err = db.Query(sqlStatement, itemId, listId)
	return err
}

// AddItemToList
// adds an item and then links it to a list
func AddItemToList(db *sql.DB, listId string, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	itemInserted, err = AddItem(db, item)
	if err != nil {
		return itemInserted, err
	}
	err = AssignItemToLinkList(db, itemInserted.Id, listId)
	return itemInserted, err
}

// RemoveItemFromList
// unlinks an item then removes it
func RemoveItemFromList(db *sql.DB, itemId string, listId string) (err error) {
	err = RemoveItemFromLinkList(db, itemId, listId)
	if err != nil {
		return err
	}
	return RemoveItem(db, itemId)
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
	sqlStatement := `select count(*) from shopping_item_to_list where listId = $1`
	rows, err := db.Query(sqlStatement, listId)
	if err != nil {
		return count, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)

	return count, err
}