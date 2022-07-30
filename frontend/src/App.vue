<template>
  <v-layout>
    <v-app :theme="store.appTheme">
      <!-- 底部 App-bar -->
      <v-app-bar class="app-bar" color="#A7535A" absolute location="bottom"
        :height="store.appBarHeight" :data-wails-drag="true">
        <template v-slot:prepend>
          <!-- <v-app-bar-nav-icon variant="text" @click.stop="onToggleMenu = !onToggleMenu">
          </v-app-bar-nav-icon> -->
          <img class="ml-n3 my-0 pa-0"
            :style="`height:` + (store.appBarHeight + 2) + `px; background: lightblue; `"
            src="@/assets/images/logo.svg" />
        </template>
        <v-app-bar-title class="ml-5">{{ store.appStatusBar }}</v-app-bar-title>
        <v-spacer></v-spacer>
        <template v-slot:append>
          <v-btn variant="text"
            :icon="locale == `en` ? `mdi-translate-off` : `mdi-translate`"
            @click="onclickToggleLanguage">
          </v-btn>
          <v-btn variant="text" icon="mdi-invert-colors" @click="onclickToggleTheme">
          </v-btn>
          <v-btn variant="text" icon="mdi-magnify" @click="onclickMagnify">
          </v-btn>
          <v-btn variant="text" icon="mdi-github" @click="onclickOpenGithub"></v-btn>
        </template>
      </v-app-bar>

      <!-- 左侧 APP 导航栏 -->
      <v-navigation-drawer :v-model="true" rail rail-width="55" permanent>
        <v-list :selected="listSelected" nav>
          <!-- <v-list-subheader class="ma-0 pa-0"> </v-list-subheader> -->
          <v-list-item v-for="(menu, i) in router.getRoutes().filter((v) => { return v.meta.location == `top` }).sort((a, b) => {
            return (a ? a.meta.sort as number : 0) - (b ? b.meta.sort as number : 0);
          })" :key="i" :value="menu" active-color="#A7535A" density="comfortable"
            :to="menu.path" @click="onclickMenuListItem(menu)">
            <v-tooltip location="end">
              <template v-slot:activator="{ props }">
                <!-- v-list-item-avastar 无法调整 icon 尺寸 -->
                <v-icon v-bind="props" :icon="menu.meta.icon as string" size="large">
                </v-icon>
              </template>
              <span>{{ t(`nav.${menu.name as string}`) }}</span>
            </v-tooltip>
          </v-list-item>
        </v-list>

        <template v-slot:append>
          <v-list :selected="listSelected" nav>
            <v-list-item v-for="(menu, i) in router.getRoutes().filter((v) => { return v.meta.location == `bottom` }).sort((a, b) => {
              return (a ? a.meta.sort as number : 0) - (b ? b.meta.sort as number : 0);
            })" :key="i" :value="menu" active-color="#A7535A" density="comfortable"
              :to="menu.path" @click="onclickMenuListItem(menu)">
              <v-tooltip location="end">
                <template v-slot:activator="{ props }">
                  <v-icon v-bind="props" :icon="menu.meta.icon as string" size="large">
                  </v-icon>
                </template>
                <span>{{ t(`nav.${menu.name as string}`) }}</span>
              </v-tooltip>
            </v-list-item>
          </v-list>
        </template>
      </v-navigation-drawer>

      <!-- 主窗口 -->
      <v-main>
        <v-sheet class="ma-0 pa-0 overflow-y-auto " :height="store.mainWindowHeight"
          :color="store.appTheme == 'dark' ? '#101010' : 'grey-lighten-4'">
          <router-view />
        </v-sheet>
      </v-main>
    </v-app>
  </v-layout>
</template>

<script lang="ts" setup>
// about vue
import { onMounted, ref, watch } from 'vue'
import { useRouter, useRoute, RouteRecordRaw } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useI18n } from 'vue-i18n'
// about wails
import { OpenGithub, SetIntSetting } from '../wailsjs/go/imes/Api'
import { SysInfo } from '../wailsjs/go/main/App'
// about app
import { useBaseStore } from './stores/index'
import * as runtime from '../wailsjs/runtime/runtime'
import { target } from '../wailsjs/go/models'

const router = useRouter() // router 是管理器，可以 addRoute、removeRoute、getRoutes、push...
const route = useRoute() // route 是一个响应式对象，
const store = useBaseStore()
const { t, locale } = useI18n({ useScope: 'global' })

