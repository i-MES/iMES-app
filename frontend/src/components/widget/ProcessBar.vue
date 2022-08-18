<template>
  <div class="PB_Container"
    :style="'height:' + height + 'px; line-height:' + height + 'px;'">
    <div v-for="(v, i) in values" :key="i"
      :style="'width:' + (100 * v / total) + '%; background-color:' + colors[i]">
      {{ v == 0 ? '' : v }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
const props = withDefaults(
  defineProps<{
    values: number[],
    colors: string[],
    height: number
  }>(),
  {
    values: () => [0, 0, 0], // 数组被转换成函数了
    colors: () => ['#43A047', '#FF5252', 'grey'],
    height: 15
  })

const total = computed(() => {
  let t = 0
  props.values.forEach((v) => t += v)
  return t
})
</script> 

<style scoped>
.PB_Container {
  width: 100%;
  display: inline-flex;
  text-align: center;
  overflow: hidden;
}

.PB_Container div {
  float: left;
}
</style>
