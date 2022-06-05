<template>
  <v-layout>
    <v-app :theme="store.appTheme">
      <!-- 底部 App-bar -->
      <v-app-bar class="app-bar" color="#1c7bc9" position="bottom"
        :height="store.appBarHeight" :data-wails-drag="true">
        <template v-slot:prepend>
          <v-app-bar-nav-icon variant="text" @click.stop="onToggleMenu = !onToggleMenu">
          </v-app-bar-nav-icon>
        </template>
        <v-icon :icon="activeMenuIcon"> </v-icon>
        <v-app-bar-title class="ml-5">{{ store.appStatusBar }}</v-app-bar-title>
        <v-spacer></v-spacer>
        <template v-slot:append>
          <v-btn variant="text" icon="mdi-translate" @click="onclickToggleLanguage">
          </v-btn>
          <v-btn variant="text" icon="mdi-invert-colors" @click="onclickToggleTheme">
          </v-btn>
          <v-btn variant="text" icon="mdi-magnify" @click="onclickMagnify">
          </v-btn>
          <v-btn variant="text" icon="mdi-github" @click="onclickOpenGithub"></v-btn>
          <v-btn variant="text" icon="mdi-dots-vertical"
            @click.stop="toggleMore = !toggleMore" :disabled="disableToggleMore">
          </v-btn>
        </template>
      </v-app-bar>

      <!-- 左侧 APP 导航栏 -->
      <v-navigation-drawer v-model="onToggleMenu" temporar rail expand-on-hover>
        <v-list>
          <v-list-subheader class="h3">{{ t('nav.mainmenu') }}</v-list-subheader>
          <v-list-item v-for="(menu, i) in router.getRoutes().sort((a, b) => {
            return (a ? a.meta.sort as number : 0) - (b ? b.meta.sort as number : 0);
          })" :key="i" :value="menu" active-color="primary" density="comfortable"
            :to="menu.path" @click="onclickMenuListItem(menu.meta.icon)">
            <v-list-item-avatar start>
              <v-icon
                :icon="menu.meta.icon ? menu.meta.icon as string : 'mid-arrow-all'">
              </v-icon>
            </v-list-item-avatar>
            <v-list-item-title>
              {{ t(menu.name ? `nav.${menu.name as string}` : 'nav.home') }}
            </v-list-item-title>
          </v-list-item>
          <v-spacer></v-spacer>
        </v-list>

        <template v-slot:append>
          <div class="pa-2">
            <v-btn block> Logout </v-btn>
          </div>
        </template>
      </v-navigation-drawer>

      <!-- 右侧导航栏 -->
      <v-navigation-drawer width="160" v-model="toggleMore" position="right">
        <template v-slot:prepend> </template>
        <app-logo logoheight="80px" />
        <template v-slot:append> </template>
      </v-navigation-drawer>

      <!-- 主窗口 -->
      <v-main>
        <router-view />
      </v-main>
    </v-app>
  </v-layout>
</template>

<script lang="ts" setup>
// about vue
import { onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useI18n } from 'vue-i18n'
// about wails
import { WindowMinimise, Quit } from '../wailsjs/runtime'
import { OpenGithub } from '../wailsjs/go/imes/Api'
import { SysInfo } from '../wailsjs/go/main/App'
// about app
import { useBaseStore } from './stores/index'
import AppLogo from './components/AppLogo.vue'

const router = useRouter() // router 是管理器，可以 addRoute、removeRoute、getRoutes、push...
const route = useRoute() // route 是一个响应式对象，
const store = useBaseStore()
const { t, availableLocales, locale } = useI18n({ useScope: 'global' })

//********** 底部 App-bar 相关 **********/
const toggleMore = ref(false)
// change theme
const onclickToggleTheme = () => {
  store.appTheme = store.appTheme === 'light' ? 'dark' : 'light'
  console.log(store.appTheme)
}
// change i18n
const languages = availableLocales
const onclickLanguageHandle = (item: string) => {
  // 所有语言列表切换
  item !== locale.value ? (locale.value = item) : false
}
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
const onToggleMenu = ref(true)
const disableToggleMore = ref(true)
const activeMenuIcon = ref('mdi-home')
const onclickMenuListItem = (val: string | unknown) => {
  if (val) {
    // console.log(val);
    activeMenuIcon.value = val as string
  }
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
const onclickMinimise = () => {
  WindowMinimise()
}
// close app
const onclickQuit = () => {
  Quit()
}

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
      activeMenuIcon.value = val.meta.icon as string
    }
  })

  store.availableHeight = display.height.value - store.appBarHeight

  // store.initConfig()
  // 加载已保存的数据
  store.syncTestProductions()
  store.syncTestStages()
  store.syncTestStation()
  store.syncTestEntity()

  SysInfo().then(
    (info) => {
      if (info) {
        store.sysInfo = info
      }
    }
  )
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
  store.availableHeight = display.height.value - store.appBarHeight
  console.log('store.availableHeight changed: ', store.availableHeight)
  store.appStatusBar.width = display.width.value
  breakpoint()
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

.app-bar {
  align-items: center;
  bottom: 0;
  justify-content: center;
  opacity: 0.9;
  position: absolute;
  width: 100%;
  z-index: 1;
}
</style>
