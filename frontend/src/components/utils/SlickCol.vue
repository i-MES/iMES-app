<template>
  <v-col :style="'width:' + colspercent + '%'">
    <slot />
  </v-col>
</template>

<script lang="ts">
import { ElementMixin } from 'vue-slicksort'

export default {
  mixins: [ElementMixin],
  props: {
    colspercent: {
      type: Number,
      required: true,
    },
  },
}
</script>
<style lang="scss" scoped>
</style>
