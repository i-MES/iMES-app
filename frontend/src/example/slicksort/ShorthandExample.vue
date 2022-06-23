<template>
  <v-container>
    <v-toolbar class="mb-5" title="基于 SlickList"></v-toolbar>
    <SlickList axis="y" v-model:list="fruits" helper-class="slicksort-helper">
      <template #item="{ item }">
        <v-card class="my-1 elevation-10" :title="item">
        </v-card>
      </template>
    </SlickList>
  </v-container>
</template>

<script>
import { SlickList } from 'vue-slicksort'

export default {
  components: {
    SlickList,
  },
  data() {
    return {
      fruits: ['Apples', 'Bananas', 'Cherries', 'Dragon Fruit', 'Elderberries'],
    }
  },
}
</script>
