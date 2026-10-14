<template>
  <ClientOnly>
    <!-- 无数据时的加载空态 -->
    <div v-if="loading && isEmpty" class="flex items-center justify-center gap-2" :style="containerHeight">
      <UIcon name="i-lucide-loader-2" class="w-8 h-8 animate-spin text-muted" />
      <span class="text-sm text-muted">加载中...</span>
    </div>

    <!-- 已有数据时保留表格，刷新时叠半透明白遮罩 + loading，避免表格卸载/高度塌陷闪屏 -->
    <div class="relative">
      <UTable
        v-if="!loading || !isEmpty"
        ref="uTableRef"
        sticky="header"
        :data="props.data"
        :columns="tableColumns"
        :style="containerHeight"
        :ui="tableUi"
        :get-sub-rows="resolvedGetSubRows"
        v-model:row-selection="rowSelection"
        v-model:column-pinning="columnPinning"
        v-model:expanded="expandedModel"
      >
      <template #__selection__-header="{ table }">
        <UCheckbox
          :model-value="
            table.getIsSomeRowsSelected()
              ? 'indeterminate'
              : table.getIsAllRowsSelected()
          "
          @update:model-value="(v: any) => table.toggleAllRowsSelected(!!v)"
        />
      </template>

      <template #__selection__-cell="{ row }">
        <UCheckbox
          :model-value="row.getIsSelected()"
          @update:model-value="(v: any) => row.toggleSelected(!!v)"
        />
      </template>

      <template #__index__-header>序号</template>
      <template #__index__-cell="{ row }">#{{ row.original.id }}</template>

      <!-- 自动展开列：在 expand-key 指定的列上渲染缩进 + chevron 按钮 -->
      <template v-if="props.expandKey" #[`${props.expandKey}-cell`]="{ row }">
        <div class="flex items-center gap-1" :style="{ paddingLeft: `${row.depth * props.expandIndent}px` }">
          <UButton
            v-if="row.getCanExpand()"
            :icon="row.getIsExpanded() ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
            color="neutral"
            variant="ghost"
            size="xs"
            class="shrink-0"
            tabindex="-1"
            @click="row.toggleExpanded()"
          />
          <span v-else class="inline-block w-5 shrink-0" />
          <slot :name="`${props.expandKey}-cell-content`" :row="row">
            {{ row.original[props.expandKey] }}
          </slot>
        </div>
      </template>

      <!-- 树形展开时 UTable 默认会在 expanded 行下渲染空白行，用空模板抑制 -->
      <template #expanded />

      <!-- 固定宽度列：超长内容省略号 + UTooltip 悬浮显示全文 -->
      <template v-for="col in ellipsisColumns" :key="col.id ?? (col as any).accessorKey" #[`${col.accessorKey}-cell`]="{ getValue }">
        <UTooltip :text="String(getValue?.() ?? '')" arrow>
          <span class="block truncate">{{ getValue?.() !== null && getValue?.() !== undefined ? getValue?.() : ' ' }}</span>
        </UTooltip>
      </template>

      <template v-for="(_, name) in $slots" #[name]="slotData">
        <slot :name="name" v-bind="slotData" />
      </template>
      </UTable>

      <!-- 刷新遮罩：仅当已有数据且正在刷新时显示；表格本体保留在遮罩下层 -->
      <div v-if="loading && !isEmpty" class="absolute inset-0 z-10 flex items-center justify-center bg-white/60 dark:bg-black/40" aria-hidden="true">
        <UIcon name="i-lucide-loader-2" class="w-6 h-6 animate-spin text-primary" />
      </div>
    </div>

    <!-- 分页 -->
    <div
      v-if="showPagination && (!loading || !isEmpty)"
      ref="paginationRef"
      class="flex gap-2 mt-4"
      :class="ALIGN_CLASS_MAP[mergedPaginationOptions?.align as keyof typeof ALIGN_CLASS_MAP]"
    >
      <UPagination
        :page="pagination?.page"
        :items-per-page="pagination?.pageSize"
        :total="pagination?.total"
        @update:page="handleCurrentChange"
      />

      <USelect
        :model-value="currentPageSize"
        :items="pageSizeOptions"
        :content="{ side: 'top' }"
        @update:model-value="handleSizeChange"
      />
    </div>

    <template #fallback>
      <div class="flex items-center justify-center p-8">
        <span class="text-sm text-muted">加载中...</span>
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, useTemplateRef, watchEffect, onMounted, onUnmounted, useSlots } from 'vue'
import { useResizeObserver, useWindowSize } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import type { TableColumn } from '@nuxt/ui'
import { useTableStore } from '../stores/table'
import { useTableHeight } from '../composables/useTableHeight'
import { useCommon } from '../composables/useCommon'
import type { TablePagination, ColumnOption } from '../types/table'

