import { useAuthStore } from '../stores/auth'

/**
 * Admin 页面中间件 —— 校验登录态。
 *
 * 未登录时重定向到 /login 并携带 redirect 回跳地址。
 * 服务端跳过（HTTP 层认证由后端 AdminAuth 中间件保证）。
 */
export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return

  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
})
