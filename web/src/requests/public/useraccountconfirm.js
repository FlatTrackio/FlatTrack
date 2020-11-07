/*
  useraccountconfirm
    complete account registration
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

// GetTokenValid
// returns if a confirm token is available
function GetTokenValid (id) {
  return Request({
    url: `/api/user/confirm/${id}`,
    method: 'GET'
  }, false, true)
}

// PostUserConfirm
// posts to confirm an account
function PostUserConfirm (id, secret, phoneNumber, birthday, password) {
  return Request({
    url: `/api/user/confirm/${id}`,
    method: 'POST',
    params: {
      secret
    },
    data: {
      phoneNumber,
      birthday,
      password
    }
  }, false, true)
}

export default { GetTokenValid, PostUserConfirm }
