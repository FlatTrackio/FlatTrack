import common from '@/frontend/common/common'
import login from '@/frontend/requests/public/login'
import registration from '@/frontend/requests/public/registration'

export default [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/frontend/views/authenticated/home.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '*',
    name: 'Unknown Page',
    component: () => import('@/frontend/views/global/unknown-page.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/flat',
    name: 'My Flat',
    component: () => import('@/frontend/views/authenticated/flat.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/about-flattrack',
    name: 'About FlatTrack',
    component: () => import('@/frontend/views/authenticated/about-flattrack.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/account',
    name: 'Account',
    component: () => import('@/frontend/views/authenticated/account-home.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/account/profile',
    name: 'Account Profile',
    component: () => import('@/frontend/views/authenticated/account-profile.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/account/security',
    name: 'Account Security',
    component: () => import('@/frontend/views/authenticated/account-security.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps',
    name: 'Apps',
    component: () => import('@/frontend/views/authenticated/apps.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/flatmates',
    name: 'My Flatmates',
    component: () => import('@/frontend/views/authenticated/flatmates.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list',
    name: 'Shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list/new',
    name: 'New shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list-new.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list/list',
    redirect: {
      name: 'Shopping list'
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list/list/:id',
    name: 'View shopping list',
    component: () => import('@/frontend/views/authenticated/shopping-list-view.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list/list/:id/new',
    name: 'New shopping list item',
    component: () => import('@/frontend/views/authenticated/shopping-list-item-new.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/apps/shopping-list/list/:listId/item/:itemId',
    name: 'View shopping list item',
    component: () => import('@/frontend/views/authenticated/shopping-list-item-view.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/admin',
    name: 'Admin home',
    component: () => import('@/frontend/views/admin/home.vue'),
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
    }
  },
  {
    path: '/admin/settings',
    name: 'Admin settings',
    component: () => import('@/frontend/views/admin/settings.vue'),
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
    }
  },
  {
    path: '/admin/accounts',
    name: 'Admin accounts',
    component: () => import('@/frontend/views/admin/accounts.vue'),
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
    }
  },
  {
    path: '/admin/accounts/new',
    name: 'Admin new account',
    component: () => import('@/frontend/views/admin/accounts-new.vue'),
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
    }
  },
  {
    path: '/admin/accounts/edit',
    redirect: {
      name: 'Admin accounts'
    },
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
    }
  },
  {
    path: '/admin/accounts/edit/:id',
    name: 'View user account',
    component: () => import('@/frontend/views/admin/account-edit.vue'),
    meta: {
      requiresAuth: true,
      requiresGroup: 'admin'
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
    meta: {
      requiresNoAuth: true
    }
  }
]
