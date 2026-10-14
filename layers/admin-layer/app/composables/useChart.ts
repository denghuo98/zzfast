/**
 * useChart — ECharts 图表管理 composable
 *
 * 提供 ECharts 图表的完整生命周期管理：初始化、更新、销毁、暗黑自动适配、响应式 resize。
 * 分为三个导出：
 * - useChartOps() — 主题色板（读 CSS var --ui-primary）
 * - useChart(options) — 低层级图表生命周期
 * - useChartComponent(options) — 高层抽象：props → generateOptions → watch → theme regenerate
 *
 * 适配 Nuxt UI v4：isDark 来源 useColorMode()、resize 用 ResizeObserver（覆盖侧栏折叠）。
 */

import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount, onUnmounted } from 'vue'
import { useResizeObserver } from '@vueuse/core'
import { echarts, type EChartsOption } from '../utils/echarts'
import type { BaseChartProps, ChartThemeConfig, UseChartOptions } from '../types/chart'

// ─── SSR-safe CSS 变量读取 ─────────────────────────────

/** 读取 CSS 自定义属性值。SSR 时返回空字符串，仅客户端执行实际读取 */
export function getCssVar(name: string): string {
  if (typeof document === 'undefined') return ''
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

// ─── 主题色板 ──────────────────────────────────────────

export function useChartOps(): ChartThemeConfig {
  return {
    chartHeight: '16rem',
    fontSize: 13,
    fontColor: '#999',
    themeColor: getCssVar('--ui-primary') || '#00704A',
    colors: [
      getCssVar('--ui-primary') || '#00704A',
      '#4ABEFF',
      '#FFAF20',
      '#FA8A6C',
      '#A48DF3',
      '#14DEBA',
      '#FF80C8',
      '#38C0FC',
    ],
  }
}

// ─── useChart（低层级）──────────────────────────────────

const RESIZE_DEBOUNCE_DELAY = 100

export function useChart(options: UseChartOptions = {}) {
  const { initOptions, initDelay = 0, threshold = 0.1, autoTheme = true } = options

  // --- 暗黑状态（直接读 html.dark class，Nuxt UI / @nuxtjs/color-mode 以此为真值）---
  const isDark = ref(false)
  let darkObserver: MutationObserver | null = null

  if (typeof document !== 'undefined') {
    const syncDark = () => {
      isDark.value = document.documentElement.classList.contains('dark')
    }
    syncDark()
    darkObserver = new MutationObserver(syncDark)
    darkObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  }

  onBeforeUnmount(() => {
    darkObserver?.disconnect()
    darkObserver = null
  })

  const chartRef = ref<HTMLElement>()
  let chart: echarts.ECharts | null = null
  let resizeTimeoutId: ReturnType<typeof setTimeout> | null = null
  let resizeFrameId: number | null = null
  let isDestroyed = false
  let emptyStateDiv: HTMLElement | null = null

  // 清理定时器
  const clearTimers = () => {
    if (resizeTimeoutId) { clearTimeout(resizeTimeoutId); resizeTimeoutId = null }
    if (resizeFrameId) { cancelAnimationFrame(resizeFrameId); resizeFrameId = null }
  }

  // rAF resize
  const requestAnimationResize = () => {
    if (resizeFrameId) cancelAnimationFrame(resizeFrameId)
    resizeFrameId = requestAnimationFrame(() => { handleResize(); resizeFrameId = null })
  }

  // debounce resize（窗口 resize）
  const debouncedResize = () => {
    if (resizeTimeoutId) clearTimeout(resizeTimeoutId)
    resizeTimeoutId = setTimeout(() => { requestAnimationResize(); resizeTimeoutId = null }, RESIZE_DEBOUNCE_DELAY)
  }

  // --- ResizeObserver（chart 容器自身尺寸变化，覆盖侧栏折叠等任何容器变化）---
  let resizeStop: (() => void) | null = null
  let observerSetup = false

  const setupResizeObserver = () => {
    if (observerSetup || !chartRef.value) return
    observerSetup = true
    if (useResizeObserver) {
      resizeStop = useResizeObserver(chartRef, () => { requestAnimationResize() })
    }
  }

  // --- 主题 watcher ---
  let themeStopHandle: (() => void) | null = null

  const setupThemeWatcher = () => {
    if (!autoTheme) return
    themeStopHandle = watch(isDark, () => {
      emptyStateManager.updateStyle()
      if (chart && !isDestroyed) {
        requestAnimationFrame(() => {
          if (chart && !isDestroyed) {
            const currentOptions = chart.getOption()
            if (currentOptions) updateChart(currentOptions as EChartsOption)
          }
        })
      }
    })
  }

  const cleanupThemeWatcher = () => { themeStopHandle?.(); themeStopHandle = null }

  // --- 样式 cache（暗黑切换时清空）---
  const createLineStyle = (color: string, width = 1, type?: 'solid' | 'dashed') => ({
    color, width, ...(type && { type }),
  })

  const styleCache: Record<string, any> = { lastDarkValue: isDark.value }

  const clearStyleCache = () => {
    styleCache.axisLine = null; styleCache.splitLine = null; styleCache.axisLabel = null
    styleCache.lastDarkValue = isDark.value
  }

  const getAxisLineStyle = (show = true) => {
    if (styleCache.lastDarkValue !== isDark.value) clearStyleCache()
    if (!styleCache.axisLine) {
      styleCache.axisLine = { show, lineStyle: createLineStyle(isDark.value ? '#444' : '#EDEDED') }
    }
    return styleCache.axisLine
  }

  const getSplitLineStyle = (show = true) => {
    if (styleCache.lastDarkValue !== isDark.value) clearStyleCache()
    if (!styleCache.splitLine) {
      styleCache.splitLine = { show, lineStyle: createLineStyle(isDark.value ? '#444' : '#EDEDED', 1, 'dashed') }
    }
    return styleCache.splitLine
  }

  const getAxisLabelStyle = (show = true) => {
    if (styleCache.lastDarkValue !== isDark.value) clearStyleCache()
    if (!styleCache.axisLabel) {
      const { fontColor, fontSize } = useChartOps()
      styleCache.axisLabel = { show, color: fontColor, fontSize }
    }
    return styleCache.axisLabel
  }

  const getAxisTickStyle = () => ({ show: false })

  const getAnimationConfig = (animationDelay = 50, animationDuration = 1500) => ({
    animationDelay: (idx: number) => idx * animationDelay + 200,
    animationDuration: (idx: number) => animationDuration - idx * 50,
    animationEasing: 'quarticOut' as const,
  })

  const getTooltipStyle = (trigger: 'item' | 'axis' = 'axis', customOptions: any = {}) => ({
    trigger,
    backgroundColor: isDark.value ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.9)',
    borderColor: isDark.value ? '#333' : '#ddd',
    borderWidth: 1,
    textStyle: { color: isDark.value ? '#fff' : '#333' },
    ...customOptions,
  })

  const getLegendStyle = (position: 'bottom' | 'top' | 'left' | 'right' = 'bottom', customOptions: any = {}) => {
    const base = {
      textStyle: { color: isDark.value ? '#fff' : '#333' },
      itemWidth: 12, itemHeight: 12, itemGap: 20,
      ...customOptions,
    }
    switch (position) {
      case 'bottom': return { ...base, bottom: 0, left: 'center', orient: 'horizontal', icon: 'roundRect' }
      case 'top': return { ...base, top: 0, left: 'center', orient: 'horizontal', icon: 'roundRect' }
      case 'left': return { ...base, left: 0, top: 'center', orient: 'vertical', icon: 'roundRect' }
      case 'right': return { ...base, right: 0, top: 'center', orient: 'vertical', icon: 'roundRect' }
      default: return base
    }
  }

  const getGridWithLegend = (showLegend: boolean, legendPosition = 'bottom' as const, baseGrid: any = {}) => {
    const defaults = { top: 15, right: 15, bottom: 8, left: 0, containLabel: true, ...baseGrid }
    if (!showLegend) return defaults
    switch (legendPosition) {
      case 'bottom': return { ...defaults, bottom: 40 }
      case 'top': return { ...defaults, top: 40 }
      case 'left': return { ...defaults, left: 120 }
      case 'right': return { ...defaults, right: 120 }
      default: return defaults
    }
  }

  // --- 空状态管理器 ---
  const emptyStateManager = {
    create: () => {
      if (!chartRef.value || emptyStateDiv) return
      emptyStateDiv = document.createElement('div')
      emptyStateDiv.style.cssText = `position:absolute;top:0;left:0;right:0;bottom:0;display:flex;flex-direction:column;align-items:center;justify-content:center;font-size:12px;color:${isDark.value ? '#666' : '#999'};background:transparent;z-index:10`
      emptyStateDiv.innerHTML = '<span>暂无数据</span>'
      const pos = getComputedStyle(chartRef.value).position
      if (pos !== 'relative' && pos !== 'absolute') chartRef.value.style.position = 'relative'
      chartRef.value.appendChild(emptyStateDiv)
    },
    remove: () => {
      if (emptyStateDiv && chartRef.value) { chartRef.value.removeChild(emptyStateDiv); emptyStateDiv = null }
    },
    updateStyle: () => {
      if (emptyStateDiv) emptyStateDiv.style.color = isDark.value ? '#666' : '#999'
    },
  }

  // --- 核心方法 ---

  const performChartInit = (opts: EChartsOption) => {
    if (!chart && chartRef.value && !isDestroyed) {
      chart = echarts.init(chartRef.value)
      setupResizeObserver()
      setupThemeWatcher()
    }
    if (chart && !isDestroyed) { chart.setOption(opts) }
  }

  const initChart = (opts: EChartsOption = {}, isEmpty = false) => {
    if (!chartRef.value || isDestroyed) return
    const merged = { ...initOptions, ...opts }
    try {
      if (isEmpty) { chart?.clear(); emptyStateManager.create(); return }
      emptyStateManager.remove()
      // 始终直接初始化（跳过 IntersectionObserver 懒加载，简化逻辑）
      if (initDelay > 0) setTimeout(() => performChartInit(merged), initDelay)
      else performChartInit(merged)
    } catch (error) { console.error('[useChart] init error:', error) }
  }

  const updateChart = (opts: EChartsOption) => {
    if (isDestroyed) return
    try {
      if (!chart) { initChart(opts); return }
      chart.setOption(opts)
    } catch (error) { console.error('[useChart] update error:', error) }
  }

  const handleResize = () => { if (chart && !isDestroyed) { try { chart.resize() } catch {} } }

  const destroyChart = () => {
    isDestroyed = true
    if (chart) { try { chart.dispose() } catch {} finally { chart = null } }
    cleanupThemeWatcher()
    resizeStop?.()
    emptyStateManager.remove()
    clearTimers()
    clearStyleCache()
  }

  const getChartInstance = () => chart
  const isChartInitialized = () => chart !== null

  // --- 生命周期 ---
  onMounted(() => {
    window.addEventListener('resize', debouncedResize)
  })
  onBeforeUnmount(() => {
    window.removeEventListener('resize', debouncedResize)
  })
  onUnmounted(() => destroyChart())

  return {
    isDark, chartRef, initChart, updateChart, handleResize, destroyChart,
    getChartInstance, isChartInitialized, emptyStateManager,
    getAxisLineStyle, getSplitLineStyle, getAxisLabelStyle, getAxisTickStyle,
    getAnimationConfig, getTooltipStyle, getLegendStyle, useChartOps, getGridWithLegend,
  }
}

