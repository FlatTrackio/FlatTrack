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
function PostShoppingList (name, notes, templateId, notObtained) {
  return Request({
    url: '/api/apps/shoppinglist/lists',
    method: 'POST',
    data: {
      name,
      notes,
      templateId
    },
    params: {
      notObtained
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
function GetShoppingListItems (id) {
  return Request({
    url: `/api/apps/shoppinglist/lists/${id}/items`,
    method: 'GET'
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
// adds to the shopping list
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

function GetShoppingListItemTags () {
  return Request({
    url: '/api/apps/shoppinglist/tags',
    method: 'GET'
  })
}

export default { GetShoppingLists, GetShoppingList, PostShoppingList, PatchShoppingList, PatchShoppingListCompleted, DeleteShoppingList, GetShoppingListItems, GetShoppingListItem, PostShoppingListItem, DeleteShoppingListItem, GetShoppingListItemTags, PatchShoppingListItem, PatchShoppingListItemObtained }