defineOptions({ name: 'MTable' })

const { width } = useWindowSize()
const uTableRef = useTemplateRef('uTableRef')
const paginationRef = ref<HTMLElement>()
const tableHeaderRef = ref<HTMLElement>()
const tableStore = useTableStore()
const { isBorder, isZebra, tableSize, isHeaderBackground } = storeToRefs(tableStore)
const { scrollToTop: scrollPageToTop } = useCommon()

interface PaginationOptions {
  pageSizes?: number[]
  align?: 'left' | 'center' | 'right'
  hideOnSinglePage?: boolean
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
}

interface Props<T = Record<string, any>> {
  loading?: boolean
  columns?: TableColumn<T>[]
  pagination?: TablePagination
  paginationOptions?: PaginationOptions
  emptyHeight?: string
  /** 是否开启 MTableHeader 高度检测。Header 在卡片内部时设为 false */
  showTableHeader?: boolean
  data?: T[]
  /** 树形子行取值函数。默认读 row.children。传入后启用树形嵌套展开 */
  getSubRows?: (row: any) => any[] | undefined
  /** 自动渲染展开按钮的列 accessorKey。设置后该列自带缩进 + chevron */
  expandKey?: string
  /** expandKey 列每层级缩进像素，默认 24 */
  expandIndent?: number
}

const props = withDefaults(defineProps<Props<any>>(), {
  columns: () => [],
  emptyHeight: '100%',
  showTableHeader: true,
  expandIndent: 24,
})

const DEFAULT_PAGINATION: PaginationOptions = {
  pageSizes: [10, 20, 30, 50, 100],
  align: 'center',
  hideOnSinglePage: false,
  size: 'md',
}

const ALIGN_CLASS_MAP = { left: 'justify-start', center: 'justify-center', right: 'justify-end' }

const mergedPaginationOptions = computed(() => ({ ...DEFAULT_PAGINATION, ...props.paginationOptions }))

const pageSizeOptions = computed(() =>
  mergedPaginationOptions.value.pageSizes?.map(s => ({ label: `${s}条/页`, value: s })),
)

// 仅无数据的加载空态才去边框；已有数据刷新（遮罩态）保留边框避免布局跳动
const border = computed(() => (props.loading && isEmpty.value) ? false : isBorder.value)
const stripe = computed(() => isZebra.value)
const headerBackground = computed(() => isHeaderBackground.value)
const size = computed(() => tableSize.value)

const slots = useSlots()

/// MTable 内部已渲染插槽的列 key，ellipsis 自动启用时需跳过以避免同名插槽冲突
const reservedCellSlots = computed(() => {
  const keys = new Set(['__selection__', '__index__'])
  if (props.expandKey) keys.add(props.expandKey)
  return keys
})

/// 自动启用超长省略号(ellipsis)的列：有固定宽度 + 纯文本内容(无父插槽覆盖) + 非特殊列
/// 可通过 ColumnOption.ellipsis 显式控制：true 强制启用、false 强制关闭
const ellipsisColumns = computed(() =>
  (tableColumns.value ?? []).filter(col => {
    const ext = col as ColumnOption
    const key = ext.accessorKey as string
    if (!key) return false
    if (ext.width == null && ext.minWidth == null) return false
    // ellipsis 显式值优先；未设置时自动判断(无父插槽则启用)
    const use = ext.ellipsis ?? !slots[`${key}-cell`]
    if (!use) return false
    return !reservedCellSlots.value.has(key)
  }),
)

const tableUi = computed(() => {
  // m-table 无条件挂在 root 上，供组件内非 scoped CSS 定位固定列样式
  const rootClasses = ['m-table']
  const ui: Record<string, string> = { thead: '' }
  if (headerBackground.value) {
    ui.thead = '[&>tr]:bg-elevated/50 [&>tr]:after:content-none'
    rootClasses.push('m-table-header-bg')
  }
  if (border.value) { rootClasses.push('border', 'border-default', 'rounded-md'); ui.thead += ' rounded-t-md'; ui.tbody = 'rounded-b-md' }
  if (stripe.value) {
    ui.tbody = (ui.tbody || '') + ' [&>tr:nth-child(even)]:bg-elevated/50'
    rootClasses.push('m-table-stripe')
  }
  ui.root = rootClasses.join(' ')
  // th 始终不换行
  switch (size.value) {
    case 'sm': ui.th = 'whitespace-nowrap px-2 py-1.5 text-xs'; ui.td = 'px-2 py-1.5 text-xs'; break
    case 'lg': ui.th = 'whitespace-nowrap px-6 py-5 text-default text-base'; ui.td = 'px-6 py-5 text-base'; break
    default:
      // 基于 Nuxt UI v4 UTable 原始默认值修改：th 换行→nowrap、text-sm→text-base；
      ui.th = 'px-4 py-3.5 text-highlighted text-left rtl:text-right font-semibold whitespace-nowrap [&:has([role=checkbox])]:pe-0'
      ui.td = 'p-4 text-default whitespace-nowrap [&:has([role=checkbox])]:pe-0'
      break
  }
  return ui
})

