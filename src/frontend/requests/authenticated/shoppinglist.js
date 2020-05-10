/*
  shoppinglist
    manage shopping lists
*/

import Request from '@/frontend/requests/requests'

// GetShoppingLists
// returns a list of all shopping lists
function GetShoppingLists () {
  return Request({
    url: '/api/apps/shoppinglist/lists',
    method: 'GET'
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
function GetShoppingListItems (id, sortBy) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/items`,
    method: 'GET',
    params: {
      sortBy
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
function PostShoppingListItem (id, name, notes, price, quantity, tag) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/items`,
    method: 'POST',
    data: {
      name,
      notes,
      price,
      quantity,
      tag
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
  DeleteShoppingListItem,
  GetShoppingListItemTags,
  UpdateShoppingListItemTag,
  PatchShoppingListItem,
  UpdateShoppingListItem,
  PatchShoppingListItemObtained
}
