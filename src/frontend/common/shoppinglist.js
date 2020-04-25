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
        price: responseList[item].price || 0
      }

      list = [...list, newItem]
    } else {
      var currentListPosition = list.length - 1
      var currentSubListItems = list[currentListPosition].items

      list[currentListPosition].items = [...currentSubListItems, responseList[item]]
      list[currentListPosition].price += (responseList[item].price || 0)
    }
  }
  return list
}

// GetShoppingListFromCache
// given an id returns a list if available
function GetShoppingListFromCache (id) {
  return JSON.parse(localStorage.getItem(`shoppinglist.list.${id}.items`) || [])
}

// WriteShoppingListToCache
// given an id writes a list to the cache
function WriteShoppingListToCache (id, items) {
  localStorage.setItem(`shoppinglist.list.${id}.items`, JSON.stringify(items || []))
}

// GetShoppingListAutoRefresh
// returns if the shopping list should be auto refresh
function GetShoppingListAutoRefresh () {
  return localStorage.getItem('shoppinglist.autorefresh') || true
}

export default { RestructureShoppingListToTags, GetShoppingListAutoRefresh, GetShoppingListFromCache, WriteShoppingListToCache }
