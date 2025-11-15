import common from "@/common/common";
import cani from "@/requests/authenticated/can-i";
import login from "@/requests/public/login";

// requireAuthToken
// given an no auth token redirect to the login page
async function requireAuthToken(to, from, next) {
  login
    .GetUserAuth()
    .then((resp) => {
      if (resp.data.data !== true) {
        next({ name: "Login", query: { redirect: to.fullPath } });
      } else {
        next();
      }
    })
    .catch((err) => {
      next({ name: "Login", query: { redirect: to.fullPath } });
    });
}

// requireNoAuthToken
// given an auth token, redirect to the home page
async function requireNoAuthToken(to, from, next) {
  login
    .GetUserAuth(false)
    .then((resp) => {
      console.log({ resp });
      if (resp.data.data !== true) {
        next();
      } else {
        next({ name: "Home" });
      }
    })
    .catch((err) => {
      console.log({ err });
      next();
    });
}

async function requireGroup(to, from, next) {
  cani
    .GetCanIgroup(to.meta.requiresGroup)
    .then((resp) => {
      if (resp.data.data === true) {
        next();
      } else {
        next(from.path);
      }
    })
    .catch(() => {
      next(from.path);
    });
}

function isPublicRoute(to) {
  return to.meta.requiresAuth !== true;
}

export default {
  requireGroup,
  isPublicRoute,
  requireAuthToken,
  requireNoAuthToken,
};
