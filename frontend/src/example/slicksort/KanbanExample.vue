<template>
  <v-container>
    <v-toolbar class="mb-5" title="基于 SlickList & SlickItem"></v-toolbar>

    <SlickList v-model:list="kanban.columns" axis="x" lock-axis="x"
      class="column-container" use-drag-handle useWindowAsScrollContainer
      helper-class="slicksort-helper">
      <SlickItem v-for="(col, i) in kanban.columns" :key="col.id" :index="i">
        <v-card class="mx-2 py-3 px-1 elevation-5">
          <DragHandle>{{ col.name }}({{ col.items.length }})</DragHandle>
          <SlickList v-model:list="col.items" axis="y" lock-axis="y" :group="col.group"
            @sort-start="onSortStart($event)" helper-class="slicksort-helper">
            <SlickItem v-for="(item, j) in col.items" :key="item.id" :index="j">
              <v-card class="my-2 elevation-10" :title="item.value">
              </v-card>
            </SlickItem>
          </SlickList>
        </v-card>
      </SlickItem>
    </SlickList>
  </v-container>
</template>

<script>
import { reactive } from 'vue'
import { SlickList, SlickItem } from 'vue-slicksort'
import DragHandle from './components/DragHandle.vue'
import { stringsToItems } from './utils'

export default {
  components: {
    SlickList,
    SlickItem,
    DragHandle,
  },
  setup() {
    const kanban = reactive({
      name: 'Kanban Example',
      columns: [
        {
          id: 'backlog',
          name: 'backlog',
          group: 'items',
          items: stringsToItems(['Vue 4 support', 'TypeScript 2 support']),
        },
        {
          id: 'todos',
          name: 'To Do',
          group: 'items',
          items: stringsToItems(['Cancel dragging', 'Multiselect', 'Accessibility', 'UI Tests', 'Keyboard navigation']),
        },
        {
          id: 'inprogress',
          name: 'In Progress',
          group: 'items',
          items: stringsToItems(['Fix drag settling jankiness']),
        },
        {
          id: 'review',
          name: 'Review',
          group: 'items',
          items: stringsToItems([]),
        },
        {
          id: 'done',
          name: 'Done',
          group: 'items',
          items: stringsToItems(['Vue 3 support', 'TypeScript support']),
        },
      ],
    })

    const onSortStart = (event) => {
      console.log(event)
    }

    return {
      kanban,
      onSortStart
    }
  },
}
</script>

<style lang="scss" scoped>
.column-container {
  display: flex;
  align-items: start;
}

// .kanban-column {
//   width: 400px;
//   margin: 5px;
//   padding: 5px;
//   // background: #eee;
// }

// .kanban-list {
//   max-height: 500px;
//   overflow: auto;
//   display: flex;
//   flex-wrap: wrap;
// }

// @media screen and (max-width: 768px) {
//   .column-container {
//     flex-direction: column;
//     align-items: stretch;
//   }

//   .kanban-column {
//     width: auto;
//   }
// }


.kanban-list-item {
  width: calc(50% - 10px);
  margin-top: 5px;

  .kanban-list-item-inner {
    min-height: 60px;
    // padding: 10px 15px;
    // border-radius: 10px;
    // background: white;
    box-shadow: 0 2px 0 rgba(0, 0, 0, 0.1);
    cursor: grab;
    transition: background 0.2s, transform 0.2s;
  }

  &.kanban-helper .kanban-list-item-inner {
    transform: rotate(10deg);
    background: #9b51e0;
    color: white;
  }
}
</style>
