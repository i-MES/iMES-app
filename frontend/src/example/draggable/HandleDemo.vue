<template>
  <v-row justify="center">
    <v-col cols="6">
      <div>
        <draggable tag="ul" :list="list" class="list-group" handle=".handle"
          item-key="name">
          <template #item="{ element }">
            <li class="list-group-item">
              <v-btn class=" handle"></v-btn>
              <v-input type="text" class="form-control" v-model="element.text">{{
                  element.name
              }}
              </v-input>
            </li>
          </template>
        </draggable>
      </div>
    </v-col>
  </v-row>

</template>

<script>
let id = 3
import draggable from 'vuedraggable'
export default {
  name: 'HandleDemo',
  display: 'Handle',
  instruction: 'Drag using the handle icon',
  order: 5,
  components: {
    draggable
  },
  data() {
    return {
      list: [
        { name: 'John', text: '', id: 0 },
        { name: 'Joao', text: '', id: 1 },
        { name: 'Jean', text: '', id: 2 }
      ],
      dragging: false
    }
  },
  computed: {
    draggingInfo() {
      return this.dragging ? 'under drag' : ''
    }
  },
  methods: {
    removeAt(idx) {
      this.list.splice(idx, 1)
    },
    add: function () {
      id++
      this.list.push({ name: 'Juan ' + id, id, text: '' })
    }
  }
}
</script>
<style scoped>
.button {
  margin-top: 35px;
}

.handle {
  float: left;
  padding-top: 8px;
  padding-bottom: 8px;
}

.close {
  float: right;
  padding-top: 8px;
  padding-bottom: 8px;
}

input {
  display: inline-block;
  width: 50%;
}

.text {
  margin: 20px;
}
</style>
