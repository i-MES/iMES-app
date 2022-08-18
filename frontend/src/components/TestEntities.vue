<template>
  <v-container class="mt-10">
    <v-row justify="center">
      <v-col v-for="e in store.testEntities" :key="e.ip.toString()" :cols="defcols">
        <v-card :elevation="5" @click="onclickEntity(e.id)"
          :color="store.appTheme == 'dark' ? store.darkmaincolor : store.lightmaincolor">
          <!-- <v-card-avatar></v-card-avatar> -->
          <template v-slot:title>{{ e.ip.toString().replaceAll(',', '.') }}
          </template>
          <template v-slot:subtitle>code: {{ e.code }}<br />tags:{{
              e.tags
          }}</template>
          <template v-slot:text>
            <v-row>
              <v-col cols="12"
                :class="timer && (store.tisNum(e.id).started > 0) ? 'animate__animated animate__flash' : ''">
                <process-bar
                  :values="[store.tisNum(e.id).pass, store.tisNum(e.id).ng, (store.testItemCounter - store.tisNum(e.id).pass - store.tisNum(e.id).ng)]" />
              </v-col>
            </v-row>
          </template>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useBaseStore } from '../stores/index'
import ProcessBar from './widget/ProcessBar.vue'
const store = useBaseStore()

withDefaults(
  defineProps<{
    defcols: number,
  }>(),
  {
    defcols: 3
  }
)

const onclickEntity = (id: string) => {
  console.log(id)
  store.activedTestEntityId = id
  store.TEsNotTE = false
}
const timer = ref(true)
setInterval(() => {
  timer.value = !timer.value
}, 3000)

</script>

<style>
</style>
