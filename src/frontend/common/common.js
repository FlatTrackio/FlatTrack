/*
  common
    commonly used JS functions
*/

function getAuthToken () {
  return localStorage.getItem('authToken')
}

export default {
  getAuthToken
}
