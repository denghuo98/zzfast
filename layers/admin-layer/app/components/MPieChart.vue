<!-- 饼图/环形图：radius 控制，环形中心可显示文字 -->
<template>
    <div ref="chartRef" class="relative w-full" :style="{ height: props.height }" />
</template>

<script setup lang="ts">
import type { EChartsOption } from '../utils/echarts'
import { useChartOps, useChartComponent } from '../composables/useChart'
import type { PieChartProps } from '../types/chart'

defineOptions({ name: 'MPieChart' })

const props = withDefaults(defineProps<PieChartProps>(), {
  height: '16rem',
  loading: false,
  isEmpty: false,
  colors: () => undefined as string[] | undefined,
  data: () => [],
  radius: () => ['50%', '80%'],
  borderRadius: 0,
  centerText: '',
  showLabel: false,
  showTooltip: true,
  showLegend: false,
  legendPosition: 'right',
})

const resolvedColors = computed(() => props.colors?.length ? props.colors : useChartOps().colors)

const { chartRef, isDark, getAnimationConfig, getTooltipStyle, getLegendStyle } = useChartComponent<PieChartProps>({
  props,
  checkEmpty: () => !props.data?.length || props.data.every(item => item.value === 0),
  watchSources: [() => props.data, () => props.centerText],
  generateOptions: (): EChartsOption => {
    const getCenter = (): [string, string] => {
      if (!props.showLegend) return ['50%', '50%']
      switch (props.legendPosition) {
        case 'left': return ['60%', '50%']
        case 'right': return ['40%', '50%']
        case 'top': return ['50%', '60%']
        case 'bottom': return ['50%', '40%']
        default: return ['50%', '50%']
      }
    }
    const center = getCenter()
    const opt: EChartsOption = {
      tooltip: props.showTooltip ? getTooltipStyle('item', { formatter: '{b}: {c} ({d}%)' }) : undefined,
      legend: props.showLegend ? getLegendStyle(props.legendPosition) : undefined,
      series: [{
        name: '数据占比', type: 'pie',
        radius: props.radius,
        center,
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: props.borderRadius,
          borderColor: isDark.value ? '#2c2c2c' : '#fff',
          borderWidth: 0,
        },
        label: {
          show: props.showLabel,
          formatter: '{b}\n{d}%',
          position: 'outside',
          color: isDark.value ? '#ccc' : '#999',
          fontSize: 12,
        },
        emphasis: { label: { show: false, fontSize: 14, fontWeight: 'bold' } },
        labelLine: { show: props.showLabel, length: 15, length2: 25, smooth: true },
        data: props.data,
        color: resolvedColors.value,
        ...getAnimationConfig(),
        animationType: 'expansion',
      }],
    }
    if (props.centerText) {
      opt.title = {
        text: props.centerText,
        left: center[0], top: center[1],
        textAlign: 'center', textVerticalAlign: 'middle',
        textStyle: { fontSize: 18, fontWeight: 500, color: isDark.value ? '#999' : '#ADB0BC' },
      }
    }
    return opt
  },
})
</script>
