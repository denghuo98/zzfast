/**
 * useTable — 表格数据管理
 *
 * 功能：数据请求 / 分页 / 搜索防抖 / 缓存 / 刷新策略
 */
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, readonly, watch, type Ref } from 'vue'
import { useWindowSize } from '@vueuse/core'
import { useTableColumns } from './useTableColumns'
import type { ColumnOption, TablePagination } from '../types/table'
import type { PaginatedResponse } from '../types/api'
import { TableCache, CacheInvalidationStrategy } from '../utils/table/tableCache'
import {
  type TableError,
  defaultResponseAdapter,
  extractTableData,
  updatePaginationFromResponse,
  createSmartDebounce,
  createErrorHandler,
} from '../utils/table/tableUtils'
import { tableConfig } from '../utils/table/tableConfig'

// 类型推导
type InferApiParams<T> = T extends (params: infer P) => any ? P : never
type InferApiResponse<T> = T extends (params: any) => Promise<infer R> ? R : never
type InferRecordType<T> =
  T extends PaginatedResponse<infer U> ? U
    : T extends { data: (infer U)[]; total?: number } ? U
      : never

export interface UseTableConfig<
  TApiFn extends (params: any) => Promise<any> = (params: any) => Promise<any>,
  TRecord = InferRecordType<InferApiResponse<TApiFn>>,
  TParams = InferApiParams<TApiFn>,
  TResponse = InferApiResponse<TApiFn>,
> {
  core: {
    apiFn: TApiFn
    apiParams?: Partial<TParams>
    excludeParams?: string[]
    immediate?: boolean
    columnsFactory?: () => ColumnOption<TRecord>[]
    paginationKey?: { current?: string; size?: string }
  }
  transform?: {
    dataTransformer?: (data: TRecord[]) => TRecord[]
    responseAdapter?: (response: TResponse) => PaginatedResponse<TRecord>
  }
  performance?: {
    enableCache?: boolean
    cacheTime?: number
    debounceTime?: number
    maxCacheSize?: number
  }
  hooks?: {
    onSuccess?: (data: TRecord[], response: PaginatedResponse<TRecord>) => void
    onError?: (error: any) => void
    onCacheHit?: (data: TRecord[], response: PaginatedResponse<TRecord>) => void
    resetFormCallback?: () => void
  }
  /** 表格交互控制 */
  table?: {
    /** MTable 的模板引用，通过它直接访问 TanStack tableApi */
    tableRef?: Ref<any>
    /** 启用列固定功能 */
    enableColumnPinning?: boolean
  }
  debug?: { enableLog?: boolean }
}

