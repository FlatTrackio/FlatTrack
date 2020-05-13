/*
  login
    manage JWTs saved
*/

import Request from '@/frontend/requests/requests'

// GetUserAuth
// validate JWT
function GetUserAuth (redirect) {
  return new Promise((resolve, reject) => {
    Request({
      url: '/api/user/auth',
      method: 'GET'
    }, redirect).then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    }, false, true)
  })
}

// PostUserAuth
// login and return a JWT
function PostUserAuth (credentials) {
  return Request({
    url: '/api/user/auth',
    method: 'POST',
    data: credentials
  }, false, true)
}

export default { GetUserAuth, PostUserAuth }
