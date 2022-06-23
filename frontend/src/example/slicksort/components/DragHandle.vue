<template>
  <v-btn size="small" variant="plain" prepend-icon="mdi-dialpad">
    <slot></slot>
  </v-btn>
</template>

<script>
// Makeshift handle because the directive doesnt work with vitepress
export default {
  mounted() {
    this.$el.sortableHandle = true
  },
}
</script>

<style scoped>
</style>
