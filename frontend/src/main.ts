import { createApp } from "vue";
import { createPinia } from 'pinia'
import App from "./App.vue";
import router from "./router";
import { createI18n } from "vue-i18n";
import zhHans from "./i18n/zh-Hans.json";
import en from "./i18n/en.json";
import vuetify from "./plugins/vuetify"

// Register global common components
import publicComponents from "./components/public/";

const i18n = createI18n({
  locale: "zh-Hans",
  fallbackLocale: "zh-Hans",
  messages: {
    en: en,
    "zh-Hans": zhHans,
  },
});

const pinia = createPinia();


createApp(App)
  .use(publicComponents)
  .use(pinia)
  .use(router)
  .use(i18n)
  .use(vuetify)
  .mount("#app");