/// 将 ColumnOption.width / minWidth 解析为 CSS 宽度值：数字转 px，字符串原样保留
const resolveWidth = (v: string | number | undefined): string | undefined => {
  if (v == null) return undefined
  return typeof v === 'number' ? `${v}px` : String(v)
}

/// 固定区/滚动区交界列（antd rc-table fixUtil 语义）：
/// leftLast = 显示序中最后一个左固定列且右邻未固定；rightFirst = 第一个右固定列且左邻未固定。
/// 交界列挂标记 class，供 ::after 阴影条定位
const pinBoundary = computed(() => {
  const pinning = columnPinning.value
  const leftSet = new Set<string>(pinning?.left ?? [])
  const rightSet = new Set<string>(pinning?.right ?? [])
  const cols = props.columns ?? []
  let leftLast: string | undefined
  let rightFirst: string | undefined
  for (let i = 0; i < cols.length; i++) {
    const key = (cols[i] as any).accessorKey as string
    if (leftSet.has(key) && !leftSet.has((cols[i + 1] as any)?.accessorKey)) leftLast = key
    if (rightFirst === undefined && rightSet.has(key) && !rightSet.has((cols[i - 1] as any)?.accessorKey)) rightFirst = key
  }
  return { leftLast, rightFirst }
})

/// 列宽映射：width → th/td 内联 width + max-width（table-layout: auto 下 max-width 是硬上限，
/// 能阻止列被内容撑开）；minWidth → 内联 min-width（列下限）。
/// 数值型 width 同时同步为 TanStack 的 size/minSize，供列固定偏移量等计算使用
const tableColumns = computed(() =>
  props.columns?.map(col => {
    const ext = col as ColumnOption
    const width = resolveWidth(ext.width)
    const minWidth = resolveWidth(ext.minWidth)
    const numWidth = typeof ext.width === 'number' ? ext.width : undefined
    const cellStyle = {
      ...(width ? { width, minWidth: width, maxWidth: width } : {}),
      ...(minWidth ? { minWidth } : {}),
    }
    const key = ext.accessorKey as string
    const boundaryClass = key === pinBoundary.value.leftLast
      ? 'm-cell-fix-left-last'
      : key === pinBoundary.value.rightFirst ? 'm-cell-fix-right-first' : undefined
    const cellClass = {
      th: [ext.meta?.class?.th, boundaryClass].filter(Boolean).join(' ') || undefined,
      td: [ext.meta?.class?.td, boundaryClass].filter(Boolean).join(' ') || undefined,
    }
    return {
      ...col,
      size: col.size ?? numWidth,
      minSize: col.minSize ?? (col.size != null ? Number(col.size) : undefined) ?? numWidth,
      meta: {
        ...ext.meta,
        class: cellClass,
        style: {
          ...ext.meta?.style,
          th: { ...ext.meta?.style?.th, ...cellStyle },
          td: { ...ext.meta?.style?.td, ...cellStyle },
        },
      },
    }
  }) ?? [],
)

const isEmpty = computed(() => !props.data?.length)

const paginationHeight = ref(0)
const tableHeaderHeight = ref(0)
const searchFormHeight = ref(0)

useResizeObserver(paginationRef, (entries) => {
  const entry = entries[0]
  if (entry) requestAnimationFrame(() => { paginationHeight.value = entry.contentRect.height })
})

useResizeObserver(tableHeaderRef, (entries) => {
  const entry = entries[0]
  if (entry) requestAnimationFrame(() => { tableHeaderHeight.value = entry.contentRect.height })
})

const searchFormRef = ref<HTMLElement>()
useResizeObserver(searchFormRef, (entries) => {
  const entry = entries[0]
  if (entry) requestAnimationFrame(() => { searchFormHeight.value = entry.contentRect.height })
})

