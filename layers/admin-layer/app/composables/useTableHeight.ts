/**
 * useTableHeight — 表格容器高度自动计算
 * 纯计算逻辑，SSR safe
 */
import { computed, type Ref } from 'vue'

interface TableHeightOptions {
  showTableHeader: Ref<boolean>
  showPagination?: Ref<boolean>
  paginationHeight: Ref<number>
  tableHeaderHeight: Ref<number>
  paginationSpacing: Ref<number>
  /** Auto-detected via #m-search-form; used by consumer for page-level layout */
  searchFormHeight?: Ref<number>
}

const DEFAULT_HEADER_HEIGHT = 48
const DEFAULT_PAGINATION_HEIGHT = 44

export function useTableHeight(options: TableHeightOptions) {
  const containerHeight = computed(() => {
    const { showTableHeader, showPagination, paginationHeight, tableHeaderHeight, paginationSpacing } = options

    let offset = 0
    if (showTableHeader.value) {
      offset += tableHeaderHeight.value || DEFAULT_HEADER_HEIGHT
    }
    // 预留分页高度：即便尚未测量到（首帧/懒渲染）或测量偏小，也按默认值预留，
    // 避免表格算高偏大把分页条挤出容器被裁掉
    if (showPagination?.value ?? paginationHeight.value > 0) {
      offset += Math.max(paginationHeight.value, DEFAULT_PAGINATION_HEIGHT) + paginationSpacing.value
    }

    if (offset === 0) {
      return { height: '100%', minHeight: '0' }
    }
    return { height: `calc(100% - ${offset}px)`, minHeight: '100px' }
  })

  return { containerHeight }
}
