import Vue from 'vue'
import VueRouter from 'vue-router'
import routes from './routes'
import routerCommon from './common'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  }
})

router.beforeEach((to, from, next) => {
  if (typeof to.name !== 'undefined') {
    document.title = `FlatTrack | ${to.name}`
  }
  if (to.matched.some(route => route.meta.requiresAuth)) {
    routerCommon.requireAuthToken(to, from, next)
  }
  if (to.matched.some(route => route.meta.requiresNoAuth === true)) {
    routerCommon.requireNoAuthToken(to, from, next)
  }
  if (to.matched.some(route => route.meta.requiresGroup)) {
    console.log({ to })
    routerCommon.requireGroup(to, from, next)
  }
  next()
})

export default router
