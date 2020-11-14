/*
  flatInfo
    fetch info
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

// GetFlatName
// gets the name of the flat
function GetFlatName () {
  return Request({
    url: '/api/system/flatName',
    method: 'GET'
  })
}

// GetFlatNotes
// gets the notes of the flat
function GetFlatNotes () {
  return Request({
    url: '/api/flat/info',
    method: 'GET'
  })
}

export default {
  GetFlatName,
  GetFlatNotes
}
