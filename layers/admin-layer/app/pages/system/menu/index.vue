<script setup lang="ts">
import type { ColumnOption } from '#layers/admin-layer/app/types/table'
import type { AdminMenuItem } from '../../../types/admin-api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const toast = useToast()

const list = ref<AdminMenuItem[]>([])
const loading = ref(false)

async function fetchList() {
  loading.value = true
  try {
    const res = await api.get<{ list: AdminMenuItem[] }>('/admin/menu/list', { page: 1, pageSize: 1000 })
    list.value = res.list ?? []
  } catch (e: any) {
    toast.add({ title: '加载失败', description: e.message, color: 'error' })
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)

// 扁平化菜单树（带缩进），用于表格展示与父级下拉
const flatRows = computed(() => flatten(list.value))

function flatten(items: AdminMenuItem[], depth = 0): Array<AdminMenuItem & { indentTitle: string }> {
  return items.flatMap((item) => {
    const row = { ...item, indentTitle: (depth > 0 ? '　'.repeat(depth) : '') + item.title }
    return item.children?.length ? [row, ...flatten(item.children, depth + 1)] : [row]
  })
}

const columns: ColumnOption[] = [
  { accessorKey: 'indentTitle', header: '菜单名称' },
  { accessorKey: 'type', header: '类型', useSlot: true, width: 80 },
  { accessorKey: 'name', header: '路由别名' },
  { accessorKey: 'path', header: '路由地址' },
  { accessorKey: 'icon', header: '图标', useSlot: true, width: 80 },
  { accessorKey: 'sort', header: '排序', width: 70 },
  { accessorKey: 'status', header: '状态', useSlot: true, width: 80 },
  { accessorKey: 'actions', header: '操作', useSlot: true, width: 200, fixed: 'right' },
]

const typeMap: Record<number, { label: string; color: string }> = {
  1: { label: '目录', color: 'primary' },
  2: { label: '菜单', color: 'info' },
  3: { label: '按钮', color: 'neutral' },
}

// ==================== 编辑弹窗 ====================

const modalOpen = ref(false)
const editingId = ref<number>()
const saving = ref(false)
const formState = reactive({
  pid: 0,
  title: '',
  name: '',
  path: '',
  icon: '',
  type: 2,
  redirect: '',
  component: '',
  sort: 0,
  status: 1,
  permissions: '',
  remark: '',
})

const parentOptions = computed(() => [
  { label: '根目录', value: 0 },
  ...flatRows.value
    .filter(r => r.type !== 3)
    .map(r => ({ label: r.indentTitle, value: r.id })),
])

function openCreate(pid = 0) {
  editingId.value = undefined
  Object.assign(formState, { pid, title: '', name: '', path: '', icon: '', type: 2, redirect: '', component: '', sort: 0, status: 1, permissions: '', remark: '' })
  modalOpen.value = true
}

function openEdit(row: AdminMenuItem) {
  editingId.value = row.id
  Object.assign(formState, {
    pid: row.pid,
    title: row.title,
    name: row.name,
    path: row.path,
    icon: row.icon,
    type: row.type,
    redirect: row.redirect,
    component: row.component,
    sort: row.sort,
    status: row.status,
    permissions: row.permissions,
    remark: row.remark,
  })
  modalOpen.value = true
}

async function save() {
  saving.value = true
  try {
    await api.post('/admin/menu/edit', {
      id: editingId.value,
      pid: formState.pid,
      title: formState.title,
      name: formState.name,
      path: formState.path,
      icon: formState.icon,
      type: formState.type,
      redirect: formState.redirect,
      component: formState.component,
      sort: formState.sort,
      status: formState.status,
      permissions: formState.permissions,
      remark: formState.remark,
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

async function remove(row: AdminMenuItem) {
  await api.post('/admin/menu/delete', { id: row.id }).catch((e: any) => {
    toast.add({ title: '删除失败', description: e.message, color: 'error' })
    throw e
  })
  toast.add({ title: '已删除', color: 'success' })
  fetchList()
}
</script>

<template>
  <div class="flex flex-col gap-4 h-full">
    <UCard :ui="{ body: 'p-0 h-full' }" class="flex-1 min-h-50 overflow-hidden">
      <MTableHeader :loading="loading" @refresh="fetchList">
        <template #left>
          <span class="text-sm font-medium">菜单列表</span>
        </template>
        <template #right>
          <UButton size="xs" color="primary" icon="i-lucide-plus" @click="openCreate(0)">新增菜单</UButton>
        </template>
      </MTableHeader>

      <MTable :loading="loading" :columns="columns" :data="flatRows">
        <template #type-data="{ row }">
          <UBadge :color="typeMap[row.original.type]?.color || 'neutral'" variant="subtle">
            {{ typeMap[row.original.type]?.label }}
          </UBadge>
        </template>
        <template #icon-data="{ row }">
          <UIcon v-if="row.original.icon" :name="row.original.icon" class="w-4 h-4" />
        </template>
        <template #status-data="{ row }">
          <UBadge :color="row.original.status === 1 ? 'success' : 'neutral'" variant="subtle">
            {{ row.original.status === 1 ? '启用' : '禁用' }}
          </UBadge>
        </template>
        <template #actions-data="{ row }">
          <div class="flex gap-1">
            <UButton v-if="row.original.type !== 3" size="xs" variant="ghost" color="info" @click="openCreate(row.original.id)">新增子项</UButton>
            <UButton size="xs" variant="ghost" color="primary" @click="openEdit(row.original)">编辑</UButton>
            <UButton size="xs" variant="ghost" color="error" @click="remove(row.original)">删除</UButton>
          </div>
        </template>
      </MTable>
    </UCard>

    <UModal v-model:open="modalOpen" :title="editingId ? '编辑菜单' : '新增菜单'">
      <template #body>
        <div class="space-y-4">
          <UFormField label="上级菜单">
            <USelect v-model="formState.pid" :items="parentOptions" />
          </UFormField>
          <UFormField label="菜单名称" required>
            <UInput v-model="formState.title" placeholder="菜单名称" />
          </UFormField>
          <UFormField label="类型">
            <USelect v-model="formState.type" :items="[{ label: '目录', value: 1 }, { label: '菜单', value: 2 }, { label: '按钮', value: 3 }]" />
          </UFormField>
          <UFormField label="路由别名" required>
            <UInput v-model="formState.name" placeholder="路由别名" />
          </UFormField>
          <UFormField v-if="formState.type !== 3" label="路由地址" required>
            <UInput v-model="formState.path" placeholder="如 /system/user" />
          </UFormField>
          <UFormField v-if="formState.type === 2" label="组件路径">
            <UInput v-model="formState.component" placeholder="组件路径" />
          </UFormField>
          <UFormField label="图标">
            <UInput v-model="formState.icon" placeholder="如 i-lucide-users" />
          </UFormField>
          <UFormField label="权限标识">
            <UInput v-model="formState.permissions" placeholder="权限标识" />
          </UFormField>
          <UFormField label="排序">
            <UInput v-model.number="formState.sort" type="number" placeholder="排序号" />
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
