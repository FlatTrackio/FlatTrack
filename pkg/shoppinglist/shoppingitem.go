/*
  shoppinglist
    item
      manage shopping list items
*/

package shoppinglist

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/imdario/mergo"

	"gitlab.com/flattrack/flattrack/pkg/types"
)

// ValidateShoppingListItem ...
// given a shopping list item, return it's validity
func ValidateShoppingListItem(db *sql.DB, item types.ShoppingItemSpec) (valid bool, err error) {
	if len(item.Name) == 0 || len(item.Name) >= 30 || item.Name == "" {
		return valid, fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if item.Notes != "" && len(item.Notes) >= 40 {
		return valid, fmt.Errorf("Unable to save shopping list notes, as they are too long")
	}
	if item.Tag != "" && len(item.Tag) == 0 || len(item.Tag) >= 30 {
		return valid, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	if item.Quantity < 1 {
		return valid, fmt.Errorf("Item quantity must be at least one")
	}
	if item.TemplateID != "" {
		list, err := GetShoppingList(db, item.TemplateID)
		if err != nil || list.ID == "" {
			return valid, fmt.Errorf("Unable to find list to use as template from provided id")
		}
	}
	return true, err
}

// GetShoppingListItems ...
// returns a list of items on a shopping list
func GetShoppingListItems(db *sql.DB, listID string, options types.ShoppingItemOptions) (items []types.ShoppingItemSpec, err error) {
	// sort by tags
	var obtained sql.NullBool
	obtained.Scan(options.Selector.Obtained)

	sqlQueryValues := []interface{}{
		listID,
	}
	sqlStatement := `select * from shopping_item where listId = $1`
	if options.Selector.Obtained != "" {
		sqlStatement += ` and obtained = $2`
		sqlQueryValues = append(sqlQueryValues, obtained)
	}
	if options.SortBy == types.ShoppingItemSortByHighestPrice {
		sqlStatement += ` order by price desc, name asc`
	} else if options.SortBy == types.ShoppingItemSortByHighestQuantity {
		sqlStatement += ` order by quantity desc, name desc`
	} else if options.SortBy == types.ShoppingItemSortByLowestPrice {
		sqlStatement += ` order by price asc, name asc`
	} else if options.SortBy == types.ShoppingItemSortByLowestQuantity {
		sqlStatement += ` order by quantity asc, name desc`
	} else if options.SortBy == types.ShoppingItemSortByRecentlyAdded {
		sqlStatement += ` order by creationTimestamp desc`
	} else if options.SortBy == types.ShoppingItemSortByRecentlyUpdated {
		sqlStatement += ` order by modificationTimestamp desc`
	} else if options.SortBy == types.ShoppingItemSortByLastAdded {
		sqlStatement += ` order by creationTimestamp asc`
	} else if options.SortBy == types.ShoppingItemSortByLastUpdated {
		sqlStatement += ` order by modificationTimestamp asc`
	} else if options.SortBy == types.ShoppingItemSortByAlphabeticalDescending {
		sqlStatement += ` order by name asc`
	} else if options.SortBy == types.ShoppingItemSortByAlphabeticalAscending {
		sqlStatement += ` order by name desc`
	} else {
		sqlStatement += ` order by tag asc, name asc`
	}
	rows, err := db.Query(sqlStatement, sqlQueryValues...)
	if err != nil {
		log.Println(err)
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err := GetItemObjectFromRows(rows)
		if err != nil {
			return items, err
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

// AddItemToList ...
// adds a new item
func AddItemToList(db *sql.DB, listID string, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	valid, err := ValidateShoppingListItem(db, item)
	if !valid || err != nil {
		return itemInserted, err
	}

	item.AuthorLast = item.Author

	sqlStatement := `insert into shopping_item (listId, name, price, quantity, notes, author, authorLast, tag, obtained, templateId)
                         values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                         returning *`
	rows, err := db.Query(sqlStatement, listID, item.Name, item.Price, item.Quantity, item.Notes, item.Author, item.AuthorLast, item.Tag, item.Obtained, &item.TemplateID)
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
	rows.Scan(&item.ID, &item.ListID, &item.Name, &item.Price, &item.Quantity, &item.Notes, &item.Obtained, &item.Tag, &item.Author, &item.AuthorLast, &item.CreationTimestamp, &item.ModificationTimestamp, &item.DeletionTimestamp, &item.TemplateID)
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
