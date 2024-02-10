import common from '@/common/common'
import cani from '@/requests/authenticated/can-i'
import login from '@/requests/public/login'

// requireAuthToken
// given an no auth token redirect to the login page
function requireAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
  if (
    typeof authToken === 'undefined' ||
    authToken === null ||
    authToken === ''
  ) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }
  next()
}

// requireNoAuthToken
// given an auth token, redirect to the home page
function requireNoAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
  if (
    typeof authToken === 'undefined' ||
    authToken === null ||
    authToken === ''
  ) {
    next()
    return
  }
  login.GetUserAuth(false).then(() => {
    window.location.href = '/'
  })
}

function requireGroup (to, from, next) {
  cani
    .GetCanIgroup(to.meta.requiresGroup)
    .then((resp) => {
      if (resp.data.data === true) {
        next()
      } else {
        next(from.path)
      }
    })
    .catch(() => {
      next(from.path)
    })
}

function isPublicRoute (to) {
  return to.meta.requiresAuth !== true
}

export default {
  requireAuthToken,
  requireNoAuthToken,
  requireGroup,
  isPublicRoute
}
