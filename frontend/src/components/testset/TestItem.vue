<template>
  <v-expansion-panel>
    <v-expansion-panel-title>
      <template v-slot:default>
        <span style="max-width:100%" class="text-truncate">{{ ti.title
        }}</span>
      </template>
    </v-expansion-panel-title>
    <v-expansion-panel-text>
      {{ ti.title }}<br /><br />
      {{ ti.filename }}
    </v-expansion-panel-text>
    <div>
      <v-progress-linear v-model="progressdata" :indeterminate="isrunning">
      </v-progress-linear>
    </div>
  </v-expansion-panel>
</template>

<script lang="ts" setup>
import { ref, onBeforeUnmount, onMounted } from 'vue'
import { target } from '../../../wailsjs/go/models'

const props = defineProps<{
  teid: string,
  tgid: string,
  tcid: string,
  ti: target.TestItem
}>()
const progressdata = ref(0)
const isrunning = ref(false)

// watch(progressdata, (n) => {
//   if (n < 100) return
//   progressdata.value = 0
// })
// const interval = setInterval(() => {
//   progressdata.value += Math.random() * (15 - 5) + 5
// }, 2000)

onMounted(() => {
  console.log('++ TestItem:', props.ti)
  window.runtime.EventsOn('testitemstatus', (tistatus: target.TestItemStatus) => {
    // console.log('receive event', props.tc.title, props.tc.id, tistatus)
    if (props.ti.id == tistatus.testitemid) {
      console.log('==', props.ti.title, tistatus.status)
      if (tistatus.status == 'started') {
        isrunning.value = true
      } else if (tistatus.status == 'finished') {
        isrunning.value = false
        progressdata.value = 100
      }
    } else {
      // 同 group、同 tc 的 ti 需要……
    }
  })
  window.runtime.EventsOn('clearprocessbar', (ids) => {
    if (ids.teid == props.teid && ids.tgid == props.tgid) {
      console.log('EventsOn clearprocessbar')
      progressdata.value = 0
    }
  })
})

onBeforeUnmount(() => {
  // clearInterval(interval)
  console.log('-- TestItem:', props.ti)
  // window.runtime.EventsOff('testitemstatus')
})
</script>

<style lang="scss" scoped>
</style>