<template>
  <v-container>
    <v-toolbar>SimpleGroupExample</v-toolbar>
    <v-row class="mt-5">
      <v-col cols="6">
        <div class="list-wrapper">
          <h4>{{ shelf.name }}</h4>
          <code>group: 'groceries'</code>
          <SortableList axis="y" group="groceries" :accept="shelf.accept"
            helper-class="slicksort-helper" :block="shelf.block"
            v-model:list="shelf.items">
            <SortableItem v-for="(item, index) in shelf.items" :key="index"
              :index="index" :item="item" />
          </SortableList>
        </div>
      </v-col>
      <v-col cols="6">
        <div class="list-wrapper">
          <h4>{{ cart.name }}</h4>
          <code>group: 'groceries'</code>
          <SortableList axis="y" group="groceries" :accept="cart.accept"
            :block="cart.block" helper-class="slicksort-helper"
            v-model:list="cart.items">
            <SortableItem v-for="(item, index) in cart.items" :key="index"
              :index="index" :item="item" />
          </SortableList>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ref } from 'vue'
import { random, } from './utils'
import SortableList from './components/SortableList.vue'
import SortableItem from './components/SortableItem.vue'

let id = 100

const colors = ['#eb5757', '#9b51e1', '#58cbf2']
const fruits = ['Apples', 'Bananas', 'Cherries', 'Dragon Fruit']
const veggies = ['Potatoes', 'Broccoli']

export default {
  name: 'SimpleGroupExample',
  components: {
    SortableItem,
    SortableList,
  },
  setup() {
    const shelf = ref({
      id: id++,
      name: 'Shelf',
      items: fruits.map((value) => {
        return {
          value,
          height: random(49, 100),
          background: colors[0],
          id: id++,
        }
      }),
    })

    const cart = ref({
      id: id++,
      name: 'Cart',
      items: veggies.map((value) => {
        return {
          value,
          height: random(49, 120),
          background: colors[1],
          id: id++,
        }
      }),
    })

    return {
      shelf,
      cart,
    }
  },
}
</script>

<style scoped>
.groups-example {
  display: flex;
  position: relative;
}
</style>
