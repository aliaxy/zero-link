<template>
  <v-chart class="chart" :option="option" autoresize />
</template>

<script setup>
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps({
  items: { type: Array, default: () => [] },
})

const option = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: '#ffffff',
    borderColor: 'rgba(60,60,67,0.12)',
    borderWidth: 1,
    textStyle: { color: '#1d1d1f', fontSize: 13 },
    extraCssText: 'box-shadow: 0 4px 12px rgba(0,0,0,0.08); border-radius: 8px;',
  },
  legend: {
    top: 0,
    right: 0,
    textStyle: { color: '#6e6e73', fontSize: 12 },
    itemWidth: 14,
    itemHeight: 4,
    borderRadius: 2,
  },
  grid: { top: 40, right: 16, bottom: 24, left: 40, containLabel: true },
  xAxis: {
    type: 'category',
    data: props.items.map((i) => i.stat_date),
    axisLine: { lineStyle: { color: 'rgba(60,60,67,0.12)' } },
    axisTick: { show: false },
    axisLabel: { color: '#aeaeb2', fontSize: 11 },
  },
  yAxis: [
    {
      type: 'value',
      name: 'PV',
      nameTextStyle: { color: '#aeaeb2', fontSize: 11 },
      axisLabel: { color: '#aeaeb2', fontSize: 11 },
      splitLine: { lineStyle: { color: 'rgba(60,60,67,0.08)' } },
    },
    {
      type: 'value',
      name: 'UV',
      nameTextStyle: { color: '#aeaeb2', fontSize: 11 },
      axisLabel: { color: '#aeaeb2', fontSize: 11 },
      splitLine: { show: false },
    },
  ],
  series: [
    {
      name: 'Page Views',
      type: 'line',
      data: props.items.map((i) => i.pv),
      smooth: true,
      symbol: 'circle',
      symbolSize: 5,
      lineStyle: { color: '#007aff', width: 2 },
      itemStyle: { color: '#007aff' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(0,122,255,0.12)' },
            { offset: 1, color: 'rgba(0,122,255,0)' },
          ],
        },
      },
    },
    {
      name: 'Unique Visitors',
      type: 'line',
      yAxisIndex: 1,
      data: props.items.map((i) => i.uv),
      smooth: true,
      symbol: 'circle',
      symbolSize: 5,
      lineStyle: { color: '#34c759', width: 2 },
      itemStyle: { color: '#34c759' },
    },
  ],
}))
</script>

<style scoped>
.chart {
  width: 100%;
  height: 280px;
}
</style>
