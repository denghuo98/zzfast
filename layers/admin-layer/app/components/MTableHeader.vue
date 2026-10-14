<!-- 表格头部，包含表格大小、刷新、全屏、列设置、其他设置 -->
<template>
  <div id="m-table-header" class="flex justify-between py-[6px]">
    <div class="flex flex-wrap gap-2">
      <slot name="left" />
    </div>

    <div class="flex flex-wrap gap-2">
      <slot name="right" />

      <!-- 刷新按钮 -->
      <UButton
        v-if="shouldShow('refresh')"
        :loading="loading && isManualRefresh"
        icon="i-lucide-refresh-cw"
        variant="soft"
        color="neutral"
        size="xs"
        @click="refresh"
      />

      <!-- 密度 -->
      <UDropdownMenu v-if="shouldShow('size')" :items="tableSizeDropDownItems">
        <UButton icon="i-lucide-ratio" variant="soft" color="neutral" size="xs" />
      </UDropdownMenu>

      <!-- 列设置 -->
      <UPopover v-if="shouldShow('columns')" placement="bottom" trigger="click">
        <UButton icon="i-lucide-columns-3" variant="soft" color="neutral" size="xs" />
        <template #content>
          <div class="p-1">
            <div
              v-for="column in columnChecks"
              :key="column.accessorKey"
              class="flex items-center gap-2"
              :class="{ 'fixed-column': column.fixed }"
            >
              <UCheckbox v-model="column.checked" :label="column.header" />
            </div>
          </div>
        </template>
      </UPopover>

      <!-- 其他设置 -->
      <UPopover
        v-if="shouldShow('settings')"
        placement="bottom"
        trigger="click"
      >
        <UButton icon="i-lucide-settings" variant="soft" color="neutral" size="xs" />
        <template #content>
          <div class="p-2 flex flex-col gap-2">
            <UCheckbox v-if="showZebra" v-model="isZebra" label="斑马纹" />
            <UCheckbox v-if="showBorder" v-model="isBorder" label="边框" />
            <UCheckbox v-if="showHeaderBackground" v-model="isHeaderBackground" label="表头背景" />
          </div>
        </template>
      </UPopover>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { TableSizeEnum } from '../enums/formEnum'
import { useTableStore } from '../stores/table'
import type { ColumnOption } from '../types/table'

defineOptions({ name: 'MTableHeader' })

interface Props {
  showZebra?: boolean
  showBorder?: boolean
  showHeaderBackground?: boolean
  /** 组件布局，子组件名用逗号分隔 */
  layout?: string
  /** 加载中 */
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showZebra: true,
  showBorder: true,
  showHeaderBackground: true,
  layout: 'refresh,size,columns,settings',
})

const columnChecks = defineModel<ColumnOption[]>('columns', {
  default: () => [],
})

const emit = defineEmits<{
  refresh: []
}>()

const tableSizeOptions = [
  { value: TableSizeEnum.SMALL, label: '紧凑' },
  { value: TableSizeEnum.DEFAULT, label: '默认' },
  { value: TableSizeEnum.LARGE, label: '宽松' },
]

const tableSizeDropDownItems = computed(() =>
  tableSizeOptions.map(item => ({
    label: item.label,
    value: item.value,
    type: 'checkbox' as const,
    checked: tableSize.value === item.value,
    onUpdateChecked: (checked: boolean) => {
      if (checked) handleTableSizeChange(item.value)
    },
  })),
)

const tableStore = useTableStore()
const { tableSize, isZebra, isBorder, isHeaderBackground } = storeToRefs(tableStore)

const layoutItems = computed(() => props.layout.split(',').map(item => item.trim()))

const shouldShow = (componentName: string) => layoutItems.value.includes(componentName)

const refresh = () => {
  isManualRefresh.value = true
  emit('refresh')
}

const handleTableSizeChange = (command: TableSizeEnum) => {
  tableStore.setTableSize(command)
}

const isManualRefresh = ref(false)
</script>
