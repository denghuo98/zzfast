<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">图表</h1>
      <p class="text-muted-foreground mt-1">基于 ECharts 的图表组件展示：柱状 / 折线 / 饼图 / 散点 / 雷达。</p>
    </div>

    <!-- 柱状图（单组） -->
    <UCard>
      <template #header>
        <h3 class="font-semibold">柱状图 — 单组数据</h3>
      </template>
      <MBarChart height="16rem" :data="singleBarData" :x-axis-data="months" />
    </UCard>

    <!-- 柱状图（多组 + 图例） -->
    <UCard>
      <template #header>
        <h3 class="font-semibold">柱状图 — 多组数据 + 堆叠 + 图例</h3>
      </template>
      <MBarChart height="18rem" :data="multiBarData" :x-axis-data="months" bar-width="26%" :stack="true" show-legend />
    </UCard>

    <!-- 柱状图（多组并排） -->
    <UCard>
      <template #header>
        <h3 class="font-semibold">柱状图 — 多组并排</h3>
      </template>
      <MBarChart height="16rem" :data="groupedBarData" :x-axis-data="months" bar-width="26%" show-legend />
    </UCard>

    <!-- 折线图（平滑 + 面积渐变） -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <UCard>
        <template #header>
          <h3 class="font-semibold">折线图 — 平滑 + 面积渐变</h3>
        </template>
        <MLineChart height="16rem" :data="lineData" :x-axis-data="months" smooth show-area-color />
      </UCard>

      <!-- 折线图（多组） -->
      <UCard>
        <template #header>
          <h3 class="font-semibold">折线图 — 多组 + 图例</h3>
        </template>
        <MLineChart height="16rem" :data="multiLineData" :x-axis-data="months" show-legend />
      </UCard>
    </div>

    <!-- 饼图 / 环形图 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <UCard>
        <template #header>
          <h3 class="font-semibold">环形图 — 中心文字</h3>
        </template>
        <MPieChart height="16rem" :data="pieData" :radius="['50%','80%']" center-text="¥ 300,458" show-legend />
      </UCard>

      <UCard>
        <template #header>
          <h3 class="font-semibold">饼图 — 显示标签</h3>
        </template>
        <MPieChart height="16rem" :data="pieData" :radius="['0%','70%']" show-label show-legend />
      </UCard>
    </div>

    <!-- 散点图 + 雷达图 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <UCard>
        <template #header>
          <h3 class="font-semibold">散点图</h3>
        </template>
        <MScatterChart height="16rem" :data="scatterData" />
      </UCard>

      <UCard>
        <template #header>
          <h3 class="font-semibold">雷达图</h3>
        </template>
        <MRadarChart height="16rem" :indicator="radarIndicator" :data="radarData" show-legend />
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const months = ['一月', '二月', '三月', '四月', '五月', '六月']

const singleBarData = [120, 200, 150, 80, 70, 110]

const multiBarData = [
  { name: '线上销售', data: [120, 200, 150, 80, 70, 110], stack: 'total' },
  { name: '线下销售', data: [80, 130, 100, 60, 90, 120], stack: 'total' },
  { name: '渠道分销', data: [50, 70, 60, 40, 30, 50], stack: 'total' },
]

const groupedBarData = [
  { name: '2025 年', data: [120, 200, 150, 80, 70, 110] },
  { name: '2026 年', data: [150, 180, 170, 110, 90, 140] },
]

const lineData = [58, 85, 90, 78, 95, 110]

const multiLineData = [
  { name: '线上', data: [120, 200, 150, 80, 70, 110] },
  { name: '线下', data: [80, 130, 100, 60, 90, 120] },
]

const pieData = [
  { name: '线上销售', value: 1048 },
  { name: '线下门店', value: 735 },
  { name: '渠道分销', value: 580 },
  { name: '海外市场', value: 484 },
  { name: '其他', value: 300 },
]

const scatterData = [
  { value: [10, 20] }, { value: [15, 35] }, { value: [25, 15] },
  { value: [30, 40] }, { value: [35, 25] }, { value: [40, 50] },
  { value: [45, 30] }, { value: [50, 60] }, { value: [55, 40] },
  { value: [60, 55] }, { value: [65, 45] }, { value: [70, 70] },
]

const radarIndicator = [
  { name: '销售', max: 100 },
  { name: '管理', max: 100 },
  { name: '技术', max: 100 },
  { name: '客服', max: 100 },
  { name: '研发', max: 100 },
  { name: '市场', max: 100 },
]

const radarData = [
  { name: '预算分配', value: [80, 90, 85, 70, 95, 75] },
  { name: '实际开销', value: [70, 85, 80, 65, 90, 70] },
]
</script>
