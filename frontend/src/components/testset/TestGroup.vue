<template>
  <v-sheet onselectstart="return flase"
    :color="store.appTheme == 'dark' ? 'blue-grey-darken-2' : 'blue-grey-lighten-3'">
    <v-toolbar height="20">
      <v-toolbar-title>{{ tg.title }}</v-toolbar-title>
      <v-btn @click="starttestgroup(tg)" icon="mdi-arrow-right-bold-circle-outline">
      </v-btn>
      <v-btn @click="stoptestgroup(tg)" icon="mdi-stop-circle-outline">
      </v-btn>
    </v-toolbar>
    <slot />
  </v-sheet>
</template>

<script lang="ts">
import { ContainerMixin } from 'vue-slicksort'
export default {
  mixins: [ContainerMixin]
}
</script>

<script lang="ts" setup>
import { testset } from '../../../wailsjs/go/models'
import { TestGroupStart } from '../../../wailsjs/go/imes/Api'
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()
defineProps<{
  tg: {
    type: testset.TestGroup,
    required: true,
  },
}>()
const starttestgroup = (tg: testset.TestGroup) => {
  console.log(tg)
  TestGroupStart(tg)
  // store.testGroupById(id)?.testItems.forEach(
  //   async (ti) => {
  //     console.log(ti)
  //     await TestItemStart(ti)
  //   }
  // )
}
const stoptestgroup = (tg: testset.TestGroup) => {
  console.log(tg)
}
</script>

<style lang="scss" scoped>
</style>
