import common from '@/frontend/common/common'
import cani from '@/frontend/requests/authenticated/can-i'
import login from '@/frontend/requests/public/login'
import registration from '@/frontend/requests/public/registration'

// requireAuthToken
// given an no auth token redirect to the login page
function requireAuthToken (to, from, next) {
  var authToken = common.GetAuthToken()
  if (typeof authToken === 'undefined' || authToken === null || authToken === '') {
    next('/login')
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
  window.location.href = '/'
}

function checkForAdminGroup (to, from, next) {
  cani.GetCanIgroup('admin').then(resp => {
    if (resp.data.spec === true) {
      next()
    } else {
      next(from.path)
    }
  }).catch(() => {
    next(from.path)
  })
}

export default [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/frontend/views/authenticated/home.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '*',
    name: 'Unknown Page',
    component: () => import('@/frontend/views/global/unknown-page.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/flat',
    name: 'My Flat',
    component: () => import('@/frontend/views/authenticated/flat.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/frontend/views/authenticated/profile.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps',
    name: 'Apps',
    component: () => import('@/frontend/views/authenticated/apps.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/flatmates',
    name: 'My Flatmates',
    component: () => import('@/frontend/views/authenticated/flatmates.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list',
    name: 'Shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list/new',
    name: 'New shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list-new.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list/list',
    redirect: {
      name: 'Shopping list'
    },
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list/list/:id',
    name: 'View shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list-view.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list/list/:id/new',
    name: 'New shopping list item',
    component: () => import('@/frontend/views/authenticated/shopping-list-item-new.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/apps/shopping-list/list/:listId/item/:itemId',
    name: 'View shopping list item',
    component: () => import('@/frontend/views/authenticated/shopping-list-item-view.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/admin',
    name: 'Admin home',
    component: () => import('@/frontend/views/admin/home.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
      checkForAdminGroup(to, from, next)
    }
  },
  {
    path: '/admin/accounts',
    name: 'Admin accounts',
    component: () => import('@/frontend/views/admin/accounts.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
      checkForAdminGroup(to, from, next)
    }
  },
  {
    path: '/admin/accounts/new',
    name: 'Admin new account',
    component: () => import('@/frontend/views/admin/accounts-new.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
      checkForAdminGroup(to, from, next)
    }
  },
  {
    path: '/admin/accounts/edit',
    redirect: {
      name: 'Admin accounts'
    },
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
    }
  },
  {
    path: '/admin/accounts/edit/:id',
    name: 'View user account',
    component: () => import('@/frontend/views/admin/account-edit.vue'),
    beforeEnter: (to, from, next) => {
      requireAuthToken(to, from, next)
      checkForAdminGroup(to, from, next)
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/frontend/views/public/login.vue'),
    beforeEnter: (to, from, next) => {
      var instanceRegistered
      var hasAuthToken
      var validAuthToken
      var nextRoute

      function handleRedirections () {
        if (instanceRegistered && validAuthToken) {
          nextRoute = '/'
        } else if (!instanceRegistered) {
          nextRoute = '/setup'
        } else if (!(hasAuthToken || validAuthToken)) {
          nextRoute = null
        }

        if (!(nextRoute === null || nextRoute === '')) {
          next(nextRoute)
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
        validAuthToken = resp.data.data === true
        handleRedirections()
      }).catch(err => {
        validAuthToken = err.response.data.data === true
        handleRedirections()
        next()
      })
    }
  },
  {
    path: '/forgot-password',
    name: 'Forgot password',
    component: () => import('@/frontend/views/public/forgot-password.vue')
  },
  // {
  //   path: '/admin',
  //   name: 'admin',
  //   component: () => import('@/frontend/views/admin/home.vue')
  // },
  {
    path: '/setup',
    name: 'Set up',
    component: () => import('@/frontend/views/public/setup.vue'),
    beforeEnter: (to, from, next) => {
      registration.GetInstanceRegistered().then(resp => {
        if (resp.data.data === true) {
          next('/')
          return
        }
        next()
      }).catch(() => {
        next()
      })
    }
  },
  {
    path: '/useraccountconfirm/:id',
    name: 'User account confirm',
    component: () => import('@/frontend/views/public/useraccountconfirm.vue'),
    beforeEnter: (to, from, next) => {
      requireNoAuthToken(to, from, next)
    }
  }
]
