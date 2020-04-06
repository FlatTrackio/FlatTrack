/*
  shoppinglist
    manage shopping lists
*/

import Request from '@/frontend/requests/requests'

// GetShoppingLists
// returns a list of all shopping lists
function GetShoppingLists () {
  return Request({
    url: '/api/apps/shoppinglist',
    method: 'GET'
  })
}

// GetShoppingList
// returns a shopping list
function GetShoppingList (id) {
  return Request({
    url: `/api/apps/shoppinglist/${id}`,
    method: 'GET'
  })
}

// PostShoppingList
// given a name and optional notes, create a shopping list
function PostShoppingList (name, notes) {
  return Request({
    url: '/api/apps/shoppinglist',
    method: 'POST',
    data: {
      name,
      notes
    }
  })
}

// DeleteShoppingList
// deletes a shopping list
function DeleteShoppingList (id) {
  return Request({
    url: `/api/apps/shoppinglist/${id}`,
    method: 'DELETE'
  })
}

// TODO
// GetShoppingListItems
// returns shopping list items by id
function GetShoppingListItems (id) {
  return Request({
    url: `/api/apps/shoppinglist/${id}/items`,
    method: 'GET'
  })
}

// TODO
// PostShoppingListItems
// adds to the shopping list
function PostShoppingListItems (id) {
  return Request({
    url: `/api/apps/shoppinglist/${id}/items`,
    method: 'POST'
  })
}

// TODO
// DeleteShoppingListItem
// adds to the shopping list
function DeleteShoppingListItem (id, itemId) {
  return Request({
    url: `/api/apps/shoppinglist/${id}/items/${itemId}`,
    method: 'POST'
  })
}

export default { GetShoppingLists, GetShoppingList, PostShoppingList, GetShoppingListItems, PostShoppingListItems, DeleteShoppingListItem }
