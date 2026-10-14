<script setup lang="ts">
import * as z from 'zod'
import type { ColumnOption } from '#layers/admin-layer/app/types/table'
import type { AdminLoginLogItem } from '../../../types/admin-api'
import { ADMIN_TOKEN_COOKIE } from '../../../composables/useApi'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const toast = useToast()

const searchSchema = z.object({
  username: z.string().optional().meta({ title: '账号' }),
  loginIp: z.string().optional().meta({ title: '登录 IP' }),
  status: z.enum(['全部', '成功', '失败']).optional().meta({ title: '状态' }),
})

const searchForm = reactive<z.infer<typeof searchSchema>>({ username: '', loginIp: '', status: '全部' })

function statusValue(v: string | undefined) {
  return v === '成功' ? 1 : v === '失败' ? 2 : undefined
}

const list = ref<AdminLoginLogItem[]>([])
const loading = ref(false)
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get<{ list: AdminLoginLogItem[]; totalCount: number }>('/admin/system/login-log/list', {
      page: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username || undefined,
      loginIp: searchForm.loginIp || undefined,
      status: statusValue(searchForm.status),
    })
    list.value = res.list ?? []
    pagination.total = res.totalCount ?? 0
  } catch (e: any) {
    toast.add({ title: '加载失败', description: e.message, color: 'error' })
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)

function doSearch() { pagination.page = 1; fetchList() }
function doReset() { searchForm.username = ''; searchForm.loginIp = ''; searchForm.status = '全部'; pagination.page = 1; fetchList() }
function handleSizeChange(size: number) { pagination.pageSize = size; pagination.page = 1; fetchList() }
function handleCurrentChange(page: number) { pagination.page = page; fetchList() }

const columns: ColumnOption<AdminLoginLogItem>[] = [
  { accessorKey: 'id', header: 'ID', width: 70 },
  { accessorKey: 'username', header: '账号' },
  { accessorKey: 'loginIp', header: '登录 IP' },
  { accessorKey: 'province', header: '位置', useSlot: true },
  { accessorKey: 'status', header: '状态', useSlot: true, width: 90 },
  { accessorKey: 'os', header: '系统', useSlot: true, width: 120 },
  { accessorKey: 'loginAt', header: '登录时间', width: 170 },
  { accessorKey: 'actions', header: '操作', useSlot: true, width: 80, fixed: 'right' },
]

async function remove(row: AdminLoginLogItem) {
  await api.del('/admin/system/login-log/delete', { id: row.id }).catch((e: any) => {
    toast.add({ title: '删除失败', description: e.message, color: 'error' })
    throw e
  })
  toast.add({ title: '已删除', color: 'success' })
  fetchList()
}

async function exportLog() {
  try {
    const token = useCookie<string>(ADMIN_TOKEN_COOKIE)
    const blob = await $fetch<Blob>('/admin/system/login-log/export', {
      method: 'GET',
      params: {
        username: searchForm.username || undefined,
        loginIp: searchForm.loginIp || undefined,
        status: statusValue(searchForm.status),
      },
      headers: token.value ? { Authorization: `Bearer ${token.value}` } : {},
      responseType: 'blob',
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'login-log.xlsx'
    a.click()
    URL.revokeObjectURL(url)
  } catch (e: any) {
    toast.add({ title: '导出失败', description: e.message, color: 'error' })
  }
}
</script>

<template>
  <div class="flex flex-col gap-4 h-full">
    <MFormSearch :schema="searchSchema" :state="searchForm" @search="doSearch" @reset="doReset" />

    <UCard :ui="{ body: 'p-0 h-full' }" class="flex-1 min-h-50 overflow-hidden">
      <MTableHeader :loading="loading" @refresh="fetchList">
        <template #left>
          <span class="text-sm font-medium">登录日志</span>
        </template>
        <template #right>
          <UButton size="xs" color="neutral" variant="outline" icon="i-lucide-download" @click="exportLog">导出</UButton>
        </template>
      </MTableHeader>

      <MTable
        :loading="loading"
        :columns="columns"
        :data="list"
        :pagination="pagination"
        :pagination-options="{ pageSizes: [10, 20, 50] }"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #province-data="{ row }">
          {{ [row.original.province, row.original.city].filter(Boolean).join(' ') || '—' }}
        </template>
        <template #status-data="{ row }">
          <UBadge :color="row.original.status === 1 ? 'success' : 'error'" variant="subtle">
            {{ row.original.status === 1 ? '成功' : '失败' }}
          </UBadge>
        </template>
        <template #os-data="{ row }">
          {{ [row.original.os, row.original.browser].filter(Boolean).join(' / ') || '—' }}
        </template>
        <template #actions-data="{ row }">
          <UButton size="xs" variant="ghost" color="error" @click="remove(row.original)">删除</UButton>
        </template>
      </MTable>
    </UCard>
  </div>
</template>
