import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { createI18n } from "vue-i18n";
import zhHans from "./i18n/zh-Hans.json";
import en from "./i18n/en.json";

// Register global common components
import publicComponents from "./components/public/";

const i18n = createI18n({
  locale: "en",
  fallbackLocale: "en",
  messages: {
    en: en,
    "zh-Hans": zhHans,
  },
});

const app = createApp(App);

app.use(publicComponents);

app.use(router).use(i18n).mount("#app");
