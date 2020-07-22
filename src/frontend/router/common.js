import common from '@/frontend/common/common'
import cani from '@/frontend/requests/authenticated/can-i'

// requireAuthToken
// given an no auth token redirect to the login page
function requireAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
  if (typeof authToken === 'undefined' || authToken === null || authToken === '') {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  next()
}

// requireNoAuthToken
// given an auth token, redirect to the home page
function requireNoAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
  if (typeof authToken === 'undefined' || authToken === null || authToken === '') {
    next()
    return
  }
  this.$router.push({ name: 'Home' })
}

function requireGroup (to, from, next) {
  cani.GetCanIgroup(to.meta.requiresGroup).then(resp => {
    if (resp.data.data === true) {
      next()
    } else {
      next(from.path)
    }
  }).catch(() => {
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
