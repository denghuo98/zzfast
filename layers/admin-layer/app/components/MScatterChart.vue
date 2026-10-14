<!-- 散点图 -->
<template>
    <div ref="chartRef" class="relative w-full" :style="{ height: props.height }" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsOption } from '../utils/echarts'
import { useChartOps, useChartComponent } from '../composables/useChart'
import type { ScatterChartProps } from '../types/chart'

defineOptions({ name: 'MScatterChart' })

const props = withDefaults(defineProps<ScatterChartProps>(), {
  height: '16rem',
  loading: false,
  isEmpty: false,
  colors: () => undefined as string[] | undefined,
  data: () => [],
  symbolSize: 14,
  showAxisLabel: true,
  showAxisLine: true,
  showSplitLine: true,
  showTooltip: true,
  showLegend: false,
  legendPosition: 'bottom',
})

const resolvedColors = computed(() => props.colors?.length ? props.colors : useChartOps().colors)

const {
  chartRef, isDark,
  getAxisLineStyle, getAxisLabelStyle, getAxisTickStyle,
  getSplitLineStyle, getTooltipStyle, getLegendStyle, getGridWithLegend,
} = useChartComponent<ScatterChartProps>({
  props,
  checkEmpty: () => !props.data?.length,
  watchSources: [() => props.data],
  generateOptions: (): EChartsOption => ({
    grid: getGridWithLegend(!!props.showLegend, props.legendPosition, { top: 15, right: 15, left: 0 }),
    tooltip: props.showTooltip ? getTooltipStyle('item', { formatter: (p: any) => `X: ${p.value[0]}<br/>Y: ${p.value[1]}` }) : undefined,
    legend: props.showLegend ? getLegendStyle(props.legendPosition) : undefined,
    xAxis: {
      type: 'value',
      axisTick: getAxisTickStyle(),
      axisLine: getAxisLineStyle(props.showAxisLine),
      axisLabel: getAxisLabelStyle(props.showAxisLabel),
      splitLine: getSplitLineStyle(props.showSplitLine),
    },
    yAxis: {
      type: 'value',
      axisLabel: getAxisLabelStyle(props.showAxisLabel),
      axisLine: getAxisLineStyle(props.showAxisLine),
      splitLine: getSplitLineStyle(props.showSplitLine),
    },
    series: [{
      type: 'scatter', data: props.data, symbolSize: props.symbolSize,
      color: resolvedColors.value[0],
      itemStyle: {
        shadowBlur: 10,
        shadowColor: isDark.value ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)',
      },
    }],
  }),
})
</script>
