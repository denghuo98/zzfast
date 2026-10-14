<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">分析页</h1>
      <p class="text-muted-foreground mt-1">数据分析与统计图表展示。</p>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UCard v-for="stat in stats" :key="stat.label">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-primary/10">
            <UIcon :name="stat.icon" class="size-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stat.value }}</p>
            <p class="text-sm text-muted-foreground">{{ stat.label }}</p>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <UCard class="lg:col-span-2">
        <template #header>
          <h3 class="font-semibold">周度趋势（Weekly Trend）</h3>
        </template>
        <MLineChart height="14rem" :data="weeklyTrend" :x-axis-data="days" smooth show-area-color />
      </UCard>

      <UCard>
        <template #header>
          <h3 class="font-semibold">流量来源（Traffic Sources）</h3>
        </template>
        <MPieChart height="14rem" :data="trafficSources" :radius="['50%','80%']" show-legend />
      </UCard>
    </div>

    <UCard>
      <template #header>
        <h3 class="font-semibold">转化漏斗（Conversion Funnel）</h3>
      </template>
      <MBarChart height="16rem" :data="funnelData" :x-axis-data="funnelStages" bar-width="30%" :stack="true" show-legend />
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const stats = [
  { label: '总用户', value: '12,480', icon: 'i-lucide-users' },
  { label: '活跃会话', value: '342', icon: 'i-lucide-activity' },
  { label: '收入', value: '¥ 48,320', icon: 'i-lucide-dollar-sign' },
  { label: '待解决问题', value: '23', icon: 'i-lucide-alert-circle' },
]

const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const weeklyTrend = [150, 230, 224, 218, 235, 160, 120]

const trafficSources = [
  { name: '直接访问', value: 335 },
  { name: '搜索引擎', value: 210 },
  { name: '社交媒体', value: 154 },
  { name: '邮件营销', value: 100 },
  { name: '外部链接', value: 85 },
]

const funnelStages = ['曝光', '点击', '访问', '咨询', '下单']
const funnelData = [
  { name: '完成转化', data: [500, 350, 200, 120, 75], stack: 'total' },
  { name: '流失', data: [0, 150, 150, 80, 45], stack: 'total' },
]
</script>
