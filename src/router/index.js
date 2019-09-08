import Vue from 'vue'
import Router from 'vue-router'
import home from '@/components/home'
import unknownPage from '@/components/unknown-page'
import tasks from '@/components/tasks'
import task from '@/components/task'
import aboutFlatTrack from '@/components/about-flattrack'
import aboutFlat from '@/components/about-flat'
import highFives from '@/components/high-fives'
import members from '@/components/members'
import noticeboard from '@/components/noticeboard'
import recipes from '@/components/recipes'
import sharedCalendar from '@/components/shared-calendar'
import shoppingList from '@/components/shopping-list'
import adminHome from '@/components/admin-home'
import adminConfigureFeatures from '@/components/admin-configure-features'
import adminManageMembers from '@/components/admin-manage-members'
import adminManageMember from '@/components/admin-manage-member'

Vue.use(Router)

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
      path: '/aboutflattrack',
      name: 'aboutflattrack',
      component: aboutFlatTrack
    },
    {
      path: '/aboutflat',
      name: 'aboutflat',
      component: aboutFlat
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
