<template>
  <v-toolbar class="entity-toolbar" :height="store.toolbarheight">
    <v-row class="ma-0 pa-0">
      <v-col cols="3">
        <v-select class="mt-3" filled label="产品" dense hide-details
          v-model="selectedProd" hideDetails
          :items="store.testProductions.map((v, _) => v.id + '-' + v.title)">
        </v-select>
      </v-col>
      <v-col cols="3">
        <v-select class="mt-3" filled label="工序" dense hide-details
          v-model="selectedStage" :items="stages">
        </v-select>
      </v-col>
    </v-row>
    <template v-slot:append>
      <v-btn variant="text" icon="mdi-view-module" @click="onclickViewModule">
      </v-btn>
    </template>
  </v-toolbar>
  <v-sheet class="ma-0 pt-10 overflow-y-auto"
    :height="store.availableHeight - store.logHeight - store.toolbarheight">
    <TE v-if="store.TEorTI" />
    <TI v-else />
  </v-sheet>
  <v-sheet class="ma-0 pa-0 overflow-y-auto "
    :height="store.toolbarheight + store.logHeight">
    <TL />
  </v-sheet>
</template>

<script lang="ts" setup>
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import TI from '../components/TestItem.vue'
import TE from '../components/TestEntity.vue'
import TL from '../components/TestLog.vue'

const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

// Toolbar 相关
const autoLoadConfig = ref(true)
const selectedProd = ref()
const selectedStage = ref()
const planab = ref('a')
if (false) {
  store.initConfig()
}
// 加载已保存的数据
store.syncTestProductions()
store.syncTestStages()
store.syncTestStation()
store.syncTestEntity()
store.syncTestItem()

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
    store.testStageByProductionId(pid).map((v, _) => (v.id + '-' + v.title)).forEach(
      (n) => stages.push(n)
    )
  }
)
watch(
  () => selectedStage.value,
  (nv) => {
    var sid = Number((nv as string).split('-')[0])
    store.activedTestStageId = sid
  }
)

const onclickViewModule = () => {
  store.TEorTI = true
}
</script>

