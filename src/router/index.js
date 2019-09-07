import Vue from 'vue'
import Router from 'vue-router'
import home from '@/components/home'
import tasks from '@/components/tasks'
import task from '@/components/task'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: home
    },
    {
      path: '/',
      name: 'tasks',
      component: tasks
    },
    {
      path: '/task',
      name: 'task',
      component: task
    }
  ]
})
