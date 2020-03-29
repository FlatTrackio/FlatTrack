/*
  registration
    register the FlatTrack instance
*/

import Request from '@/frontend/requests/requests'

// PostAdminRegister
// login and return a JWT
function PostAdminRegister (form) {
  return Request({
    url: '/api/admin/register',
    method: 'POST',
    data: form
  })
}

export default { PostAdminRegister }
