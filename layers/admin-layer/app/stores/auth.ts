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

export const useAuthStore = defineStore('admin-auth', () => {
  const token = useCookie<string>(ADMIN_TOKEN_COOKIE, { maxAge: 86400, path: '/' })
  const userInfo = ref<AdminUserInfo | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  function setToken(value: string) {
    token.value = value
  }

  function setUserInfo(info: AdminUserInfo | null) {
    userInfo.value = info
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
  async function fetchUserInfo() {
    if (!token.value) return null
    const api = useApi()
    const info = await api.get<AdminUserInfo>('/admin/member/info')
    setUserInfo(info)
    return info
  }

  /** 退出登录 */
  async function logout() {
    const api = useApi()
    try {
      await api.post('/admin/site/account/logout')
    } catch {
      // 忽略登出接口错误，本地状态照常清理
    }
    token.value = null
    setUserInfo(null)
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
    logout,
  }
})
