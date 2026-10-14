<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">基础表格</h1>
      <p class="text-muted-foreground mt-1">使用 MTable 展示基础数据表格，以及树形嵌套展开。</p>
    </div>

    <!-- 搜索 + 平铺表格 -->
    <UCard>
      <template #header>
        <span class="text-sm font-medium">搜索与平铺表格</span>
      </template>

      <div class="flex items-center gap-3 mb-4">
        <UInput v-model="search" icon="i-lucide-search" placeholder="搜索..." class="max-w-xs" />
      </div>

      <MTable :data="filteredData" :columns="flatColumns" />
    </UCard>

    <!-- 列固定（Column Pinning） -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">列固定</span>
          <span class="text-xs text-muted-foreground">
            点击列头 Pin 按钮切换固定状态；横向滚动查看固定效果
          </span>
        </div>
      </template>

      <MTable
        :data="pinningData"
        :columns="pinningColumns"
        v-model:column-pinning="pinningState"
      />

      <p class="text-xs text-muted-foreground mt-3">
        提示：左侧固定 Name 和 Email 两列，右侧固定操作列；缩窄浏览器窗口或调整下方容器宽度可观察横向滚动效果。
      </p>
    </UCard>

    <!-- 树形嵌套展开 -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">树形嵌套展开</span>
          <span class="text-xs text-muted-foreground">expand-key + get-sub-rows</span>
        </div>
      </template>

      <MTable
        :data="treeData"
        :columns="treeColumns"
        expand-key="name"
        :get-sub-rows="(row: any) => row.children"
        v-model:expanded="treeExpanded"
      >
        <template #role-cell="{ row }">
          <UBadge
            :color="row.original.role === '部门' ? 'primary' : row.original.role === '团队' ? 'info' : 'neutral'"
            variant="subtle"
          >
            {{ row.original.role }}
          </UBadge>
        </template>
        <template #status-cell="{ row }">
          <UBadge
            :color="row.original.status === '活跃' ? 'success' : 'neutral'"
            variant="subtle"
          >
            {{ row.original.status }}
          </UBadge>
        </template>
        <template #memberCount-cell="{ row }">
          <span class="text-xs text-muted-foreground">
            {{ row.original.memberCount ?? row.original.children?.length ?? '-' }} 人
          </span>
        </template>
      </MTable>
    </UCard>

    <!-- 行详情展开 -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">行详情展开</span>
          <span class="text-xs text-muted-foreground">#expanded slot</span>
        </div>
      </template>

      <MTable
        :data="detailData"
        :columns="detailColumns"
        v-model:expanded="detailExpanded"
      >
        <template #name-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton
              :icon="row.getIsExpanded() ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
              color="neutral"
              variant="ghost"
              size="xs"
              tabindex="-1"
              @click="row.toggleExpanded()"
            />
            <span>{{ row.original.name }}</span>
          </div>
        </template>
        <template #expanded="{ row }">
          <pre class="text-xs bg-elevated/50 rounded p-3 overflow-auto">{{ JSON.stringify(row.original, null, 2) }}</pre>
        </template>
      </MTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { Column } from '@tanstack/vue-table'

definePageMeta({ layout: 'admin' })

const UButton = resolveComponent('UButton')

const search = ref('')

// ======== 平铺表格 ========

const flatColumns = [
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'email', header: 'Email' },
  { accessorKey: 'role', header: 'Role' },
  // 固定宽度测试列：width 应映射为 80px，与其余自动宽度列形成对比
  { accessorKey: 'status', header: 'Status', width: 80 },
]

const data = ref([
  { name: 'Alice Johnson', email: 'alice@example.com', role: 'Admin', status: 'Active' },
  { name: 'Bob Smith', email: 'bob@example.com', role: 'Editor', status: 'Pending review by security team for access escalation' },
  { name: 'Charlie Brown', email: 'charlie@example.com', role: 'User', status: 'Inactive' },
  { name: 'Diana Prince', email: 'diana@example.com', role: 'Editor', status: 'Active' },
  { name: 'Eve Williams', email: 'eve@example.com', role: 'User', status: 'Active' },
])

const filteredData = computed(() => {
  if (!search.value) return data.value
  const q = search.value.toLowerCase()
  return data.value.filter(row =>
    row.name.toLowerCase().includes(q) || row.email.toLowerCase().includes(q),
  )
})

// ======== 列固定（Column Pinning） ========

/// 列固定 header 辅助函数：渲染 UButton 切换 pin 状态(参照 Nuxt UI 官方示例)
function pinningHeader(column: Column<any>, label: string, position: 'left' | 'right' = 'left') {
  const isPinned = column.getIsPinned()
  return h(UButton, {
    color: 'neutral',
    variant: 'ghost',
    label,
    icon: isPinned ? 'i-lucide-pin-off' : 'i-lucide-pin',
    class: '-mx-2.5',
    size: 'xs' as const,
    onClick() {
      column.pin(isPinned === position ? false : position)
    },
  })
}

