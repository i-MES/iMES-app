<template>
  <v-row>
    <slot />
  </v-row>
</template>

<script >
import { ContainerMixin } from 'vue-slicksort'
export default {
  mixins: [ContainerMixin],
}
</script>

<style lang="scss" scoped>
</style>
