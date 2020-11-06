/*
  flatmates
    manage user accounts
*/

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
