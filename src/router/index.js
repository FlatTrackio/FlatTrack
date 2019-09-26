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
import member from '@/components/authenticated/member'
import noticeboard from '@/components/authenticated/noticeboard'
import recipes from '@/components/authenticated/recipes'
import sharedCalendar from '@/components/authenticated/shared-calendar'
import shoppingList from '@/components/authenticated/shopping-list'
import adminHome from '@/components/admin/home'
import adminConfigureFeatures from '@/components/admin/configure-features'
import adminManageMembers from '@/components/admin/manage-members'
import adminManageMember from '@/components/admin/manage-member'

Vue.use(Router)
Vue.component('home', () => import('../components/home.vue'))

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: home
    },
    {
      path: '*',
      name: 'unknown-page',
      component: unknownPage
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: tasks
    },
    {
      path: '/tasks/view',
      name: 'task-view',
      component: task
    },
    {
      path: '/about/flattrack',
      name: 'aboutflattrack',
      component: aboutFlatTrack,
      alias: '/aboutflattrack'
    },
    {
      path: '/about/flat',
      name: 'aboutflat',
      component: aboutFlat,
      alias: '/aboutflat'
    },
    {
      path: '/high-fives',
      name: 'highfives',
      component: highFives
    },
    {
      path: '/members',
      name: 'members',
      component: members
    },
    {
      path: '/members/u',
      name: 'member',
      component: member
    },
    {
      path: '/tasks/view',
      name: 'task-view',
      component: task
    },
    {
      path: '/noticeboard',
      name: 'noticeboard',
      component: noticeboard
    },
    {
      path: '/recipes',
      name: 'recipes',
      component: recipes
    },
    {
      path: '/shared-calendar',
      name: 'shared-calendar',
      component: sharedCalendar
    },
    {
      path: '/shopping-list',
      name: 'shopping-list',
      component: shoppingList
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
      component: adminHome
    },
    {
      path: '/admin/features',
      name: 'admin-configure-features',
      component: adminConfigureFeatures
    },
    {
      path: '/admin/members',
      name: 'admin-manage-members',
      component: adminManageMembers
    },
    {
      path: '/admin/members/u',
      name: 'admin-manage-member',
      component: adminManageMember
    }
  ]
})
