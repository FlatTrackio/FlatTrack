/*
  login
    manage JWTs saved
*/

import Request from '@/frontend/requests/requests'

// GetUserAuth
// validate JWT
function GetUserAuth () {
  return new Promise((resolve, reject) => {
    Request({
      url: '/api/user/auth',
      method: 'GET'
    }).then(res => resolve())
      .catch(err => {
        if (err.response.status === 401) {
          reject(err)
        }
      })
  })
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

export { GetUserAuth, PostUserAuth }
