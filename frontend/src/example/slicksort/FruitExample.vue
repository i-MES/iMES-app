<template>
  <v-container>
    <v-toolbar>FruitExample (axis:{{ axis }}, showHandle:{{ showHandle }})</v-toolbar>
    <SortableList axis="y" v-model:list="fruits" class="example-list pa-4"
      helper-class="slicksort-helper" :class="{ horizontal: axis === 'x' }">
      <SortableItem v-for="(fruit, i) in fruits" :key="fruit" :index="i"
        :item="{ value: fruit }" :showHandle="showHandle" />
    </SortableList>
    <v-card title="Fruits">
      <v-card-content>
        {{ fruits }}
      </v-card-content>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, withDefaults, defineProps } from 'vue'
import SortableList from './components/SortableList.vue'
import SortableItem from './components/SortableItem.vue'

// 使用 ref 而不是 reactive，官方 FAQ：
// https://vue-slicksort.netlify.app/troubleshooting.html#v-model-is-not-working
const fruits = ref(['Apples', 'Bananas', 'Cherries', 'Dragon Fruit', 'Elderberries'])

withDefaults(
  defineProps<{
    axis: string,
    showHandle: boolean,
  }>(),
  {
    axis: 'y',
    showHandle: false,
  }
)
</script>

<!-- 
<script>
import SortableList from './components/SortableList.vue'
import SortableItem from './components/SortableItem.vue'

export default {
  components: {
    SortableList,
    SortableItem,
  },
  props: {
    axis: {
      type: String,
      default: 'y',
    },
    showHandle: Boolean,
  },
  data() {
    return {
      fruits: ['Apples', 'Bananas', 'Cherries', 'Dragon Fruit', 'Elderberries'],
    }
  },
}
</script> -->

<style scoped>
</style>