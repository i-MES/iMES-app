<template>
  <v-layout>
    <v-app :theme="store.apptheme">
      <v-app-bar color="#1c7bc9" position="bottom" height="30" data-wails-drag>
        <v-app-bar-nav-icon variant="text"
          @click.stop="onToggleDrawer = !onToggleDrawer">
        </v-app-bar-nav-icon>
        <v-btn variant="text" :icon="activeMenuIcon"> </v-btn>
        <v-spacer></v-spacer>
        <v-btn variant="text" icon="mdi-translate" @click="onclickToggleLanguage">
        </v-btn>
        <v-btn variant="text" icon="mdi-invert-colors" @click="onclickToggleTheme">
        </v-btn>
        <v-btn variant="text" icon="mdi-magnify"></v-btn>
        <v-btn variant="text" icon="mdi-view-module" @click="onclickViewModule"></v-btn>
        <v-btn variant="text" icon="mdi-github" @click="onclickOpenGithub"></v-btn>
        <v-btn variant="text" icon="mdi-dots-vertical"></v-btn>
      </v-app-bar>

      <v-navigation-drawer v-model="onToggleDrawer" temporary>
        <v-list>
          <v-list-subheader class="h3">{{
              t("nav.mainmenu")
          }}</v-list-subheader>
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

      <v-main>
        <router-view />
      </v-main>
    </v-app>
  </v-layout>
</template>

<script lang="ts" setup>
// vue about
import { onMounted, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useDisplay } from "vuetify";
import { useI18n } from "vue-i18n";
// wails about
import { WindowMinimise, Quit } from "../wailsjs/runtime";
import * as imes from "../wailsjs/go/imes/Middleware";
import * as app from "../wailsjs/go/main/App";
// app about
import { useBaseStore } from "./stores/index";

const router = useRouter(); // router 是管理器，可以 addRoute、removeRoute、getRoutes、push...
const route = useRoute(); // route 是一个响应式对象，
const store = useBaseStore();
const { t, availableLocales, locale } = useI18n({ useScope: "global" });

// vuetify's display info
const display = useDisplay();
onMounted(() => {
  console.log(display.height.value);
  console.log(display.width.value);
  console.log(display.mobile.value);
  console.log(display.platform.value);
  // 默认导航到页面
  router.push({
    name: "dashboard",
  });
  // console.log(router.getRoutes());
});

// 导航栏 menu 菜单相关
const onToggleDrawer = ref(true);
const activeMenuIcon = ref("mdi-view-dashboard");
const onclickMenuListItem = (val: string | unknown) => {
  if (val) {
    console.log(val);
    activeMenuIcon.value = val as string;
  }
};
watch(
  () => route.path,
  (newPath) => {
    console.log("watching route.path:", newPath);
  }
); watch(
  () => route.params,
  (v) => {
    console.log("watching route.params:", v);
  }
);

// change theme
const onclickToggleTheme = () => {
  store.apptheme = store.apptheme === "light" ? "dark" : "light";
  console.log(store.apptheme);
};

// change i18n
const languages = availableLocales;
const onclickLanguageHandle = (item: string) => {
  // 所有语言列表切换
  item !== locale.value ? (locale.value = item) : false;
};
const onclickToggleLanguage = () => {
  // 中英文切换
  locale.value
    ? locale.value == "en"
      ? (locale.value = "zh-Hans")
      : (locale.value = "en")
    : false;
};
const onclickViewModule = () => {
  store.addCounter();
  store.anothercounter += 1;
};
const onclickOpenGithub = () => {
  imes.OpenGithub();
};
// hide window
const onclickMinimise = () => {
  WindowMinimise();
};
// close app
const onclickQuit = () => {
  Quit();
};
</script>

<style lang="scss">
@import url("./assets/css/reset.css");
@import url("./assets/css/font.css");

#app {
  position: relative;
  height: 100%;
  overflow: hidden;
}
</style>
