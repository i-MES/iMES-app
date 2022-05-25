<template>
  <v-toolbar class="entity-toolbar" :height="store.toolbarheight">
    <v-toolbar-title>{{ t('testpage.testentity-overall') }}</v-toolbar-title>
    <v-spacer></v-spacer>
    <v-text-field hide-details append-icon="mdi-magnify"> </v-text-field>
  </v-toolbar>
  <v-container class="fill-height width-100 mt-10">
    <v-row>
      <v-col v-for="entity in entities" :key="entity.id" cols="4">
        <v-card :elevation="5"
          :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
          <v-card-avatar></v-card-avatar>
          <template v-slot:title>{{ entity.title }}</template>
          <template v-slot:subtitle>{{ entity.desc }}</template>
          <template v-slot:text>{{ entity.id }}- {{ entity.title }} - {{ entity.desc
          }}</template>
          <v-card-actions>

          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import { useBaseStore } from '../stores/index'
import { useI18n } from 'vue-i18n'
import { imes } from '../../wailsjs/go/models'
import { GetActivedTestEntity, TestItemStart } from '../../wailsjs/go/imes/Api'

const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()
var entities = reactive<imes.TestEntity[]>([])

// const timer = setInterval(() => {
//   TestItemStart(1).then((val) => {
//     console.log('测试项启动：', val ? '成功' : '失败')
//   })
// }, 3000)

onMounted(() => {
  GetActivedTestEntity().then(
    (_entites) => {
      _entites.forEach((e) => entities.push(e))
      console.log(entities)
    }
  )
})
onUnmounted(() => {
  // clearInterval(timer)
})
</script>

<style>
.entity-toolbar {
  top: 0;
  position: absolute;
  width: 100%;
  /* opacity: 0.95; */
  z-index: 1;
}
</style>
