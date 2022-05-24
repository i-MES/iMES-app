<template>
  <v-layout>
    <v-app :theme="store.appTheme">
      <!-- 底部 App-bar -->
      <v-app-bar color="#1c7bc9" position="bottom" :height="store.appBarHeight"
        data-wails-drag>
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
          <v-btn variant="text" icon="mdi-view-module" @click="onclickViewModule">
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
            <v-list-item-title
              v-text="t(menu.name ? `nav.${menu.name as string}` : 'nav.home')">
            </v-list-item-title>
          </v-list-item>
        </v-list>

        <template v-slot:append>
          <div class="pa-2">
            <v-btn block> Logout </v-btn>
          </div>
        </template>
      </v-navigation-drawer>

      <!-- 右侧产品导航栏 -->
      <v-navigation-drawer width="160" v-model="toggleMore" position="right">
        <template v-slot:prepend>
          <div class="pa-2">
            <v-btn block v-if="planab === 'a'" @click="onclickLoadConfig">{{
                t("nav.loadconfig")
            }}</v-btn>
            <v-file-input v-if="planab === 'b'" label="File input" outlined dense>
            </v-file-input>
          </div>
        </template>
        <v-select filled label="产品" dense hide-details v-model="selectedProd"
          :items="store.testProductions.map((v, _) => v.id + '-' + v.title)">
        </v-select>
        <v-select filled label="机型" dense hide-details
          :items="store.teststeps.map((v, _) => v.id + v.title)">
        </v-select>
        <v-select filled label="工序" dense hide-details :items="['组装', '高温', '写版本']">
        </v-select>
        <template v-slot:append>
          <div class="px-2 ">
            <v-switch v-model="store.testPageViewModel" color="success"
              density="compact" @change="onchangeViewModel"
              :label="`${t('testpage.view-model')}: ${store.testPageViewModel ? t('testpage.view-model-single') : t('testpage.view-model-all')}`">
            </v-switch>
          </div>
        </template>
      </v-navigation-drawer>
      <v-main>
        <router-view />
      </v-main>
    </v-app>
  </v-layout>
</template>

<script lang="ts" setup>
// about vue
import { onMounted, ref, watch, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useI18n } from 'vue-i18n'
// about wails
import { WindowMinimise, Quit } from '../wailsjs/runtime'
import { OpenGithub } from '../wailsjs/go/imes/Api'
// about app
import { useBaseStore } from './stores/index'

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
}
const onclickViewModule = () => {
  store.addCounter()
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

//********** 右侧产品导航栏相关 **********/
const selectedProd = ref()
const planab = ref('a')
const onclickLoadConfig = () => {
  if (planab.value == 'a') {
    // 方案1：后端加载
    store.loadTestProductions()
  } else {
    // 方案2：前端加载
  }
}
watch(
  () => selectedProd.value,
  (nv) => {
    var _id: number = nv.split('-')[0] as number
    var _title: string = nv.split('-')[1]
    console.log('selectedProd changed: ', _id, _title)
    // 用户重新选择了产品
    var _tp = store.testProductionById(_id)
    if (_tp) {
      store.appStatusBar.Prod = _tp.title
      // 加载对应产品的工序
      store.loadSteps()
    }
  }
)
const onchangeViewModel = () => {
  console.log("view model value changed:", viewmodel.value ? '单机' : '全体')
}

//********** 其他 **********/
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
  const display = useDisplay()
  console.log("vuetify's display - height: ", display.height.value)
  console.log("vuetify's display - width: ", display.width.value)
  console.log("vuetify's display - mobile: ", display.mobile.value)
  console.log("vuetify's display - platform: ", display.platform.value)

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
})
</script>

<style lang="scss">
@import url('./assets/css/reset.css');
@import url('./assets/css/font.css');

#app {
  position: relative;
  height: 100%;
  overflow: hidden;
}
</style>
