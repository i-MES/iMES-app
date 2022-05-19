<template>
  <v-hover v-slot="{ isHovering, props }" open-delay="600" close-delay="600">
    <v-card class="mx-0 log-card" color="grey-lighten-4" v-bind="props">
      <v-toolbar :height="store.toolbarheight">
        <v-toolbar-title>{{ t("testpage.testitem-log") }}</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-text-field hide-details :height="store.toolbarheight - 10"
          prepend-icon="mdi-magnify">
        </v-text-field>
        <v-btn icon="mdi-upload-outline"
          @click="logHeight + 50 > store.tiPageAvilableHeight ? (logHeight = store.tiPageAvilableHeight) : (logHeight += 50)">
        </v-btn>
        <v-btn icon="mdi-download-outline"
          @click="logHeight - 50 < 0 ? (logHeight = 0) : (logHeight -= 50)">
        </v-btn>
        <v-btn :icon="logHeightMaxed ? 'mdi-download-multiple' : 'mdi-upload-multiple'"
          @click="logHeightMaxed = !logHeightMaxed">
        </v-btn>
        <v-btn :icon="sticky ? 'mdi-pin' : 'mdi-pin-off'" @click="sticky = !sticky">
        </v-btn>
      </v-toolbar>

      <v-table v-if="isHovering || sticky" density="compact" :height="logHeight">
        <thead>
          <tr>
            <th class="text-left" width="8%">No.</th>
            <th class="text-left" width="30%">Name</th>
            <th class="text-left">Calories</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(log, idx) in store.testitemsLogs" :key="log.timestamp">
            <td>{{ idx + 1 }}</td>
            <td>{{ DateTime.fromSeconds(log.timestamp).toFormat('yyyy-MM-dd HH:MM:ss')
            }}</td>
            <td>{{ log.message }}</td>
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </v-hover>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useBaseStore } from '../stores/index'
import { TestItemStart } from '../../wailsjs/go/imes/Middleware'
import { useI18n } from 'vue-i18n'
import { DateTime } from 'luxon'

const { t } = useI18n({ useScope: 'global' })

const store = useBaseStore()
const logHeight = ref(200) // card - toolbar height, only log height
const logHeightBak = ref(0)
const logHeightMaxed = ref(false)
const sticky = ref(false)
const stickyBak = ref(false)

const timer = setInterval(() => {
  TestItemStart(1).then((val) => {
    console.log('测试项启动：', val ? '成功' : '失败')
  })
}, 3000)

onMounted(() => {
  store.appStatusBar.logHeight = logHeight.value
})
onUnmounted(() => {
  clearInterval(timer)
})

watch(
  () => logHeight.value,
  (nv, ov) => {
    store.tiPageAvilableHeight -= (nv - ov)
    store.appStatusBar.logHeight = nv
  }
)
watch(logHeightMaxed, (nv) => {
  if (nv) {
    // max log window
    logHeightBak.value = logHeight.value
    logHeight.value = store.tiPageAvilableHeight
    stickyBak.value = sticky.value
    sticky.value = true
  } else {
    // retreave log window
    logHeight.value = logHeightBak.value
    sticky.value = stickyBak.value
  }
})

window.runtime.EventsOn('testitemlog', (data) => {
  store.testitemsLogs.push(data)
})
</script>

<style>
.log-card {
  align-items: center;
  bottom: 0;
  justify-content: center;
  opacity: 0.95;
  position: absolute;
  width: 100%;
  z-index: 1;
}
</style>
