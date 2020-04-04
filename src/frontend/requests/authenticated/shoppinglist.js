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
    url: '/api/apps/shoppinglist/{id}',
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

export default { GetShoppingLists, GetShoppingList, PostShoppingList }
