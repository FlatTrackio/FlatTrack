import Vue from 'vue'
import Router from 'vue-router'
import home from '@/components/authenticated/home'
import login from '@/components/public/login'
import forgotPassword from '@/components/public/forgot-password'
import unknownPage from '@/components/global/unknown-page'
import tasks from '@/components/authenticated/tasks'
import task from '@/components/authenticated/task'
import aboutFlatTrack from '@/components/authenticated/about-flattrack'
import aboutFlat from '@/components/authenticated/about-flat'
import highFives from '@/components/authenticated/high-fives'
import members from '@/components/authenticated/members'
import noticeboard from '@/components/authenticated/noticeboard'
import recipes from '@/components/authenticated/recipes'
import sharedCalendar from '@/components/authenticated/shared-calendar'
import shoppingList from '@/components/authenticated/shopping-list'
import adminHome from '@/components/admin/home'
import adminConfigureFeatures from '@/components/admin/configure-features'
import adminManageMembers from '@/components/admin/manage-members'
import adminManageMember from '@/components/admin/manage-member'
import adminManageTasks from '@/components/admin/manage-tasks'
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

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: home,
      beforeEnter: authenticatedRoute
    },
    {
      path: '*',
      name: 'unknown-page',
      component: unknownPage
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: tasks,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/tasks/view',
      name: 'task-view',
      component: task,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/about/flattrack',
      name: 'aboutflattrack',
      component: aboutFlatTrack,
      alias: '/aboutflattrack',
      beforeEnter: authenticatedRoute
    },
    {
      path: '/about/flat',
      name: 'aboutflat',
      component: aboutFlat,
      alias: '/aboutflat',
      beforeEnter: authenticatedRoute
    },
    {
      path: '/high-fives',
      name: 'highfives',
      component: highFives,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/members',
      name: 'members',
      component: members,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/tasks/view',
      name: 'task-view',
      component: task,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/noticeboard',
      name: 'noticeboard',
      component: noticeboard,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/recipes',
      name: 'recipes',
      component: recipes,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/shared-calendar',
      name: 'shared-calendar',
      component: sharedCalendar,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/shopping-list',
      name: 'shopping-list',
      component: shoppingList,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/login',
      name: 'login',
      component: login
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: forgotPassword
    },
    {
      path: '/admin',
      name: 'admin',
      component: adminHome,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/features',
      name: 'admin-configure-features',
      component: adminConfigureFeatures,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/members',
      name: 'admin-manage-members',
      component: adminManageMembers,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/members/u',
      name: 'admin-manage-member',
      component: adminManageMember,
      beforeEnter: authenticatedRoute
    },
    {
      path: '/admin/tasks',
      name: 'admin-manage-tasks',
      component: adminManageTasks,
      beforeEnter: authenticatedRoute
    }
  ]
})
