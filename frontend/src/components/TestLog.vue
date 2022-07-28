<template>
  <v-sheet>
    <v-toolbar :height="store.toolbarheight"
      style="position: fixed; width:calc(100% - 55px)">
      <v-toolbar-title>{{ t('testpage.testitem-log') }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-text-field hide-details append-icon="mdi-magnify"> </v-text-field>
      <v-btn :icon="logHeightMaxed ? 'mdi-download-multiple' : 'mdi-upload-multiple'"
        @click="onclickMax">
      </v-btn>
    </v-toolbar>

    <v-table class="pt-10" density="compact">
      <thead>
        <tr>
          <th class="text-left" width="8%">No.</th>
          <th class="text-left" width="30%">时间戳</th>
          <th class="text-left" width="30%">TestItemID</th>
          <th class="text-left">Flag</th>
          <th class="text-left">Message</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(st, idx) in store.LastestTIStatus[store.activedTestEntityId]"
          :key="st.timestamp">
          <td>{{ idx + 1 }}</td>
          <td>{{ DateTime.fromSeconds(st.timestamp).toFormat('yyyy-MM-dd HH:MM:ss')
          }}</td>
          <td>{{ st.testitemid }}</td>
          <td>{{ st.status == 'pass' ? "√" : st.status === 'ng' ? "×" : "" }}</td>
          <td>{{ st.status }}</td>
        </tr>
      </tbody>
    </v-table>
  </v-sheet>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useBaseStore } from '../stores/index'
// import { TestItemStart } from '../../wailsjs/go/imes/Api'
import { useI18n } from 'vue-i18n'
import { DateTime } from 'luxon'
import * as runtime from '../../wailsjs/runtime/runtime'

// const props = defineProps<{
//   mainWindowHeight: number
// }>()
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()
const logHeightMaxed = ref(false)
// const pinBak = ref(false)

const onclickMax = () => {
  //   if (!logHeightMaxed.value) {
  //     // max log window
  //     logHeightBak.value = logHeight.value
  //     logHeight.value = store.mainWindowHeight - store.toolbarheight
  //     pinBak.value = pin.value
  //     pin.value = true
  //   } else {
  //     // retreave log window
  //     logHeight.value = logHeightBak.value
  //     pin.value = pinBak.value
  //   }
  //   logHeightMaxed.value = !logHeightMaxed.value
}
// watch(
//   () => pin.value,
//   (nv) => {
//     if (nv) {
//       store.logHeight = logHeight.value
//     } else {
//       store.logHeight = 0
//     }
//   }
// )
// const timer = setInterval(() => {
// TestItemStart(1).then((val) => {
//   console.log('测试项启动：', val ? '成功' : '失败')
// })
// }, 3000)

onMounted(() => {
  // store.appStatusBar.logHeight = logHeight.value
  runtime.EventsOn('testitemlog', (data) => {
    store.testitemsLogs.push(data)
  })
})
onUnmounted(() => {
  // clearInterval(timer)
  runtime.EventsOff('testitemlog')
})



</script>

<style>
.log-card {
  align-items: center;
  bottom: 0;
  justify-content: center;
  opacity: 0.9;
  position: absolute;
  width: 100%;
  z-index: 1;
}
</style>
