/**
 * useApi —— 管理后台请求封装。
 *
 * 约定后端 zzfast 通过 ResponseHandler 中间件统一返回：
 *   { code, message?, data?, error?, timestamp, traceID }
 * code === 0 表示成功，此时返回 data；否则抛出错误。
 *
 * token 通过 Authorization: Bearer <token> 传递（后端 ztoken.GetAuthorization 会去掉 Bearer 前缀）。
 */
import { useCookie } from '#app'
import type { UseFetchOptions } from 'nuxt/app'

export const ADMIN_TOKEN_COOKIE = 'admin_token'

export interface ApiResponse<T = unknown> {
  code: number
  message?: string
  data?: T
  error?: unknown
  timestamp?: number
  traceID?: string
}

export function useApi() {
  const token = useCookie<string>(ADMIN_TOKEN_COOKIE)

  async function request<T>(url: string, options: UseFetchOptions<T> = {}): Promise<T> {
    const headers: Record<string, string> = { ...(options.headers as Record<string, string>) }
    if (token.value) {
      headers.Authorization = `Bearer ${token.value}`
    }

    const res = await $fetch<ApiResponse<T>>(url, {
      ...options,
      headers,
    })

    if (res.code !== 0) {
      const message = res.message || (typeof res.error === 'string' ? res.error : '请求失败')
      throw new Error(message)
    }

    return res.data as T
  }

  return {
    request,
    get: <T>(url: string, params?: Record<string, unknown>) =>
      request<T>(url, { method: 'GET', params }),
    post: <T>(url: string, body?: unknown) =>
      request<T>(url, { method: 'POST', body }),
    del: <T>(url: string, params?: Record<string, unknown>) =>
      request<T>(url, { method: 'DELETE', params }),
  }
}
