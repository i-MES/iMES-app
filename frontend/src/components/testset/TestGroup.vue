<template>
  <v-sheet onselectstart="return false"
    :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
    <v-toolbar height="20">
      <v-toolbar-title>Group_{{ tg.title == '' ? tg.id.substr(0, 5) : tg.title }}
      </v-toolbar-title>
      <v-btn @click="starttestgroup(tg)" icon="mdi-arrow-right-bold-circle-outline"
        :disabled="store.testGroupsStatus[tgidx].disableStartBtn">
      </v-btn>
      <v-btn @click="stoptestgroup(tg)" icon="mdi-stop-circle-outline"
        :disabled="store.testGroupsStatus[tgidx].disableStopBtn">
      </v-btn>
      <v-btn @click="store.newTestGroup(tg.id)"
        :disabled="store.testGroupsStatus[tgidx].disableNewBtn"
        icon="mdi-format-horizontal-align-right">
      </v-btn>
    </v-toolbar>
    <div style="min-height:100px">
      <slot />
    </div>
  </v-sheet>
</template>

<script lang="ts">
import { ContainerMixin } from 'vue-slicksort'
import { watch } from 'fs'
export default {
  mixins: [ContainerMixin]
}
</script>

<script lang="ts" setup>
import { onMounted, onBeforeUnmount } from 'vue'
import { target } from '../../../wailsjs/go/models'
import { RunTestGroup } from '../../../wailsjs/go/imes/Api'
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()
const props = defineProps<{
  tg: target.TestGroup,
  tgidx: number
}>()
onMounted(() => {
  console.log('++++++ TestGroup: ', props.tg)
  console.log(props.tgidx)
  console.log(store.testGroupsStatus[props.tgidx])

  // 置初始值
  if (props.tg.testclasses.length == 0) {
    store.testGroupsStatus[props.tgidx].disableStartBtn = true
    store.testGroupsStatus[props.tgidx].disableStopBtn = true
  }

  window.runtime.EventsOn('testclasscreated', (tgid) => {
    if (tgid == props.tg.id) {
      // 在当前 group 中创建了 testclass
      store.testGroupsStatus[props.tgidx].disableStartBtn = false
      store.testGroupsStatus[props.tgidx].disableStopBtn = false
    }
  })
  window.runtime.EventsOn('testclassdeleted', (tgid) => {
    if (tgid == props.tg.id) {
      // 在当前 group 中删除了 testclass
      if (props.tg.testclasses.length == 0) {
        store.testGroupsStatus[props.tgidx].disableStartBtn = true
        store.testGroupsStatus[props.tgidx].disableStopBtn = true
      }
    }
  })
  window.runtime.EventsOn('testgroupfinished', (tgid) => {
    if (tgid == props.tg.id) {
      // 后端组测试完毕
      store.testGroupsStatus[props.tgidx].disableStartBtn = false
      store.testGroupsStatus[props.tgidx].disableStopBtn = true
      store.testGroupsStatus[props.tgidx].disableNewBtn = false
    }
  })
})
onBeforeUnmount(() => {
  console.log('------ TestGroup: ', props.tg)
})

// 执行组测试
const starttestgroup = (tg: target.TestGroup) => {
  console.log(tg)
  store.testGroupsStatus[props.tgidx].disableStartBtn = true
  store.testGroupsStatus[props.tgidx].disableStopBtn = false
  store.testGroupsStatus[props.tgidx].disableNewBtn = true
  // 只有 activedTestEntityId 才会被用户点击
  RunTestGroup(store.activedTestEntityId, tg)
  // 将 grop 内所有 ti 的滚动条清零
  window.runtime.EventsEmit('clearprocessbar', {
    teid: store.activedTestEntityId,
    tgid: tg.id,
  })
}
// 停止组测试
const stoptestgroup = (tg: target.TestGroup) => {
  console.log(tg)
}
</script>

<style lang="scss" scoped>
</style>
