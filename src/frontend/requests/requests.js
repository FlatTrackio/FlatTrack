import axios from 'axios'
import common from '@/frontend/common/common'

const siteSubPath = document.head.querySelector('[name~=sitesubpath][content]').content || window.location.host

function redirectToLogin (redirect) {
  if (redirect === false) {
    return
  }
  if (window.location.pathname !== siteSubPath + '/login') {
    // TODO direct by name instead of URL
    window.location.pathname = siteSubPath + '/login'
  }
}

function Request (request, redirect = true, publicRoute = false) {
  return new Promise((resolve, reject) => {
    var authToken = common.GetAuthToken()
    // if there is no token, and the request is a public route
    if ((typeof authToken === 'undefined' || authToken === null || authToken === '') && publicRoute !== true) {
      redirectToLogin(redirect)
    }
    request.headers = {
      Accept: 'application/json'
    }
    if (publicRoute !== true) {
      request.headers.Authorization = `bearer ${authToken}`
    }
    request.baseURL = window.location.origin + siteSubPath
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
