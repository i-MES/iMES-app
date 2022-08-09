<template>
  <v-row dense :class="store.enableTGWrap ? '' : 'flex-nowrap'">
    <slot />
  </v-row>
</template>

<script lang="ts">
import { ContainerMixin } from 'vue-slicksort'
export default {
  mixins: [ContainerMixin],
}
</script>
<script lang="ts" setup>
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()
</script>
<style lang="scss" scoped>
</style>
