<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold">卡片</h1>
      <p class="text-muted-foreground mt-1">对齐 art-design-pro 的卡片模板：统计 / 进度 / 图表 / 数据列表 / 图片。</p>
    </div>

    <!-- 1. 统计卡片（文字） -->
    <div>
      <h2 class="text-lg font-semibold mb-3">统计卡片（文字）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="card in statsCards" :key="card.id">
          <div class="flex items-start justify-between">
            <div>
              <p class="text-sm text-muted-foreground">{{ card.title }}</p>
              <p class="text-xl font-bold mt-1">{{ card.count?.toLocaleString() }}</p>
              <p class="text-xs text-muted-foreground mt-1">{{ card.description }}</p>
            </div>
            <div class="size-10 rounded-lg flex items-center justify-center shrink-0" :class="card.iconBg">
              <UIcon :name="card.icon" class="size-5" :class="card.iconColor" />
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 2. 统计卡片（数字滚动） -->
    <div>
      <h2 class="text-lg font-semibold mb-3">统计卡片（数字滚动）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="card in statsCards" :key="'count-' + card.id">
          <div class="flex items-start justify-between">
            <div>
              <p class="text-sm text-muted-foreground">{{ card.title }}</p>
              <p class="text-xl font-bold mt-1 tabular-nums">
                <MCountTo :target="card.count" :duration="2000" separator="," />
              </p>
              <p class="text-xs text-muted-foreground mt-1">{{ card.description }}</p>
            </div>
            <div class="size-10 rounded-lg flex items-center justify-center shrink-0" :class="card.iconBg">
              <UIcon :name="card.icon" class="size-5" :class="card.iconColor" />
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 3. 统计卡片（自定义样式） -->
    <div>
      <h2 class="text-lg font-semibold mb-3">统计卡片（自定义样式）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="card in statsCards" :key="'custom-' + card.id" :class="card.boxBg">
          <div class="flex items-start justify-between">
            <div>
              <p class="text-sm" :class="card.textColor">{{ card.title }}</p>
              <p class="text-xl font-bold mt-1" :class="card.textColor">{{ card.count?.toLocaleString() }}</p>
              <p class="text-xs mt-1 opacity-70" :class="card.textColor">{{ card.description }}</p>
            </div>
            <div class="size-12 rounded-lg flex items-center justify-center shrink-0 text-3xl" :class="card.customIconBg">
              <UIcon :name="card.icon" class="size-7" :class="card.customIconColor" />
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 4. 进度卡片 -->
    <div>
      <h2 class="text-lg font-semibold mb-3">进度卡片</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="proj in progressCards" :key="proj.id">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium">{{ proj.title }}</span>
            <span class="text-sm text-muted-foreground">{{ proj.percentage }}%</span>
          </div>
          <UProgress :value="proj.percentage" :color="proj.uiColor as any" size="sm" />
        </UCard>
      </div>
    </div>

    <!-- 5. 进度卡片（icon） -->
    <div>
      <h2 class="text-lg font-semibold mb-3">进度卡片（icon）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="proj in progressCards" :key="'icon-' + proj.id">
          <div class="flex items-center gap-3 mb-3">
            <div class="size-9 rounded-lg flex items-center justify-center shrink-0" :class="proj.iconBg">
              <UIcon :name="proj.icon" class="size-5" :class="proj.iconColor" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{{ proj.title }}</p>
              <p class="text-xs text-muted-foreground">{{ proj.percentage }}%</p>
            </div>
          </div>
          <UProgress :value="proj.percentage" :color="proj.uiColor as any" size="sm" />
        </UCard>
      </div>
    </div>

    <!-- 6. 图表卡片（小图表）— 左上数值文案 + 右上 chart + 左下环比 + 右下时间 -->
    <div>
      <h2 class="text-lg font-semibold mb-3">图表卡片（小图表）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0' }">
          <div class="relative" style="height:9.5rem">
            <div class="p-4">
              <p class="text-xl font-semibold leading-tight">2,545</p>
              <p class="text-xs text-muted-foreground mt-0.5">新用户</p>
            </div>
            <div class="absolute top-3 right-3 w-[40%] -mr-4">
              <MLineChart height="3.5rem" :data="[120,132,101,134,90,230,210]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" smooth />
            </div>
            <p class="absolute bottom-3 left-4 text-xs text-success">+1.2%</p>
            <p class="absolute bottom-3 right-4 text-xs" style="color:#9ca3af">过去7天</p>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0' }">
          <div class="relative" style="height:9.5rem">
            <div class="p-4">
              <p class="text-xl font-semibold leading-tight">15,480</p>
              <p class="text-xs text-muted-foreground mt-0.5">浏览量</p>
            </div>
            <div class="absolute top-3 right-3 w-[40%] -mr-4">
              <MBarChart height="3.5rem" :data="[120,100,150,140,90,120,130]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" bar-width="45%" />
            </div>
            <p class="absolute bottom-3 left-4 text-xs text-error">-4.15%</p>
            <p class="absolute bottom-3 right-4 text-xs" style="color:#9ca3af">过去14天</p>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0' }">
          <div class="relative" style="height:9.5rem">
            <div class="p-4">
              <p class="text-xl font-semibold leading-tight">2,545</p>
              <p class="text-xs text-muted-foreground mt-0.5">粉丝数</p>
            </div>
            <div class="absolute top-3 right-3 w-[40%] -mr-4">
              <MLineChart height="3.5rem" :data="[150,180,160,200,180,220,240]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" smooth show-area-color />
            </div>
            <p class="absolute bottom-3 left-4 text-xs text-success">+1.2%</p>
            <p class="absolute bottom-3 right-4 text-xs" style="color:#9ca3af">过去30天</p>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0' }">
          <div class="flex items-center" style="height:9.5rem">
            <div class="flex-1 p-4">
              <p class="text-lg font-semibold">粉丝量</p>
              <p class="text-xl font-semibold mt-1">36,358</p>
              <p class="text-xs text-success mt-1">+18% <span class="text-muted-foreground">较去年</span></p>
              <div class="flex gap-3 mt-2 text-xs text-muted-foreground">
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-primary" />2022</span>
                <span class="flex items-center gap-1"><span class="size-2 rounded-full" style="background:#e6e8f7" />2021</span>
              </div>
            </div>
            <div class="w-28 h-full flex items-center justify-center pr-2">
              <MPieChart height="6rem" :data="[{name:'2022',value:50},{name:'2021',value:40}]" :radius="['70%','90%']" :show-tooltip="false" />
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 7. 图表卡片（大图表）— 顶部数值+环比，底部 chart 撑满 -->
    <div>
      <h2 class="text-lg font-semibold mb-3">图表卡片（大图表）</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0 flex flex-col' }" style="height:11rem">
          <div class="flex items-start justify-between px-5 pt-4 shrink-0">
            <div>
              <p class="text-xl font-semibold leading-tight">2,545</p>
              <p class="text-xs text-muted-foreground mt-0.5">新用户</p>
            </div>
            <p class="text-xs text-success">+1.2%</p>
          </div>
          <div class="flex-1 relative -mr-4">
            <div class="absolute inset-0">
              <MLineChart height="100%" :data="[120,132,101,134,90,230,210]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" smooth />
            </div>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0 flex flex-col' }" style="height:11rem">
          <div class="flex items-start justify-between px-5 pt-4 shrink-0">
            <div>
              <p class="text-xl font-semibold leading-tight">15,480</p>
              <p class="text-xs text-muted-foreground mt-0.5">浏览量</p>
            </div>
            <p class="text-xs text-error">-4.15%</p>
          </div>
          <div class="flex-1 relative -mr-4">
            <div class="absolute inset-0">
              <MBarChart height="100%" :data="[120,100,150,140,90,120,130,110]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" bar-width="50%" />
            </div>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0 flex flex-col' }" style="height:11rem">
          <div class="flex items-start justify-between px-5 pt-4 shrink-0">
            <div>
              <p class="text-xl font-semibold leading-tight">2,545</p>
              <p class="text-xs text-muted-foreground mt-0.5">粉丝数</p>
            </div>
            <p class="text-xs text-success">+1.2%</p>
          </div>
          <div class="flex-1 relative -mr-4">
            <div class="absolute inset-0">
              <MLineChart height="100%" :data="[150,180,160,200,180,220,240]" :show-axis-label="false" :show-axis-line="false" :show-split-line="false" :show-tooltip="false" smooth show-area-color />
            </div>
          </div>
        </UCard>
        <UCard :ui="{ root: 'overflow-hidden', body: '!p-0' }" style="height:11rem">
          <div class="flex h-full items-center">
            <div class="flex-1 p-4">
              <p class="text-lg font-semibold">粉丝量</p>
              <p class="text-xl font-semibold mt-1">36,358</p>
              <p class="text-xs text-error mt-1">-18% <span class="text-muted-foreground">较11月</span></p>
              <div class="flex gap-3 mt-1.5 text-xs text-muted-foreground">
                <span class="flex items-center gap-1"><span class="size-2 rounded-full bg-primary" />12月</span>
                <span class="flex items-center gap-1"><span class="size-2 rounded-full" style="background:#e6e8f7" />11月</span>
              </div>
            </div>
            <div class="w-36 h-full flex items-center justify-center pr-2">
              <MPieChart height="7rem" :data="[{name:'12月',value:70},{name:'11月',value:30}]" :radius="['60%','85%']" :show-tooltip="false" />
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 8. 数据列表卡片 -->
    <div>
      <h2 class="text-lg font-semibold mb-3">数据列表卡片</h2>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <UCard>
          <template #header><h3 class="font-semibold text-sm">待办事项</h3></template>
          <p class="text-xs text-muted-foreground mb-3">今日待处理任务</p>
          <div class="space-y-2">
            <div v-for="item in dataList" :key="item.title" class="flex items-center gap-3 py-1.5">
              <div class="size-8 rounded-lg flex items-center justify-center shrink-0" :class="item.bg">
                <UIcon :name="item.icon" class="size-4" :class="item.color" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm truncate">{{ item.title }}</p>
                <p class="text-xs text-muted-foreground">{{ item.time }}</p>
              </div>
              <UBadge size="xs" variant="soft" :color="item.statusColor as any">{{ item.status }}</UBadge>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header><h3 class="font-semibold text-sm">最近活动</h3></template>
          <p class="text-xs text-muted-foreground mb-3">近期活动列表</p>
          <div class="space-y-2">
            <div v-for="item in dataList.slice(0, 4)" :key="item.title" class="flex items-center gap-3 py-1.5">
              <div class="size-8 rounded-lg flex items-center justify-center shrink-0" :class="item.bg">
                <UIcon :name="item.icon" class="size-4" :class="item.color" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm truncate">{{ item.title }}</p>
                <p class="text-xs text-muted-foreground">{{ item.time }}</p>
              </div>
              <UBadge size="xs" variant="soft" :color="item.statusColor as any">{{ item.status }}</UBadge>
            </div>
          </div>
          <UButton variant="ghost" size="sm" block class="mt-2">查看更多</UButton>
        </UCard>

        <UCard>
          <template #header><h3 class="font-semibold text-sm">最近交易</h3></template>
          <p class="text-xs text-muted-foreground mb-3">{{ new Date().toLocaleDateString('zh-CN') }}</p>
          <div class="space-y-3">
            <div v-for="tx in timelineData" :key="tx.content" class="flex gap-3">
              <div class="flex flex-col items-center">
                <span class="size-2.5 rounded-full shrink-0 mt-1.5" :style="{ backgroundColor: tx.color }" />
                <div v-if="tx !== timelineData[timelineData.length - 1]" class="w-px flex-1 my-1" :style="{ backgroundColor: tx.color, opacity: 0.3 }" />
              </div>
              <div class="pb-2 flex-1">
                <p class="text-xs text-muted-foreground">{{ tx.time }}</p>
                <p class="text-sm mt-0.5">{{ tx.content }}</p>
                <p v-if="tx.code" class="text-xs text-muted-foreground mt-0.5">{{ tx.code }}</p>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- 9. 图片卡片 -->
    <div>
      <h2 class="text-lg font-semibold mb-3">图片卡片</h2>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard v-for="card in imageCards" :key="card.id" class="cursor-pointer hover:shadow-md transition-shadow" :ui="{ root: 'overflow-hidden !p-0', body: '!p-0' }">
          <img :src="card.imageUrl" :alt="card.title" class="w-full h-40 object-cover" />
          <div class="p-3 space-y-1.5">
            <UBadge size="xs" variant="soft">{{ card.category }}</UBadge>
            <h4 class="text-sm font-semibold line-clamp-2">{{ card.title }}</h4>
            <div class="flex items-center gap-3 text-xs text-muted-foreground">
              <span class="flex items-center gap-1"><UIcon name="i-lucide-clock" class="size-3" />{{ card.readTime }}</span>
              <span class="flex items-center gap-1"><UIcon name="i-lucide-eye" class="size-3" />{{ card.views?.toLocaleString() }}</span>
              <span class="flex items-center gap-1"><UIcon name="i-lucide-message-circle" class="size-3" />{{ card.comments }}</span>
            </div>
            <p class="text-xs text-muted-foreground">{{ card.date }}</p>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const statsCards = [
  { id: 1, title: '销售产品', count: 1235, description: '鞋子、牛仔裤、派对服装、手表', icon: 'i-lucide-bar-chart-3', iconBg: 'bg-info/10', iconColor: 'text-info', boxBg: '', customIconBg: 'bg-primary/10', customIconColor: 'text-primary', textColor: 'text-primary' },
  { id: 2, title: '活跃用户', count: 5000, description: '日活跃用户超过 5,000+', icon: 'i-lucide-user-round-search', iconBg: 'bg-warning/10', iconColor: 'text-warning', boxBg: '', customIconBg: 'bg-warning/10', customIconColor: 'text-warning', textColor: 'text-warning' },
  { id: 3, title: '总收入', count: 35000, description: '月收入超过 ¥350,000+', icon: 'i-lucide-banknote', iconBg: 'bg-secondary/10', iconColor: 'text-secondary', boxBg: '', customIconBg: 'bg-secondary/10', customIconColor: 'text-secondary', textColor: 'text-secondary' },
  { id: 4, title: '客户评价', count: 4800, description: '平均评分 4.8/5', icon: 'i-lucide-message-circle', iconBg: 'bg-error/10', iconColor: 'text-error', boxBg: '', customIconBg: 'bg-error/10', customIconColor: 'text-error', textColor: 'text-error' },
]

