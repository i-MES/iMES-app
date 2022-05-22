<template>
  <v-navigation-drawer width="180" v-model="toggleStep" expand-on-hover rail position="right">
    <v-select
      filled
      label="产品"
      dense
      hide-details
      v-model="selectedProd"
      :items="store.testProductions.map((v, idx) => v.id + v.title)"
    ></v-select>
    <v-select filled label="机型" dense hide-details :items="['Foo', 'Bar']"></v-select>
    <v-select filled label="工序" dense hide-details :items="['组装', '高温', '写版本']"></v-select>
    <v-select filled label="视图" dense hide-details :items="['单体', '全体']"></v-select>
    <template v-slot:append>
      <div class="pa-2">
        <v-btn @click="onclickLoadConfig">LoadConfig</v-btn>
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { useBaseStore } from '../stores/index'

const props = withDefaults(
  defineProps<{
    toggleStep: boolean
  }>(),
  {
    toggleStep: false,
  }
)
const store = useBaseStore()
const selectedProd = ref(0)

const onclickLoadConfig = () => {
  store.loadTestProductions()
}

watch(
  () => selectedProd.value,
  (nv) => {
    console.log('selectedProd changed: ', nv)
    // 用户重新选择了产品
    var _tp = store.testProductionById(nv)
    if (_tp) {
      store.appStatusBar.Production = _tp.title
      // 加载对应产品的工序
      store.loadSteps()
    }
  }
)
</script>