// ─── useChartComponent（高层抽象）──────────────────────────

interface UseChartComponentOptions<T extends BaseChartProps> {
  props: T
  generateOptions: () => EChartsOption
  checkEmpty?: () => boolean
  watchSources?: (() => any)[]
  chartOptions?: UseChartOptions
}

export function useChartComponent<T extends BaseChartProps>(opts: UseChartComponentOptions<T>) {
  const { props, generateOptions, checkEmpty, watchSources = [], chartOptions = {} } = opts
  const chart = useChart(chartOptions)
  const { chartRef, initChart, isDark, emptyStateManager } = chart

  const isEmpty = computed(() => {
    if (props.isEmpty) return true
    if (checkEmpty) return checkEmpty()
    return false
  })

  const updateChart = () => {
    nextTick(() => {
      if (isEmpty.value) {
        chart.getChartInstance()?.clear()
        emptyStateManager.create()
      } else {
        emptyStateManager.remove()
        initChart(generateOptions())
      }
    })
  }

  const stopHandles: (() => void)[] = []

  const setupWatchers = () => {
    if (watchSources.length > 0) {
      stopHandles.push(watch(watchSources, updateChart, { deep: true }))
    }
    stopHandles.push(watch(isDark, () => {
      emptyStateManager.updateStyle()
      updateChart()
    }))
  }

  const setupLifecycle = () => {
    onMounted(() => {
      updateChart()
    })
    onBeforeUnmount(() => {
      stopHandles.forEach(s => s()); stopHandles.length = 0
      emptyStateManager.remove()
    })
  }

  setupWatchers()
  setupLifecycle()

  return { ...chart, isEmpty, updateChart }
}
