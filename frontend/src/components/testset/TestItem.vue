<template>
  <v-expansion-panel>
    <v-expansion-panel-title>
      <template v-slot:default>
        <span style="max-width:100%" class="text-truncate">{{ ti.docstr ? ti.docstr :
            ti.title
        }}</span>
      </template>
    </v-expansion-panel-title>
    <v-expansion-panel-text>
      TestItem<br />
      desc: {{ ti.desc }}<br />
      modulepath: {{ ti.modulepath }}<br />
      modulename: {{ ti.modulename }}<br />
      funcname: {{ ti.funcname }}<br />
      args: {{ ti.args }}<br />
      docstr: {{ ti.docstr }}<br />
    </v-expansion-panel-text>
    <div>
      <v-progress-linear v-model="progressdata" :indeterminate="isrunning"
        :color="pcolor">
      </v-progress-linear>
    </div>
  </v-expansion-panel>
</template>

<script lang="ts" setup>
import { ref, onBeforeUnmount, onMounted } from 'vue'
import { target } from '../../../wailsjs/go/models'
import * as runtime from '../../../wailsjs/runtime/runtime'

const props = defineProps<{
  teid: string,
  tgid: string,
  tcid: string,
  ti: target.TestItem
}>()
const progressdata = ref(0)
const isrunning = ref(false)
const pcolor = ref('')

// watch(progressdata, (n) => {
//   if (n < 100) return
//   progressdata.value = 0
// })
// const interval = setInterval(() => {
//   progressdata.value += Math.random() * (15 - 5) + 5
// }, 2000)

onMounted(() => {
  console.log('++ TestItem:', props.ti)
  runtime.EventsOn('testitemstatus', (tistatus: target.TestItemStatus) => {
    console.log('receive event', tistatus)
    if (props.ti.id == tistatus.testitemid) {
      console.log('==', props.ti.title, tistatus.status)
      if (tistatus.status == 'started') {
        isrunning.value = true
        pcolor.value = ''
      } else if (tistatus.status == 'pass') {
        isrunning.value = false
        progressdata.value = 100
        pcolor.value = 'green'
      } else if (tistatus.status == 'ng') {
        isrunning.value = false
        progressdata.value = 100
        pcolor.value = 'red'
      }
    } else {
      // 同 group、同 tc 的 ti 需要……
    }
  })
  runtime.EventsOn('clearprocessbar', (ids) => {
    if (ids.teid == props.teid && ids.tgid == props.tgid) {
      console.log('EventsOn clearprocessbar')
      progressdata.value = 0
    }
  })
})

onBeforeUnmount(() => {
  // clearInterval(interval)
  console.log('-- TestItem:', props.ti)
  // runtime.EventsOff('testitemstatus')
})
</script>

<style lang="scss" scoped>
</style>