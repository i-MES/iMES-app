<template>
  <v-container class="fill-height width-100">
    <v-row>
      <v-col v-for="entity in store.testEntities" :key="entity.id" cols="6" md="4"
        lg="3" xl="2">
        <v-card :elevation="5" @click="onclickEntity(entity.id)"
          :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
          <!-- <v-card-avatar></v-card-avatar> -->
          <template v-slot:title>{{ entity.ip }}</template>
          <template v-slot:subtitle>{{ entity.title }} - {{ entity.desc }}</template>
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
      <v-col cols="6" md="4" lg="3" xl="2">
        <AddEntity />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { useBaseStore } from '../stores/index'
import { useI18n } from 'vue-i18n'
import AddEntity from './public/AddEntity.vue'
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()
const onclickEntity = (id) => {
  console.log(id)
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
