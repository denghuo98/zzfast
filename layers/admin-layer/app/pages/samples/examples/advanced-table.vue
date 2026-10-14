<script setup lang="ts">
import * as z from 'zod'
import type { ColumnOption } from '#layers/admin-layer/app/types/table'

definePageMeta({ layout: 'admin' })

const toast = useToast()

// ==================== 搜索表单 Schema ====================

const searchSchema = z.object({
  keyword: z.string().optional().meta({
    title: '关键词',
    hoverTip: '搜索用户名称、邮箱或角色',
  }),
  role: z.enum(['全部', '管理员', '编辑', '普通用户']).optional().meta({ title: '角色' }),
  status: z.enum(['全部', '活跃', '待审核', '禁用']).optional().meta({ title: '状态' }),
  dateFrom: z.string().optional().meta({ title: '开始日期' }),
  dateTo: z.string().optional().meta({ title: '结束日期' }),
})

const searchParams = reactive<z.infer<typeof searchSchema>>({
  keyword: '',
  role: '全部',
  status: '全部',
  dateFrom: '',
  dateTo: '',
})

// ==================== Mock 数据 ====================

interface UserRecord {
  id: number
  name: string
  email: string
  role: string
  status: string
  lastLogin: string
}

const mockUsers: UserRecord[] = Array.from({ length: 26 }, (_, i) => ({
  id: i + 1,
  name: `用户 ${String.fromCodePoint(65 + (i % 26))}`,
  email: `user${i + 1}@example.com`,
  role: ['管理员', '编辑', '普通用户'][i % 3],
  status: ['活跃', '待审核', '禁用'][i % 3],
  lastLogin: `2026-06-${String(i + 1).padStart(2, '0')} ${String(10 + (i % 14)).padStart(2, '0')}:${String(i % 60).padStart(2, '0')}`,
}))

// ==================== 搜索逻辑 ====================

const loading = ref(false)
const filteredData = ref<UserRecord[]>(mockUsers)
const currentPage = ref(1)
const pageSize = ref(10)

function filterUsers(data: z.infer<typeof searchSchema>): UserRecord[] {
  let list = [...mockUsers]
  if (data.keyword) {
    const q = data.keyword.toLowerCase()
    list = list.filter(u =>
      u.name.includes(q) || u.email.includes(q) || u.role.includes(q),
    )
  }
  if (data.role && data.role !== '全部') {
    list = list.filter(u => u.role === data.role)
  }
  if (data.status && data.status !== '全部') {
    list = list.filter(u => u.status === data.status)
  }
  return list
}

function doSearch(data: z.infer<typeof searchSchema>) {
  loading.value = true
  // 模拟异步搜索
  setTimeout(() => {
    filteredData.value = filterUsers(data)
    currentPage.value = 1
    loading.value = false
  }, 300)
}

function addUser() {
  const nextId = mockUsers.length + 1
  mockUsers.push({
    id: nextId,
    name: `新用户 ${nextId}`,
    email: `user${nextId}@example.com`,
    role: '普通用户',
    status: '活跃',
    lastLogin: '—',
  })

  // 追加后按当前搜索条件同步刷新，并跳到最后一页让新用户可见
  filteredData.value = filterUsers(searchParams as z.infer<typeof searchSchema>)
  currentPage.value = Math.max(1, Math.ceil(filteredData.value.length / pageSize.value))
  toast.add({ title: '已添加', description: `已模拟添加「新用户 ${nextId}」` })
}

function doReset() {
  searchParams.keyword = ''
  searchParams.role = '全部'
  searchParams.status = '全部'
  searchParams.dateFrom = ''
  searchParams.dateTo = ''
  filteredData.value = mockUsers
  currentPage.value = 1
}

// ==================== 表格列定义 ====================

const columnsFactory = (): ColumnOption[] => [
  { accessorKey: '__selection__', type: 'selection', header: '' },
  { accessorKey: '__index__', type: 'index', header: '序号', width: 80 },
  { accessorKey: 'name', header: '名称' },
  { accessorKey: 'email', header: '邮箱' },
  { accessorKey: 'role', header: '角色', useSlot: true },
  { accessorKey: 'status', header: '状态', useSlot: true },
  { accessorKey: 'lastLogin', header: '最后登录' },
]

const columns = computed(() => columnsFactory())

// ==================== 分页 ====================

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  total: filteredData.value.length,
}))

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredData.value.slice(start, start + pageSize.value)
})

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(page: number) {
  currentPage.value = page
}

function refreshData() {
  currentPage.value = 1
  doSearch(searchParams as z.infer<typeof searchSchema>)
}

const selectedRows = ref<any[]>([])
</script>

<template>
  <div class="flex flex-col gap-4 h-full">
    <!-- 搜索表单 -->
    <MFormSearch
      :schema="searchSchema"
      :state="searchParams"
      @search="doSearch"
      @reset="doReset"
    />

    <!-- 表格 + 选中信息容器：flex-1 填充剩余空间 -->
    <div class="flex-1 min-h-0 flex flex-col gap-4">
      <UCard :ui="{ body: 'p-0 h-full' }" class="flex-1 min-h-50 overflow-hidden">
        <MTableHeader :loading="loading" @refresh="refreshData">
          <template #left>
            <span class="text-sm font-medium">用户列表</span>
          </template>
          <template #right>
            <UButton size="xs" color="primary" icon="i-lucide-plus" @click="addUser">
              添加用户
            </UButton>
          </template>
        </MTableHeader>

        <MTable
          :loading="loading"
          :columns="columns"
          :data="paginatedData"
          :pagination="pagination"
          :show-table-header="true"
          :pagination-options="{ pageSizes: [10, 20, 50] }"
          @pagination:size-change="handleSizeChange"
          @pagination:current-change="handleCurrentChange"
          @selection-change="(s: any[]) => selectedRows = s"
        >
          <template #role-data="{ row }">
            <UBadge
              :color="row.original.role === '管理员' ? 'primary' : row.original.role === '编辑' ? 'info' : 'neutral'"
              variant="subtle"
            >
              {{ row.original.role }}
            </UBadge>
          </template>
          <template #status-data="{ row }">
            <UBadge
              :color="row.original.status === '活跃' ? 'success' : row.original.status === '待审核' ? 'warning' : 'neutral'"
              variant="subtle"
            >
              {{ row.original.status }}
            </UBadge>
          </template>
        </MTable>
      </UCard>

      <!-- 选中行信息 -->
      <UCard v-if="selectedRows.length" class="text-sm shrink-0">
        <p class="text-muted-foreground">
          已选择 {{ selectedRows.length }} 行：
        </p>
        <p v-for="(row, i) in selectedRows" :key="i" class="mt-1">
          #{{ row.original.id }} {{ row.original.name }} — {{ row.original.email }}
        </p>
      </UCard>
    </div>
  </div>
</template>
