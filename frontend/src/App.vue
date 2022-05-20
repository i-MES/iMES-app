<template>
  <v-layout>
    <v-app :theme="store.appTheme">
      <v-app-bar color="#1c7bc9" position="bottom" :height="store.appBarHeight"
        data-wails-drag>
        <v-app-bar-nav-icon variant="text" @click.stop="onToggleMenu = !onToggleMenu">
        </v-app-bar-nav-icon>
        <v-icon :icon="activeMenuIcon"> </v-icon>
        <a class="ml-5">{{ store.appStatusBar }}</a>
        <v-spacer></v-spacer>
        <v-btn variant="text" icon="mdi-translate" @click="onclickToggleLanguage">
        </v-btn>
        <v-btn variant="text" icon="mdi-invert-colors" @click="onclickToggleTheme">
        </v-btn>
        <v-btn variant="text" icon="mdi-magnify"></v-btn>
        <v-btn variant="text" icon="mdi-view-module" @click="onclickViewModule"></v-btn>
        <v-btn variant="text" icon="mdi-github" @click="onclickOpenGithub"></v-btn>
        <v-btn variant="text" icon="mdi-dots-vertical"
          @click.stop="onToggleStep = !onToggleStep" :disabled="disableToggleStep">
        </v-btn>
      </v-app-bar>

      <v-navigation-drawer v-model="onToggleMenu" temporary>
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
              v-text="t(menu.name ? `nav.${menu.name as string}` : 'home')">
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-navigation-drawer>

      <v-navigation-drawer v-model="onToggleStep" expand-on-hover rail position="right"
        v-if="!disableToggleStep">
        <v-dialog v-model="dialog" scrollable max-width="300px">
          <template v-slot:activator="{ isActive, props }">
            <v-btn color="primary" dark v-bind="props">
              {{ selectedProd ? store.testProduction(selectedProd).title : '选择产品' }}
            </v-btn>
          </template>
          <v-card>
            <v-card-title>选择待测产品</v-card-title>
            <v-divider></v-divider>
            <v-card-text>
              <v-radio-group v-model="selectedProd" column>
                <v-radio v-for="(tp, idx) in store.testProductions"
                  :label="`${tp.title}(${tp.desc})`" :value="tp.id">
                </v-radio>
              </v-radio-group>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions class="d-inline-flex">
              <v-btn color="blue darken-1" text @click="dialog = false">
                取消
              </v-btn>
              <v-btn color="blue darken-1" text @click="dialog = false">
                确定
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-list density="compact" nav>
          <v-list-item v-for="(step, i) in store.teststeps.sort((a, b) => {
            return a.sequence - b.sequence
          })" :title="step.title" :value="step.id"
            @click="onclickTestStep(step.id, step.title)">
          </v-list-item>
        </v-list>
      </v-navigation-drawer>

      <v-main>
        <router-view />
      </v-main>
    </v-app>
  </v-layout>
</template>

<script lang="ts" setup>
// vue about
import { onMounted, ref, watch, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useI18n } from 'vue-i18n'
// wails about
import { WindowMinimise, Quit } from '../wailsjs/runtime'
import * as imes from '../wailsjs/go/imes/Middleware'
// app about
import { useBaseStore } from './stores/index'

const router = useRouter() // router 是管理器，可以 addRoute、removeRoute、getRoutes、push...
const route = useRoute() // route 是一个响应式对象，
const store = useBaseStore()
const { t, availableLocales, locale } = useI18n({ useScope: 'global' })

const dialog = ref(false)
const selectedProd = ref(0)
watch(
  () => selectedProd.value,
  (nv) => {
    // 用户重新选择了产品
    store.appStatusBar.Production = store.testProduction(nv).title
    // 加载对应产品的工序
    store.loadSteps()
  }
)

const onclickTestStep = (id: number, title: string) => {
  store.activeTestStepId = id
  store.appStatusBar.TestStep = title

}

// 导航栏 menu 菜单相关
const onToggleMenu = ref(true)
const onToggleStep = ref(false)
const disableToggleStep = ref(true)
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
      disableToggleStep.value = false
    } else {
      disableToggleStep.value = true
    }
  }
)

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
const onclickViewModule = () => {
  store.addCounter()
}
const onclickOpenGithub = () => {
  imes.OpenGithub()
}
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

  store.loadTestProductions()
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
