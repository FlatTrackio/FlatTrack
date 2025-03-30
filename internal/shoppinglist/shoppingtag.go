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

type ShoppingTagManager struct {
	manager *Manager
	db      *sql.DB
}

func (m *Manager) ShoppingTag() *ShoppingTagManager {
	return &ShoppingTagManager{
		manager: m,
		db:      m.db,
	}
}

// Validate ...
// given a shopping tag, return it's validity
func (m *ShoppingTagManager) Validate(tag types.ShoppingTag) (valid bool, err error) {
	if tag.Name != "" && len(tag.Name) == 0 || len(tag.Name) >= 30 {
		return false, ErrInvalidShoppingItemTag
	}
	// TODO check if one already exists with that name
	return true, err
}

// Create ...
// adds a new tag to be used in lists
func (m *ShoppingTagManager) Create(newTag types.ShoppingTag) (tag types.ShoppingTag, err error) {
	valid, err := m.Validate(newTag)
	if !valid || err != nil {
		return types.ShoppingTag{}, err
	}
	newTag.AuthorLast = newTag.Author
	sqlStatement := `insert into shopping_list_tag (name, author, authorLast)
                         values ($1, $2, $3)
                         returning *`
	rows, err := m.db.Query(sqlStatement, newTag.Name, newTag.Author, newTag.AuthorLast)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()

	tag, err = getTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
}

// ListTagsInList ...
// returns a list of tags used in items in a list
func (m *ShoppingTagManager) ListTagsInList(listID string) (tags []string, err error) {
	sqlStatement := `select distinct tag from shopping_item where listId = $1 order by tag`
	rows, err := m.db.Query(sqlStatement, listID)
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

// GetInList ...
// returns a tags used in items in a list
func (m *ShoppingTagManager) GetInList(listID string, tag string) (tagInDB string, err error) {
	sqlStatement := `select tag from shopping_item where listId = $1 and tag = $2`
	rows, err := m.db.Query(sqlStatement, listID, tag)
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

// UpdateInList ...
// updates a tag's name in a list
func (m *ShoppingTagManager) UpdateInList(listID string, tag string, tagUpdate string) (tagNew string, err error) {
	tagInDB, err := m.GetInList(listID, tag)
	if tagInDB == "" || err != nil {
		return "", ErrFailedToFindShoppingTagToUpdate
	}
	valid, err := m.manager.ShoppingTag().Validate(types.ShoppingTag{
		Name: tagUpdate,
	})
	if err != nil {
		return "", err
	}
	if !valid {
		return "", ErrInvalidShoppingItemTag
	}
	sqlStatement := `update shopping_item set tag = $3 where listId = $1 and tag = $2 returning tag`
	rows, err := m.db.Query(sqlStatement, listID, tag, tagUpdate)
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

// Get ...
// returns a tag, given an id
func (m *ShoppingTagManager) Get(id string) (tag types.ShoppingTag, err error) {
	sqlStatement := `select * from shopping_list_tag where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	tag, err = getTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
}

// List ...
// returns a list of all tags used in items across lists
func (m *ShoppingTagManager) List(options types.ShoppingTagOptions) (tags []types.ShoppingTag, err error) {
	sqlStatement := `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name asc`
	switch options.SortBy {
	case types.ShoppingTagSortByRecentlyUpdated:
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by modificationTimestamp desc`
	case types.ShoppingTagSortByLastUpdated:
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by modificationTimestamp asc`
	case types.ShoppingTagSortByLastAdded:
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by creationTimestamp asc`
	case types.ShoppingTagSortByAlphabeticalDescending:
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name asc`
	case types.ShoppingTagSortByAlphabeticalAscending:
		sqlStatement = `select * from shopping_list_tag
                         where deletionTimestamp = 0
	                 order by name desc`
	}
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return []types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		tag, err := getTagObjectFromRows(rows)
		if err != nil {
			return []types.ShoppingTag{}, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// UpdateShoppingTag ...
// updates a tag's name
func (m *ShoppingTagManager) Update(id string, tag types.ShoppingTag) (tagUpdated types.ShoppingTag, err error) {
	tagInDB, err := m.Get(id)
	if tagInDB.ID == "" || err != nil {
		return types.ShoppingTag{}, ErrFailedToFindShoppingTagToUpdate
	}
	valid, err := m.Validate(tag)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	if !valid {
		return types.ShoppingTag{}, ErrInvalidShoppingItemTag
	}
	if tag.Name != "" && len(tag.Name) == 0 || len(tag.Name) > 30 {
		return types.ShoppingTag{}, fmt.Errorf("Unable to use the provided tag, as it is either empty or too long or too short")
	}
	sqlStatement := `update shopping_list_tag set name = $2, authorLast = $3, modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int where id = $1 returning *`
	rows, err := m.db.Query(sqlStatement, id, tag.Name, tag.AuthorLast)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	tagUpdated, err = getTagObjectFromRows(rows)
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tagUpdated, nil
}

// Delete ...
// deletes a shopping tag
func (m *ShoppingTagManager) Delete(id string) (err error) {
	sqlStatement := `delete from shopping_list_tag where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
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

// getTagObjectFromRows ...
// returns a shopping tag object from rows
func getTagObjectFromRows(rows *sql.Rows) (tag types.ShoppingTag, err error) {
	if err := rows.Scan(&tag.ID, &tag.Name, &tag.Author, &tag.AuthorLast, &tag.CreationTimestamp, &tag.ModificationTimestamp, &tag.DeletionTimestamp); err != nil {
		return types.ShoppingTag{}, err
	}
	err = rows.Err()
	if err != nil {
		return types.ShoppingTag{}, err
	}
	return tag, nil
}
