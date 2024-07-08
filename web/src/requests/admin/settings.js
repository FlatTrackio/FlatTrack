/*
  settings
    manage admin settings
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

// PostFlatName
// changes the FlatName
function PostFlatName (flatName) {
  return Request({
    url: `/api/admin/settings/flatName`,
    method: 'POST',
    data: {
      flatName
    }
  })
}

// PutFlatNotes
// changes the FlatNotes
function PutFlatNotes (notes) {
  return Request({
    url: `/api/admin/settings/flatNotes`,
    method: 'PUT',
    data: {
      notes
    }
  })
}

// GetCostsWriteRequireGroupAdmin
// changes the setting for CostsWriteRequireGroupAdmin
function GetCostsWriteRequireGroupAdmin () {
  return Request({
    url: `/api/admin/settings/costsWriteRequireGroupAdmin`,
    method: 'GET'
  })
}

// PutCostsWriteRequireGroupAdmin
// changes the setting for CostsWriteRequireGroupAdmin
function PutCostsWriteRequireGroupAdmin (requireAdmin) {
  return Request({
    url: `/api/admin/settings/costsWriteRequireGroupAdmin`,
    method: 'PUT',
    data: {
      requireAdmin
    }
  })
}

export default {
  PostFlatName,
  PutFlatNotes,

  GetCostsWriteRequireGroupAdmin,
  PutCostsWriteRequireGroupAdmin
}
