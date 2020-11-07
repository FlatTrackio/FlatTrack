/*
  shoppinglist
    helper functions for displaying the shopping list data
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

// RestructureShoppingListToTags
// returns items structured by tags
function RestructureShoppingListToTags (responseList) {
  var currentTag = ''
  var list = []
  responseList.forEach(item => {
    if (currentTag !== item.tag) {
      currentTag = item.tag
      var newItem = {
        tag: currentTag || 'Untagged',
        items: [item],
        price: item.price * item.quantity || 0
      }
      list = [...list, newItem]
      return
    }
    var currentListPosition = list.length - 1
    var currentSubListItems = list[currentListPosition].items
    list[currentListPosition].items = [...currentSubListItems, item]
    list[currentListPosition].price += (item.price * item.quantity || 0)
  })
  return list
}

// GetShoppingListFromCache
// given an id returns a list if available
function GetShoppingListFromCache (id) {
  var items = localStorage.getItem(`shoppinglist.list.${id}.items`)
  return JSON.parse(items) || []
}

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
  return localStorage.getItem('shoppinglist.autorefresh')
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