const pinningData = ref([
  { name: 'Alice', email: 'alice@corp.com', role: 'Admin', dept: '技术中心', location: '杭州', status: 'Active', joinDate: '2024-03-15' },
  { name: 'Bob', email: 'bob.smith@corp.com', role: 'Editor', dept: '产品中心', location: '北京', status: 'Active', joinDate: '2024-06-01' },
  { name: 'Charlie', email: 'charlie.brown@corp.com', role: 'Viewer', dept: '设计部', location: '上海', status: 'Inactive', joinDate: '2023-11-20' },
  { name: 'Diana', email: 'diana.prince@corp.com', role: 'Admin', dept: '运营部', location: '深圳', status: 'Active', joinDate: '2025-01-10' },
  { name: 'Edward', email: 'edward.elric@corp.com', role: 'Editor', dept: '技术中心', location: '杭州', status: 'Active', joinDate: '2024-09-05' },
  { name: 'Fiona', email: 'fiona.green@corp.com', role: 'Viewer', dept: '市场部', location: '广州', status: 'Active', joinDate: '2025-04-18' },
  { name: 'George', email: 'george.w@corp.com', role: 'Editor', dept: '技术中心', location: '成都', status: 'Inactive', joinDate: '2023-07-30' },
])

/// pinning 列使用 width(非 size)，确保视觉宽度与 TanStack 内部宽度一致，
/// 消除 table-layout: auto 下列固定间隙。width 会被 MTable 映射为 max-width 硬上限。
const pinningColumns = [
  { accessorKey: 'name', header: ({ column }: { column: Column<any> }) => pinningHeader(column, 'Name'), width: 120 },
  { accessorKey: 'email', header: ({ column }: { column: Column<any> }) => pinningHeader(column, 'Email'), width: 220 },
  { accessorKey: 'role', header: ({ column }: { column: Column<any> }) => pinningHeader(column, 'Role'), width: 100 },
  { accessorKey: 'dept', header: ({ column }: { column: Column<any> }) => pinningHeader(column, '部门'), width: 130 },
  { accessorKey: 'location', header: ({ column }: { column: Column<any> }) => pinningHeader(column, '城市'), width: 100 },
  { accessorKey: 'status', header: ({ column }: { column: Column<any> }) => pinningHeader(column, 'Status'), width: 100 },
  { accessorKey: 'joinDate', header: ({ column }: { column: Column<any> }) => pinningHeader(column, '入职'), width: 130 },
  { accessorKey: 'actions', header: '操作', width: 120, cell: () => {
    const btn = resolveComponent('UButton')
    return h('div', { class: 'flex gap-1' }, [
      h(btn, { color: 'neutral', variant: 'ghost', icon: 'i-lucide-pencil', size: 'xs' as const }),
      h(btn, { color: 'neutral', variant: 'ghost', icon: 'i-lucide-trash', size: 'xs' as const }),
    ])
  }},
]

const pinningState = ref({
  left: ['name', 'email'],
  right: ['actions'],
})

// ======== 树形嵌套 ========

interface TreeRow {
  name: string
  email: string
  role: string
  status: string
  memberCount?: number
  children?: TreeRow[]
}

const treeData: TreeRow[] = [
  {
    name: '技术中心',
    email: 'tech@example.com',
    role: '部门',
    status: '活跃',
    memberCount: 12,
    children: [
      {
        name: '前端组',
        email: 'fe@example.com',
        role: '团队',
        status: '活跃',
        memberCount: 5,
        children: [
          { name: '张三', email: 'zhangsan@example.com', role: '工程师', status: '活跃' },
          { name: '李四', email: 'lisi@example.com', role: '工程师', status: '活跃' },
          { name: '王五', email: 'wangwu@example.com', role: '实习生', status: '活跃' },
        ],
      },
      {
        name: '后端组',
        email: 'be@example.com',
        role: '团队',
        status: '活跃',
        memberCount: 4,
        children: [
          { name: '赵六', email: 'zhaoliu@example.com', role: '工程师', status: '活跃' },
        ],
      },
    ],
  },
  {
    name: '产品中心',
    email: 'product@example.com',
    role: '部门',
    status: '活跃',
    memberCount: 6,
    children: [
      { name: '孙七', email: 'sunqi@example.com', role: '产品经理', status: '活跃' },
      { name: '周八', email: 'zhouba@example.com', role: '设计师', status: '活跃' },
    ],
  },
]

const treeColumns = [
  { accessorKey: 'name', header: '名称' },
  { accessorKey: 'email', header: '邮箱' },
  { accessorKey: 'role', header: '角色' },
  { accessorKey: 'status', header: '状态' },
  { accessorKey: 'memberCount', header: '人数' },
]

const treeExpanded = ref<Record<string, boolean>>({ '0': true })

// ======== 行详情展开 ========

const detailData = data.value.slice(0, 3)
const detailColumns = [
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'email', header: 'Email' },
  { accessorKey: 'role', header: 'Role' },
]

const detailExpanded = ref<Record<string, boolean>>({})
</script>
