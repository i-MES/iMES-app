<template>
  <v-sheet onselectstart="return false"
    :color="store.appTheme == 'dark' ? store.darkmaincolor : store.lightmaincolor">
    <v-toolbar height="40"
      :style="'background:' + (store.appTheme == 'dark' ? store.darkmaincolor : store.lightmaincolor)">
      <drag-handle />
      <v-toolbar-title>
        G-{{ tg.title == '' ? tg.id.substring(0, 4) : tg.title }}
      </v-toolbar-title>
      <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
        @click="starttestgroup()" :disabled="disableBtnRunGroup">
        <v-icon>mdi-arrow-right-bold-circle-outline</v-icon>
      </v-btn>
      <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
        @click="stoptestgroup()" :disabled="disableBtnStopGroup">
        <v-icon>mdi-stop-circle-outline</v-icon>
      </v-btn>
      <v-menu open-on-hover>
        <template v-slot:activator="{ props }">
          <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
            v-bind="props">
            <v-icon>mdi-dots-horizontal</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item density="compact" active-color="primary"
            @click="store.NewTestGroup(tg.id)" :disabled="disableBtnNewGroup">
            <v-list-item-avatar>
              <v-icon icon="mdi-format-horizontal-align-right"></v-icon>
            </v-list-item-avatar>
            <v-list-item-title>新建组</v-list-item-title>
          </v-list-item>
          <v-list-item density="compact" active-color="primary"
            @click="store.DelTestGroup(tg.id)" :disabled="disableBtnDelGroup">
            <v-list-item-avatar>
              <v-icon icon="mdi-delete-circle-outline"></v-icon>
            </v-list-item-avatar>
            <v-list-item-title>删除组</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-toolbar>
    <v-progress-linear model-value="0" color="light-green darken-4" height="10" striped>
    </v-progress-linear>
    <div style="min-height:100px">
      <slot />
    </div>
  </v-sheet>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { target } from '../../../wailsjs/go/models'
import { RunTestGroup } from '../../../wailsjs/go/imes/Api'
import { useBaseStore } from '../../stores/index'
import DragHandle from '../utils/DragHandle.vue'
import * as runtime from '../../../wailsjs/runtime/runtime'

const props = defineProps<{
  tg: target.TestGroup,
}>()
const store = useBaseStore()
const disableBtnRunGroup = ref(true)
const disableBtnStopGroup = ref(true)
const disableBtnNewGroup = ref(false)
const disableBtnDelGroup = ref(true)

// watch(
//   () => props.tg.testclasses,
//   (nv) => {
//     if (nv.length == 0) {
//       disableBtnRunGroup.value = true
//       disableBtnStopGroup.value = true
//       disableBtnNewGroup.value = false
//       disableBtnDelGroup.value = false
//     } else {
//       disableBtnRunGroup.value = false
//       disableBtnStopGroup.value = true
//       disableBtnNewGroup.value = false
//       disableBtnDelGroup.value = true
//     }
//   }
// )
onMounted(() => {
  console.log('++++++ TestGroup: ', props.tg)

  // 置初始值
  if (props.tg.testclasses.length == 0) {
    disableBtnRunGroup.value = true
    disableBtnStopGroup.value = true
    disableBtnNewGroup.value = false
    disableBtnDelGroup.value = false
  } else {
    disableBtnRunGroup.value = false
    disableBtnStopGroup.value = true
    disableBtnNewGroup.value = false
    disableBtnDelGroup.value = true
  }
  runtime.EventsOn('testclasscreated', (tgid: string) => {
    if (tgid == props.tg.id) {
      // 在当前 group 中放入了 testclass：使能测试
      disableBtnRunGroup.value = false
      disableBtnStopGroup.value = true
      disableBtnDelGroup.value = true
    }
  })
  runtime.EventsOn('testclassdeleted', (tgid: string) => {
    if (tgid == props.tg.id) {
      // 在当前 group 中移走了 testclass：去使能测试
      if (props.tg.testclasses.length == 0) {
        disableBtnRunGroup.value = true
        disableBtnStopGroup.value = true
        disableBtnDelGroup.value = false
      }
    }
  })
  runtime.EventsOn('testgroupfinished', (tgid: string) => {
    if (tgid == props.tg.id) {
      // 后端（go）组测试完毕：使能测试
      disableBtnRunGroup.value = false
      disableBtnStopGroup.value = true
      disableBtnNewGroup.value = false
    }
  })
  runtime.EventsOn('startallgroup', () => {
    starttestgroup()
  })
})
onBeforeUnmount(() => {
  console.log('------ TestGroup: ', props.tg)
})

// 执行组测试
const starttestgroup = () => {
  console.log(props.tg)
  disableBtnRunGroup.value = true
  disableBtnStopGroup.value = false
  disableBtnNewGroup.value = true
  disableBtnDelGroup.value = true
  // 清空本 group 的所有 item 的 status
  if (store.testEntitiesTIStatus[store.activedTestEntityId]) {
    store.testEntitiesTIStatus[store.activedTestEntityId].forEach((tis) => {
      if (tis.testgroupid == props.tg.id) {
        tis.status = ''
      }
    })
  }
  // 只有 activedTestEntityId 才会被用户点击
  RunTestGroup(store.activedTestEntityId, props.tg)
  // 将 group 内所有 ti 的滚动条清零
  runtime.EventsEmit('clearprocessbar', {
    teid: store.activedTestEntityId,
    tgid: props.tg.id,
  })
}
// 停止组测试
const stoptestgroup = () => {
  console.log('stop group', props.tg)
}
</script>

<style lang="scss" scoped>
</style>
