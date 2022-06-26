<template>
  <div class="mt-5">
    <DragHandle v-if="showHandle" />
    <div class="mt-2 ml-2">
      {{ tc.title }}
    </div>
    <v-expansion-panels class="px-2 mt-1" multiple>
      <test-item v-for="(ti, i) in tc.testitems" :key="ti.title" :id="i" :teid="teid"
        :tgid="tgid" :tcid="tc.id" :ti="ti" />
    </v-expansion-panels>
  </div>
</template>

<script lang="ts">
import { ElementMixin } from 'vue-slicksort'
import DragHandle from '../../example/slicksort/components/DragHandle.vue'

export default {
  mixins: [ElementMixin],
  components: {
    DragHandle,
  }
}
</script>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { target } from '../../../wailsjs/go/models'
import TestItem from './TestItem.vue'

const props = defineProps<{
  teid: string,
  tgid: string,
  tc: target.TestClass,
}>()

const showHandle = ref(false)


onMounted(() => {
  console.log('++++ TestClass: ', props.tc, props.tgid)
  window.runtime.EventsEmit('testclasscreated', props.tgid)
})
onBeforeUnmount(() => {
  console.log('---- TestClass: ', props.tc, props.tgid)
  window.runtime.EventsEmit('testclassdeleted', props.tgid)
})
</script>

<style lang="scss" scoped>
</style>
