<template>
  <v-icon start icon="mdi-cursor-move" size="small"></v-icon>
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
