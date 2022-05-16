import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import User from "../views/User.vue"
import Test from "../views/Test.vue"
import Settings from "../views/Settings.vue"
import Dashboard from "../views/Dashboard.vue"

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home", // 与 i18n/x.json 中 nav.xxx 中自动对应
    component: Home,
    meta: {
      sort: 1,
      icon: 'mdi-home'
    },
    children: [

    ]
  },
  {
    path: "/test",
    name: "test",
    component: Test,
    meta: {
      sort: 2,
      icon: 'mdi-arrow-right-bold-circle-outline'
    }
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: Dashboard,
    meta: {
      sort: 3,
      icon: 'mdi-view-dashboard'
    }
  },
  {
    path: "/user/:id",
    name: "user",
    component: User,
    meta: {
      sort: 4,
      icon: 'mdi-account'
    },
  },
  {
    path: "/settings",
    name: "settings",
    component: Settings,
    meta: {
      sort: 5,
      icon: 'mdi-cog-outline'
    },
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
      sort: 6,
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
