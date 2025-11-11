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
	"log"
	"time"

	"github.com/imdario/mergo"
	"github.com/lib/pq"

	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	ErrFailedToAddItemToShoppingListFromTemplate = fmt.Errorf("Failed to add new item to new shopping list from template")
	ErrFailedToCreateShoppingList                = fmt.Errorf("Failed to create shopping list")
	ErrFailedToFindShoppingTagToUpdate           = fmt.Errorf("Unable to find tag to update")
	ErrFailedToGetExistingShoppingList           = fmt.Errorf("Failed to get existing shopping list")
	ErrFailedToGetItemsFromShoppingList          = fmt.Errorf("Failed to get items from shopping list")
	ErrFailedToPatchShoppingList                 = fmt.Errorf("Failed to patch shopping list")
	ErrFailedToRemoveAllItemsFromList            = fmt.Errorf("Failed to remove all items from list")
	ErrFailedToUpdateShoppingItemFields          = fmt.Errorf("Failed to update fields in the item")
	ErrInvalidItemQuantityMustBeOne              = fmt.Errorf("Unable to use item quantity must be at least one")
	ErrInvalidShoppingItemName                   = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	ErrInvalidShoppingItemTag                    = fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	ErrInvalidShoppingListNotes                  = fmt.Errorf("Unable to save shopping list notes, as they are too long")
	ErrInvalidShoppingItemNotes                  = fmt.Errorf("Unable to save shopping item notes, as they are too long")
	ErrShoppingListByIDNotFoundForTemplate       = fmt.Errorf("Unable to find list to use as template from provided id")
)

type Manager struct {
	db              *sql.DB
	settingsManager *settings.Manager
}

func NewManager(db *sql.DB, settingsManager *settings.Manager) *Manager {
	return &Manager{
		db:              db,
		settingsManager: settingsManager,
	}
}

type ShoppingListManager struct {
	manager *Manager
	db      *sql.DB
}

func (m *Manager) ShoppingList() *ShoppingListManager {
	return &ShoppingListManager{
		manager: m,
		db:      m.db,
	}
}

// Validate ...
// given a shopping list, return it's validity
func (m *ShoppingListManager) Validate(shoppingList types.ShoppingListSpec) (valid bool, err error) {
	if len(shoppingList.Name) == 0 || len(shoppingList.Name) >= 30 || shoppingList.Name == "" {
		return false, ErrInvalidShoppingItemName
	}
	if shoppingList.Notes != "" && len(shoppingList.Notes) > 100 {
		return false, ErrInvalidShoppingListNotes
	}
	if shoppingList.TemplateID != "" {
		list, err := m.Get(shoppingList.TemplateID)
		if err != nil || list.ID == "" {
			return false, ErrShoppingListByIDNotFoundForTemplate
		}
	}
	return true, nil
}

