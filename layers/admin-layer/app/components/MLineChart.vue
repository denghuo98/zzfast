<!-- 折线图：支持单组/多组，平滑曲线 + 面积渐变 -->
<template>
    <div ref="chartRef" class="relative w-full" :style="{ height: props.height }" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { graphic, type EChartsOption } from '../utils/echarts'
import { useChartOps, useChartComponent, getCssVar } from '../composables/useChart'
import type { LineChartProps, LineDataItem } from '../types/chart'

defineOptions({ name: 'MLineChart' })

const props = withDefaults(defineProps<LineChartProps>(), {
  height: '16rem',
  loading: false,
  isEmpty: false,
  colors: () => undefined as string[] | undefined,
  data: () => [0, 0, 0, 0, 0, 0, 0],
  xAxisData: () => [],
  lineWidth: 2.5,
  showAreaColor: false,
  smooth: true,
  symbol: 'none',
  symbolSize: 6,
  animationDelay: 200,
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

const primaryColor = computed(() => getCssVar('--ui-primary') || '#00704A')

const getColor = (custom?: string, index?: number): string => {
  if (custom) return custom
  if (index !== undefined) return resolvedColors.value[index % resolvedColors.value.length]
  return primaryColor.value
}

/** hex → rgba 工具 */
function hexToRgba(hex: string, alpha: number): string {
  const clean = hex.replace('#', '')
  const full = clean.length === 3 ? clean.split('').map(c => c + c).join('') : clean
  const r = parseInt(full.substring(0, 2), 16)
  const g = parseInt(full.substring(2, 4), 16)
  const b = parseInt(full.substring(4, 6), 16)
  return `rgba(${r},${g},${b},${Math.max(0, Math.min(1, alpha))})`
}

const generateAreaStyle = (item: LineDataItem, color: string) => {
  const ac = item.areaStyle || {}
  if (ac.custom) return ac.custom
  return {
    color: new graphic.LinearGradient(0, 0, 0, 1, [
      { offset: 0, color: hexToRgba(color, ac.startOpacity ?? 0.2) },
      { offset: 1, color: hexToRgba(color, ac.endOpacity ?? 0.02) },
    ]),
  }
}

const createSeriesItem = (cfg: { name?: string; data: number[]; color: string; smooth?: boolean; symbol?: string; symbolSize?: number; lineWidth?: number; areaStyle?: any }) => ({
  name: cfg.name, data: cfg.data, type: 'line' as const,
  color: cfg.color,
  smooth: cfg.smooth ?? props.smooth,
  symbol: (cfg.symbol ?? props.symbol) as any,
  symbolSize: cfg.symbolSize ?? props.symbolSize,
  lineStyle: { width: cfg.lineWidth ?? props.lineWidth, color: cfg.color },
  areaStyle: cfg.areaStyle,
  emphasis: { focus: 'series' as const, lineStyle: { width: (cfg.lineWidth ?? props.lineWidth) + 1 } },
})

const {
  chartRef,
  getAxisLineStyle, getAxisLabelStyle, getAxisTickStyle,
  getSplitLineStyle, getTooltipStyle, getLegendStyle, getGridWithLegend,
} = useChartComponent<LineChartProps>({
  props,
  checkEmpty: () => {
    if (!isMultipleData.value) {
      const sd = props.data as number[]
      return !sd.length || sd.every(v => v === 0)
    }
    const md = props.data as LineDataItem[]
    return !md.length || md.every(item => !item.data?.length || item.data.every(v => v === 0))
  },
  watchSources: [() => props.data, () => props.xAxisData, () => props.colors],
  generateOptions: (): EChartsOption => {
    const opt: EChartsOption = {
      grid: getGridWithLegend(!!props.showLegend && isMultipleData.value, props.legendPosition, { top: 15, right: 15, left: 0 }),
      tooltip: props.showTooltip ? getTooltipStyle() : undefined,
      xAxis: {
        type: 'category', boundaryGap: false, data: props.xAxisData,
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
      const md = props.data as LineDataItem[]
      opt.series = md.map((item, i) => {
        const c = getColor(resolvedColors.value[i], i)
        return createSeriesItem({
          name: item.name, data: item.data, color: c,
          smooth: item.smooth, symbol: item.symbol, lineWidth: item.lineWidth,
          areaStyle: (item.showAreaColor ?? props.showAreaColor) ? generateAreaStyle(item, c) : undefined,
        })
      })
    } else {
      const sd = props.data as number[]
      const c = getColor(resolvedColors.value[0])
      opt.series = [createSeriesItem({
        data: sd, color: c,
        areaStyle: props.showAreaColor ? generateAreaStyle({ data: sd } as any, c) : undefined,
      })]
    }
    return opt
  },
})
</script>
