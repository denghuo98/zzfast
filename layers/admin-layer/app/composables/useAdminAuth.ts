import type { InjectionKey } from 'vue'

/**
 * admin-layer 所需的认证能力抽象。
 *
 * auth-module 会通过 provide 注入真实实现；
 * 未注入时（admin-layer 单独使用）返回全放行 fallback，
 * 保证 UI 组件和中间件不会报错。
 *
 * 使用 Symbol.for 让 auth-module 与 admin-layer 在运行时可以共享同一个 key，
 * 而不需要在包级别互相依赖。
 */
export interface AdminAuthAdapter {
  /** 当前是否已登录 */
  isAuthenticated: import('vue').ComputedRef<boolean> | (() => boolean)
  /** 判断用户是否拥有指定角色 */
  hasRole: (role: string) => boolean
  /** 判断用户是否拥有指定权限 */
  hasPermission: (perm: string) => boolean
  /** 当前登录用户信息 */
  userInfo: import('vue').ComputedRef<{ name?: string; fullname?: string; [key: string]: unknown } | null>
  /** 登出 */
  logout: () => void | Promise<void>
  /** 当前访问令牌（用于 AdminProjectSelect 等组件） */
  token: () => string | undefined
}

export const AdminAuthInjectionKey: InjectionKey<AdminAuthAdapter> = Symbol.for('admin-auth')

export function useAdminAuth(): AdminAuthAdapter {
  const injected = inject(AdminAuthInjectionKey, null)
  if (injected) return injected

  // 未注入 auth 实现时的默认 fallback：全放行
  const userInfo = computed(() => null)
  return {
    isAuthenticated: () => true,
    hasRole: () => true,
    hasPermission: () => true,
    userInfo,
    logout: () => {},
    token: () => undefined,
  }
}
