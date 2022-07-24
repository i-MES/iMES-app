import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import vuetify from '@vuetify/vite-plugin'

import pkg from './package.json'
import { createHtmlPlugin } from 'vite-plugin-html'
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
    vuetify({
      autoImport: true,
    }),
    createHtmlPlugin({
      minify: false,
      // entry: "/src/main.ts",
      template: 'index.html',
      inject: {
        data: {
          title: `${pkg.name}`,
        },
      },
    }),
  ],
  resolve: {
    alias: {
      // 这样 ts 中 import 才可以使用 "@/..."
      // __dirname 若提示找不到，需 yarn add @types/node --save-dev
      '@': path.resolve(__dirname, 'src'),
      // 否则 warnning：You are running the esm-bundler build of vue-i18n....
      'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
    },
  },
  build: {
    outDir: './dist',
    rollupOptions: {
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
      },
    },
  },
})
