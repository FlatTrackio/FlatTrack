import common from "@/common/common";
import login from "@/requests/public/login";
import healthz from "@/requests/public/healthz";
import registration from "@/requests/public/registration";

export default [
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/authenticated/home.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/:catchAll(.*)",
    name: "Unknown Page",
    component: () => import("@/views/global/unknown-page.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/flat",
    name: "My Flat",
    component: () => import("@/views/authenticated/flat.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/about-flattrack",
    name: "About FlatTrack",
    component: () => import("@/views/authenticated/about-flattrack.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/account",
    name: "Account",
    component: () => import("@/views/authenticated/account-home.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/account/profile",
    name: "Account Profile",
    component: () => import("@/views/authenticated/account-profile.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/account/security",
    name: "Account Security",
    component: () => import("@/views/authenticated/account-security.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/account/settings",
    name: "Account Settings",
    component: () => import("@/views/authenticated/settings.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps",
    name: "Apps",
    component: () => import("@/views/authenticated/apps.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/flatmates",
    name: "My Flatmates",
    component: () => import("@/views/authenticated/flatmates.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/shopping-list",
    name: "Shopping list",
    component: () => import("@/views/authenticated/shopping-list.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/shopping-list/new",
    name: "New shopping list",
    component: () => import("@/views/authenticated/shopping-list-new.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/shopping-list/list",
    redirect: {
      name: "Shopping list",
    },
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/shopping-list/list/:id",
    name: "View shopping list",
    component: () => import("@/views/authenticated/shopping-list-view.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/apps/shopping-list/tags",
    name: "Manage shopping tags",
    component: () => import("@/views/authenticated/shopping-list-tags.vue"),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/admin",
    name: "Admin home",
    component: () => import("@/views/admin/home.vue"),
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/admin/settings",
    name: "Admin settings",
    component: () => import("@/views/admin/settings.vue"),
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/admin/accounts",
    name: "Admin accounts",
    component: () => import("@/views/admin/accounts.vue"),
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/admin/accounts/new",
    name: "Admin new account",
    component: () => import("@/views/admin/accounts-new.vue"),
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/admin/accounts/edit",
    redirect: {
      name: "Admin accounts",
    },
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/admin/accounts/edit/:id",
    name: "View user account",
    component: () => import("@/views/admin/account-edit.vue"),
    meta: {
      requiresAuth: true,
      requiresGroup: "admin",
    },
  },
  {
    path: "/useraccountconfirm/:id",
    name: "User account confirm",
    component: () => import("@/views/public/useraccountconfirm.vue"),
    meta: {
      requiresNoAuth: true,
    },
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/public/login.vue"),
    beforeEnter: (to, from, next) => {
      let instanceRegistered;
      registration
        .GetInstanceRegistered()
        .then((resp) => {
          instanceRegistered = resp.data.data === true;
          return login.GetUserAuth(false);
        })
        .then((resp) => {
          next({ name: "Home" });
        })
        .catch((err) => {
          if (instanceRegistered !== true) {
            // TODO fix this route not being served
            next({ name: "Set up" });
            return;
          }
          next();
        });
    },
  },
  {
    path: "/forgot-password",
    name: "Forgot password",
    component: () => import("@/views/public/forgot-password.vue"),
    meta: {
      requiresNoAuth: true,
    },
  },
  {
    path: "/setup",
    name: "Set up",
    component: () => import("@/views/public/setup.vue"),
    beforeEnter: (to, from, next) => {
      healthz
        .GetHealthz()
        .then((resp) => {
          return registration.GetInstanceRegistered();
        })
        .then((resp) => {
          if (resp.data.data === true) {
            next({ name: "Home" });
            return;
          }
          next();
        })
        .catch((err) => {
          if (
            err.config.url === "/_healthz" &&
            err.response.data.data === false
          ) {
            next({ name: "Unavailable" });
            return;
          }
          next();
        });
    },
  },
  {
    path: "/unavailable",
    name: "Unavailable",
    component: () => import("@/views/public/unavailable.vue"),
    beforeEnter: (to, from, next) => {
      healthz
        .GetHealthz()
        .then((resp) => {
          if (resp.data.data === true) {
            next("/");
            return;
          }
          next();
        })
        .catch(() => {
          next();
        });
    },
  },
];
