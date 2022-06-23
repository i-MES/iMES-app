<template>
  <v-container class="pa-0">
    <v-row class="mx-n12 my-2">
      <v-col v-for="(comp, i) in Comps" :key="i" cols="12">
        <v-card class="mx-10">
          <component :is="comp"></component>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
// import HandleDemo from '../example/draggable/handle.vue'
import FruitExample from '../example/slicksort/FruitExample.vue'
import GroupExample from '../example/slicksort/GroupExample.vue'
import SimpleGroupExample from '../example/slicksort/SimpleGroupExample.vue'
import LongListExample from '../example/slicksort/LongListExample.vue'
import ShorthandExample from '../example/slicksort/ShorthandExample.vue'
import KanbanExample from '../example/slicksort/KanbanExample.vue'

const Comps = [ShorthandExample, KanbanExample, FruitExample, SimpleGroupExample, GroupExample, LongListExample]
</script>