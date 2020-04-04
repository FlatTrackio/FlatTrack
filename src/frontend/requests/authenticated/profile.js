/*
  profile
    manage your account
*/

import Request from '@/frontend/requests/requests'

// GetProfile
// returns the profile of the authenticated account
function GetProfile () {
  return Request({
    url: '/api/user/profile',
    method: 'GET'
  })
}

// PostProfile
// updates the profile of the authenticated account
function PostProfile (id) {
  return Request({
    url: `/api/user/profile`,
    method: 'POST'
  })
}

export default { GetProfile, PostProfile }
