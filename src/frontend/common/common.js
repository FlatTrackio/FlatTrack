/*
  common
    commonly used JS functions
*/

// getAuthToken
// return the JWT from localStorage
function getAuthToken () {
  return localStorage.getItem('authToken')
}

// deleteAuthToken
// remove the JWT from localStorage
function deleteAuthToken () {
  return localStorage.removeItem('authToken')
}

export default {
  getAuthToken,
  deleteAuthToken
}
