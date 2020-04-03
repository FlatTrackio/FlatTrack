import common from '@/frontend/common/common'
import login from '@/frontend/requests/public/login'
import registration from '@/frontend/requests/public/registration'

function checkForAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
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
    path: '/flat',
    name: 'flat',
    component: () => import('@/frontend/views/authenticated/flat.vue'),
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

      function handleRedirections () {
        if (instanceRegistered && validAuthToken) {
          nextRoute = '/'
        } else if (!(hasAuthToken || validAuthToken)) {
          nextRoute = null
        } else if (!instanceRegistered) {
          nextRoute = '/setup'
        }

        console.log({ validAuthToken, instanceRegistered, hasAuthToken, nextRoute })

        if (!(nextRoute === null || nextRoute === '')) {
          window.location.pathname = nextRoute
        } else {
          next()
        }
      }
      // check that the instance is set up
      registration.GetInstanceRegistered().then(resp => {
        instanceRegistered = resp.data.data === true
      }).then(() => {
        // check if the authToken in localStorage isn't empty
        var authToken = common.GetAuthToken()
        hasAuthToken = (!(typeof authToken === 'undefined' || authToken === null || authToken === ''))
        // check if authToken is valid
        return login.GetUserAuth(false)
      }).then(resp => {
        console.log({ resp })
        validAuthToken = resp.data.data === true
        handleRedirections()
      }).catch(err => {
        console.log({ err })
        validAuthToken = err.response.data.data === true
        handleRedirections()
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
