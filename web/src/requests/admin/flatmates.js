/*
  flatmates
    manage user accounts
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

// PostFlatmate
// creates a user account
function PostFlatmate (names, email, phoneNumber, birthday, groups, password) {
  return Request({
    url: `/api/admin/users`,
    method: 'POST',
    data: {
      names,
      email,
      phoneNumber,
      birthday,
      groups,
      password
    }
  })
}

// PatchFlatmate
// patches a user account
function PatchFlatmate (id, names, email, phoneNumber, birthday, groups, password) {
  return Request({
    url: `/api/admin/users/${id}`,
    method: 'PATCH',
    data: {
      names,
      email,
      phoneNumber,
      birthday,
      groups,
      password
    }
  })
}

// PatchFlatmateDisabled
// patches a user account disabled field
function PatchFlatmateDisabled (id, disabled) {
  return Request({
    url: `/api/admin/users/${id}/disabled`,
    method: 'PATCH',
    data: {
      disabled
    }
  })
}

// PatchFlatmate
// updates a user account
function UpdateFlatmate (id, names, email, phoneNumber, birthday, groups, password) {
  return Request({
    url: `/api/admin/users/${id}`,
    method: 'PUT',
    data: {
      names,
      email,
      phoneNumber,
      birthday,
      groups,
      password
    }
  })
}

// DeleteFlatmate
// removes a user account by id
function DeleteFlatmate (id) {
  return Request({
    url: `/api/admin/users/${id}`,
    method: 'DELETE'
  })
}

// GetUserAccountConfirms
// returns a list of confirm tokens and secrets
function GetUserAccountConfirms (userId) {
  return Request({
    url: `/api/admin/useraccountconfirms`,
    method: 'GET',
    params: {
      userId
    }
  })
}

export default {
  PostFlatmate,
  PatchFlatmate,
  PatchFlatmateDisabled,
  DeleteFlatmate,
  GetUserAccountConfirms
}
