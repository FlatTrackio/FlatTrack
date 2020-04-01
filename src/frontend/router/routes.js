import common from '@/frontend/common/common'
import login from '@/frontend/requests/public/login'
import registration from '@/frontend/requests/public/registration'

function checkForAuthToken (to, from, next) {
  var authToken = common.getAuthToken()
  if (typeof authToken === 'undefined' || authToken === null) {
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
      var authToken = common.getAuthToken()
      if (authToken === 'undefined' || authToken === null) {
        next()
        return
      }
      login.GetUserAuth().then(resp => {
        if (resp.data.data === true) {
          window.location.href = '/'
        }
        next()
      }).catch(() => {
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
        if (resp.data.data !== false) {
          window.location.href = '/'
        }
        next()
      }).catch(err => {
        console.log({ err })
        next()
      })
    }
  }
]
