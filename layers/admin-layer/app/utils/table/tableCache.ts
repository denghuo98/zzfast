/**
 * 表格缓存管理 — 基于 ohash 的 LRU 缓存
 */
import { hash } from 'ohash'
import type { PaginatedResponse } from '../../types/api'

export enum CacheInvalidationStrategy {
  CLEAR_ALL = 'clear_all',
  CLEAR_CURRENT = 'clear_current',
  CLEAR_PAGINATION = 'clear_pagination',
  KEEP_ALL = 'keep_all',
}

export interface CacheItem<T> {
  data: T[]
  response: PaginatedResponse<T>
  timestamp: number
  params: string
  tags: Set<string>
  accessCount: number
  lastAccessTime: number
}

export class TableCache<T> {
  private cache = new Map<string, CacheItem<T>>()
  private cacheTime: number
  private maxSize: number
  private enableLog: boolean

  constructor(cacheTime = 5 * 60 * 1000, maxSize = 50, enableLog = false) {
    this.cacheTime = cacheTime
    this.maxSize = maxSize
    this.enableLog = enableLog
  }

  private log(message: string, ...args: unknown[]) {
    if (this.enableLog) console.log(`[TableCache] ${message}`, ...args)
  }

  private generateKey(params: unknown): string {
    return hash(params)
  }

  private generateTags(params: Record<string, unknown>): Set<string> {
    const tags = new Set<string>()
    const searchKeys = Object.keys(params).filter(
      key => !['current', 'size', 'total'].includes(key)
        && params[key] !== undefined && params[key] !== '' && params[key] !== null,
    )
    if (searchKeys.length > 0) {
      tags.add(`search:${searchKeys.map(k => `${k}:${String(params[k])}`).join('|')}`)
    } else {
      tags.add('search:default')
    }
    tags.add(`pagination:${params.size || 10}`)
    tags.add('pagination')
    return tags
  }

  private evictLRU(): void {
    if (this.cache.size <= this.maxSize) return
    let lruKey = ''
    let minAccess = Infinity
    let oldest = Infinity
    for (const [key, item] of this.cache.entries()) {
      if (item.accessCount < minAccess || (item.accessCount === minAccess && item.lastAccessTime < oldest)) {
        lruKey = key
        minAccess = item.accessCount
        oldest = item.lastAccessTime
      }
    }
    if (lruKey) {
      this.cache.delete(lruKey)
      this.log(`LRU evicted: ${lruKey}`)
    }
  }

  set(params: unknown, data: T[], response: PaginatedResponse<T>): void {
    const key = this.generateKey(params)
    const now = Date.now()
    this.evictLRU()
    this.cache.set(key, {
      data, response, timestamp: now, params: key,
      tags: this.generateTags(params as Record<string, unknown>),
      accessCount: 1, lastAccessTime: now,
    })
  }

  get(params: unknown): CacheItem<T> | null {
    const key = this.generateKey(params)
    const item = this.cache.get(key)
    if (!item) return null
    if (Date.now() - item.timestamp > this.cacheTime) {
      this.cache.delete(key)
      return null
    }
    item.accessCount++
    item.lastAccessTime = Date.now()
    return item
  }

  clearCurrentSearch(params: unknown): number {
    const key = this.generateKey(params)
    return this.cache.delete(key) ? 1 : 0
  }

  clearPagination(): number {
    return this.clearByTags(['pagination'])
  }

  clearByTags(tags: string[]): number {
    let count = 0
    for (const [key, item] of this.cache.entries()) {
      if (tags.some(t => [...item.tags].some(it => it.includes(t)))) {
        this.cache.delete(key)
        count++
      }
    }
    return count
  }

  clear(): void { this.cache.clear() }

  getStats() {
    let totalSize = 0; let totalAccess = 0
    for (const item of this.cache.values()) {
      totalSize += JSON.stringify(item.data).length
      totalAccess += item.accessCount
    }
    return {
      total: this.cache.size,
      size: `${(totalSize / 1024).toFixed(2)}KB`,
      hitRate: `${this.cache.size > 0 ? (totalAccess / this.cache.size).toFixed(1) : '0'} avg hits`,
    }
  }

  cleanupExpired(): number {
    let count = 0
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now - item.timestamp > this.cacheTime) { this.cache.delete(key); count++ }
    }
    return count
  }
}
