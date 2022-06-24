<template>
  <v-card class="ma-0">

    <!-- 顶部 toolbar -->
    <v-toolbar class="entity-toolbar" :height="store.toolbarheight">
      <v-row class="ma-0 pa-0" align="center" align-content="center">
        <v-col cols="3">
          <v-select class="mt-3" filled :label="t('nav.production')" dense hide-details
            v-model="selectedProd"
            :items="store.testProductions.map((v, _) => v.id + '-' + v.title)">
          </v-select>
        </v-col>
        <v-col cols="3">
          <v-select class="mt-3" filled :label="t('nav.stage')" dense hide-details
            v-model="selectedStage" :items="stages">
          </v-select>
        </v-col>
      </v-row>
      <template v-slot:append>
        <v-btn variant="text" @click="zoomOut" icon="mdi-magnify-minus-outline"></v-btn>
        <v-btn variant="text" @click="zoomIn" icon="mdi-magnify-plus-outline"></v-btn>
        <add-entity />
        <v-btn variant="text" icon="mdi-view-module"
          @click="store.TEorTI = !store.TEorTI">
        </v-btn>
      </template>
    </v-toolbar>

    <!-- TestSet 主窗口 -->
    <v-sheet class="ma-0 pt-5 overflow-y-auto"
      :height="store.mainWindowHeight - store.logHeight - store.toolbarheight">
      <test-entity v-if="store.TEorTI" :defcols="defCols" />
      <test-set v-else />
    </v-sheet>

    <!-- TestLog 可扩展窗口 -->
    <v-sheet class="ma-0 pa-0 overflow-y-auto"
      :height="store.toolbarheight + store.logHeight">
      <test-log />
    </v-sheet>
  </v-card>
</template>

<script lang="ts" setup>
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import TestSet from '../components/TestSet.vue'
import TestEntity from '../components/TestEntity.vue'
import TestLog from '../components/TestLog.vue'
import AddEntity from '../components/forms/AddEntity.vue'

const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

// Toolbar 相关
const selectedProd = ref()
const selectedStage = ref()

const stages = reactive([])
watch(
  () => selectedProd.value,
  (nv) => {
    var pid = Number((nv as string).split('-')[0])
    store.activedProductionId = pid
    // 联动 stage 工序选择栏
    var i = stages.length
    while (i--) {
      stages.splice(i, 1)
    }
    store
      .testStageByProductionId(pid)
      .map((v, _) => v.id + '-' + v.title)
      .forEach((n) => stages.push(n))
  }
)
watch(
  () => selectedStage.value,
  (nv) => {
    var sid = Number((nv as string).split('-')[0])
    store.activedTestStageId = sid
  }
)

const defCols = ref(3)
const zoomOut = () => {
  defCols.value = (defCols.value - 1 < 2) ? 2 : defCols.value - 1
  // console.log(defCols.value)
}
const zoomIn = () => {
  defCols.value = (defCols.value + 1 > 12) ? 12 : defCols.value + 1
  // console.log(defCols.value)
}

</script>
