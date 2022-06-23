<template>
  <ul class="example-list">
    <slot />
  </ul>
</template>

<script >
import { ContainerMixin } from 'vue-slicksort'
export default {
  mixins: [ContainerMixin],
}
</script>

<style lang="scss" scoped>
.example-list {
  margin: 0;
  padding: 0;
  overflow: auto;
}
</style>
