import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import Login from "./views/Login.vue";
import Logout from "./views/Logout.vue";
import store from "./store";

Vue.use(Router);

const router = new Router({
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/",
      name: "home",
      component: Home
    },
    {
      path: "/login",
      name: "login",
      component: Login
    },
    {
      path: "/logout",
      name: "logout",
      component: Logout
    }
  ]
});

router.beforeEach((to, from, next) => {
  if (to.name === "logout") {
    next();
    return;
  }

  if (store.state.user) {
    if (to.name === "login") {
      next("/");
      return;
    }

    next();
    return;
  }

  if (to.name === "login") {
    next();
    return;
  }

  next({
    path: "/login",
    query: {
      nextURL: to.fullPath
    }
  });
});

export default router;
