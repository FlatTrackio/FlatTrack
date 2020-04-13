/*
  shoppinglist
    helper functions for displaying the shopping list data
*/

function RestructureShoppingListToTags (responseList) {
  var currentTag = ''
  var list = []
  for (var item in responseList) {
    if (currentTag !== responseList[item].tag) {
      currentTag = responseList[item].tag
      var newItem = {
        tag: currentTag || 'Untagged',
        items: [responseList[item]]
      }

      list = [...list, newItem]
    } else {
      var currentListPosition = list.length - 1
      var currentSubListItems = list[currentListPosition].items

      list[currentListPosition].items = [...currentSubListItems, responseList[item]]
    }
  }
  return list
}

export default { RestructureShoppingListToTags }
