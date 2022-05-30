<template>
  <v-container class="fill-height width-100">
    <v-row justify  ="center">
      <v-col v-for="entity in store.testEntities" :key="entity.ip.toString()"
        :cols="defcols">
        <v-card :elevation="5" @click="onclickEntity(entity.ip.toString())"
          :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
          <!-- <v-card-avatar></v-card-avatar> -->
          <template v-slot:title>{{ entity.ip.toString().replaceAll(',', '.')
          }}</template>
          <template v-slot:subtitle>code: {{ entity.code }}<br />tags:{{
              entity.tags
          }}</template>
          <template v-slot:text>
            <v-row>
              <v-col cols="12">
                状态：测试中……
              </v-col>
            </v-row>
          </template>
          <v-card-actions>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { useBaseStore } from '../stores/index'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

const props = withDefaults(
  defineProps<{
    defcols: number,
  }>(),
  {
    defcols: 3
  }
)

const onclickEntity = (ip) => {
  console.log(ip)
  store.activedTestEntityIp = ip
  store.TEorTI = false
}


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
