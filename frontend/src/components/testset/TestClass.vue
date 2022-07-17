<template>
  <div class="mt-5">
    <DragHandle v-if="showHandle" />
    <v-expansion-panels v-bind="props" class="px-2 mt-1" multiple>
      <v-tooltip v-if="store.enableTCTooltip" activator="parent" location="top"
        transition="slide-y-reverse-transition">
        <h2>TestClass</h2>
        title: {{ tc.title }}<br />
        parametrizes: {{ tc.parametrizes }}<br />
        fixtures: {{ tc.fixtures }}<br />
        modulepath: {{ tc.modulepath }}<br />
        modulename: {{ tc.modulename }}<br />
      </v-tooltip>
      <test-item v-for="(ti, i) in tc.testitems" :key="ti.title" :id="i" :teid="teid"
        :tgid="tgid" :tcid="tc.id" :ti="ti" />
    </v-expansion-panels>
  </div>
</template>


<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { target } from '../../../wailsjs/go/models'
import TestItem from './TestItem.vue'
import DragHandle from '../utils/DragHandle.vue'
import * as runtime from '../../../wailsjs/runtime/runtime'
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()

const props = defineProps<{
  teid: string,
  tgid: string,
  tc: target.TestClass,
}>()

const showHandle = ref(false)

onMounted(() => {
  console.log('++++ TestClass: ', props.tc, props.tgid)
  runtime.EventsEmit('testclasscreated', props.tgid)
})
onBeforeUnmount(() => {
  console.log('---- TestClass: ', props.tc, props.tgid)
  runtime.EventsEmit('testclassdeleted', props.tgid)
})
</script>

<style lang="scss" scoped>
</style>