// List ...
// returns a list of all shopping lists (name, notes, author, etc...)
func (m *ShoppingListManager) List(options types.ShoppingListOptions) (shoppingLists []types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list where deletionTimestamp = 0 `
	fields := []interface{}{}

	if options.SortBy == types.ShoppingListSortByTemplated {
		sqlStatement = `with popularity as (
                          select id, (select count(*) from shopping_list where templateid = c.id) as tally from shopping_list c)
                        select id, name, notes, author, authorlast, completed, creationtimestamp, modificationtimestamp, deletiontimestamp, templateid, total_tag_exclude
                        from shopping_list
                        join popularity using(id) where deletiontimestamp = 0 `
	}

	if options.Selector.ModificationTimestampBefore != 0 {
		sqlStatement += fmt.Sprintf(`and modificationTimestamp < $%v `, len(fields)+1)
		fields = append(fields, options.Selector.ModificationTimestampBefore)
	}
	if options.Selector.CreationTimestampBefore != 0 {
		sqlStatement += fmt.Sprintf(`and creationTimestamp < $%v `, len(fields)+1)
		fields = append(fields, options.Selector.CreationTimestampBefore)
	}
	if options.Selector.ModificationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and modificationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.ModificationTimestampAfter)
	}
	if options.Selector.CreationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and creationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.CreationTimestampAfter)
	}

	switch options.SortBy {
	case types.ShoppingListSortByRecentlyUpdated:
		sqlStatement += `order by modificationTimestamp desc `
	case types.ShoppingListSortByLastUpdated:
		sqlStatement += `order by modificationTimestamp asc `
	case types.ShoppingListSortByRecentlyAdded:
		sqlStatement += `order by creationTimestamp asc `
	case types.ShoppingListSortByLastAdded:
		sqlStatement += `order by creationTimestamp asc `
	case types.ShoppingListSortByAlphabeticalDescending:
		sqlStatement += `order by name asc `
	case types.ShoppingListSortByAlphabeticalAscending:
		sqlStatement += `order by name desc `
	case types.ShoppingListSortByTemplated:
		sqlStatement += `order by popularity.tally desc, shopping_list.creationtimestamp desc `
	default:
		sqlStatement += `order by creationTimestamp desc `
	}

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(`limit $%v `, len(fields)+1)
		fields = append(fields, options.Limit)
	}

	rows, err := m.db.Query(sqlStatement, fields...)
	if err != nil {
		return []types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		shoppingList, err := getListObjectFromRows(rows)
		if err != nil {
			return []types.ShoppingListSpec{}, err
		}
		shoppingList.Count, err = m.manager.ShoppingItem().GetListItemCount(shoppingList.ID)
		if err != nil {
			return []types.ShoppingListSpec{}, err
		}

		if options.Selector.Completed == "true" && !shoppingList.Completed {
			continue
		} else if options.Selector.Completed == "false" && shoppingList.Completed {
			continue
		}
		shoppingLists = append(shoppingLists, shoppingList)
	}
	return shoppingLists, nil
}

// Get ...
// returns a given shopping list, by it's ID
func (m *ShoppingListManager) Get(listID string) (shoppingList types.ShoppingListSpec, err error) {
	sqlStatement := `select * from shopping_list where id = $1 and deletionTimestamp = 0`
	rows, err := m.db.Query(sqlStatement, listID)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	shoppingList, err = getListObjectFromRows(rows)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	shoppingList.Count, err = m.manager.ShoppingItem().GetListItemCount(shoppingList.ID)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	return shoppingList, nil
}

// Create ...
// creates a shopping list for adding items to
func (m *ShoppingListManager) Create(shoppingList types.ShoppingListSpec, options types.ShoppingItemOptions) (shoppingListInserted types.ShoppingListSpec, err error) {
	valid, err := m.Validate(shoppingList)
	if !valid || err != nil {
		return types.ShoppingListSpec{}, err
	}

	if shoppingList.TemplateID != "" {
		templateList, err := m.Get(shoppingList.TemplateID)
		if err != nil {
			return types.ShoppingListSpec{}, err
		}
		shoppingList.TotalTagExclude = templateList.TotalTagExclude
	} else {
		shoppingList.TotalTagExclude = []string{}
	}

	shoppingList.AuthorLast = shoppingList.Author
	shoppingList.Completed = false

	sqlStatement := `insert into shopping_list (name, notes, author, authorLast, completed, templateId, total_tag_exclude)
                         values ($1, $2, $3, $4, $5, $6, $7)
                         returning *`
	rows, err := m.db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.Author, shoppingList.AuthorLast, shoppingList.Completed, shoppingList.TemplateID, pq.Array(shoppingList.TotalTagExclude))
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	shoppingListInserted, err = getListObjectFromRows(rows)
	if err != nil || shoppingListInserted.ID == "" {
		log.Printf("error getting list object from rows: %v\n", err)
		return types.ShoppingListSpec{}, ErrFailedToCreateShoppingList
	}
	log.Printf("%+v\n", shoppingListInserted)
	if shoppingList.TemplateID == "" {
		return shoppingListInserted, nil
	}

	// if using other list as a template
	shoppingListItems, err := m.manager.ShoppingItem().List(shoppingList.TemplateID, options)
	if err != nil {
		if err := m.Delete(shoppingListInserted.ID); err != nil {
			return types.ShoppingListSpec{}, err
		}
		return types.ShoppingListSpec{}, ErrFailedToGetItemsFromShoppingList
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
		_, err := m.manager.ShoppingItem().AddItemToList(shoppingListInserted.ID, newItem)
		if err != nil {
			log.Printf("error adding item to list: %v\n", err)
			if err := m.Delete(shoppingListInserted.ID); err != nil {
				return types.ShoppingListSpec{}, err
			}
			return types.ShoppingListSpec{}, ErrFailedToAddItemToShoppingListFromTemplate
		}
	}
	return shoppingListInserted, nil
}

// Patch ...
// patches a shopping list
func (m *ShoppingListManager) Patch(listID string, shoppingList types.ShoppingListSpec) (shoppingListPatched types.ShoppingListSpec, err error) {
	existingList, err := m.Get(listID)
	if err != nil || existingList.ID == "" {
		return types.ShoppingListSpec{}, ErrFailedToGetExistingShoppingList
	}
	err = mergo.Merge(&shoppingList, existingList)
	if err != nil {
		return types.ShoppingListSpec{}, ErrFailedToUpdateShoppingItemFields
	}
	valid, err := m.Validate(existingList)
	if !valid || err != nil {
		return types.ShoppingListSpec{}, err
	}

	sqlStatement := `update shopping_list set name = $1, notes = $2, authorLast = $3, completed = $4, total_tag_exclude = $5, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $6
                         returning *`
	rows, err := m.db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.AuthorLast, shoppingList.Completed, pq.Array(shoppingList.TotalTagExclude), listID)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	shoppingListPatched, err = getListObjectFromRows(rows)
	if err != nil || shoppingListPatched.ID == "" {
		log.Printf("error getting shopping list from rows: %v", err)
		return types.ShoppingListSpec{}, ErrFailedToPatchShoppingList
	}
	return shoppingListPatched, nil
}

// UpdateShoppingList ...
// updates a shopping list
func (m *ShoppingListManager) Update(listID string, shoppingList types.ShoppingListSpec) (shoppingListUpdated types.ShoppingListSpec, err error) {
	valid, err := m.Validate(shoppingList)
	if !valid || err != nil {
		return types.ShoppingListSpec{}, err
	}

	sqlStatement := `update shopping_list set name = $1, notes = $2, authorLast = $3, completed = $4, total_tag_exclude = $5::text[], modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $6
                         returning *`
	rows, err := m.db.Query(sqlStatement, shoppingList.Name, shoppingList.Notes, shoppingList.AuthorLast, shoppingList.Completed, pq.Array(shoppingList.TotalTagExclude), listID)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	shoppingListUpdated, err = getListObjectFromRows(rows)
	if err != nil || shoppingListUpdated.ID == "" {
		log.Printf("error getting shopping list from rows: %v", err)
		return types.ShoppingListSpec{}, ErrFailedToCreateShoppingList
	}
	return shoppingListUpdated, nil
}

// SetListCompleted ...
// updates the list's completed field
func (m *ShoppingListManager) SetListCompleted(listID string, completed bool, userID string) (list types.ShoppingListSpec, err error) {
	sqlStatement := `update shopping_list set completed = $1 where id = $2 returning *`
	rows, err := m.db.Query(sqlStatement, completed, listID)
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		list, err = getListObjectFromRows(rows)
		if err != nil {
			return types.ShoppingListSpec{}, err
		}
	}

	return list, nil
}

// getListObjectFromRows ...
// returns a shopping list object from rows
func getListObjectFromRows(rows *sql.Rows) (list types.ShoppingListSpec, err error) {
	if err := rows.Scan(&list.ID, &list.Name, &list.Notes, &list.Author, &list.AuthorLast, &list.Completed, &list.CreationTimestamp, &list.ModificationTimestamp, &list.DeletionTimestamp, &list.TemplateID, pq.Array(&list.TotalTagExclude)); err != nil {
		return types.ShoppingListSpec{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.ShoppingListSpec{}, err
	}
	return list, nil
}

// DeleteShoppingList ...
// deletes a shopping list, given a shopping list Id
func (m *ShoppingListManager) Delete(listID string) (err error) {
	err = m.manager.ShoppingItem().DeleteAll(listID)
	if err != nil {
		return ErrFailedToRemoveAllItemsFromList
	}
	sqlStatement := `delete from shopping_list where id = $1`
	rows, err := m.db.Query(sqlStatement, listID)
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

// GetListCount ...
// returns a count lists
func (m *ShoppingListManager) GetListCount() (count int, err error) {
	sqlStatement := `select count(*) from shopping_list`
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return 0, err
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

// DeleteCleanup ...
// cleans up shopping lists older than policy
func (m *ShoppingListManager) DeleteCleanup() (string, func() error) {
	return types.CronTabScheduleShoppingListCleanup, func() error {
		policy, err := m.manager.settingsManager.GetShoppingListKeepPolicy()
		if err != nil {
			return err
		}
		timestamp := time.Now()
		limit := -1
		switch policy {
		case types.ShoppingListKeepPolicyThreeMonths:
			timestamp = timestamp.AddDate(0, -3, 0)
		case types.ShoppingListKeepPolicySixMonths:
			timestamp = timestamp.AddDate(0, -6, 0)
		case types.ShoppingListKeepPolicyOneYear:
			timestamp = timestamp.AddDate(-1, 0, 0)
		case types.ShoppingListKeepPolicyTwoYears:
			timestamp = timestamp.AddDate(-2, 0, 0)
		case types.ShoppingListKeepPolicyLast10:
			limit = 10
		case types.ShoppingListKeepPolicyLast50:
			limit = 50
		case types.ShoppingListKeepPolicyLast100:
			limit = 100
		default:
			return nil
		}
		lists, err := m.List(types.ShoppingListOptions{
			Selector: types.ShoppingListSelector{
				CreationTimestampBefore: timestamp.Unix(),
			},
		})
		if err != nil {
			return err
		}
		if len(lists) == 0 {
			return nil
		}
		if limit != -1 {
			removeAmount := len(lists) - limit
			if len(lists) < limit {
				removeAmount = 0
			}
			lists = lists[:removeAmount]
		}
		log.Printf("[cleanup] Deleting old lists with policy %v\n", policy)
		for _, list := range lists {
			if err := m.Delete(list.ID); err != nil {
				return err
			}
		}
		return nil
	}
}