const PAGINATION_SPACING = computed(() => 16)
const { containerHeight } = useTableHeight({
  showTableHeader: computed(() => props.showTableHeader),
  paginationHeight,
  tableHeaderHeight,
  searchFormHeight,
  paginationSpacing: PAGINATION_SPACING,
})

const showPagination = computed(() => props.pagination && !isEmpty.value)

// USelect 需要明确的 computed 绑定以确保响应式更新
const currentPageSize = computed(() => props.pagination?.pageSize ?? 10)

const emit = defineEmits<{
  (e: 'pagination:size-change', val: number): void
  (e: 'pagination:current-change', val: number): void
  (e: 'selection-change', selection: any[]): void
}>()

const columnPinning = defineModel<Record<string, string[]>>('columnPinning')

// 受控展开状态（树形子行 + 行详情共用）
const expandedModel = defineModel<Record<string, boolean>>('expanded', { default: () => ({}) })

// 树形展开时 UTable 会为每个 expanded 行渲染空的 #expanded 详情行，
// 用 JS 在 DOM 渲染后定位并隐藏这些空行（CSS `visibility:collapse` 对 <tr> 最可靠）。
function hideEmptyGapRows() {
  nextTick(() => {
    const el = (uTableRef.value as any)?.$el as HTMLElement | undefined
    if (!el) return
    el.querySelectorAll('tr[data-slot="tr"]:not([data-expanded]):not([data-selected])').forEach((tr) => {
      const td = tr.querySelector('td[colspan]')
      // 仅隐藏空内容行（无子节点或仅含空白/注释），避免误杀合法展开内容
      if (td && td.childNodes.length === 0 || (td && td.textContent?.trim() === '')) {
        ;(tr as HTMLElement).style.visibility = 'collapse'
      }
    })
  })
}

// 每次数据或展开状态变化后执行
watch([() => props.data, () => expandedModel.value], hideEmptyGapRows, { deep: false, flush: 'post' })
watch(() => props.loading, (v) => { if (!v) hideEmptyGapRows() })
onMounted(hideEmptyGapRows)

// 滚动状态追踪：为固定列边缘阴影检测横滚位置，并标记两轴是否溢出（用于按轴关闭弹性滚动）
let scrollCleanupFn: (() => void) | null = null

/// setAttribute 即便写同值也会触发 style invalidation，滚动每帧写属性会导致整表 reflow 卡顿，
/// 因此仅在值变化时写 DOM
function setAttrIfChanged(el: HTMLElement, name: string, value: string | null) {
  const prev = el.getAttribute(name)
  if (value === null) {
    if (prev !== null) el.removeAttribute(name)
    return
  }
  if (prev !== value) el.setAttribute(name, value)
}

/// 与 antd rc-table 相同的 ping 语义：该侧是否存在被滚动隐藏的内容
/// data-overflow-x / -y = 该轴当前是否真的可滚动（仅在可滚动的轴上关闭 overscroll，避免阻断页面滚动链）
function updateScrollState(root: HTMLElement) {
  const { scrollLeft, scrollWidth, clientWidth, scrollHeight, clientHeight } = root
  const overflowX = scrollWidth - clientWidth > 1
  setAttrIfChanged(root, 'data-overflow-x', String(overflowX))
  setAttrIfChanged(root, 'data-overflow-y', String(scrollHeight - clientHeight > 1))
  if (!overflowX) {
    setAttrIfChanged(root, 'data-ping-left', null)
    setAttrIfChanged(root, 'data-ping-right', null)
    return
  }
  setAttrIfChanged(root, 'data-ping-left', String(scrollLeft > 1))
  setAttrIfChanged(root, 'data-ping-right', String(scrollLeft + clientWidth < scrollWidth - 1))
}

function bindScrollShadow() {
  try {
    const root = (uTableRef.value as any)?.$el as HTMLElement | undefined
    if (!root) return
    const handler = () => updateScrollState(root)
    handler()
    root.addEventListener('scroll', handler, { passive: true })
    // 窗口缩放 / 数据变化导致溢出状态变化时复算
    const ro = new ResizeObserver(() => handler())
    ro.observe(root)
    scrollCleanupFn = () => {
      root.removeEventListener('scroll', handler)
      ro.disconnect()
      ;['data-ping-left', 'data-ping-right', 'data-overflow-x', 'data-overflow-y'].forEach(a => root.removeAttribute(a))
    }
  } catch {}
}

onMounted(() => {
  if (!props.loading) nextTick(bindScrollShadow)
})

watch(() => props.loading, (loading) => {
  if (loading) {
    scrollCleanupFn?.()
    scrollCleanupFn = null
  } else {
    nextTick(bindScrollShadow)
  }
})

