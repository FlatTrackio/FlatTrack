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

function GetNotes () {
  return new Promise((resolve, reject) => {
    resolve({
      resp: {
        data: {
          spec: 'Hiiii'
        }
      }
    })
  })
}

export default {
  GetView,
  GetNotes
}
