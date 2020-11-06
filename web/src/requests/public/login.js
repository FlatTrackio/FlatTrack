/*
  login
    manage JWTs saved
*/

import Request from '@/requests/requests'

// GetUserAuth
// validate JWT
function GetUserAuth (redirect) {
  return Request({
    url: '/api/user/auth',
    method: 'GET'
  }, redirect, true)
}

// PostUserAuth
// login and return a JWT
function PostUserAuth (email, password) {
  return Request({
    url: '/api/user/auth',
    method: 'POST',
    data: {
      email,
      password
    }
  }, false, true)
}

export default {
  GetUserAuth,
  PostUserAuth
}
