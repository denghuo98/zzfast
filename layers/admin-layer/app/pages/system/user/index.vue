<script setup lang="ts">
import * as z from 'zod'
import type { ColumnOption } from '#layers/admin-layer/app/types/table'
import type { AdminMemberItem, AdminRoleItem } from '../../../types/admin-api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const toast = useToast()

// ==================== 搜索 ====================

const searchSchema = z.object({
  username: z.string().optional().meta({ title: '账号' }),
  realName: z.string().optional().meta({ title: '姓名' }),
  status: z.enum(['全部', '启用', '禁用']).optional().meta({ title: '状态' }),
})

const searchForm = reactive<z.infer<typeof searchSchema>>({ username: '', realName: '', status: '全部' })

function statusValue(v: string | undefined) {
  return v === '启用' ? 1 : v === '禁用' ? 2 : undefined
}

// ==================== 列表 ====================

const list = ref<AdminMemberItem[]>([])
const loading = ref(false)
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get<{ list: AdminMemberItem[]; totalCount: number }>('/admin/member/list', {
      page: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username || undefined,
      realName: searchForm.realName || undefined,
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

function doSearch() {
  pagination.page = 1
  fetchList()
}

function doReset() {
  searchForm.username = ''
  searchForm.realName = ''
  searchForm.status = '全部'
  pagination.page = 1
  fetchList()
}

function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  fetchList()
}

function handleCurrentChange(page: number) {
  pagination.page = page
  fetchList()
}

// ==================== 列定义 ====================

const columns: ColumnOption<AdminMemberItem>[] = [
  { accessorKey: 'id', header: 'ID', width: 70 },
  { accessorKey: 'username', header: '账号' },
  { accessorKey: 'realName', header: '姓名' },
  { accessorKey: 'email', header: '邮箱' },
  { accessorKey: 'mobile', header: '手机号' },
  { accessorKey: 'roles', header: '角色', useSlot: true },
  { accessorKey: 'status', header: '状态', useSlot: true, width: 90 },
  { accessorKey: 'createdAt', header: '创建时间', width: 170 },
  { accessorKey: 'actions', header: '操作', useSlot: true, width: 200, fixed: 'right' },
]

// ==================== 编辑弹窗 ====================

const modalOpen = ref(false)
const editingId = ref<number>()
const formState = reactive({
  username: '',
  realName: '',
  password: '',
  roleIds: [] as number[],
  email: '',
  mobile: '',
  remark: '',
  status: 1,
})

const roleOptions = ref<AdminRoleItem[]>([])
const saving = ref(false)

async function loadRoles() {
  try {
    const res = await api.get<{ list: AdminRoleItem[] }>('/admin/role/list', { page: 1, pageSize: 1000 })
    roleOptions.value = res.list ?? []
  } catch {}
}

function openCreate() {
  editingId.value = undefined
  Object.assign(formState, { username: '', realName: '', password: '', roleIds: [], email: '', mobile: '', remark: '', status: 1 })
  modalOpen.value = true
}

function openEdit(row: AdminMemberItem) {
  editingId.value = row.id
  Object.assign(formState, {
    username: row.username,
    realName: row.realName,
    password: '',
    roleIds: [...(row.roleIds ?? [])],
    email: row.email,
    mobile: row.mobile,
    remark: row.remark,
    status: row.status,
  })
  modalOpen.value = true
}

async function save() {
  saving.value = true
  try {
    await api.post('/admin/member/edit', {
      id: editingId.value,
      username: formState.username,
      realName: formState.realName,
      password: formState.password || undefined,
      roleIds: formState.roleIds,
      email: formState.email,
      mobile: formState.mobile,
      remark: formState.remark,
      status: formState.status,
    })
    toast.add({ title: '保存成功', color: 'success' })
    modalOpen.value = false
    fetchList()
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e.message, color: 'error' })
  } finally {
    saving.value = false
  }
}

async function remove(row: AdminMemberItem) {
  await api.post('/admin/member/delete', { id: row.id }).catch((e: any) => {
    toast.add({ title: '删除失败', description: e.message, color: 'error' })
    throw e
  })
  toast.add({ title: '已删除', color: 'success' })
  fetchList()
}

async function resetPassword(row: AdminMemberItem) {
  await api.post('/admin/member/reset-pwd', { id: row.id })
  toast.add({ title: '已重置为默认密码', color: 'success' })
}
</script>

<template>
  <div class="flex flex-col gap-4 h-full">
    <MFormSearch :schema="searchSchema" :state="searchForm" @search="doSearch" @reset="doReset" />

    <UCard :ui="{ body: 'p-0 h-full' }" class="flex-1 min-h-50 overflow-hidden">
      <MTableHeader :loading="loading" @refresh="fetchList">
        <template #left>
          <span class="text-sm font-medium">用户列表</span>
        </template>
        <template #right>
          <UButton size="xs" color="primary" icon="i-lucide-plus" @click="openCreate">新增用户</UButton>
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
        <template #roles-data="{ row }">
          <div class="flex flex-wrap gap-1">
            <UBadge v-for="r in row.original.roles" :key="r.id" color="neutral" variant="subtle">{{ r.name }}</UBadge>
          </div>
        </template>
        <template #status-data="{ row }">
          <UBadge :color="row.original.status === 1 ? 'success' : 'neutral'" variant="subtle">
            {{ row.original.status === 1 ? '启用' : '禁用' }}
          </UBadge>
        </template>
        <template #actions-data="{ row }">
          <div class="flex gap-1">
            <UButton size="xs" variant="ghost" color="primary" @click="openEdit(row.original)">编辑</UButton>
            <UButton size="xs" variant="ghost" color="warning" @click="resetPassword(row.original)">重置密码</UButton>
            <UButton size="xs" variant="ghost" color="error" @click="remove(row.original)">删除</UButton>
          </div>
        </template>
      </MTable>
    </UCard>

    <UModal v-model:open="modalOpen" :title="editingId ? '编辑用户' : '新增用户'">
      <template #body>
        <div class="space-y-4">
          <UFormField label="账号" required>
            <UInput v-model="formState.username" placeholder="登录账号" :disabled="!!editingId" />
          </UFormField>
          <UFormField label="姓名">
            <UInput v-model="formState.realName" placeholder="真实姓名" />
          </UFormField>
          <UFormField v-if="!editingId" label="密码" required>
            <UInput v-model="formState.password" type="password" placeholder="登录密码" />
          </UFormField>
          <UFormField label="角色">
            <USelectMenu v-model="formState.roleIds" :items="roleOptions.map(r => ({ label: r.name, value: r.id }))" multiple placeholder="选择角色" />
          </UFormField>
          <UFormField label="邮箱">
            <UInput v-model="formState.email" placeholder="邮箱" />
          </UFormField>
          <UFormField label="手机号">
            <UInput v-model="formState.mobile" placeholder="手机号" />
          </UFormField>
          <UFormField label="备注">
            <UTextarea v-model="formState.remark" placeholder="备注" :rows="2" />
          </UFormField>
          <UFormField label="状态">
            <USelect v-model="formState.status" :items="[{ label: '启用', value: 1 }, { label: '禁用', value: 2 }]" />
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="ghost" @click="modalOpen = false">取消</UButton>
          <UButton color="primary" :loading="saving" @click="save">保存</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>
