import axios from 'axios'

function getAuthToken () {
  return localStorage.getItem('authToken')
}

function Request (request) {
  return new Promise((resolve, reject) => {
    var authToken = getAuthToken()
    if (typeof authToken === 'undefined') {
      reject(401) // eslint-disable-line prefer-promise-reject-errors
      window.location.href = '/login'
    }
    request.headers = {
      Authorization: 'Bearer ' + authToken
    }
    axios(request)
      .then(resp => resolve(resp))
      .catch(err => {
        if (err.response.status === 401) {
          localStorage.authToken = ''
          window.location.href = '/login'
        }
        reject(err)
      })
  })
}

export default Request
