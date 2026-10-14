<!-- 雷达图 -->
<template>
    <div ref="chartRef" class="relative w-full" :style="{ height: props.height }" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsOption } from '../utils/echarts'
import { useChartOps, useChartComponent } from '../composables/useChart'
import type { RadarChartProps } from '../types/chart'

defineOptions({ name: 'MRadarChart' })

const props = withDefaults(defineProps<RadarChartProps>(), {
  height: '16rem',
  loading: false,
  isEmpty: false,
  colors: () => undefined as string[] | undefined,
  indicator: () => [],
  data: () => [],
  showTooltip: true,
  showLegend: false,
  legendPosition: 'bottom',
})

const resolvedColors = computed(() => props.colors?.length ? props.colors : useChartOps().colors)

const {
  chartRef, isDark, getAnimationConfig,
  getTooltipStyle, getLegendStyle,
} = useChartComponent<RadarChartProps>({
  props,
  checkEmpty: () => !props.data?.length || props.data.every(d => !d.value?.length || d.value.every(v => v === 0)),
  watchSources: [() => props.data, () => props.indicator],
  generateOptions: (): EChartsOption => ({
    tooltip: props.showTooltip ? getTooltipStyle('item') : undefined,
    legend: props.showLegend ? getLegendStyle(props.legendPosition) : undefined,
    radar: {
      center: ['50%', '50%'],
      radius: '70%',
      indicator: props.indicator,
      axisName: { color: isDark.value ? '#ccc' : '#999', fontSize: 12 },
      splitLine: { lineStyle: { color: isDark.value ? '#444' : '#e6e6e6' } },
      axisLine: { lineStyle: { color: isDark.value ? '#444' : '#e6e6e6' } },
      splitArea: {
        areaStyle: {
          color: isDark.value
            ? ['rgba(255,255,255,0.02)', 'rgba(255,255,255,0.04)']
            : ['rgba(0,0,0,0.02)', 'rgba(0,0,0,0.04)'],
        },
      },
    },
    series: (props.data ?? []).map((item, i) => ({
      name: item.name,
      type: 'radar' as const,
      data: [{ value: item.value }],
      color: resolvedColors.value[i % resolvedColors.value.length],
      areaStyle: { opacity: 0.1 },
      emphasis: { areaStyle: { opacity: 0.25 }, lineStyle: { width: 3 } },
      ...getAnimationConfig(200, 1800),
    })),
  }),
})
</script>
