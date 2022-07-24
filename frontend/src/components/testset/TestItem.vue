<template>
  <v-expansion-panel>
    <v-expansion-panel-title>
      <template v-slot:default>
        <span style="max-width:100%" class="text-truncate">{{ ti.desc ? ti.desc :
            ti.title
        }}</span>
      </template>
    </v-expansion-panel-title>
    <v-expansion-panel-text>
      <ul>
        <li>title: {{ ti.title }}</li>
        <li>desc: {{ ti.desc }}</li>
        <li>args: {{ ti.args }}</li>
        <li>status: {{ statusmsg }}</li>
      </ul>
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
import { useBaseStore } from '../../stores/index'

const store = useBaseStore()

const props = defineProps<{
  teid: string,
  tgid: string,
  tcid: string,
  ti: target.TestItem
}>()
const progressdata = ref(0)
const isrunning = ref(false)
const pcolor = ref('')
const statusmsg = ref('')

// watch(progressdata, (n) => {
//   if (n < 100) return
//   progressdata.value = 0
// })
// const interval = setInterval(() => {
//   progressdata.value += Math.random() * (15 - 5) + 5
// }, 2000)

const setstatus = (tis: target.TestItemStatus) => {
  if (tis.status == 'started') {
    isrunning.value = true
    pcolor.value = ''
  } else if (tis.status == 'pass') {
    isrunning.value = false
    progressdata.value = 100
    pcolor.value = 'green'
  } else if (tis.status == 'ng') {
    isrunning.value = false
    progressdata.value = 100
    pcolor.value = 'red'
  }
}

var tis = store.LastestTIStatusById(props.teid, props.tgid, props.tcid, props.ti.id)
if (tis) {
  console.log('=-=-=-=-=-=-=')
  setstatus(tis)
} else {
  console.error(tis)
}

onMounted(() => {
  // console.log('++ TestItem:', props.ti)
  // 注册状态响应函数
  runtime.EventsOn('testitemstatus', (tistatus: target.TestItemStatus) => {
    // console.log('receive event', tistatus)
    if (props.ti.id == tistatus.testitemid && props.teid == tistatus.testentityid
      && props.tgid == tistatus.testgroupid && props.tcid == tistatus.testclassid) {
      setstatus(tistatus)
    }
  })
  runtime.EventsOn('clearprocessbar', ({ teid, tgid }) => {
    if (teid == props.teid && tgid == props.tgid) {
      console.log('EventsOn clearprocessbar')
      progressdata.value = 0
    }
  })
})

onBeforeUnmount(() => {
  // clearInterval(interval)
  // console.log('-- TestItem:', props.ti)
  // runtime.EventsOff('testitemstatus')
})
</script>

<style lang="scss" scoped>
</style>