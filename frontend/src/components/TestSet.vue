<template>
  <v-container>
    <v-tabs class="sticky" centered v-model="activeTab" color="deep-purple-accent-4">
      <v-tab v-for="e in store.testEntities" :key="e.ip.toString()"
        :value="e.ip.toString()"> {{ e.ip.toString().replaceAll(',', '.') }}
      </v-tab>
    </v-tabs>
    <v-window v-model="activeTab">
      <v-window-item v-for="e in store.testEntities" :key="e.ip.toString()"
        :value="e.ip.toString()">
        <v-container class="fill-height width-100 mt-10">
          <v-row>
            <v-col v-for="tg in store.testGroups" :key="tg.title"
              :cols="12 / store.testGroups.length">
              <test-group v-model:list="tg.testclasses" axis="y" group="tg"
                :distance="10" helper-class="slicksort-helper" :tg="tg">
                <test-class v-for="(tc, i) in tg.testclasses" :key="tc.id" :index="i"
                  :tc="tc" />
              </test-group>
            </v-col>
          </v-row>
        </v-container>
      </v-window-item>
    </v-window>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted, } from 'vue'
import { useBaseStore } from '../stores/index'
import TestGroup from './testset/TestGroup.vue'
import TestClass from './testset/TestClass.vue'

const store = useBaseStore()
const activeTab = ref(store.activedTestEntityIp)
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
  store.syncTestSet()
  console.log(store.testGroups)
})
</script>

<style>
</style>