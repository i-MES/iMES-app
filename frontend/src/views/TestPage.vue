<template>
  <v-card class="ma-0">

    <!-- 顶部 toolbar -->
    <v-toolbar class="entity-toolbar" :height="store.toolbarheight">
      <v-row class="ma-0 pa-0" align="center" align-content="center">
        <v-col cols="4">
          <v-select class="mt-3" filled :label="t('nav.production')" dense hide-details
            v-model="selectedProd"
            :items="store.testProductions.map((v, _) => v.id + '-' + v.title)">
          </v-select>
        </v-col>
        <v-col cols="4">
          <v-select class="mt-3" filled :label="t('nav.stage')" dense hide-details
            v-model="selectedStage" :items="stages">
          </v-select>
        </v-col>
      </v-row>

      <template v-slot:append>

        <!-- TestEntities 专用 Button -->
        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" v-if="store.TEsNotTE" @click="zoomOut">
              <v-icon>mdi-magnify-minus-outline</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.zoomout') }}</span>
        </v-tooltip>

        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" v-if="store.TEsNotTE" @click="zoomIn">
              <v-icon>mdi-magnify-plus-outline</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.zoomin') }}</span>
        </v-tooltip>

        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <add-entity v-if="store.TEsNotTE" v-bind="props" />
          </template>
          <span>{{ t('test.addentity') }}</span>
        </v-tooltip>

        <!-- TestEntity 专用 Button -->
        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" v-if="!store.TEsNotTE"
              @click="store.LoadTestGroup('src', true)">
              <v-badge v-if="testgroupsrcnewer" dot color="error">
                <v-icon>mdi-file-document-edit</v-icon>
              </v-badge>
              <v-icon v-else>mdi-file-document-edit</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.loadtestgroup') }}</span>
        </v-tooltip>

        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" v-if="!store.TEsNotTE"
              @click="store.canSortTestClass = !store.canSortTestClass">
              <v-icon>{{ store.canSortTestClass ? `mdi-hand-back-left` :
                  `mdi-hand-back-left-off`
              }}</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.cansorttestclass') }}</span>
        </v-tooltip>

        <!-- TestEntities & TestEntity 都可用 Button -->
        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" @click="startallgroup">
              <v-icon>mdi-arrow-right-bold-circle-outline</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.startallgroup') }}</span>
        </v-tooltip>

        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" @click="stopallgroup">
              <v-icon>mdi-stop-circle-outline</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.stopallgroup') }}</span>
        </v-tooltip>

        <v-tooltip location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
              v-bind="props" @click="store.TEsNotTE = !store.TEsNotTE">
              <v-icon>mdi-view-module</v-icon>
            </v-btn>
          </template>
          <span>{{ t('test.viewmode') }}</span>
        </v-tooltip>
      </template>
    </v-toolbar>

    <!-- TestSet 主窗口 -->
    <v-sheet class="ma-0 pt-10 overflow-y-auto"
      :height="store.mainWindowHeight - store.logHeight - store.toolbarheight">
      <test-entities v-if="store.TEsNotTE" :defcols="defCols" />
      <test-entity v-else />
    </v-sheet>

    <!-- TestLog 可扩展窗口 -->
    <v-sheet class="ma-0 pa-0 overflow-y-auto"
      :height="store.toolbarheight + store.logHeight">
      <test-log />
    </v-sheet>
  </v-card>
</template>

<script lang="ts" setup>
import { onUnmounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import TestEntities from '../components/TestEntities.vue'
import TestEntity from '../components/TestEntity.vue'
import TestLog from '../components/TestLog.vue'
import AddEntity from '../components/forms/AddEntity.vue'
import { StopTestGroupSyncMonitor } from '../../wailsjs/go/imes/Api'

const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

// Toolbar 相关
const selectedProd = ref()
const selectedStage = ref()
const testgroupsrcnewer = ref(false)

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
const startallgroup = () => {
  if (store.TEsNotTE) {
    console.log('run all entity\'s all group')
  } else {
    console.log('run activity entity\'s all group')
    window.runtime.EventsEmit('startallgroup')
  }
}
const stopallgroup = () => {
  if (store.TEsNotTE) {
    console.log('stop all entity\'s all group')
  } else {
    console.log('stop activity entity\'s all group')
  }
}
window.runtime.EventsOn('testgroupmonitor', (x: string) => {
  if (x == 'srcnewer') {
    testgroupsrcnewer.value = true
    store.LoadTestGroup('config', true)
  } else {
    testgroupsrcnewer.value = false
  }
})

onUnmounted(() => {
  StopTestGroupSyncMonitor()
})
</script>

<style>
.entity-toolbar {
  top: 0;
  position: absolute;
  width: 100%;
  /* opacity: 0.95; */
  z-index: 10000;
}
</style>
