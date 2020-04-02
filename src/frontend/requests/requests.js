import axios from 'axios'

function getAuthToken () {
  return localStorage.getItem('authToken')
}

function redirectToLogin (redirect) {
  if (redirect === false) {
    return
  }
  if (window.location.pathname !== '/login') {
    console.log('redirecting to /login')
    window.location.href = '/login'
  }
}

function Request (request, redirect) {
  return new Promise((resolve, reject) => {
    var authToken = getAuthToken()
    if (typeof authToken === 'undefined' || authToken === null || authToken === '') {
      redirectToLogin(redirect)
    }
    request.headers = {
      Authorization: 'bearer ' + authToken
    }
    axios(request)
      .then(resp => resolve(resp))
      .catch(err => {
        if (err.response.status === 401) {
          redirectToLogin()
        }
        reject(err)
      })
  })
}

export default Request
