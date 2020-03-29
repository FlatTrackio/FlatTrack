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
    path: '/flatmates',
    name: 'flatmates',
    component: () => import('@/frontend/views/authenticated/flatmates.vue')
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
  // {
  //   path: '/admin',
  //   name: 'admin',
  //   component: () => import('@/frontend/views/admin/home.vue')
  // },
  {
    path: '/setup',
    name: 'setup',
    component: () => import('@/frontend/views/public/setup.vue')
  }
]
