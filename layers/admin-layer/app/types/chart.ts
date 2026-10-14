/**
 * 图表组件类型定义
 *
 * 覆盖柱状/折线/饼图/散点/雷达共 5 种图表类型的 Props 与数据项。
 * 所有 Props 均扩展 BaseChartProps（统一 height/loading/empty/colors），
 * AxisDisplayProps（轴线显隐）与 InteractionProps（tooltip/legend）。
 */

import type { EChartsOption } from 'echarts'

// ─── 基础 ─────────────────────────────────────────────

/** 图例位置 */
export type LegendPosition = 'bottom' | 'top' | 'left' | 'right'

/** 数据点标记形状 */
export type SymbolType =
  | 'circle' | 'rect' | 'roundRect' | 'triangle' | 'diamond' | 'pin' | 'arrow' | 'none'

// ─── 主题 ─────────────────────────────────────────────

export interface ChartThemeConfig {
  chartHeight: string
  fontSize: number
  fontColor: string
  /** 品牌主色（读 CSS var --ui-primary） */
  themeColor: string
  /** 默认色板 */
  colors: string[]
}

export interface UseChartOptions {
  initOptions?: EChartsOption
  initDelay?: number
  threshold?: number
  autoTheme?: boolean
}

// ─── 共享 Props ────────────────────────────────────────

export interface BaseChartProps {
  height?: string
  loading?: boolean
  isEmpty?: boolean
  colors?: string[]
}

export interface AxisDisplayProps {
  showAxisLabel?: boolean
  showAxisLine?: boolean
  showSplitLine?: boolean
}

export interface InteractionProps {
  showTooltip?: boolean
  showLegend?: boolean
  legendPosition?: LegendPosition
}

// ─── 柱状图 ────────────────────────────────────────────

export interface BarDataItem {
  name: string
  data: number[]
  barWidth?: string | number
  stack?: string
}

export interface BarChartProps extends BaseChartProps, AxisDisplayProps, InteractionProps {
  data: number[] | BarDataItem[]
  xAxisData?: string[]
  barWidth?: string | number
  stack?: boolean
  borderRadius?: number | number[]
}

// ─── 折线图 ────────────────────────────────────────────

export interface LineDataItem {
  name: string
  data: number[]
  lineWidth?: number
  showAreaColor?: boolean
  areaStyle?: {
    startOpacity?: number
    endOpacity?: number
    custom?: any
  }
  smooth?: boolean
  symbol?: SymbolType
  symbolSize?: number
}

export interface LineChartProps extends BaseChartProps, AxisDisplayProps, InteractionProps {
  data: number[] | LineDataItem[]
  xAxisData?: string[]
  lineWidth?: number
  showAreaColor?: boolean
  smooth?: boolean
  symbol?: SymbolType
  symbolSize?: number
  animationDelay?: number
}

// ─── 饼图/环形图 ───────────────────────────────────────

export interface PieDataItem {
  value: number
  name: string
}

export interface PieChartProps extends BaseChartProps, InteractionProps {
  data: PieDataItem[]
  radius?: string[]
  borderRadius?: number
  centerText?: string
  showLabel?: boolean
}

// ─── 散点图 ────────────────────────────────────────────

export interface ScatterDataItem {
  value: number[]
}

export interface ScatterChartProps extends BaseChartProps, AxisDisplayProps, InteractionProps {
  data?: ScatterDataItem[]
  symbolSize?: number
}

// ─── 雷达图 ────────────────────────────────────────────

export interface RadarDataItem {
  name: string
  value: number[]
}

export interface RadarChartProps extends BaseChartProps, InteractionProps {
  indicator?: { name: string; max: number }[]
  data?: RadarDataItem[]
}
