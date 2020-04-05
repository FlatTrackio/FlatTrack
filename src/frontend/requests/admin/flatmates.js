/*
  flatmates
    manage user accounts
*/

import Request from '@/frontend/requests/requests'

// PostFlatmate
// creates a user account
function PostFlatmate (form) {
  return Request({
    url: `/api/admin/users`,
    method: 'POST',
    data: form
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

export default { PostFlatmate, DeleteFlatmate }
