import common from '@/frontend/common/common'
import login from '@/frontend/requests/public/login'
import registration from '@/frontend/requests/public/registration'

function checkForAuthToken (to, from, next) {
  var authToken = common.getAuthToken()
  if (typeof authToken === 'undefined' || authToken === null || authToken === '') {
    window.location.href = '/login'
  }
  next()
}

export default [
  {
    path: '/',
    name: 'home',
    component: () => import('@/frontend/views/authenticated/home.vue'),
    beforeEnter: (to, from, next) => {
      checkForAuthToken(to, from, next)
    }
  },
  {
    path: '*',
    name: 'unknown-page',
    component: () => import('@/frontend/views/global/unknown-page.vue'),
    beforeEnter: (to, from, next) => {
      checkForAuthToken(to, from, next)
    }
  },
  {
    path: '/apps',
    name: 'apps',
    component: () => import('@/frontend/views/authenticated/apps.vue'),
    beforeEnter: (to, from, next) => {
      checkForAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/flatmates',
    name: 'flatmates',
    component: () => import('@/frontend/views/authenticated/flatmates.vue'),
    beforeEnter: (to, from, next) => {
      checkForAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list',
    name: 'shopping-list',
    component: () => import('@/frontend/views/authenticated/shopping-list.vue'),
    beforeEnter: (to, from, next) => {
      checkForAuthToken(to, from, next)
    }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/frontend/views/public/login.vue'),
    beforeEnter: (to, from, next) => {
      var instanceRegistered
      var hasAuthToken
      var validAuthToken
      var nextRoute = ''
      // check that the instance is set up
      registration.GetInstanceRegistered().then(resp => {
        instanceRegistered = resp.data.data === true
      }).then(resp => {
        // check if the authToken in localStorage isn't empty
        var authToken = common.getAuthToken()
        hasAuthToken = (!(typeof authToken === 'undefined' || authToken === null || authToken === ''))
        // check if authToken is valid
        login.GetUserAuth(false).then(resp => {
          return resp
        }).catch(() => {
          return false
        })
      }).then(resp => {
        if (typeof resp === 'undefined') {
          validAuthToken = false
        } else {
          validAuthToken = resp.data.data === true
        }
        if (instanceRegistered && validAuthToken) {
          nextRoute = '/'
        } else if (!hasAuthToken && instanceRegistered) {
          nextRoute = null
        } else if (!instanceRegistered) {
          nextRoute = '/setup'
        }

        if (nextRoute !== null) {
          window.location.pathname = nextRoute
        } else {
          next()
        }
      }).catch(err => {
        console.log({ err })
        next()
      })
    }
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('@/frontend/views/public/forgot-password.vue')
  },
  // {
  //   path: '/admin',
  //   name: 'admin',
  //   component: () => import('@/frontend/views/admin/home.vue')
  // },
  {
    path: '/setup',
    name: 'setup',
    component: () => import('@/frontend/views/public/setup.vue'),
    beforeEnter: (to, from, next) => {
      registration.GetInstanceRegistered().then(resp => {
        if (resp.data.data === true) {
          window.location.pathname = '/'
        }
        next()
      }).catch(() => {
        next()
      })
    }
  }
]
