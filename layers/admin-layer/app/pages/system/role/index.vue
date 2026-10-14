<script setup lang="ts">
import * as z from 'zod'
import type { ColumnOption } from '#layers/admin-layer/app/types/table'
import type { AdminMenuItem, AdminRoleItem } from '../../../types/admin-api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const toast = useToast()

const searchSchema = z.object({
  name: z.string().optional().meta({ title: '角色名称' }),
  key: z.string().optional().meta({ title: '角色标识' }),
  status: z.enum(['全部', '启用', '禁用']).optional().meta({ title: '状态' }),
})

const searchForm = reactive<z.infer<typeof searchSchema>>({ name: '', key: '', status: '全部' })

function statusValue(v: string | undefined) {
  return v === '启用' ? 1 : v === '禁用' ? 2 : undefined
}

const list = ref<AdminRoleItem[]>([])
const loading = ref(false)
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get<{ list: AdminRoleItem[]; totalCount: number }>('/admin/role/list', {
      page: pagination.page,
      pageSize: pagination.pageSize,
      name: searchForm.name || undefined,
      key: searchForm.key || undefined,
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
function doReset() { searchForm.name = ''; searchForm.key = ''; searchForm.status = '全部'; pagination.page = 1; fetchList() }
function handleSizeChange(size: number) { pagination.pageSize = size; pagination.page = 1; fetchList() }
function handleCurrentChange(page: number) { pagination.page = page; fetchList() }

const columns: ColumnOption<AdminRoleItem>[] = [
  { accessorKey: 'id', header: 'ID', width: 70 },
  { accessorKey: 'name', header: '角色名称' },
  { accessorKey: 'key', header: '角色标识' },
  { accessorKey: 'remark', header: '备注' },
  { accessorKey: 'status', header: '状态', useSlot: true, width: 90 },
  { accessorKey: 'createdAt', header: '创建时间', width: 170 },
  { accessorKey: 'actions', header: '操作', useSlot: true, width: 140, fixed: 'right' },
]

// ==================== 编辑弹窗 ====================

const modalOpen = ref(false)
const editingId = ref<number>()
const saving = ref(false)
const formState = reactive({
  name: '',
  key: '',
  remark: '',
  sort: 0,
  status: 1,
  menuIds: [] as number[],
})

const menuTree = ref<AdminMenuItem[]>([])
const menuOptions = computed(() => flattenMenu(menuTree.value))

function flattenMenu(items: AdminMenuItem[], depth = 0): { label: string; value: number }[] {
  return items.flatMap((item) => {
    const prefix = depth > 0 ? '　'.repeat(depth - 1) + '└ ' : ''
    const current = [{ label: prefix + item.title, value: item.id }]
    return item.children?.length ? [...current, ...flattenMenu(item.children, depth + 1)] : current
  })
}

async function loadMenus() {
  try {
    const res = await api.get<{ list: AdminMenuItem[] }>('/admin/menu/list', { page: 1, pageSize: 1000 })
    menuTree.value = res.list ?? []
  } catch {}
}

function openCreate() {
  editingId.value = undefined
  Object.assign(formState, { name: '', key: '', remark: '', sort: 0, status: 1, menuIds: [] })
  modalOpen.value = true
}

function openEdit(row: AdminRoleItem) {
  editingId.value = row.id
  Object.assign(formState, {
    name: row.name,
    key: row.key,
    remark: row.remark,
    sort: row.sort,
    status: row.status,
    menuIds: [...(row.menuIds ?? [])],
  })
  modalOpen.value = true
}

async function save() {
  saving.value = true
  try {
    await api.post('/admin/role/edit', {
      id: editingId.value,
      name: formState.name,
      key: formState.key,
      remark: formState.remark,
      sort: formState.sort,
      status: formState.status,
      menuIds: formState.menuIds,
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

async function remove(row: AdminRoleItem) {
  await api.post('/admin/role/delete', { id: row.id }).catch((e: any) => {
    toast.add({ title: '删除失败', description: e.message, color: 'error' })
    throw e
  })
  toast.add({ title: '已删除', color: 'success' })
  fetchList()
}
</script>

<template>
  <div class="flex flex-col gap-4 h-full">
    <MFormSearch :schema="searchSchema" :state="searchForm" @search="doSearch" @reset="doReset" />

    <UCard :ui="{ body: 'p-0 h-full' }" class="flex-1 min-h-50 overflow-hidden">
      <MTableHeader :loading="loading" @refresh="fetchList">
        <template #left>
          <span class="text-sm font-medium">角色列表</span>
        </template>
        <template #right>
          <UButton size="xs" color="primary" icon="i-lucide-plus" @click="openCreate">新增角色</UButton>
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
        <template #status-data="{ row }">
          <UBadge :color="row.original.status === 1 ? 'success' : 'neutral'" variant="subtle">
            {{ row.original.status === 1 ? '启用' : '禁用' }}
          </UBadge>
        </template>
        <template #actions-data="{ row }">
          <div class="flex gap-1">
            <UButton size="xs" variant="ghost" color="primary" @click="openEdit(row.original)">编辑</UButton>
            <UButton size="xs" variant="ghost" color="error" @click="remove(row.original)">删除</UButton>
          </div>
        </template>
      </MTable>
    </UCard>

    <UModal v-model:open="modalOpen" :title="editingId ? '编辑角色' : '新增角色'" @open="loadMenus">
      <template #body>
        <div class="space-y-4">
          <UFormField label="角色名称" required>
            <UInput v-model="formState.name" placeholder="角色名称" />
          </UFormField>
          <UFormField label="角色标识" required>
            <UInput v-model="formState.key" placeholder="角色标识（如 admin）" />
          </UFormField>
          <UFormField label="菜单权限">
            <USelectMenu v-model="formState.menuIds" :items="menuOptions" multiple placeholder="选择菜单权限" />
          </UFormField>
          <UFormField label="排序">
            <UInput v-model.number="formState.sort" type="number" placeholder="排序号" />
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
