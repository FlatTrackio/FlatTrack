/*
  shoppinglist
    helper functions for displaying the shopping list data
*/

// RestructureShoppingListToTags
// returns items structured by tags
function RestructureShoppingListToTags (responseList) {
  var currentTag = ''
  var list = []
  for (var item in responseList) {
    if (currentTag !== responseList[item].tag) {
      currentTag = responseList[item].tag
      var newItem = {
        tag: currentTag || 'Untagged',
        items: [responseList[item]],
        price: responseList[item].price * responseList[item].quantity || 0
      }

      list = [...list, newItem]
    } else {
      var currentListPosition = list.length - 1
      var currentSubListItems = list[currentListPosition].items

      list[currentListPosition].items = [...currentSubListItems, responseList[item]]
      list[currentListPosition].price += (responseList[item].price * responseList[item].quantity || 0)
    }
  }
  return list
}

// GetShoppingListFromCache
// given an id returns a list if available
function GetShoppingListFromCache (id) {
  var items = localStorage.getItem(`shoppinglist.list.${id}.items`)
  return JSON.parse(items || [])
}

// WriteShoppingListToCache
// given an id writes a list to the cache
function WriteShoppingListToCache (id, items) {
  localStorage.setItem(`shoppinglist.list.${id}.items`, JSON.stringify(items || []))
}

// DeleteShoppingListFromCache
// given an id deletes a cached list
function DeleteShoppingListFromCache (id, items) {
  localStorage.removeItem(`shoppinglist.list.${id}.items`)
}

// returns if the shopping list should auto refresh
function GetShoppingListAutoRefresh () {
  return Boolean(localStorage.getItem('shoppinglist.autorefresh')) || true
}

// WriteShoppingListAutoRefresh
// writes if the shopping list should auto refresh
function WriteShoppingListAutoRefresh (autorefresh) {
  return localStorage.setItem('shoppinglist.autorefresh', autorefresh)
}

// GetShoppingListSortBy
// returns how the shopping list should sort by
function GetShoppingListSortBy () {
  return localStorage.getItem('shoppinglist.sortBy') || 'tags'
}

// WriteShoppingListSortBy
// writes how the shopping list should sort by
function WriteShoppingListSortBy (sortBy) {
  return localStorage.setItem('shoppinglist.sortBy', sortBy)
}

// GetShoppingListObtainedFilter
// returns the shopping list obtained filter
function GetShoppingListObtainedFilter (id) {
  return Number(localStorage.getItem(`shoppinglist.list.${id}.obtainedFilter`)) || 0
}

// WriteShoppingListObtainedFilter
// writes the shopping list obtained filter
function WriteShoppingListObtainedFilter (id, state) {
  return localStorage.setItem(`shoppinglist.list.${id}.obtainedFilter`, state)
}

// GetShoppingListSearch
// returns the last search from the current shopping list
function GetShoppingListSearch (id) {
  return sessionStorage.getItem(`shoppinglist.list.${id}.search`) || ''
}

// WriteShoppingListAutoRefresh
// writes the last search string of the current shopping list
function WriteShoppingListSearch (id, search) {
  return sessionStorage.setItem(`shoppinglist.list.${id}.search`, search)
}

export default {
  RestructureShoppingListToTags,
  GetShoppingListAutoRefresh,
  GetShoppingListFromCache,
  WriteShoppingListToCache,
  DeleteShoppingListFromCache,
  GetShoppingListSortBy,
  WriteShoppingListSortBy,
  GetShoppingListObtainedFilter,
  WriteShoppingListObtainedFilter,
  GetShoppingListSearch,
  WriteShoppingListSearch
}
