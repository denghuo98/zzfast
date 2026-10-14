<!-- 柱状图：支持单组/多组/堆叠，垂直渐变 -->
<template>
  <div ref="chartRef" class="relative w-full" :style="{ height: props.height }" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { graphic, type EChartsOption } from '../utils/echarts'
import { useChartOps, useChartComponent, getCssVar } from '../composables/useChart'
import type { BarChartProps, BarDataItem } from '../types/chart'

defineOptions({ name: 'MBarChart' })

const props = withDefaults(defineProps<BarChartProps>(), {
  height: '16rem',
  loading: false,
  isEmpty: false,
  colors: () => undefined as string[] | undefined,
  borderRadius: 4,
  data: () => [0, 0, 0, 0, 0, 0, 0],
  xAxisData: () => [],
  barWidth: '40%',
  stack: false,
  showAxisLabel: true,
  showAxisLine: true,
  showSplitLine: true,
  showTooltip: true,
  showLegend: false,
  legendPosition: 'bottom',
})

const resolvedColors = computed(() => props.colors?.length ? props.colors : useChartOps().colors)

const isMultipleData = computed(() =>
  Array.isArray(props.data) && props.data.length > 0 && typeof props.data[0] === 'object' && 'name' in props.data[0],
)

/** 是否为堆叠模式（stack=true 且有多组数据） */
const isStacked = computed(() => props.stack && isMultipleData.value)

const getColor = (customColor?: string, index?: number) => {
  if (customColor) return customColor
  if (index !== undefined) return resolvedColors.value[index % resolvedColors.value.length]
  return new graphic.LinearGradient(0, 0, 0, 1, [
    { offset: 0, color: getCssVar('--ui-primary') || '#00704A' },
    { offset: 1, color: getCssVar('--ui-primary') || '#005739' },
  ])
}

/**
 * 堆叠柱状图 borderRadius 策略：
 * - 最底层：仅底部圆角 [0, 0, br, br]
 * - 中间层：无圆角 0
 * - 最顶层：仅顶部圆角 [br, br, 0, 0]
 * 非堆叠模式下所有柱子四个角都保留 borderRadius。
 */
function stackBorderRadius(index: number, total: number): number | number[] {
  if (total <= 1) return props.borderRadius
  if (index === 0) return [0, 0, props.borderRadius, props.borderRadius]   // bottom
  if (index === total - 1) return [props.borderRadius, props.borderRadius, 0, 0] // top
  return 0 // middle
}

const createSeriesItem = (config: {
  name?: string; data: number[]; color?: any; barWidth?: string | number; stack?: string; borderRadius?: number | number[]
}) => ({
  name: config.name,
  data: config.data,
  type: 'bar' as const,
  stack: config.stack,
  itemStyle: {
    borderRadius: config.borderRadius ?? props.borderRadius,
    color: typeof config.color === 'string'
      ? new graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: config.color }, { offset: 1, color: config.color }])
      : config.color,
  },
  barWidth: config.barWidth || props.barWidth,
  ...getAnimationConfig(),
})

const {
  chartRef,
  getAxisLineStyle, getAxisLabelStyle, getAxisTickStyle,
  getSplitLineStyle, getAnimationConfig,
  getTooltipStyle, getLegendStyle, getGridWithLegend,
} = useChartComponent<BarChartProps>({
  props,
  checkEmpty: () => {
    if (Array.isArray(props.data) && typeof props.data[0] === 'number') {
      const sd = props.data as number[]
      return !sd.length || sd.every(v => v === 0)
    }
    if (Array.isArray(props.data) && typeof props.data[0] === 'object') {
      const md = props.data as BarDataItem[]
      return !md.length || md.every(item => !item.data?.length || item.data.every(v => v === 0))
    }
    return true
  },
  watchSources: [() => props.data, () => props.xAxisData, () => props.colors],
  generateOptions: (): EChartsOption => {
    const opt: EChartsOption = {
      grid: getGridWithLegend(!!props.showLegend && isMultipleData.value, props.legendPosition, { top: 15, right: 0, left: 0 }),
      tooltip: props.showTooltip ? getTooltipStyle() : undefined,
      xAxis: {
        type: 'category', data: props.xAxisData,
        axisTick: getAxisTickStyle(),
        axisLine: getAxisLineStyle(props.showAxisLine),
        axisLabel: getAxisLabelStyle(props.showAxisLabel),
      },
      yAxis: {
        type: 'value',
        axisLabel: getAxisLabelStyle(props.showAxisLabel),
        axisLine: getAxisLineStyle(props.showAxisLine),
        splitLine: getSplitLineStyle(props.showSplitLine),
      },
    }
    if (props.showLegend && isMultipleData.value) opt.legend = getLegendStyle(props.legendPosition)
    if (isMultipleData.value) {
      const md = props.data as BarDataItem[]
      opt.series = md.map((item, i) => createSeriesItem({
        name: item.name, data: item.data,
        color: props.colors?.[i] ?? getColor(undefined, i),
        barWidth: item.barWidth,
        stack: props.stack ? item.stack || 'total' : undefined,
        borderRadius: isStacked.value ? stackBorderRadius(i, md.length) : undefined,
      }))
    } else {
      opt.series = [createSeriesItem({ data: props.data as number[], color: getColor() })]
    }
    return opt
  },
})
</script>
