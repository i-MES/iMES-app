<template>
  <v-container>
    <div class="echarts" id="echartsid"></div>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted} from 'vue'
import { useI18n } from 'vue-i18n'

import * as echarts from 'echarts/core'
import {
  BarChart,
  // 系列类型的定义后缀都为 SeriesOption
  BarSeriesOption,
  LineSeriesOption,
} from 'echarts/charts'
import {
  TitleComponent,
  // 组件类型的定义后缀都为 ComponentOption
  TitleComponentOption,
  TooltipComponent,
  TooltipComponentOption,
  GridComponent,
  GridComponentOption,
  // 数据集组件
  DatasetComponent,
  DatasetComponentOption,
  // 内置数据转换器组件 (filter, sort)
  TransformComponent,
} from 'echarts/components'
import { LabelLayout, UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'

// 通过 ComposeOption 来组合出一个只有必须组件和图表的 Option 类型
type ECOption = echarts.ComposeOption<
  | BarSeriesOption
  | LineSeriesOption
  | TitleComponentOption
  | TooltipComponentOption
  | GridComponentOption
  | DatasetComponentOption
>
echarts.use([
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  BarChart,
  LabelLayout,
  UniversalTransition,
  CanvasRenderer,
])

const { t } = useI18n({ useScope: 'global' })
console.log('-=-=')
onMounted(() => {
  var id = document.getElementById('echartsid')
  console.log(id)
  if (id) {
    var myChart = echarts.init(id)
    myChart.setOption({
      title: {
        text: '直通率',
      },
      tooltip: {},
      xAxis: {
        data: ['工位1', '工位2', '工位3', '工位4', '工位5', '工位6'],
      },
      yAxis: {},
      series: [
        {
          name: '销量',
          type: 'bar',
          data: [5, 20, 36, 10, 10, 20],
        },
      ],
    })
  }
})
</script>
<style>
.echarts {
  width: 350px;
  height: 300px;
}
</style>
