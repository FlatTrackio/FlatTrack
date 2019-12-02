import Vue from 'vue'
import Router from 'vue-router'
// import axios from 'axios'

Vue.use(Router)
Vue.component('home', () => import('../components/home.vue'))

const authenticatedRoute = async (to, before, next) => {
  var authToken
  try {
    authToken = localStorage.getItem('authToken')
    /*
    var isValid = await axios({
      url: '/api/auth/validate',
      method: 'post',
      headers: {
        'authorization': `Bearer ${authToken}`
      }
    })
    console.log(isValid)
    */
  } catch (err) {
    return next({ path: '/login' })
  }
  return authToken ? next() : next({ path: '/login' })
}

const unauthenticatedRouteUninitialised = async (to, before, next) => {
  var authToken
  try {
    authToken = localStorage.getItem('authToken')
  } catch (err) {
    console.error(err)
  }
  return authToken ? next({ path: '/' }) : next()
}

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/components/authenticated/home'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '*',
      name: 'unknown-page',
      component: () => import('@/components/global/unknown-page')
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import(/* webpackChunkName: "tasks" */ '@/components/authenticated/tasks'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/tasks/t',
      name: 'task-view',
      component: () => import(/* webpackChunkName: "tasks" */ '@/components/authenticated/task'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/about/flattrack',
      name: 'aboutflattrack',
      component: () => import('@/components/authenticated/about-flattrack'),
      alias: '/aboutflattrack',
      beforeEnter: authenticatedRoute
    },
    {
      path: '/about/flat',
      name: 'aboutflat',
      component: () => import('@/components/authenticated/about-flat'),
      alias: '/aboutflat',
      beforeEnter: authenticatedRoute
    },
    {
      path: '/high-fives',
      name: 'highfives',
      component: () => import('@/components/authenticated/high-fives'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/members',
      name: 'members',
      component: () => import('@/components/authenticated/members'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/noticeboard',
      name: 'noticeboard',
      component: () => import('@/components/authenticated/noticeboard'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/noticeboard/p',
      name: 'noticeboard posts',
      component: () => import('@/components/authenticated/noticeboard-post'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/recipes',
      name: 'recipes',
      component: () => import('@/components/authenticated/recipes'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/shared-calendar',
      name: 'shared-calendar',
      component: () => import('@/components/authenticated/shared-calendar'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/shopping-list',
      name: 'shopping-list',
      component: () => import('@/components/authenticated/shopping-list'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/components/public/login')
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/components/public/forgot-password')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/components/admin/home'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/features',
      name: 'admin-configure-features',
      component: () => import('@/components/admin/configure-features'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/members',
      name: 'admin-manage-members',
      component: () => import(/* webpackChunkName: "admin-members" */ '@/components/admin/manage-members'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/members/u',
      name: 'admin-manage-member',
      component: () => import(/* webpackChunkName: "admin-members" */ '@/components/admin/manage-member'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/tasks',
      name: 'admin-manage-tasks',
      component: () => import(/* webpackChunkName: "admin-tasks" */ '@/components/admin/manage-tasks'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/tasks/t',
      name: 'admin-manage-task-edit',
      component: () => import(/* webpackChunkName: "admin-tasks" */ '@/components/admin/manage-task'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/entries',
      name: 'admin-manage-entries',
      component: () => import('@/components/admin/manage-entries'),
      beforeEnter: authenticatedRoute
    },
    {
      path: '/setup',
      name: 'setup',
      component: () => import('@/components/public/setup'),
      beforeEnter: unauthenticatedRouteUninitialised
    }
  ]
})
