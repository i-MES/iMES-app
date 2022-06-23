<template>
  <v-container>
    <v-toolbar>LongListExample</v-toolbar>
    <SortableList v-model:list="items" axis="y" class="long-list"
      helper-class="slicksort-helper">
      <SortableItem v-for="(item, index) in items" :key="item.id" :index="index"
        :item="item" />
    </SortableList>
  </v-container>
</template>

<script>
import { random, range } from './utils'
import SortableList from './components/SortableList.vue'
import SortableItem from './components/SortableItem.vue'

export default {
  components: {
    SortableItem,
    SortableList,
  },
  data() {
    return {
      items: range(10).map((i) => {
        return {
          value: `Item ${i}`,
          id: i,
          height: random(50, 100),
        }
      }),
    }
  },
}
</script>

<style scoped>
.long-list {
  max-height: 300px;
  overflow: auto;
}
</style>
