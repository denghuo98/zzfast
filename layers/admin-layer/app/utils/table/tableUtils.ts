/**
 * 表格工具函数 — 响应适配、防抖、错误处理
 */
import { tableConfig } from './tableConfig'
import type { TablePagination } from '../../types/table'
import type { PaginatedResponse } from '../../types/api'

export interface TableError {
  code: string
  message: string
  details?: unknown
}

// --- 响应适配 ---

function extractRecords<T>(obj: Record<string, unknown>, fields: string[]): T[] {
  for (const field of fields) {
    if (field in obj && Array.isArray(obj[field])) return obj[field] as T[]
  }
  return []
}

function extractTotal(obj: Record<string, unknown>, records: unknown[], fields: string[]): number {
  for (const field of fields) {
    if (field in obj && typeof obj[field] === 'number') return obj[field] as number
  }
  return records.length
}

function extractPagination(obj: Record<string, unknown>, data?: Record<string, unknown>) {
  const result: Partial<{ page: number; itemsPerPage: number }> = {}
  for (const src of [obj, data ?? {}]) {
    for (const field of tableConfig.currentFields) {
      if (field in src && typeof src[field] === 'number') { result.page = src[field] as number; break }
    }
    if (result.page !== undefined) break
  }
  for (const src of [obj, data ?? {}]) {
    for (const field of tableConfig.sizeFields) {
      if (field in src && typeof src[field] === 'number') { result.itemsPerPage = src[field] as number; break }
    }
    if (result.itemsPerPage !== undefined) break
  }
  return Object.keys(result).length > 0 ? result : undefined
}

export const defaultResponseAdapter = <T>(response: unknown): PaginatedResponse<T> => {
  if (!response) return { list: [], total: 0, page: 1, itemsPerPage: 10 }
  if (Array.isArray(response)) return { list: response, total: response.length, page: 1, itemsPerPage: response.length }
  if (typeof response !== 'object') return { list: [], total: 0, page: 1, itemsPerPage: 10 }

  const res = response as Record<string, unknown>
  let records = extractRecords<T>(res, tableConfig.recordFields)
  let total = extractTotal(res, records, tableConfig.totalFields)
  const pagination = extractPagination(res)

  if (records.length === 0 && 'data' in res && typeof res.data === 'object' && !Array.isArray(res.data)) {
    const data = res.data as Record<string, unknown>
    records = extractRecords<T>(data, ['list', 'records', 'items'])
    total = extractTotal(data, records, tableConfig.totalFields)
  }

  const result: PaginatedResponse<T> = { list: records, total, page: 1, itemsPerPage: 10 }
  if (pagination) Object.assign(result, pagination)
  return result
}

export const extractTableData = <T>(response: PaginatedResponse<T>): T[] =>
  Array.isArray(response.list) ? response.list : []

export const updatePaginationFromResponse = <T>(
  pagination: TablePagination,
  response: PaginatedResponse<T>,
): void => {
  pagination.page = response.page ?? pagination.page ?? 1
  pagination.pageSize = response.itemsPerPage ?? pagination.pageSize ?? 10
  pagination.total = response.total ?? pagination.total ?? 0
}

// --- 防抖 ---

export const createSmartDebounce = <T extends (...args: any[]) => Promise<any>>(fn: T, delay: number) => {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  let lastArgs: Parameters<T> | null = null
  let lastResolve: ((v: any) => void) | null = null
  let lastReject: ((r: any) => void) | null = null

  const debouncedFn = (...args: Parameters<T>): Promise<any> =>
    new Promise((resolve, reject) => {
      if (timeoutId) clearTimeout(timeoutId)
      lastArgs = args; lastResolve = resolve; lastReject = reject
      timeoutId = setTimeout(async () => {
        try { resolve(await fn(...args)) }
        catch (e) { reject(e) }
        finally { timeoutId = null; lastArgs = null; lastResolve = null; lastReject = null }
      }, delay)
    })

  ;(debouncedFn as any).cancel = () => { if (timeoutId) { clearTimeout(timeoutId); timeoutId = null } }
  ;(debouncedFn as any).flush = async () => {
    if (timeoutId && lastArgs && lastResolve && lastReject) {
      clearTimeout(timeoutId); timeoutId = null
      try { lastResolve(await fn(...lastArgs!)) }
      catch (e) { lastReject(e) }
    }
  }
  return debouncedFn as T & { cancel: () => void; flush: () => Promise<any> }
}

// --- 错误处理 ---

export const createErrorHandler = (onError?: (e: TableError) => void, _enableLog = false) =>
  (err: unknown, context: string): TableError => {
    const tableError: TableError = { code: 'UNKNOWN_ERROR', message: '未知错误', details: err }
    if (err instanceof Error) { tableError.message = err.message; tableError.code = err.name }
    else if (typeof err === 'string') tableError.message = err
    onError?.(tableError)
    return tableError
  }
