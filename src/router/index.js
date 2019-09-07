import Vue from 'vue'
import Router from 'vue-router'
import home from '@/components/home'
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

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: home
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: tasks
    },
    {
      path: '/task',
      name: 'task',
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
    }
  ]
})