//********** 底部 App-bar 相关 **********/
// change theme
const onclickToggleTheme = () => {
  store.appTheme = store.appTheme === 'light' ? 'dark' : 'light'
  console.log(store.appTheme)
}
// change i18n
// const languages = availableLocales
// const onclickLanguageHandle = (item: string) => {
//   // 所有语言列表切换
//   item !== locale.value ? (locale.value = item) : false
// }
const onclickToggleLanguage = () => {
  // 中英文切换
  locale.value ? (locale.value == 'en' ? (locale.value = 'zh-Hans') : (locale.value = 'en')) : false
}
const onclickMagnify = () => {
  console.log('onclickMagnify')
}
const onclickOpenGithub = () => {
  OpenGithub()
}

//********** 左侧 APP 导航栏相关 **********/
const disableToggleMore = ref(true)
const listSelected = ref([])
const onclickMenuListItem = (menu: RouteRecordRaw) => {
  console.log('-0-0-', listSelected.value, menu)
}
watch(
  () => route.path,
  (newPath) => {
    console.log('watching route.path:', newPath)
    if (newPath.indexOf('test') > 0) {
      disableToggleMore.value = false
      console.log('disableToggleMore watch:', newPath, 'false')
    } else {
      disableToggleMore.value = true
      console.log('disableToggleMore watch:', newPath, 'true')
    }
  }
)

//********** 右侧导航栏相关 **********/

//********** 其他 **********/
const display = useDisplay()
// hide window
// const onclickMinimise = () => {
//   WindowMinimise()
// }
// // close app
// const onclickQuit = () => {
//   Quit()
// }

onMounted(() => {
  // practise vuetify's display props
  console.log('vuetify\'s display - height: ', display.height.value)
  console.log('vuetify\'s display - width: ', display.width.value)
  console.log('vuetify\'s display - mobile: ', display.mobile.value)
  console.log('vuetify\'s display - platform: ', display.platform.value)

  // 默认导航的页面
  const _dr = store.defaultRoute
  router.push({
    name: _dr,
  })
  // console.log('===', router.getRoutes())
  router.getRoutes().forEach((val) => {
    if (val.name === _dr) {
      // activeMenuIcon.value = val.meta.icon as string
    }
  })

  store.mainWindowHeight = display.height.value - store.appBarHeight

  // store.initConfig()
  // 加载已保存的数据
  store.syncTestProductions()
  store.syncTestStages()
  store.syncTestStation()
  store.syncTestEntity()
  store.LoadStringSetting('darkmaincolor')
  store.LoadStringSetting('lightmaincolor')
  store.LoadNumberSetting('paneFirstLengthPercent')

  SysInfo().then(
    (info) => {
      if (info) {
        store.sysInfo = info
      }
    }
  )

  // 注册状态响应函数
  runtime.EventsOn('testitemstatus', (tis: target.TestItemStatus) => {
    var _tises = store.LastestTIStatus[tis.testentityid]
    if (_tises) {
      for (let idx = 0; idx < _tises.length; idx++) {
        if (_tises[idx].testgroupid == tis.testgroupid && _tises[idx].testclassid == tis.testclassid && _tises[idx].testitemid == tis.testitemid) {
          store.LastestTIStatus[tis.testentityid][idx] = tis
          console.log('更新已有 tistatus', tis.testitemid, tis.status)
          return
        }
      }
      store.LastestTIStatus[tis.testentityid].push(tis)
      console.log('添加新的 tistatus', tis.testitemid, tis.status)
    } else {
      store.LastestTIStatus[tis.testentityid] = [tis]
      console.log('创建 entity 对应的 tistatus 数组', tis.testitemid, tis.status)
    }
  })
})

const breakpoint = () => {
  if (store.sysInfo.buildtype) {
    console.log('buildtype: ', store.sysInfo.buildtype)
    if (store.sysInfo.buildtype == 'dev') {
      var w: number = display.width.value
      if (w < 600) {
        store.appStatusBar.width = 'xs-' + w
      } else if (w < 960) {
        store.appStatusBar.width = 'sm-' + w
      } else if (w < 1264) {
        store.appStatusBar.width = 'md-' + w
      } else if (w < 1904) {
        store.appStatusBar.width = 'lg-' + w
      } else {
        store.appStatusBar.width = 'xl-' + w
      }
    }
  }
}
breakpoint()
window.onresize = () => {
  store.mainWindowHeight = display.height.value - store.appBarHeight
  console.log('store.mainWindowHeight changed: ', store.mainWindowHeight)
  store.appStatusBar.width = display.width.value
  breakpoint()
  SetIntSetting('display-height', display.height.value)
  SetIntSetting('display-width', display.width.value)
}

</script>

<style lang="scss">
@import url('./assets/css/reset.css');
@import url('./assets/css/font.css');
@import url('./app.css');

#app {
  position: relative;
  height: 100%;
  overflow: hidden;
}
</style>
