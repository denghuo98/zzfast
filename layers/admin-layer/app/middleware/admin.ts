import { useAuthStore } from '../stores/auth'

/**
 * Admin 页面中间件 —— 校验登录态。
 *
 * 仅判断 cookie 是否存在并不够：后端停止或 token 过期时 cookie 仍在，
 * 因此这里会向后端请求用户信息确认会话真实有效（带超时与短期缓存）。
 * 校验失败时清空 token / 用户信息，并重定向到 /login 携带 redirect 回跳地址。
 *
 * 服务端跳过（HTTP 层认证由后端 AdminAuth 中间件保证）。
 */
export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) return

  const authStore = useAuthStore()
  const valid = await authStore.validateSession()

  if (!valid) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
})