const progressCards = [
  { id: 1, title: '完成进度', percentage: 75, uiColor: 'success', icon: 'i-lucide-arrow-up-circle', iconBg: 'bg-success/10', iconColor: 'text-success' },
  { id: 2, title: '项目进度', percentage: 65, uiColor: 'primary', icon: 'i-lucide-twitch', iconBg: 'bg-primary/10', iconColor: 'text-primary' },
  { id: 3, title: '学习进度', percentage: 45, uiColor: 'error', icon: 'i-lucide-gamepad-2', iconBg: 'bg-error/10', iconColor: 'text-error' },
  { id: 4, title: '任务进度', percentage: 90, uiColor: 'secondary', icon: 'i-lucide-flag', iconBg: 'bg-secondary/10', iconColor: 'text-secondary' },
]

const dataList = [
  { title: '新加坡之行', status: '进行中', time: '5 分钟', icon: 'i-lucide-camera', bg: 'bg-primary/10', color: 'text-primary', statusColor: 'primary' },
  { title: '归档数据', status: '进行中', time: '10 分钟', icon: 'i-lucide-bar-chart-3', bg: 'bg-secondary/10', color: 'text-secondary', statusColor: 'secondary' },
  { title: '客户会议', status: '待处理', time: '15 分钟', icon: 'i-lucide-user', bg: 'bg-warning/10', color: 'text-warning', statusColor: 'warning' },
  { title: '筛选任务团队', status: '进行中', time: '20 分钟', icon: 'i-lucide-user-round', bg: 'bg-error/10', color: 'text-error', statusColor: 'error' },
  { title: '发送信封给小王', status: '已完成', time: '20 分钟', icon: 'i-lucide-mail', bg: 'bg-success/10', color: 'text-success', statusColor: 'success' },
]

