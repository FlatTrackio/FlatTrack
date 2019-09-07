import Vue from 'vue'
import Router from 'vue-router'
import home from '@/components/home'
import tasks from '@/components/tasks'
import task from '@/components/task'
import aboutFlatTrack from '@/components/about-flattrack'
import aboutFlat from '@/components/about-flat'

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
    }
  ]
})
