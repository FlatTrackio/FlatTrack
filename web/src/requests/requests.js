// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import axios from 'axios'
import common from '@/common/common'

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
    request.headers = {
      Accept: 'application/json'
    }
    if (publicRoute !== true) {
      request.headers.Authorization = `bearer ${authToken}`
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
