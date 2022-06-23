<template>
  <v-row justify="center">
    <v-col cols="3">
      <v-card>
        <v-card-text class="my-4 text-center text-h6">
          Draggable 1
        </v-card-text>
        <draggable class="dragArea list-group" :list="list1" :clone="clone"
          :group="{ name: 'people', pull: pullFunction }" @start="start" item-key="id">
          <template #item="{ element }">
            <div class="list-group-item">
              {{ element.name }}
            </div>
          </template>
        </draggable>
      </v-card>
    </v-col>

    <v-col cols="3">
      <v-card>
        <v-card-text class="my-4 text-center text-h6">
          Draggable 2
        </v-card-text>
        <draggable class="dragArea list-group" :list="list2" group="people"
          :clone="clone" @start="start" item-key="id">
          <template #item="{ element }">
            <div class="list-group-item">
              {{ element.name }}
            </div>
          </template>
        </draggable>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import draggable from 'vuedraggable'
let idGlobal = 8
const list1 = reactive([
  { name: 'Jesus', id: 1 },
  { name: 'Paul', id: 2 },
  { name: 'Peter', id: 3 }
])
const list2 = reactive([
  { name: 'Luc', id: 5 },
  { name: 'Thomas', id: 6 },
  { name: 'John', id: 7 }
])
const controlOnStart = ref(true)
const clone = ({ name }) => {
  return { name, id: idGlobal++ }
}
const pullFunction = () => {
  return this.controlOnStart.value ? 'clone' : true
}
const start = ({ originalEvent }) => {
  this.controlOnStart.value = originalEvent.ctrlKey
}
</script>
<style scoped>
.list-group {
  min-height: 20px;
}

.list-group-item {
  cursor: move;
  min-height: 20px;
}

.list-group-item i {
  cursor: pointer;
}
</style>