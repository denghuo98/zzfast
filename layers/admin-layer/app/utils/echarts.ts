/**
 * ECharts 按需注册插件
 *
 * 仅注册 admin-layer 实际使用的图表类型和组件，减小打包体积。
 * 使用 .client.ts 后缀确保仅客户端执行，避免 SSR 引入 canvas/DOM。
 *
 * 使用：import { echarts, graphic, type EChartsOption } from '#layers/admin-layer/app/plugins/echarts'
 */

import * as echarts from 'echarts/core'

// 图表类型
import { BarChart, LineChart, PieChart, ScatterChart, RadarChart } from 'echarts/charts'

// 组件
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  MarkPointComponent,
  MarkLineComponent,
} from 'echarts/components'

// 渲染器
import { CanvasRenderer } from 'echarts/renderers'

// 注册
echarts.use([
  BarChart, LineChart, PieChart, ScatterChart, RadarChart,
  TitleComponent, TooltipComponent, GridComponent, LegendComponent,
  MarkPointComponent, MarkLineComponent,
  CanvasRenderer,
])

export { echarts }
export type { EChartsOption } from 'echarts'

/** echarts.graphic — 用于创建 LinearGradient 等 */
export const graphic = echarts.graphic
