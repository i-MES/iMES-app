<template>
  <v-col :cols="cols">
    <slot />
  </v-col>
</template>

<script lang="ts">
import { ElementMixin } from 'vue-slicksort'

export default {
  mixins: [ElementMixin],
  props: {
    cols: {
      type: Number,
      required: true,
    },
  },
}
</script>
<style lang="scss" scoped>
</style>
