/**
 * useAdminPermission — 权限演示用 composable
 *
 * 封装 adminAuth 的角色信息访问，供 permissions.vue 使用。
 * 从 useAdminAuth().userInfo 提取 roles（auth-module 的 UserInfo.groups[].code）。
 */
import { computed } from 'vue'
import type { ComputedRef } from 'vue'
import { useAdminAuth } from './useAdminAuth'

export function useAdminPermission() {
  const auth = useAdminAuth()

  /** 当前用户所有角色 code（如 ['admin@matrix', 'user@matrix']） */
  const currentRoles: ComputedRef<string[]> = computed(() => {
    const info = auth.userInfo.value as Record<string, unknown> | null
    if (!info) return []
    const groups = info.groups as { code: string }[] | undefined
    return groups?.map(g => g.code) ?? []
  })

  /** 首个角色，去掉 @ 后缀（如 'admin'） */
  const currentRole: ComputedRef<string> = computed(() => {
    const first = currentRoles.value[0]
    return first?.split('@')[0] ?? ''
  })

  /** 前 10 个角色（展示用），超出截断 */
  const displayRoles: ComputedRef<string[]> = computed(() => currentRoles.value.slice(0, 10))

  /** 被截断的角色数量（> 0 时说明有省略） */
  const remainingRoleCount: ComputedRef<number> = computed(() => Math.max(0, currentRoles.value.length - 10))

  /** 登录名 */
  const userName: ComputedRef<string> = computed(() => {
    const info = auth.userInfo.value as Record<string, unknown> | null
    return (info?.fullname ?? info?.name ?? '') as string
  })

  const hasRole = auth.hasRole

  /** 触发登出 → 跳 SSO 换账号登录 → 回跳本页 */
  function switchRole(_accountLabel: string) {
    // 1. 清除认证 cookie（同 store.clear）
    if (typeof document !== 'undefined') {
      document.cookie = 'ACCESS_TOKEN=; path=/; max-age=0'
      document.cookie = 'AUTH_USER=; path=/; max-age=0'
    }
    // 2. 跳 /login?next= 回本页，SSO 登录后 redirect 回来角色即更新
    if (typeof window !== 'undefined') {
      window.location.href = `/login?next=${encodeURIComponent(window.location.href)}`
    }
  }

  return { currentRoles, currentRole, displayRoles, remainingRoleCount, userName, hasRole, switchRole }
}
