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

// PatchProfile
// updates the profile of the authenticated account
function PatchProfile (names, email, phoneNumber, birthday, password) {
  return Request({
    url: `/api/user/profile`,
    method: 'PATCH',
    data: {
      names,
      email,
      phoneNumber,
      birthday,
      password
    }
  })
}

// ResetAuth
// revokes all JWTs
function PostAuthReset () {
  return Request({
    url: `/api/user/auth/reset`,
    method: 'POST'
  })
}

// GetProfilePicture
// get the profile picture for an account
function GetProfilePicture () {
  return Request({
    url: `/api/user/profile/picture`,
    method: 'GET'
  })
}

export default {
  GetProfile,
  PatchProfile,
  PostAuthReset,
  GetProfilePicture
}
