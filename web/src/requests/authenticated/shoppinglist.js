/*
  shoppinglist
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

import Request from '@/requests/requests'

// GetShoppingLists
// returns a list of all shopping lists
function GetShoppingLists (completed, sortBy, creationTimestampAfter, modificationTimestampAfter, limit) {
  return Request({
    url: '/api/apps/shoppinglist/lists',
    method: 'GET',
    params: {
      completed,
      sortBy,
      creationTimestampAfter,
      modificationTimestampAfter,
      limit
    }
  })
}

// GetShoppingList
// returns a shopping list
function GetShoppingList (id) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}`,
    method: 'GET'
  })
}

// PostShoppingList
// given a name and optional notes, create a shopping list
function PostShoppingList (name, notes, templateId, templateListItemSelector) {
  return Request({
    url: '/api/apps/shoppinglist/lists',
    method: 'POST',
    data: {
      name,
      notes,
      templateId
    },
    params: {
      templateListItemSelector
    }
  })
}

// PatchShoppingList
// given a name and optional notes, patch a shopping list
function PatchShoppingList (id, name, notes) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}`,
    method: 'PATCH',
    data: {
      name,
      notes
    }
  })
}

// UpdateShoppingList
// given a name and optional notes, patch a shopping list
function UpdateShoppingList (id, name, notes, completed) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}`,
    method: 'PUT',
    data: {
      name,
      notes,
      completed
    }
  })
}

// PatchShoppingListCompleted
// given a bool, patch a shopping list's completed field
function PatchShoppingListCompleted (id, completed) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/completed`,
    method: 'PATCH',
    data: {
      completed
    }
  })
}

// DeleteShoppingList
// deletes a shopping list
function DeleteShoppingList (id) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}`,
    method: 'DELETE'
  })
}

// GetShoppingListItems
// returns shopping list items by id
function GetShoppingListItems (id, sortBy, obtained) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/items`,
    method: 'GET',
    params: {
      sortBy,
      obtained
    }
  })
}

// GetShoppingListItem
// returns shopping item by id
function GetShoppingListItem (listId, itemId) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/items/${itemId}`,
    method: 'GET'
  })
}

// PostShoppingListItem
// adds to the shopping list
function PostShoppingListItem (id, name, notes, price, quantity, tag, obtained) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/items`,
    method: 'POST',
    data: {
      name,
      notes,
      price,
      quantity,
      tag,
      obtained
    }
  })
}

// PatchShoppingListItem
// patches the shopping list item
function PatchShoppingListItem (listId, itemId, name, notes, price, quantity, tag) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/items/${itemId}`,
    method: 'PATCH',
    data: {
      name,
      notes,
      price,
      quantity,
      tag
    }
  })
}

// UpdateShoppingListItem
// updates the shopping list item
function UpdateShoppingListItem (listId, itemId, name, notes, price, quantity, tag, obtained) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/items/${itemId}`,
    method: 'PUT',
    data: {
      name,
      notes,
      price,
      quantity,
      tag,
      obtained
    }
  })
}

// PatchShoppingListItemObtained
// adds to the shopping list
function PatchShoppingListItemObtained (listId, itemId, obtained) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/items/${itemId}/obtained`,
    method: 'PATCH',
    data: {
      obtained
    }
  })
}

// DeleteShoppingListItem
// adds to the shopping list
function DeleteShoppingListItem (listId, itemId) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/items/${itemId}`,
    method: 'DELETE'
  })
}

// GetShoppingListItemTags
// fetches all tags used in a list
function GetShoppingListItemTags (listId) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/tags`,
    method: 'GET'
  })
}

// GetAllShoppingListItemTags
// fetches all tags
function GetAllShoppingListItemTags () {
  return Request({
    url: `/api/apps/shoppinglist/tags`,
    method: 'GET'
  })
}

// UpdateShoppingListItemTag
// updates a tag name used in a list
function UpdateShoppingListItemTag (listId, tagName, tagNameNew) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${listId}/tags/${tagName}`,
    method: 'PUT',
    data: {
      name: tagNameNew
    }
  })
}

// PostShoppingTag
// creates a shopping tag
function PostShoppingTag (name) {
  return Request({
    url: `/api/apps/shoppinglist/tags`,
    method: 'POST',
    data: {
      name
    }
  })
}

// GetShoppingTags
// returns a list of shopping tags
function GetShoppingTags (sortBy) {
  return Request({
    url: `/api/apps/shoppinglist/tags`,
    method: 'GET',
    params: {
      sortBy
    }
  })
}

// GetShoppingTag
// returns a list of shopping tag
function GetShoppingTag (id) {
  return Request({
    url: `/api/apps/shoppinglist/tags/${id}`,
    method: 'GET'
  })
}

// UpdateShoppingTag
// updates a shopping tag name
function UpdateShoppingTag (id, name) {
  return Request({
    url: `/api/apps/shoppinglist/tags/${id}`,
    method: 'PUT',
    data: {
      name
    }
  })
}

// DeleteShoppingTag
// removes a shopping tag
function DeleteShoppingTag (id) {
  return Request({
    url: `/api/apps/shoppinglist/tags/${id}`,
    method: 'DELETE'
  })
}

// GetShoppingListNotes
// returns a notes of shopping lists
function GetShoppingListNotes () {
  return Request({
    url: `/api/apps/shoppinglist/settings/notes`,
    method: 'GET'
  })
}

// PutShoppingListNotes
// updates a notes of shopping lists
function PutShoppingListNotes (notes) {
  return Request({
    url: `/api/admin/settings/shoppingListNotes`,
    method: 'PUT',
    data: {
      notes
    }
  })
}

export default {
  GetShoppingLists,
  GetShoppingList,
  PostShoppingList,
  PatchShoppingList,
  UpdateShoppingList,
  PatchShoppingListCompleted,
  DeleteShoppingList,

  GetShoppingListItems,
  GetShoppingListItem,
  PostShoppingListItem,
  PatchShoppingListItem,
  UpdateShoppingListItem,
  PatchShoppingListItemObtained,
  DeleteShoppingListItem,

  GetAllShoppingListItemTags,
  GetShoppingListItemTags,
  UpdateShoppingListItemTag,
  PostShoppingTag,
  GetShoppingTags,
  GetShoppingTag,
  UpdateShoppingTag,
  DeleteShoppingTag,

  GetShoppingListNotes,
  PutShoppingListNotes
}
