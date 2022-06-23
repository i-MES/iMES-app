<template>
  <v-container>
    <v-toolbar>GroupExample（排好分组后获胜）</v-toolbar>
    <v-row class="mt-5">
      <div class="win-overlay" v-if="showWinScreen">
        <div>
          <h1>Winner!</h1>
          <button class="button" @click="resetList">Reset</button>
        </div>
      </div>
      <div v-for="list in lists" :key="list.id" class="list-wrapper">
        <v-card class="mx-1">
          <v-card-title>{{ list.name }}</v-card-title>
          <v-card-text>
            <code>group: '{{ list.group }}'</code>
            <br />
            <code>accept: {{ list.accept }}</code>
          </v-card-text>
          <v-card-content>
            <SortableList axis="y" :group="list.group" :accept="list.accept"
              helper-class="slicksort-helper" :block="list.block"
              v-model:list="list.items">
              <SortableItem v-for="(item, index) in list.items" :key="item.id"
                :index="index" :item="item" />
            </SortableList>
          </v-card-content>
        </v-card>
      </div>
    </v-row>
  </v-container>
</template>

<script>
import { ref, watch } from 'vue'
import { random, range } from './utils'
import SortableList from './components/SortableList.vue'
import SortableItem from './components/SortableItem.vue'

let id = 100
const colors = ['#eb5757', '#9b51e1', '#58cbf2']
const makeList = () => {
  return [
    {
      id: id++,
      name: 'List A',
      group: 'a',
      accept: ['b'],
      items: range(3).map((value) => {
        return {
          value: 'Item ' + (value + 1),
          height: random(49, 100),
          background: colors[value],
          id: id++,
        }
      }),
    },
    {
      id: id++,
      name: 'List B',
      group: 'b',
      accept: true,
      items: range(3).map((value) => {
        return {
          value: 'Item ' + (value + 1),
          height: random(49, 120),
          background: colors[value],
          id: id++,
        }
      }),
    },
    {
      id: id++,
      name: 'List C',
      group: 'c',
      accept: ['b'],
      items: range(3).map((value) => {
        return {
          value: 'Item ' + (value + 1),
          height: random(49, 120),
          background: colors[value],
          id: id++,
        }
      }),
    },
  ]
}

export default {
  name: 'GroupExample',
  components: {
    SortableItem,
    SortableList,
  },
  setup() {
    const lists = ref(makeList())

    const showWinScreen = ref(false)

    const resetList = () => {
      lists.value = makeList()
      showWinScreen.value = false
    }

    watch(
      () => lists.value,
      (newValue) => {
        const mapped = newValue.map((l) => l.items.map((i) => i.value))

        const winning = [
          ['Item 1', 'Item 1', 'Item 1'],
          ['Item 2', 'Item 2', 'Item 2'],
          ['Item 3', 'Item 3', 'Item 3'],
        ]

        if (JSON.stringify(mapped) === JSON.stringify(winning)) {
          showWinScreen.value = true
        }
      },
      {
        deep: true,
      }
    )

    return {
      lists,
      showWinScreen,
      resetList,
    }
  },
}
</script>

<style scoped>
.groups-example {
  display: flex;
  position: relative;
}

.list-wrapper {
  width: 33%;
}

.win-overlay {
  position: absolute;
  z-index: 1;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}
</style>