onUnmounted(() => { scrollCleanupFn?.(); scrollCleanupFn = null })

// getSubRows 解析：未传时默认读 row.children，使树形开箱即用
const resolvedGetSubRows = computed(() => props.getSubRows ?? ((row: any) => row?.children as any[] | undefined))

const rowSelection = ref<Record<string, boolean>>({})

watch(() => props.data, () => { rowSelection.value = {} })

watch(rowSelection, (selection) => {
  const selectedRows = Object.entries(selection)
    .filter(([, s]) => s)
    .map(([id]) => {
      const index = parseInt(id)
      return props.data?.[index] ? { original: props.data[index] } : null
    }).filter(Boolean)
  emit('selection-change', selectedRows)
}, { deep: true })

const scrollToTop = () => {
  nextTick(() => {
    if (uTableRef.value?.$el) (uTableRef.value.$el as HTMLElement).scrollTop = 0
    scrollPageToTop()
  })
}

const handleSizeChange = (val: number) => emit('pagination:size-change', val)
const handleCurrentChange = (val: number) => { emit('pagination:current-change', val); scrollToTop() }

const findPageElements = () => {
  const el = document.getElementById('m-table-header')
  tableHeaderRef.value = el ?? undefined
  const sf = document.getElementById('m-search-form')
  searchFormRef.value = sf ?? undefined
}

watchEffect(() => {
  void props.data?.length
  nextTick(() => findPageElements())
}, { flush: 'post' })

defineExpose({ scrollToTop, uTableRef })
</script>

<style>
.m-table {
  /* 固定列不透明底色 */
  --m-table-bg: var(--ui-bg);
  /* 等价于 bg-elevated/50 叠在 --ui-bg 上的合成色 —— 让固定列与同行普通列深浅一致 */
  --m-table-bg-alt: color-mix(in oklab, var(--ui-bg-elevated) 50%, var(--ui-bg));
  --m-table-shadow: rgb(0 0 0 / 0.15);
}
.dark .m-table { --m-table-shadow: rgb(0 0 0 / 0.45); }

/* ---- 取消表格自身弹性滚动：仅作用于真正可滚动的轴，避免阻断页面滚动链 ---- */
.m-table[data-overflow-x="true"] { overscroll-behavior-x: none; }
.m-table[data-overflow-y="true"] { overscroll-behavior-y: none; }

/* ---- 性能：sticky thead 默认的半透明 + backdrop-blur 在滚动每帧重采样下方内容（repaint 大户，
       横滚/页面竖滚都卡）。表格落在卡片纯色底上，改不透明后静止视觉等价 ---- */
.m-table thead {
  background-color: var(--m-table-bg);
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
}

/* ---- 固定列背景：不透明（遮挡下方滚动内容），且跟随斑马纹 / 表头背景 / 选中态 ---- */
.m-table th[data-pinned="left"], .m-table th[data-pinned="right"],
.m-table td[data-pinned="left"], .m-table td[data-pinned="right"] {
  background-color: var(--m-table-bg);
}
.m-table-header-bg th[data-pinned="left"],
.m-table-header-bg th[data-pinned="right"] { background-color: var(--m-table-bg-alt); }
.m-table-stripe tbody tr:nth-child(even) td[data-pinned="left"],
.m-table-stripe tbody tr:nth-child(even) td[data-pinned="right"] { background-color: var(--m-table-bg-alt); }
.m-table tbody tr[data-selected="true"] td[data-pinned="left"],
.m-table tbody tr[data-selected="true"] td[data-pinned="right"] { background-color: var(--m-table-bg-alt); }

/* ---- 固定区/滚动区交界阴影（antd 设计：::after 阴影条贴在交界列外侧，
       border-collapse 下单元格自身 box-shadow 绘制不可靠，独立伪元素盒子稳定绘制） ---- */
.m-table .m-cell-fix-left-last::after,
.m-table .m-cell-fix-right-first::after {
  content: "";
  position: absolute;
  top: 0;
  bottom: -1px;
  width: 30px;
  transition: box-shadow 0.3s;
  pointer-events: none;
}
.m-table .m-cell-fix-left-last::after { left: 100%; }
.m-table .m-cell-fix-right-first::after { right: 100%; }
/* ping = 该侧有被滚动隐藏的内容时显示阴影 */
.m-table[data-ping-left="true"] .m-cell-fix-left-last::after {
  box-shadow: inset 10px 0 8px -8px var(--m-table-shadow);
}
.m-table[data-ping-right="true"] .m-cell-fix-right-first::after {
  box-shadow: inset -10px 0 8px -8px var(--m-table-shadow);
}
</style>
