<template>
  <div>
    <DragHandle v-if="showHandle" />
    <div class="text-subtitle-2 ml-2">
      {{ tc.title }}
    </div>
    <v-expansion-panels class="px-2 mt-2" multiple>
      <v-expansion-panel v-for="(ti, i) in tc.testitems" :key="ti.title" :id="i">
        <v-expansion-panel-title>
          <template v-slot:default>
            <span style="max-width:100%" class="text-truncate ">{{ ti.title }}</span>
          </template>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          {{ ti.title }}<br /><br />
          {{ ti.filename }}
        </v-expansion-panel-text>
        <div>
          <v-progress-linear v-model="value" :buffer-value="bufferValue">
          </v-progress-linear>
        </div>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script lang="ts">
import { ElementMixin } from 'vue-slicksort'
import DragHandle from '../../example/slicksort/components/DragHandle.vue'
// import { testset } from '../../../wailsjs/go/models'

export default {
  mixins: [ElementMixin],
  components: {
    DragHandle,
  },
  props: {
    tc: {
      type: Object,
      required: true,
    },
    showHandle: Boolean,
  },
}
</script>

<script lang="ts" setup>
import { ref, watch, onBeforeUnmount } from 'vue'

const value = ref(10)
const bufferValue = ref(20)

const interval = setInterval(() => {
  // value.value += Math.random() * (15 - 5) + 5;
  // bufferValue.value += Math.random() * (15 - 5) + 6;
}, 2000)

watch(value, (n) => {
  if (n < 100) return
  value.value = 0
  bufferValue.value = 10
})
onBeforeUnmount(() => {
  clearInterval(interval)
})


</script>

<style lang="scss" scoped>
.button {
  margin-left: auto;
}

.example-list-item {
  list-style-type: none;
  display: flex;
  align-items: center;
  padding: 10px 20px;
  margin: 10px;
  border-radius: 10px;
  font-weight: bold;
  background: white;
  box-shadow: inset 0 0 0 3px rgba(0, 0, 0, 0.1), 1px 2px 5px rgba(0, 0, 0, 0.15);
  background: #9b51e0;
  color: white;
  line-height: 1.4;
  user-select: none;
}
</style>
