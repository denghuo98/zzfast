/**
 * useAuthStore —— 管理后台认证状态。
 *
 * token 存储在 cookie（admin_token）中，SSR 与客户端均可读取；
 * 请求封装 useApi 会读取该 cookie 并附带 Authorization 头。
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useCookie } from '#app'
import { ADMIN_TOKEN_COOKIE, useApi } from '../composables/useApi'
import { aesGcmEncrypt } from '../utils/crypto'
import type { AdminLoginResponse, AdminUserInfo } from '../types/admin-api'

/** 会话校验结果的缓存时长，避免每次路由跳转都请求后端 */
const SESSION_VALIDATE_TTL = 30_000
/** 会话校验请求超时时间，后端不可用时尽快跳转登录页 */
const SESSION_VALIDATE_TIMEOUT = 4_000

export const useAuthStore = defineStore('admin-auth', () => {
  const token = useCookie<string>(ADMIN_TOKEN_COOKIE, { maxAge: 86400, path: '/' })
  const userInfo = ref<AdminUserInfo | null>(null)
  const validatedAt = ref(0)

  const isAuthenticated = computed(() => !!token.value)

  function setToken(value: string) {
    token.value = value
    validatedAt.value = 0
  }

  function setUserInfo(info: AdminUserInfo | null) {
    userInfo.value = info
  }

  function clearSession() {
    token.value = null
    setUserInfo(null)
    validatedAt.value = 0
  }

  /** 获取登录验证码 */
  async function fetchCaptcha() {
    const api = useApi()
    return api.get<{ id: string; base64: string }>('/admin/site/captcha')
  }

  /** 登录 */
  async function login(params: { username: string; password: string; cid?: string; code?: string }) {
    const api = useApi()
    const runtimeConfig = useRuntimeConfig()
    const encryptKey = (runtimeConfig.public.admin as Record<string, unknown>)?.requestEncryptKey as string || '1234567890abcdef'

    const encryptedPassword = await aesGcmEncrypt(params.password, encryptKey)

    const res = await api.post<AdminLoginResponse>('/admin/site/account/login', {
      username: params.username,
      password: encryptedPassword,
      cid: params.cid || '',
      code: params.code || '',
    })

    setToken(res.token)
    setUserInfo({ id: res.id, username: res.username, realName: '', avatar: '', sex: 0, email: '', mobile: '' })

    await fetchUserInfo().catch(() => {})
    return res
  }

  /** 获取当前登录用户信息 */
  async function fetchUserInfo(timeout?: number) {
    if (!token.value) return null
    const api = useApi()
    const info = await api.get<AdminUserInfo>('/admin/member/info', undefined, timeout)
    setUserInfo(info)
    validatedAt.value = Date.now()
    return info
  }

  /**
   * 校验当前会话是否有效（token 存在且后端认可）。
   *
   * 校验成功会刷新用户信息；失败会清空 token 和用户信息。
   * 校验结果默认缓存 SESSION_VALIDATE_TTL，可用 force 强制刷新。
   */
  async function validateSession(force = false): Promise<boolean> {
    if (!token.value) {
      clearSession()
      return false
    }

    if (!force && userInfo.value && Date.now() - validatedAt.value < SESSION_VALIDATE_TTL) {
      return true
    }

    try {
      await fetchUserInfo(SESSION_VALIDATE_TIMEOUT)
      return true
    } catch {
      clearSession()
      return false
    }
  }

  /** 退出登录 */
  async function logout() {
    const api = useApi()
    try {
      await api.post('/admin/site/account/logout')
    } catch {
      // 忽略登出接口错误，本地状态照常清理
    }
    clearSession()
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    setToken,
    setUserInfo,
    fetchCaptcha,
    login,
    fetchUserInfo,
    validateSession,
    logout,
  }
})
