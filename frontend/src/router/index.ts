import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import AppHome from '../views/AppHome.vue'
import UserManage from '../views/UserManage.vue'
import TestPage from '../views/TestPage.vue'
import AppSettings from '../views/AppSettings.vue'
import DashBoard from '../views/DashBoard.vue'
import FooBar from '../views/FooBar.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home', // 与 i18n/x.json 中 nav.xxx 中自动对应
    component: AppHome,
    meta: {
      location: 'top',
      sort: 1,
      icon: 'mdi-home'
    }
  },
  {
    path: '/test',
    name: 'test',
    component: TestPage,
    meta: {
      location: 'top',
      sort: 2,
      icon: 'mdi-arrow-right-bold-circle-outline'
    }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashBoard,
    meta: {
      location: 'top',
      sort: 3,
      icon: 'mdi-finance'
    }
  },
  // {
  //   path: '/plguin',
  //   name: 'plugin',
  //   component: FooBar,
  //   meta: {
  //     location: 'top',
  //     sort: 4,
  //     icon: 'mdi-checkbox-blank-badge'
  //   }
  // },
  {
    path: '/user/:id',
    name: 'user',
    component: UserManage,
    meta: {
      location: 'bottom',
      sort: 1,
      icon: 'mdi-account'
    },
  },
  {
    path: '/settings',
    name: 'settings',
    component: AppSettings,
    meta: {
      location: 'bottom',
      sort: 2,
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
      location: 'bottom',
      sort: 3,
      icon: 'mdi-information-outline'
    }
  },
  // {
  //   path: '/foobar',
  //   name: 'foobar',
  //   component: FooBar,
  //   meta: {
  //     sort: 7,
  //     icon: 'mdi-test-tube'
  //   },
  // },
]

const router = createRouter({
  // vue-router 提供了 history、hash 模式
  history: createWebHashHistory(),  // hash 模式  
  routes,
})

export default router
