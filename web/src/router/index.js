import { createRouter, createWebHistory } from "vue-router";
import routes from "./routes";
import routerCommon from "./common";
import common from "../common/common";

const router = new createRouter({
  history: createWebHistory("/"),
  base: import.meta.env.BASE_URL,
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { x: 0, y: 0 };
    }
  },
});

router.beforeEach(async (to, from, next) => {
  if (typeof to.name !== "undefined") {
    document.title = `${to.name} | FlatTrack`;
  }
  if (to.meta.requiresAuth === true && common.HasAuthToken() === false) {
    await routerCommon.requireAuthToken(to, from, next);
  } else if (to.meta.requiresNoAuth === true) {
    await routerCommon.requireNoAuthToken(to, from, next);
  } else if (typeof to.meta.requiresGroup !== "undefined") {
    await routerCommon.requireGroup(to, from, next);
  } else {
    next();
  }
});

export default router;
