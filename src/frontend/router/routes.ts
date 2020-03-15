export default [
  {
    path: '/',
    name: 'home',
    component: () => import('@/frontend/views/authenticated/home.vue')
  },
  {
    path: '*',
    name: 'unknown-page',
    component: () => import('@/frontend/views/global/unknown-page.vue')
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: () => import(/* webpackChunkName: "tasks" */ '@/frontend/views/authenticated/tasks.vue')
  },
  {
    path: '/tasks/t',
    name: 'task-view',
    component: () => import(/* webpackChunkName: "tasks" */ '@/frontend/views/authenticated/task.vue')
  },
  {
    path: '/about/flattrack',
    name: 'aboutflattrack',
    component: () => import('@/frontend/views/authenticated/about-flattrack.vue'),
    alias: '/aboutflattrack'
  },
  {
    path: '/about/flat',
    name: 'aboutflat',
    component: () => import('@/frontend/views/authenticated/about-flat.vue'),
    alias: '/aboutflat'
  },
  {
    path: '/high-fives',
    name: 'highfives',
    component: () => import('@/frontend/views/authenticated/high-fives.vue')
  },
  {
    path: '/members',
    name: 'members',
    component: () => import('@/frontend/views/authenticated/members.vue')
  },
  {
    path: '/noticeboard',
    name: 'noticeboard',
    component: () => import('@/frontend/views/authenticated/noticeboard.vue')
  },
  {
    path: '/noticeboard/p',
    name: 'noticeboard posts',
    component: () => import('@/frontend/views/authenticated/noticeboard-post.vue')
  },
  {
    path: '/recipes',
    name: 'recipes',
    component: () => import('@/frontend/views/authenticated/recipes.vue')
  },
  {
    path: '/shared-calendar',
    name: 'shared-calendar',
    component: () => import('@/frontend/views/authenticated/shared-calendar.vue')
  },
  {
    path: '/shopping-list',
    name: 'shopping-list',
    component: () => import('@/frontend/views/authenticated/shopping-list.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/frontend/views/public/login.vue')
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('@/frontend/views/public/forgot-password.vue')
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/frontend/views/admin/home.vue')
  },
  {
    path: '/admin/features',
    name: 'admin-configure-features',
    component: () => import('@/frontend/views/admin/configure-features.vue')
  },
  {
    path: '/admin/members',
    name: 'admin-manage-members',
    component: () => import(/* webpackChunkName: "admin-members" */ '@/frontend/views/admin/manage-members.vue')
  },
  {
    path: '/admin/members/u',
    name: 'admin-manage-member',
    component: () => import(/* webpackChunkName: "admin-members" */ '@/frontend/views/admin/manage-member.vue')
  },
  {
    path: '/admin/tasks',
    name: 'admin-manage-tasks',
    component: () => import(/* webpackChunkName: "admin-tasks" */ '@/frontend/views/admin/manage-tasks.vue')
  },
  {
    path: '/admin/tasks/t',
    name: 'admin-manage-task-edit',
    component: () => import(/* webpackChunkName: "admin-tasks" */ '@/frontend/views/admin/manage-task.vue')
  },
  {
    path: '/admin/entries',
    name: 'admin-manage-entries',
    component: () => import('@/frontend/views/admin/manage-entries.vue')
  },
  {
    path: '/setup',
    name: 'setup',
    component: () => import('@/frontend/views/public/setup.vue')
  }
]
