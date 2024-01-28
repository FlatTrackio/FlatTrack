/*
  shoppinglist
    item
      manage shopping list items
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

	"github.com/imdario/mergo"

	"gitlab.com/flattrack/flattrack/pkg/types"
)

type ShoppingItemManager struct {
	manager *Manager
	db      *sql.DB
}

func (m *Manager) ShoppingItem() *ShoppingItemManager {
	return &ShoppingItemManager{
		manager: m,
		db:      m.db,
	}
}

// Validate ...
// given a shopping list item, return it's validity
func (m *ShoppingItemManager) Validate(item types.ShoppingItemSpec) (valid bool, err error) {
	if len(item.Name) == 0 || len(item.Name) >= 30 || item.Name == "" {
		return false, ErrInvalidShoppingItemName
	}
	if item.Notes != "" && len(item.Notes) >= 40 {
		return false, ErrInvalidShoppingItemNotes
	}
	if item.Tag != "" && len(item.Tag) == 0 || len(item.Tag) >= 30 {
		return false, ErrInvalidShoppingItemTag
	}
	if item.Quantity < 1 {
		return false, ErrInvalidItemQuantityMustBeOne
	}
	if item.TemplateID != "" {
		list, err := m.manager.ShoppingList().Get(item.TemplateID)
		if err != nil || list.ID == "" {
			return false, ErrShoppingListByIDNotFoundForTemplate
		}
	}
	return true, nil
}

// List ...
// returns a list of items on a shopping list
func (m *ShoppingItemManager) List(listID string, options types.ShoppingItemOptions) (any, error) {
	// sort by tags
	var obtained sql.NullString
	if err := obtained.Scan(options.Selector.Obtained); err != nil {
		log.Printf("error: scanning obtained; %v\n", err)
		return []types.ShoppingItemSpec{}, err
	}

	sqlQueryValues := []interface{}{listID}
	sqlStatement := `select * from shopping_item where listId = $1`
	if options.Selector.Obtained != "" && options.Selector.TemplateListItemSelector != "" {
		sqlStatement += fmt.Sprintf(` and obtained = $%v`, len(sqlQueryValues)+1)
		sqlQueryValues = append(sqlQueryValues, obtained)
	}
	switch options.Selector.TemplateListItemSelector {
	case "unobtained":
		sqlStatement += ` and obtained = false`
	case "obtained":
		sqlStatement += ` and obtained = true`
	default:
	}
	if options.SearchName != "" {
		sqlStatement += ` and name ilike '%' ||` + fmt.Sprintf("$%v", len(sqlQueryValues)+1) + `|| '%'`
		sqlQueryValues = append(sqlQueryValues, options.SearchName)
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

	rows, err := m.db.Query(sqlStatement, sqlQueryValues...)
	if err != nil {
		log.Println(err)
		return []types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	items := []types.ShoppingItemSpec{}
	for rows.Next() {
		item, err := getItemObjectFromRows(rows)
		if err != nil {
			return []types.ShoppingItemSpec{}, err
		}
		items = append(items, item)
	}
	if options.SortBy == types.ShoppingItemSortByTag {
		return m.groupItemsByTag(items), nil
	}
	return items, nil
}

func (m *ShoppingItemManager) groupItemsByTag(items []types.ShoppingItemSpec) (lists []types.ShoppingItemsGroupByTag) {
	// Group items by tag and calculate tag price
	for _, item := range items {
		foundTagInList := -1
	lists:
		for i, list := range lists {
			if list.Tag == item.Tag {
				foundTagInList = i
				break lists
			}
		}
		if foundTagInList == -1 {
			lists = []types.ShoppingItemsGroupByTag{{
				Tag:   item.Tag,
				Items: []types.ShoppingItemSpec{item},
				Price: item.Price,
			}}
		} else {
			lists[foundTagInList].Items = append(lists[foundTagInList].Items, item)
			lists[foundTagInList].Price += item.Price * float64(item.Quantity)
		}
	}
	return lists
}

// Get ...
// given an item id, return it's properties
func (m *ShoppingItemManager) Get(listid, itemID string) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `select * from shopping_item where listid = $1 and id = $2`
	rows, err := m.db.Query(sqlStatement, listid, itemID)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	item, err = getItemObjectFromRows(rows)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	return item, nil
}

// AddItemToList ...
// adds a new item
func (m *ShoppingItemManager) AddItemToList(listID string, item types.ShoppingItemSpec) (itemInserted types.ShoppingItemSpec, err error) {
	valid, err := m.Validate(item)
	if !valid || err != nil {
		return types.ShoppingItemSpec{}, err
	}

	if item.Tag == "" {
		item.Tag = "Untagged"
	}

	item.AuthorLast = item.Author

	sqlStatement := `insert into shopping_item (listId, name, price, quantity, notes, author, authorLast, tag, obtained, templateId)
                         values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                         returning *`
	rows, err := m.db.Query(sqlStatement, listID, item.Name, item.Price, item.Quantity, item.Notes, item.Author, item.AuthorLast, item.Tag, item.Obtained, &item.TemplateID)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemInserted, err = getItemObjectFromRows(rows)
		if err != nil {
			return types.ShoppingItemSpec{}, err
		}
	}
	return itemInserted, nil
}

// Patch ...
// patches a shopping item
func (m *ShoppingItemManager) Patch(listid string, itemID string, item types.ShoppingItemSpec) (itemPatched types.ShoppingItemSpec, err error) {
	existingItem, err := m.Get(listid, itemID)
	if err != nil || existingItem.ID == "" {
		return types.ShoppingItemSpec{}, ErrFailedToGetExistingShoppingList
	}
	err = mergo.Merge(&item, existingItem)
	if err != nil {
		return types.ShoppingItemSpec{}, ErrFailedToUpdateShoppingItemFields
	}

	valid, err := m.Validate(existingItem)
	if !valid || err != nil {
		return types.ShoppingItemSpec{}, err
	}

	if item.Tag == "" {
		item.Tag = "Untagged"
	}

	sqlStatement := `update shopping_item set name = $2, price = $3, quantity = $4, notes = $5, authorLast = $6, tag = $7, obtained = $8, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := m.db.Query(sqlStatement, itemID, item.Name, item.Price, item.Quantity, item.Notes, item.AuthorLast, item.Tag, item.Obtained)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemPatched, err = getItemObjectFromRows(rows)
		if err != nil {
			return types.ShoppingItemSpec{}, err
		}
	}

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: item.AuthorLast,
	}
	_, err = m.manager.ShoppingList().Patch(itemPatched.ListID, shoppingListPatch)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	return itemPatched, nil
}

// Update ...
// patches a shopping item
func (m *ShoppingItemManager) Update(listID string, itemID string, item types.ShoppingItemSpec) (itemUpdated types.ShoppingItemSpec, err error) {
	valid, err := m.Validate(item)
	if !valid || err != nil {
		return types.ShoppingItemSpec{}, err
	}

	if item.Tag == "" {
		item.Tag = "Untagged"
	}

	sqlStatement := `update shopping_item set name = $3, price = $4, quantity = $5, notes = $6, authorLast = $7, tag = $8, obtained = $9, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where listId = $1 and id = $2 returning *`
	rows, err := m.db.Query(sqlStatement, listID, itemID, item.Name, item.Price, item.Quantity, item.Notes, item.AuthorLast, item.Tag, item.Obtained)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemUpdated, err = getItemObjectFromRows(rows)
		if err != nil {
			return types.ShoppingItemSpec{}, err
		}
	}
	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: item.AuthorLast,
	}
	_, err = m.manager.ShoppingList().Patch(itemUpdated.ListID, shoppingListPatch)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	return itemUpdated, nil
}

// SetItemObtained ...
// updates the item's obtained field
func (m *ShoppingItemManager) SetItemObtained(listID string, itemID string, obtained bool, authorLast string) (item types.ShoppingItemSpec, err error) {
	sqlStatement := `update shopping_item set obtained = $3 where listId = $1 and id = $2 returning *`
	rows, err := m.db.Query(sqlStatement, listID, itemID, obtained)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		item, err = getItemObjectFromRows(rows)
		if err != nil {
			return types.ShoppingItemSpec{}, err
		}
	}

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: authorLast,
	}
	_, err = m.manager.ShoppingList().Patch(item.ListID, shoppingListPatch)
	if err != nil {
		return types.ShoppingItemSpec{}, err
	}
	return item, nil
}

// getItemObjectFromRows ...
// returns an item object from rows
func getItemObjectFromRows(rows *sql.Rows) (item types.ShoppingItemSpec, err error) {
	if err := rows.Scan(&item.ID, &item.ListID, &item.Name, &item.Price, &item.Quantity, &item.Notes, &item.Obtained, &item.Tag, &item.Author, &item.AuthorLast, &item.CreationTimestamp, &item.ModificationTimestamp, &item.DeletionTimestamp, &item.TemplateID); err != nil {
		return types.ShoppingItemSpec{}, err
	}
	if err := rows.Err(); err != nil {
		return types.ShoppingItemSpec{}, err
	}
	return item, nil
}

// Delete ...
// given an item id, remove it
func (m *ShoppingItemManager) Delete(id string, listID string, authorLast string) (err error) {
	sqlStatement := `delete from shopping_item where id = $1 and listId = $2`
	rows, err := m.db.Query(sqlStatement, id, listID)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: authorLast,
	}
	_, err = m.manager.ShoppingList().Patch(listID, shoppingListPatch)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
// given an item id, remove it
func (m *ShoppingItemManager) DeleteTagItems(listID string, tagName string, authorLast string) (err error) {
	sqlStatement := `delete from shopping_item where listId = $1 and tag = $2`
	rows, err := m.db.Query(sqlStatement, listID, tagName)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()

	shoppingListPatch := types.ShoppingListSpec{
		AuthorLast: authorLast,
	}
	_, err = m.manager.ShoppingList().Patch(listID, shoppingListPatch)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll ...
// given an item id, remove all items
// only intended to be called when deleting list
func (m *ShoppingItemManager) DeleteAll(listID string) (err error) {
	sqlStatement := `delete from shopping_item where listId = $1`
	rows, err := m.db.Query(sqlStatement, listID)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return err
}

// GetListItemCount ...
// returns a count of the items in a list
func (m *ShoppingItemManager) GetListItemCount(listID string) (count int, err error) {
	sqlStatement := `select count(*) from shopping_item where listId = $1`
	rows, err := m.db.Query(sqlStatement, listID)
	if err != nil {
		return count, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	if err := rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
