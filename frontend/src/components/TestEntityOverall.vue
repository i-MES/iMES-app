<template>
  <v-hover v-slot="{ isHovering, props }" open-delay="600" close-delay="600">
    <v-card class="mx-0 entity-card" color="grey-lighten-4" v-bind="props">
      <v-toolbar :height="store.toolbarheight">
        <v-toolbar-title>{{ t('testpage.testentity-overall') }}</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn :icon="esticky ? 'mdi-pin' : 'mdi-pin-off'" @click="esticky = !esticky"> </v-btn>
      </v-toolbar>

      <v-row justify="space-around" v-if="isHovering || esticky">
        <v-col v-for="elevation in [1, 2, 3]" :key="elevation" cols="12" md="4">
          <v-sheet class="ma-5 pa-12" color="grey lighten-3">
            <v-sheet :elevation="elevation" class="mx-auto" height="100" width="100">
              <v-btn>Entity{{ elevation }}</v-btn>
            </v-sheet>
          </v-sheet>
        </v-col>
      </v-row>
      <v-row justify="space-around" v-if="isHovering || esticky">
        <v-col v-for="elevation in [1, 2, 3]" :key="elevation" cols="12" md="4">
          <v-sheet class="ma-5 pa-12" color="grey lighten-3">
            <v-sheet :elevation="elevation" class="mx-auto" height="100" width="100">
              <v-btn>Entity{{ elevation + 3 }}</v-btn>
            </v-sheet>
          </v-sheet>
        </v-col>
      </v-row>
    </v-card>
  </v-hover>
</template>

<script lang="ts" setup>
import { onUnmounted, ref } from 'vue'
import { useBaseStore } from '../stores/index'
import { useI18n } from 'vue-i18n'
import { TestItemStart } from '../../wailsjs/go/imes/Api'

const { t } = useI18n({ useScope: 'global' })

const store = useBaseStore()
const esticky = ref(false)
const entitywindow = ref()

const timer = setInterval(() => {
  TestItemStart(1).then((val) => {
    console.log('测试项启动：', val ? '成功' : '失败')
  })
}, 3000)

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style>
.entity-card {
  align-items: center;
  top: 0;
  justify-content: center;
  opacity: 0.95;
  position: absolute;
  width: 100%;
  z-index: 2;
}
</style>
