<template>
  <v-sheet>
    <v-toolbar :height="store.toolbarheight"
      style="position: fixed; width:calc(100% - 55px)">
      <v-toolbar-title>{{ t('testpage.testitem-log') }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-text-field hide-details append-icon="mdi-magnify"> </v-text-field>
      <v-btn icon="mdi-download-outline" @click="onclickMin"> </v-btn>
      <v-btn icon="mdi-upload-outline" @click="onclickMax"> </v-btn>
    </v-toolbar>

    <v-table class="pt-10" density="compact">
      <thead>
        <tr>
          <th class="text-left">No.</th>
          <th class="text-left">时间戳</th>
          <th class="text-left">TestItemID</th>
          <th class="text-left">Flag</th>
          <th class="text-left">Message</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(st, idx) in store.testEntitiesTIStatus[store.activedTestEntityId] "
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
import { onMounted, onUnmounted } from 'vue'
import { useBaseStore } from '../stores/index'
// import { TestItemStart } from '../../wailsjs/go/imes/Api'
import { useI18n } from 'vue-i18n'
import { DateTime } from 'luxon'

// const props = defineProps<{
//   mainWindowHeight: number
// }>()
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

const onclickMin = () => {
  store.paneFirstLengthPercent = 100 - ((store.toolbarheight + 2) * 100 / store.mainWindowHeight)
}
const onclickMax = () => {
  store.paneFirstLengthPercent = store.toolbarheight * 100 / store.mainWindowHeight
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
})
onUnmounted(() => {
  // clearInterval(timer)
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
