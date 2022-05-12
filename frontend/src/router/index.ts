import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import User from "../views/User.vue"
import Dashboard from "../views/Dashboard.vue"

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home", // 与 i18n/x.json 中 nav.xxx 中自动对应
    component: Home,
    meta: {
      sort: 1,
      icon: 'mdi-clock'
    },
    children: [

    ]
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: Dashboard,
    meta: {
      sort: 2,
      icon: 'mdi-view-dashboard'
    }
  },
  {
    path: "/user/:id",
    name: "user",
    meta: {
      sort: 3,
      icon: 'mdi-account'
    },
    component: User,
  },
  {
    path: "/about",
    name: "about",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: function () {
      return import(/* webpackChunkName: "about" */ "../views/About.vue");
    },
    meta: {
      sort: 4,
      icon: 'mdi-flag'
    }
  },
];

const router = createRouter({
  // vue-router 提供了 history、hash 模式
  history: createWebHashHistory(),  // hash 模式  
  routes,
});

export default router;
