/*
  costs
    get your costs
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

// GetView
// get the costs view
function GetView () {
  return Request({
    url: '/api/apps/costs/view',
    method: 'GET'
  })
}

// GetCostsWriteRequireGroupAdmin
// get the costs setting for write requiring admin group permissions
function GetCostsWriteRequireGroupAdmin () {
  return Request({
    url: '/api/apps/costs/writeRequireGroupAdmin',
    method: 'GET'
  })
}

// GetCostItems
// get a list of cost items
function GetCostsItems () {
  return Request({
    url: '/api/apps/costs/items',
    method: 'GET'
  })
}

// GetCostsItem
// get a list of cost items
function GetCostsItem (id) {
  return Request({
    url: `/api/apps/costs/items/${id}`,
    method: 'GET'
  })
}

// PostCostsItem
// create a new cost item
function PostCostsItem (item) {
  return Request({
    url: '/api/apps/costs/items',
    method: 'POST',
    data: item
  })
}

// PutCostsItem
// create a new cost item
function PutCostsItem (id, item) {
  return Request({
    url: `/api/apps/costs/items/${id}`,
    method: 'PUT',
    data: item
  })
}

// PatchCostsItem
// create a new cost item
function PatchCostsItem (id, item) {
  return Request({
    url: `/api/apps/costs/items/${id}`,
    method: 'PATCH',
    data: item
  })
}

// DeleteCostsItem
// delete a new cost item
function DeleteCostsItem (id) {
  return Request({
    url: `/api/apps/costs/items/${id}`,
    method: 'DELETE'
  })
}

// DeleteCostsItem
// delete a new cost item
function DeleteCostsItems (ids) {
  return Request({
    url: `/api/apps/costs/items`,
    method: 'DELETE',
    data: {
      ids
    }
  })
}

// GetCostsNotes
// changes the FlatNotes
function GetCostsNotes () {
  return Request({
    url: `/api/apps/costs/settings/notes`,
    method: 'GET'
  })
}

// PutCostsNotes
// changes the FlatNotes
function PutCostsNotes (notes) {
  return Request({
    url: `/api/admin/settings/costsNotes`,
    method: 'PUT',
    data: {
      notes
    }
  })
}

export default {
  GetView,
  GetCostsWriteRequireGroupAdmin,

  GetCostsItems,
  GetCostsItem,
  PostCostsItem,
  PutCostsItem,
  PatchCostsItem,
  DeleteCostsItem,
  DeleteCostsItems,

  GetCostsNotes,
  PutCostsNotes
}
