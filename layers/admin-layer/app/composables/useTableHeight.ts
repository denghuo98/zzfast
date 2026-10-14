/**
 * useTableHeight — 表格容器高度自动计算
 * 纯计算逻辑，SSR safe
 */
import { computed, type Ref } from 'vue'

interface TableHeightOptions {
  showTableHeader: Ref<boolean>
  paginationHeight: Ref<number>
  tableHeaderHeight: Ref<number>
  paginationSpacing: Ref<number>
  /** Auto-detected via #m-search-form; used by consumer for page-level layout */
  searchFormHeight?: Ref<number>
}

const DEFAULT_HEADER_HEIGHT = 44

export function useTableHeight(options: TableHeightOptions) {
  const containerHeight = computed(() => {
    const { showTableHeader, paginationHeight, tableHeaderHeight, paginationSpacing } = options

    let offset = 0
    if (showTableHeader.value) {
      offset += tableHeaderHeight.value || DEFAULT_HEADER_HEIGHT
    }
    if (paginationHeight.value > 0) {
      offset += paginationHeight.value + paginationSpacing.value
    }

    if (offset === 0) {
      return { height: '100%', minHeight: '0' }
    }
    return { height: `calc(100% - ${offset}px)`, minHeight: '100px' }
  })

  return { containerHeight }
}
