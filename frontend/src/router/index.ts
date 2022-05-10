import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import User from "../views/User.vue"
import Login from "../views/Login.vue"
import Dashboard from "../views/Dashboard.vue"
import HelloVuetify from "../views/HelloVuetify.vue"

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
    children: [

    ]
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: function () {
      return import(/* webpackChunkName: "about" */ "../views/About.vue");
    },
  },
  {
    path: "/user/:id",
    name: "user",
    component: User,
  },
  {
    path: "/login",
    name: "login",
    component: Login,
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: Dashboard,
  },
  {
    path: "/hellovuetify",
    name: "hellovuetify",
    component: HelloVuetify,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