export function useTable<TApiFn extends (params: any) => Promise<any>>(config: UseTableConfig<TApiFn>) {
  type TRecord = InferRecordType<InferApiResponse<TApiFn>>
  type TParams = InferApiParams<TApiFn>

  const {
    core: { apiFn, apiParams = {} as Partial<TParams>, excludeParams = [], immediate = true, columnsFactory, paginationKey },
    transform: { dataTransformer, responseAdapter = defaultResponseAdapter } = {},
    performance: { enableCache = false, cacheTime = 5 * 60 * 1000, debounceTime = 300, maxCacheSize = 50 } = {},
    hooks: { onSuccess, onError, onCacheHit, resetFormCallback } = {},
    table: { tableRef, enableColumnPinning = false } = {},
    debug: { enableLog = false } = {},
  } = config

  const pageKey = paginationKey?.current || tableConfig.paginationKey.current
  const sizeKey = paginationKey?.size || tableConfig.paginationKey.size

  const cacheUpdateTrigger = ref(0)

  const logger = {
    log: (msg: string, ...args: unknown[]) => { if (enableLog) console.log(`[useTable] ${msg}`, ...args) },
    warn: (msg: string, ...args: unknown[]) => { if (enableLog) console.warn(`[useTable] ${msg}`, ...args) },
    error: (msg: string, ...args: unknown[]) => { if (enableLog) console.error(`[useTable] ${msg}`, ...args) },
  }

  const cache = enableCache ? new TableCache<TRecord>(cacheTime, maxCacheSize, enableLog) : null

  type LoadingState = 'idle' | 'loading' | 'success' | 'error'
  const loadingState = ref<LoadingState>('idle')
  const loading = computed(() => loadingState.value === 'loading')
  const error = ref<TableError | null>(null)
  const data = ref<TRecord[]>([])
  let abortController: AbortController | null = null
  let cacheCleanupTimer: ReturnType<typeof setInterval> | null = null

  const searchParams = reactive(Object.assign({ [pageKey]: 1, [sizeKey]: 10 }, apiParams || {}) as TParams)

  const pagination = reactive<TablePagination>({
    page: ((searchParams as Record<string, unknown>)[pageKey] as number) || 1,
    pageSize: ((searchParams as Record<string, unknown>)[sizeKey] as number) || 10,
    total: 0,
  })

  const { width } = useWindowSize()
  const mobilePagination = computed(() => ({ ...pagination, small: width.value < 768 }))

  const columnConfig = columnsFactory ? useTableColumns<TRecord>(columnsFactory) : null

  // ==================== 列固定 ====================

  const columnPinning = ref<Record<string, string[]>>({})

  function pinColumn(columnId: string, position: 'left' | 'right') {
    const current = { ...columnPinning.value }
    const left = (current.left || []).filter(c => c !== columnId)
    const right = (current.right || []).filter(c => c !== columnId)
    if (position === 'left') left.push(columnId)
    else right.push(columnId)
    columnPinning.value = { left, right }
  }

  function unpinColumn(columnId: string) {
    columnPinning.value = {
      left: (columnPinning.value.left || []).filter(c => c !== columnId),
      right: (columnPinning.value.right || []).filter(c => c !== columnId),
    }
  }

  function clearPinning() {
    columnPinning.value = {}
  }

  // 自动同步 columnPinning 到 TanStack tableApi
  if (enableColumnPinning) {
    watch(
      () => tableRef?.value?.uTableRef?.tableApi as any,
      (api) => {
        if (!api) return
        api.setColumnPinning(columnPinning.value)
        watch(columnPinning, (v) => api.setColumnPinning(v), { deep: true })
      },
      { immediate: true },
    )
  }

  // ==================== 数据请求 ====================

  const hasData = computed(() => data.value.length > 0)

  const cacheInfo = computed(() => {
    void cacheUpdateTrigger.value
    return cache ? cache.getStats() : { total: 0, size: '0KB', hitRate: '0 avg hits' }
  })

  const handleError = createErrorHandler(onError, enableLog)

  const clearCache = (strategy: CacheInvalidationStrategy, context?: string): void => {
    if (!cache) return
    let count = 0
    switch (strategy) {
      case CacheInvalidationStrategy.CLEAR_ALL: cache.clear(); break
      case CacheInvalidationStrategy.CLEAR_CURRENT: count = cache.clearCurrentSearch(searchParams); break
      case CacheInvalidationStrategy.CLEAR_PAGINATION: count = cache.clearPagination(); break
    }
    if (count) logger.log(`清理缓存 ${count} 条 — ${context || ''}`)
    cacheUpdateTrigger.value++
  }

  const fetchData = async (params?: Partial<TParams>, useCache = enableCache): Promise<PaginatedResponse<TRecord>> => {
    if (abortController) abortController.abort()
    const currentController = new AbortController()
    abortController = currentController
    loadingState.value = 'loading'
    error.value = null

    try {
      let requestParams = Object.assign({}, searchParams, { [pageKey]: pagination.page, [sizeKey]: pagination.pageSize }, params || {}) as TParams
      if (excludeParams.length > 0) {
        const filtered = { ...requestParams as Record<string, unknown> }
        excludeParams.forEach(k => delete filtered[k])
        requestParams = filtered as TParams
      }

      if (useCache && cache) {
        const cached = cache.get(requestParams)
        if (cached) {
          data.value = cached.data
          updatePaginationFromResponse(pagination, cached.response)
          const paramsRecord = searchParams as Record<string, unknown>
          if (paramsRecord[pageKey] !== pagination.page) paramsRecord[pageKey] = pagination.page
          if (paramsRecord[sizeKey] !== pagination.pageSize) paramsRecord[sizeKey] = pagination.pageSize
          loadingState.value = 'success'
          onCacheHit?.(cached.data, cached.response)
          return cached.response
        }
      }

      const response = await apiFn(requestParams)
      if (currentController.signal.aborted) throw new Error('请求已取消')

      const standard = responseAdapter(response)
      let tableData = extractTableData(standard)
      if (dataTransformer) tableData = dataTransformer(tableData)

      data.value = tableData
      updatePaginationFromResponse(pagination, standard)

      const paramsRecord = searchParams as Record<string, unknown>
      if (paramsRecord[pageKey] !== pagination.page) paramsRecord[pageKey] = pagination.page
      if (paramsRecord[sizeKey] !== pagination.pageSize) paramsRecord[sizeKey] = pagination.pageSize

      if (useCache && cache) { cache.set(requestParams, tableData, standard); cacheUpdateTrigger.value++ }

      loadingState.value = 'success'
      onSuccess?.(tableData, standard)
      return standard
    } catch (err) {
      if (err instanceof Error && err.message === '请求已取消') { loadingState.value = 'idle'; return { list: [], total: 0, page: 1, itemsPerPage: 10 } }
      loadingState.value = 'error'
      data.value = []
      throw handleError(err, '获取表格数据失败')
    } finally {
      if (abortController === currentController) abortController = null
    }
  }

  const getData = async (params?: Partial<TParams>) => { try { return await fetchData(params) } catch { /* handled */ } }

  const getDataByPage = async (params?: Partial<TParams>) => {
    pagination.page = 1; (searchParams as Record<string, unknown>)[pageKey] = 1
    clearCache(CacheInvalidationStrategy.CLEAR_CURRENT, '搜索')
    try { return await fetchData(params, false) } catch { /* handled */ }
  }

  const debouncedGetDataByPage = createSmartDebounce(getDataByPage, debounceTime)

  const resetSearchParams = async () => {
    debouncedGetDataByPage.cancel()
    const paramsRecord = searchParams as Record<string, unknown>
    const defaultPagination = { [pageKey]: 1, [sizeKey]: (paramsRecord[sizeKey] as number) || 10 }
    Object.keys(searchParams).forEach(k => delete paramsRecord[k])
    Object.assign(searchParams, apiParams || {}, defaultPagination)
    pagination.page = 1
    pagination.pageSize = defaultPagination[sizeKey] as number
    error.value = null
    clearCache(CacheInvalidationStrategy.CLEAR_ALL, '重置')
    await getData()
    if (resetFormCallback) { await nextTick(); resetFormCallback() }
  }

  let isCurrentChanging = false

  const handleSizeChange = async (s: number) => {
    if (s <= 0) return; debouncedGetDataByPage.cancel()
    const paramsRecord = searchParams as Record<string, unknown>
    pagination.pageSize = s; pagination.page = 1; paramsRecord[sizeKey] = s; paramsRecord[pageKey] = 1
    clearCache(CacheInvalidationStrategy.CLEAR_CURRENT, '页大小变化')
    await getData()
  }

  const handleCurrentChange = async (p: number) => {
    if (p <= 0 || isCurrentChanging || pagination.page === p) return
    try {
      isCurrentChanging = true
      const paramsRecord = searchParams as Record<string, unknown>
      pagination.page = p
      if (paramsRecord[pageKey] !== p) paramsRecord[pageKey] = p
      await getData()
    } finally { isCurrentChanging = false }
  }

  const refreshData = async () => { debouncedGetDataByPage.cancel(); clearCache(CacheInvalidationStrategy.CLEAR_ALL, '手动刷新'); await getData() }
  const refreshCreate = async () => { debouncedGetDataByPage.cancel(); pagination.page = 1; (searchParams as Record<string, unknown>)[pageKey] = 1; clearCache(CacheInvalidationStrategy.CLEAR_PAGINATION, '新增'); await getData() }
  const refreshUpdate = async () => { clearCache(CacheInvalidationStrategy.CLEAR_CURRENT, '编辑'); await getData() }
  const refreshRemove = async () => {
    const { page } = pagination; clearCache(CacheInvalidationStrategy.CLEAR_CURRENT, '删除'); await getData()
    if (data.value.length === 0 && page > 1) { pagination.page = page - 1; (searchParams as Record<string, unknown>)[pageKey] = page - 1; await getData() }
  }
  const refreshSoft = async () => { clearCache(CacheInvalidationStrategy.CLEAR_CURRENT, '软刷新'); await getData() }

  const cancelRequest = () => { abortController?.abort(); debouncedGetDataByPage.cancel() }
  const clearData = () => { data.value = []; error.value = null; clearCache(CacheInvalidationStrategy.CLEAR_ALL, '清空') }
  const clearExpiredCache = () => cache ? cache.cleanupExpired() : 0

  if (enableCache && cache) {
    cacheCleanupTimer = setInterval(() => { const c = cache.cleanupExpired(); if (c) { logger.log(`自动清理 ${c} 条过期缓存`); cacheUpdateTrigger.value++ } }, cacheTime / 2)
  }

  if (immediate) onMounted(async () => { await getData() })
  onUnmounted(() => { cancelRequest(); cache?.clear(); if (cacheCleanupTimer) clearInterval(cacheCleanupTimer) })

  return {
    data, loading: readonly(loading), error: readonly(error),
    isEmpty: computed(() => data.value.length === 0), hasData,
    pagination, paginationMobile: mobilePagination,
    handleSizeChange, handleCurrentChange,
    searchParams, resetSearchParams,
    fetchData: getData, getData: getDataByPage, getDataDebounced: debouncedGetDataByPage, clearData,
    refreshData, refreshSoft, refreshCreate, refreshUpdate, refreshRemove,
    cacheInfo, clearCache, clearExpiredCache, cancelRequest,
    columnPinning, pinColumn, unpinColumn, clearPinning,
    ...(columnConfig && {
      columns: columnConfig.columns, columnChecks: columnConfig.columnChecks,
      addColumn: columnConfig.addColumn, removeColumn: columnConfig.removeColumn,
      toggleColumn: columnConfig.toggleColumn, updateColumn: columnConfig.updateColumn,
      batchUpdateColumns: columnConfig.batchUpdateColumns, reorderColumns: columnConfig.reorderColumns,
      getColumnConfig: columnConfig.getColumnConfig, getAllColumns: columnConfig.getAllColumns,
      resetColumns: columnConfig.resetColumns,
    }),
  }
}
