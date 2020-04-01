import axios from 'axios'

function getAuthToken () {
  return localStorage.getItem('authToken')
}

function redirectToLogin (redirect) {
  if (window.location.href !== '/login' && redirect !== false) {
    window.location.href = 'login'
  }
}

function Request (request, redirect) {
  return new Promise((resolve, reject) => {
    var authToken = getAuthToken()
    if (typeof authToken === 'undefined') {
      reject(401) // eslint-disable-line prefer-promise-reject-errors
      redirectToLogin(redirect)
    }
    request.headers = {
      Authorization: 'Bearer ' + authToken
    }
    axios(request)
      .then(resp => resolve(resp))
      .catch(err => {
        // console.log({ err }, err.response)
        // if (err.response && err.response.status === 401) {
        //   delete localStorage['authToken']
        //   redirectToLogin(redirect)
        // }
        reject(err)
      })
  })
}

export default Request
