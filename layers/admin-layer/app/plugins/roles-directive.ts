/**
 * v-roles 指令（client plugin）
 *
 * 移植自 art-design-pro directives/core/roles.ts。
 * 用法：<UButton v-roles="'admin'"> 或 <UButton v-roles="['admin','super']">
 * 无权限时从 DOM 移除元素（removeChild），OR 语义。
 *
 * 注意：指令 hook 不在 Vue setup 上下文，不能直接 useAdminAuth()（内部调 inject）。
 * 因此在 plugin setup 阶段从 Nuxt app 获取 auth adapter 并闭包传入 directive。
 */

import type { Directive, DirectiveBinding } from 'vue'
import { AdminAuthInjectionKey, type AdminAuthAdapter } from '../composables/useAdminAuth'

/** fallback adapter（auth-module 未注入时全放行） */
const fallback: AdminAuthAdapter = {
  isAuthenticated: () => true,
  hasRole: () => true,
  hasPermission: () => true,
  userInfo: { value: null } as any,
  logout: () => {},
  token: () => undefined,
}

function createRolesDirective(auth: AdminAuthAdapter): Directive<HTMLElement, string | string[]> {
  function check(el: HTMLElement, binding: DirectiveBinding<string | string[]>): void {
    // 跳过 SSR（directive 仅 client 生效）
    if (typeof window === 'undefined') return
    const required = Array.isArray(binding.value) ? binding.value : [binding.value]
    const hasMatch = required.some(r => auth.hasRole(r))
    if (!hasMatch && el.parentNode) el.parentNode.removeChild(el)
  }

  return { mounted: check, updated: check }
}

export default defineNuxtPlugin((nuxtApp) => {
  // 从 Vue app 级 provides 读取 auth adapter（同 useAdminAuth 逻辑，但 directive 内无法调 inject）
  const injected = nuxtApp.vueApp._context.provides[AdminAuthInjectionKey as any] as AdminAuthAdapter | undefined
  const auth = injected ?? fallback
  nuxtApp.vueApp.directive('roles', createRolesDirective(auth))
})
