import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import AppHome from '../views/AppHome.vue'
import UserManage from '../views/UserManage.vue'
import TestPage from '../views/TestPage.vue'
import AppSettings from '../views/AppSettings.vue'
import DashBoard from '../views/DashBoard.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home', // 与 i18n/x.json 中 nav.xxx 中自动对应
    component: AppHome,
    meta: {
      sort: 1,
      icon: 'mdi-home'
    }
  },
  {
    path: '/test',
    name: 'test',
    component: TestPage,
    meta: {
      sort: 2,
      icon: 'mdi-arrow-right-bold-circle-outline'
    }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashBoard,
    meta: {
      sort: 3,
      icon: 'mdi-view-dashboard'
    }
  },
  {
    path: '/user/:id',
    name: 'user',
    component: UserManage,
    meta: {
      sort: 4,
      icon: 'mdi-account'
    },
  },
  {
    path: '/settings',
    name: 'settings',
    component: AppSettings,
    meta: {
      sort: 5,
      icon: 'mdi-cog-outline'
    },
  },
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: function () {
      return import(/* webpackChunkName: 'about' */ '../views/AppAbout.vue')
    },
    meta: {
      sort: 6,
      icon: 'mdi-flag'
    }
  },
]

const router = createRouter({
  // vue-router 提供了 history、hash 模式
  history: createWebHashHistory(),  // hash 模式  
  routes,
})

export default router
