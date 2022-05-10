<template>
  <v-layout>
    <v-app :theme="store.apptheme">
      <v-app-bar color="info" position="top" data-wails-drag>
        <v-app-bar-nav-icon variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
        <v-spacer></v-spacer>
        <v-btn variant="text" icon="mdi-translate" @click="onclickToggleLanguage"></v-btn>
        <v-btn variant="text" icon="mdi-invert-colors" @click="onclickToggleTheme"></v-btn>
        <v-btn variant="text" icon="mdi-magnify"></v-btn>
        <v-btn variant="text" icon="mdi-view-module"></v-btn>
        <v-btn variant="text" icon="mdi-dots-vertical"></v-btn>
      </v-app-bar>

      <v-navigation-drawer v-model="drawer" temporary>
        <v-list>
          <v-list-subheader>{{ t("nav.mainmenu") }}</v-list-subheader>
          <v-list-item
            v-for="(menu, i) in menus"
            :key="i"
            :value="menu"
            active-color="primary"
            density="comfortable"
            :to="menu.ref"
            @click="onclickMenuListItem(menu.ref)"
          >
            <v-list-item-avatar start>
              <v-icon :icon="menu.icon"></v-icon>
            </v-list-item-avatar>
            <v-list-item-title v-text="t(menu.title)"></v-list-item-title>
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
import { onMounted } from "vue";
import { useDisplay } from "vuetify";
import * as imes from "../wailsjs/go/imes/Middleware";
import { useBaseStore } from "./store/index";
import { reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import { WindowMinimise, Quit } from "../wailsjs/runtime";
const store = useBaseStore();
const { t, availableLocales, locale } = useI18n({ useScope: "global" });

// vuetify's display info
const display = useDisplay();
onMounted(() => {
  console.log(display.height.value);
  console.log(display.width.value);
  console.log(display.mobile.value);
  console.log(display.platform.value);
});

const drawer = ref(true);
const menus = reactive([
  {
    title: "nav.home",
    ref: "/",
    icon: 'mdi-clock'
  },
  {
    title: "nav.dashboard",
    ref: "/dashboard",
    icon: 'mdi-account'
  },
  {
    title: "nav.about",
    ref: "/about",
    icon: 'mdi-flag'
  },
])


// change theme
const onclickToggleTheme = () => {
  store.apptheme = store.apptheme === "light" ? "dark" : "light";
  console.log(store.apptheme);
};
// change i18n
const languages = availableLocales;
const onclickLanguageHandle = (item: string) => {
  item !== locale.value ? (locale.value = item) : false;
};
const onclickToggleLanguage = () => {
  locale.value
    ? locale.value == "en"
      ? (locale.value = "zh-Hans")
      : (locale.value = "en")
    : false;
};
const onclickMenuListItem = (val: string) => {
  console.log(val)
}
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
