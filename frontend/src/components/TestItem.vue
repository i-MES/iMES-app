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
            <v-col v-for="tg in store.testGroups" :key="tg.id" cols="4">
              <v-sheet
                :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
                <v-toolbar height="20">
                  <v-toolbar-title>{{ tg.title }}</v-toolbar-title>
                  <v-btn @click="starttestgroup(tg.id)"
                    icon="mdi-arrow-right-bold-circle-outline"> </v-btn>
                  <v-btn @click="stoptestgroup(tg.id)" icon="mdi-stop-circle-outline">
                  </v-btn>
                </v-toolbar>
                <v-expansion-panels class="px-2 mt-2" multiple>
                  <v-expansion-panel v-for="ti in tg.testItems" :key="ti.title"
                    :title="ti.title" :text="ti.desc">
                    <div>
                      <v-progress-linear v-model="value" :buffer-value="bufferValue">
                      </v-progress-linear>
                    </div>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-sheet>
            </v-col>
          </v-row>
        </v-container>
      </v-window-item>
    </v-window>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { useBaseStore } from '../stores/index'
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

onMounted(() => {
  store.syncTestSet()
})
onBeforeUnmount(() => {
  clearInterval(interval)
})

const starttestgroup = (id: number) => {
  console.log('tg-id: ', id)
}
const stoptestgroup = (id: number) => {
  console.log('tg-id: ', id)
}
</script>

<style>
</style>