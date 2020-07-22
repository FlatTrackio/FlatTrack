import Vue from 'vue'
import VueRouter from 'vue-router'
import routes from './routes'
import routerCommon from './common'

Vue.use(VueRouter)

const subpath = document.head.querySelector('[name~=sitesubpath][content]').content || ''

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL + subpath,
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
    document.title = `${to.name} | FlatTrack`
  }
  if (to.matched.some(route => route.meta.requiresAuth === true)) {
    routerCommon.requireAuthToken(to, from, next)
  }
  if (to.matched.some(route => route.meta.requiresNoAuth === true)) {
    routerCommon.requireNoAuthToken(to, from, next)
  }
  if (to.matched.some(route => route.meta.requiresGroup)) {
    routerCommon.requireGroup(to, from, next)
  }
  next()
})

export default router