const timelineData = [
  { time: '上午 09:30', color: '#49BEFF', content: '收到 John Doe 支付的 ¥385.90', code: '' },
  { time: '上午 10:00', color: '#369EFF', content: '新销售记录', code: 'ML-3467' },
  { time: '上午 12:00', color: '#67E8CF', content: '向 Michael 支付了 ¥64.95', code: '' },
  { time: '下午 14:30', color: '#FFC107', content: '系统维护通知', code: 'MT-2023' },
  { time: '下午 15:45', color: '#FF6969', content: '紧急订单取消提醒', code: 'OR-9876' },
  { time: '下午 17:00', color: '#67E8CF', content: '完成每日销售报表', code: '' },
]

const imageCards = [
  { id: 1, imageUrl: '/images/cover/img1.webp', title: 'AI 技术在医疗领域的创新应用与发展前景', category: '社交', readTime: '2分钟', views: 9125, comments: 3, date: '12月19日 周一' },
  { id: 2, imageUrl: '/images/cover/img2.webp', title: '大数据分析助力企业决策的实践案例', category: '技术', readTime: '3分钟', views: 7234, comments: 5, date: '12月20日 周二' },
  { id: 3, imageUrl: '/images/cover/img3.webp', title: '区块链技术在供应链管理中的应用', category: '科技', readTime: '4分钟', views: 5678, comments: 8, date: '12月21日 周三' },
  { id: 4, imageUrl: '/images/cover/img4.webp', title: '云计算技术发展趋势与未来展望', category: '云技术', readTime: '5分钟', views: 4321, comments: 6, date: '12月22日 周四' },
]
</script>
