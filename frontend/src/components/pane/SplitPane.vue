<template>
  <div ref="splitPane" class="split-pane" :class="direction"
    :style="{ flexDirection: direction as 'column', height: totalHeight + 'px' }">

    <div class="pane pane-first" :style="lengthType + ':' + paneOneLengthValue">
      <slot name="first"></slot>
    </div>

    <div class="pane-trigger" :style="lengthType + ':' + triggerLengthValue"
      @mousedown="handleMouseDown"></div>

    <div class="pane pane-second overflow-y-auto">
      <slot name="second"></slot>
    </div>

  </div>
</template>

<script lang="ts" setup>
import { ref, computed, withDefaults } from 'vue'

const props = withDefaults(
  defineProps<{
    totalHeight: number,
    direction: string,
    min?: number,
    max?: number,
    paneFirstLengthPercent?: number,
    triggerLength?: number,
  }>(),
  {
    totalHeight: 300,
    direction: 'row',
    min: 10,
    max: 90,
    paneFirstLengthPercent: 80,
    triggerLength: 4,
  }
)

const triggerLeftOffset = ref(0) // 鼠标距滑动器左(顶)侧偏移量
const splitPane = ref(null)
const lengthType = computed(() => {
  return props.direction === 'row' ? 'width' : 'height'
})

const paneOneLengthValue = computed(() => {
  return `calc(${props.paneFirstLengthPercent}% - ${props.triggerLength / 2 + 'px'})`
})

const triggerLengthValue = computed(() => {
  return props.triggerLength + 'px'
})

// 按下滑动器
const handleMouseDown = (e) => {
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)

  if (props.direction === 'row') {
    triggerLeftOffset.value = e.pageX - e.srcElement.getBoundingClientRect().left
  } else {
    triggerLeftOffset.value = e.pageY - e.srcElement.getBoundingClientRect().top
  }
}

const emit = defineEmits(['update:paneFirstLengthPercent'])
// 按下滑动器后移动鼠标
const handleMouseMove = (e) => {
  // dir() 可以打印结构化对象
  // console.dir(splitPane.value)
  // 从上面打印可以找到很多可用的属性
  // console.log('--------', splitPane.value.clientHeight)
  // element 的可用函数，参考：https://developer.mozilla.org/zh-CN/docs/Web/API/Element/getBoundingClientRect
  const clientRect = splitPane.value.getBoundingClientRect()
  let paneFirstLengthPercent = 0
  if (props.direction === 'row') {
    const offset = e.pageX - clientRect.left - triggerLeftOffset.value + props.triggerLength / 2
    paneFirstLengthPercent = (offset / clientRect.width) * 100
  } else {
    const offset = e.pageY - clientRect.top - triggerLeftOffset.value + props.triggerLength / 2
    paneFirstLengthPercent = (offset / clientRect.height) * 100
  }

  if (paneFirstLengthPercent < props.min) {
    paneFirstLengthPercent = props.min
  }
  if (paneFirstLengthPercent > props.max) {
    paneFirstLengthPercent = props.max
  }
  emit('update:paneFirstLengthPercent', paneFirstLengthPercent)
}

// 松开滑动器
const handleMouseUp = () => {
  document.removeEventListener('mousemove', handleMouseMove)
}
</script>

<style scoped lang="scss">
.split-pane {
  height: 100%;
  display: flex;

  &.row {
    .pane {
      height: 100%;
    }

    .pane-trigger {
      height: 100%;
      cursor: col-resize;
    }
  }

  &.column {
    .pane {
      width: 100%;
    }

    .pane-trigger {
      width: 100%;
      cursor: row-resize;
    }
  }

  .pane-first {
    // background: palevioletred;
  }

  .pane-trigger {
    user-select: none;
    background: #888888;
  }

  .pane-second {
    flex: 1;
    // background: turquoise;
  }
}
</style>