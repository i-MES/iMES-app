<template>
  <v-container>
    <!-- <v-tabs class="sticky" centered v-model="activeTab" color="deep-purple-accent-4">
      <v-tab v-for="e in store.testEntities" :key="e.ip.toString()"
        :value="e.ip.toString()"> {{ e.ip.toString().replaceAll(',', '.') }}
      </v-tab>
    </v-tabs>
    <v-window v-model="activeTab">
      <v-window-item v-for="te in store.testEntities" :key="te.ip.toString()"
        :value="te.ip.toString()"> -->
    <!-- <v-container class="fill-height width-100 mt-10"> -->
    <slick-row v-model:list="store.testGroups" axis="x" lock-axis="x" use-drag-handle
      useWindowAsScrollContainer helper-class="slicksort-helper"
      @update:list="store.SaveTestGroup" :accept="store.canSortTestClass">
      <slick-col v-for="(tg, idx) in store.testGroups" :key="tg.id" :index="idx"
        :cols="12 / store.testGroups.length" :disabled="!store.canSortTestClass">
        <slick-list v-model:list="tg.testclasses" axis="y" group="tg" :distance="10"
          helper-class="slicksort-helper" @update:list="store.SaveTestGroup"
          :accept="store.canSortTestClass">
          <test-group :tg="tg">
            <slick-item v-for="(tc, i) in tg.testclasses" :key="tc.id" :index="i"
              :disabled="!store.canSortTestClass">
              <test-class :teid="store.activedTestEntityId" :tgid="tg.id" :tc="tc" />
            </slick-item>
          </test-group>
        </slick-list>
      </slick-col>
    </slick-row>
    <!-- </v-container>
      </v-window-item>
    </v-window> -->
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, } from 'vue'
import { useBaseStore } from '../stores/index'
import TestGroup from './testset/TestGroup.vue'
import TestClass from './testset/TestClass.vue'
import { SlickList, SlickItem } from 'vue-slicksort'
import SlickRow from './utils/SlickRow.vue'
import SlickCol from './utils/SlickCol.vue'

const store = useBaseStore()
// const activeTab = ref(store.activedTestEntityId)
// const props = withDefaults(
//   defineProps<{
//     entityId: string,
//   }>(),
//   {
//     entityId: '127.0.0.1'
//   }
// )

onMounted(() => {
  // 加载 TestGroup、TestClass、TestItem 数据
  store.LoadTestGroup('config', true, false)
  console.log('TestEntity onMounted: ')
})
</script>

<style>
.column-container {
  display: flex;
  align-items: start;
}
</style>