import axios from 'axios'
import common from '@/frontend/common/common'

function redirectToLogin (redirect) {
  if (redirect === false) {
    return
  }
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

function Request (request, redirect = true, publicRoute = false) {
  return new Promise((resolve, reject) => {
    var authToken = common.GetAuthToken()
    // if there is no token, and the request is a public route
    if ((typeof authToken === 'undefined' || authToken === null || authToken === '') && publicRoute !== true) {
      redirectToLogin(redirect)
    }
    if (publicRoute !== true) {
      request.headers = {
        Authorization: 'bearer ' + authToken
      }
    }
    axios(request)
      .then(resp => resolve(resp))
      .catch(err => {
        if (err.response.status === 401) {
          redirectToLogin(redirect)
        }
        reject(err)
      })
  })
}

export default Request
