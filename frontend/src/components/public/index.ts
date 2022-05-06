import { App } from "vue";
import OpenLink from "./OpenLink.vue";

// Encapsulate global components as plug-ins
// 将全局组件封装为插件

export default {
  install(app: App) {
    app.component(OpenLink.name, OpenLink);
  },
};
