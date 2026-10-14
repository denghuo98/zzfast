import type { ComputedRef } from 'vue'
import { useAuthStore } from '../stores/auth'
import type { AdminUserInfo } from '../types/admin-api'

/**
 * admin-layer 认证能力抽象 —— 从内置 auth store 读取真实登录态。
 *
 * 权限控制由后端 Casbin（AdminAuth 中间件）保证，前端 hasRole/hasPermission 仅作
 * 菜单显隐的辅助判断；当前未获取角色列表，因此默认放行（true）。
 */
export interface AdminAuthAdapter {
  /** 当前是否已登录 */
  isAuthenticated: () => boolean
  /** 判断用户是否拥有指定角色 */
  hasRole: (role: string) => boolean
  /** 判断用户是否拥有指定权限 */
  hasPermission: (perm: string) => boolean
  /** 当前登录用户信息 */
  userInfo: ComputedRef<AdminUserInfo | null>
  /** 登出 */
  logout: () => void | Promise<void>
  /** 当前访问令牌 */
  token: () => string | undefined
}

export function useAdminAuth(): AdminAuthAdapter {
  const authStore = useAuthStore()

  const userInfo = computed(() => authStore.userInfo)

  return {
    isAuthenticated: () => authStore.isAuthenticated,
    hasRole: () => true,
    hasPermission: () => true,
    userInfo,
    logout: () => authStore.logout(),
    token: () => authStore.token || undefined,
  }
}
