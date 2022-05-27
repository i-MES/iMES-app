<template>
  <v-layout>
    <v-app :theme="store.appTheme">
      <!-- 底部 App-bar -->
      <v-app-bar class="app-bar" color="#1c7bc9" position="bottom"
        :height="store.appBarHeight" data-wails-drag>
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
          <v-spacer></v-spacer>
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
          <div class="pa-1">
            <v-btn :disabled="autoLoadConfig" block v-if="planab === 'a'"
              @click="onclickLoadConfig">{{
                  t("nav.loadconfig")
              }} </v-btn>
            <v-file-input v-if="planab === 'b'" label="File input" outlined dense>
            </v-file-input>
          </div>
        </template>
        <v-select filled label="产品" dense hide-details v-model="selectedProd"
          :items="store.testProductions.map((v, _) => v.id + '-' + v.title)">
        </v-select>
        <v-select filled label="工序" dense hide-details v-model="selectedStage"
          :items="stages">
        </v-select>
        <template v-slot:append>
          <div class="px-2 ">
            <v-switch v-model="autoLoadConfig" color="success" density="compact"
              :label="`自动加载`">
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
import { OpenConfigFolder } from '../wailsjs/go/imes/Api'

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
  store.TEorTI = true
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
const autoLoadConfig = ref(true)
const selectedProd = ref()
const selectedStage = ref()
const planab = ref('a')
const onclickLoadConfig = () => {
  if (planab.value == 'a') {
    // 方案1：后端加载
    OpenConfigFolder()
  } else {
    // 方案2：前端加载
  }
}

if (true) {
  store.initConfig()
}
// 加载已保存的数据
store.syncTestProductions()
store.syncTestStages()
store.syncTestStation()
store.syncTestEntity()
store.syncTestItem()

const stages = reactive([])
watch(
  () => selectedProd.value,
  (nv) => {
    var pid = Number((nv as string).split('-')[0])
    store.activedProductionId = pid
    // 联动 stage 工序选择栏
    var i = stages.length
    while (i--) {
      stages.splice(i, 1)
    }
    store.testStageByProductionId(pid).map((v, _) => (v.id + '-' + v.title)).forEach(
      (n) => stages.push(n)
    )
  }
)
watch(
  () => selectedStage.value,
  (nv) => {
    var sid = Number((nv as string).split('-')[0])
    store.activedTestStageId = sid
  }
)

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

  store.availableHeight = display.height.value - store.appBarHeight
})
window.onresize = () => {
  store.availableHeight = display.height.value - store.appBarHeight
  console.log('store.availableHeight changed: ', store.availableHeight)
}
</script>

<style lang="scss">
@import url('./assets/css/reset.css');
@import url('./assets/css/font.css');

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
