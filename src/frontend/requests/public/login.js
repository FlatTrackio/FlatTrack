/*
  login
    manage JWTs saved
*/

import Request from '@/frontend/requests/requests'

// GetUserAuth
// validate JWT
function GetUserAuth (redirect) {
  return Request({
    url: '/api/user/auth',
    method: 'GET'
  }, redirect)
}

// PostUserAuth
// login and return a JWT
function PostUserAuth (credentials) {
  return Request({
    url: '/api/user/auth',
    method: 'POST',
    data: credentials
  })
}

export default { GetUserAuth, PostUserAuth }
