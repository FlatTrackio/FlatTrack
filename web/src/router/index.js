import Vue from 'vue'
import VueRouter from 'vue-router'
import routes from './routes'
import routerCommon from './common'
import common from '../common/common'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  base: import.meta.env.BASE_URL,
  routes,
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  }
})

router.beforeEach(async (to, from, next) => {
  if (typeof to.name !== 'undefined') {
    document.title = `${to.name} | FlatTrack`
  }
  if (to.meta.requiresAuth === true && common.HasAuthToken() === false) {
    routerCommon.requireAuthToken(to, from, next)
  } else if (to.meta.requiresNoAuth === true) {
    routerCommon.requireNoAuthToken(to, from, next)
  } else if (typeof to.requiresGroup !== 'undefined') {
    routerCommon.requireGroup(to, from, next)
  } else {
    next()
  }
})

export default router
