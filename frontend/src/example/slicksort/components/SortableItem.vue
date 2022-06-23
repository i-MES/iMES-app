<template>
  <li class="example-list-item" z-index="9999"
    :style="{ height: item.height + 'px', background: item.background }">
    <DragHandle v-if="showHandle" />
    {{ item.value }}
    <button v-if="showHandle" class="button">Add to cart</button>
  </li>
</template>

<script>
import { ElementMixin } from 'vue-slicksort'
import DragHandle from './DragHandle.vue'

export default {
  mixins: [ElementMixin],
  components: {
    DragHandle,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    showHandle: Boolean,
  },
}
</script>

<style lang="scss" scoped>
.button {
  margin-left: auto;
}

.example-list-item {
  list-style-type: none;
  display: flex;
  align-items: center;
  padding: 10px 20px;
  margin: 10px;
  border-radius: 10px;
  font-weight: bold;
  background: white;
  box-shadow: inset 0 0 0 3px rgba(0, 0, 0, 0.1), 1px 2px 5px rgba(0, 0, 0, 0.15);
  background: #9b51e0;
  color: white;
  line-height: 1.4;
  user-select: none;
}
</style>
