/*
  registration
    register the FlatTrack instance
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

// GetInstanceRegistered
// determine if the instance of FlatTrack is registered
function GetInstanceRegistered () {
  return Request(
    {
      url: '/api/system/initialized',
      method: 'GET'
    },
    undefined,
    true
  )
}

// PostAdminRegister
// login and return a JWT
function PostAdminRegister (form, params) {
  console.log({ params })
  return Request(
    {
      url: '/api/admin/register',
      method: 'POST',
      data: form,
      params: {
        ...params
      }
    },
    false,
    true
  )
}

// GetInstanceRegistered
// return a list of timezones
function GetTimezones (secret) {
  return Request(
    {
      url: '/api/system/timezones',
      method: 'GET',
      params: {
        secret: secret
      }
    },
    false,
    true
  )
}

export default {
  GetInstanceRegistered,
  GetTimezones,
  PostAdminRegister
}
