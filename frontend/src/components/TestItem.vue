<template>
  <v-tabs class="sticky" centered v-model="activeTab" color="deep-purple-accent-4">
    <v-tab v-for="e in store.testEntities" :key="e.ip" :value="e.ip"> {{ e.ip }}
    </v-tab>
  </v-tabs>
  <v-window v-model="activeTab">
    <v-window-item v-for="e in store.testEntities" :key="e.ip" :value="e.ip">
      <v-container class="fill-height width-100 mt-10">
        <v-row>
          <v-col
            v-for="tg in [{ id: 1, title: 'tg1' }, { id: 2, title: 'tg2' }, { id: 3, title: 'tg3' }, { id: 4, title: 'tg4' }]"
            :key="tg.id" cols="4">
            <v-sheet
              :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
              <v-toolbar height="20">
                <v-spacer></v-spacer>
                <v-btn @click="starttestgroup(tg.id)"
                  icon="mdi-arrow-right-bold-circle-outline"> </v-btn>
                <v-btn @click="stoptestgroup(tg.id)" icon="mdi-stop-circle-outline">
                </v-btn>
              </v-toolbar>
              <v-expansion-panels class="px-2 mt-2" multiple>
                <v-expansion-panel v-for="ti in store.testitems" :key="ti.id"
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
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, onBeforeUnmount, defineProps, reactive } from 'vue'
import { useBaseStore } from '../stores/index'
import { imes } from '../../wailsjs/go/models'
import { GetActivedTestEntity } from '../../wailsjs/go/imes/Api'
const store = useBaseStore()

const activeTab = ref(1)
const props = withDefaults(
  defineProps<{
    entityId: number,
  }>(),
  {
    entityId: 1
  }
)

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

const starttestgroup = (id: number) => {
  console.log('tg-id: ', id)
}
const stoptestgroup = (id: number) => {
  console.log('tg-id: ', id)
}
</script>

<style>
</style>