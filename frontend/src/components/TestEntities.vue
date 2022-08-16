<template>
  <v-container class="mt-10">
    <v-row justify="center">
      <v-col v-for="entity in store.testEntities" :key="entity.ip.toString()"
        :cols="defcols">
        <v-card :elevation="5" @click="onclickEntity(entity.id)"
          :color="store.appTheme == 'dark' ? store.darkmaincolor : store.lightmaincolor">
          <!-- <v-card-avatar></v-card-avatar> -->
          <template v-slot:title>{{ entity.ip.toString().replaceAll(',', '.') }}
          </template>
          <template v-slot:subtitle>code: {{ entity.code }}<br />tags:{{
              entity.tags
          }}</template>
          <template v-slot:text>
            <v-row>
              <!-- <v-col cols="12">
                <div
                  :class="timer && (store.testItemStatusCounter(entity.id).started > 0) ? 'animate__animated animate__flash' : ''">
                  total: {{ store.testItemCounter }}<br />
                  started: {{ store.testItemStatusCounter(entity.id).started }}<br />
                  PASS: {{ store.testItemStatusCounter(entity.id).pass }}<br />
                  NG: {{ store.testItemStatusCounter(entity.id).ng }}<br />
                </div>
              </v-col> -->
              <v-col cols="12"
                :class="timer && (store.testItemStatusCounter(entity.id).started > 0) ? 'animate__animated animate__flash' : ''">
                <v-progress-linear v-model="store.testItemStatusCounter(entity.id).pass"
                  :max="store.testItemCounter" color="green" height="15">
                  <template v-slot:default="{}">
                    <strong>{{ store.testItemStatusCounter(entity.id).pass }}
                      /{{ store.testItemCounter }}</strong>
                  </template>
                </v-progress-linear>
                <v-progress-linear v-model="store.testItemStatusCounter(entity.id).ng"
                  reverse :max="store.testItemCounter" color="red" height="15">
                  <template v-slot:default="{}">
                    <strong>{{ store.testItemStatusCounter(entity.id).ng }}
                      /{{ store.testItemCounter }}</strong>
                  </template>
                </v-progress-linear>
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
