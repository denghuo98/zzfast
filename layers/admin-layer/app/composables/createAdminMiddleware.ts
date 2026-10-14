export interface AdminMiddlewareOptions {
  /** 判断是否已登录 */
  isAuthenticated: () => boolean
  /** 判断角色 */
  hasRole: (role: string) => boolean
  /** 登录页路径，默认 /login */
  loginPath?: string
  /** 构造跳转 URL（未登录时）。默认使用 siteUrl 或 window.location.origin */
  buildRedirectUrl?: (to: ReturnType<typeof useRoute>) => string
}

/**
 * 创建 admin 页面中间件。
 *
 * 由 auth-module 等消费者注入真实认证能力，admin-layer 自身不依赖任何 auth 实现。
 */
export function createAdminMiddleware(options: AdminMiddlewareOptions) {
  const {
    isAuthenticated,
    hasRole,
    loginPath = '/login',
    buildRedirectUrl,
  } = options

  return defineNuxtRouteMiddleware((to) => {
    // 角色检查依赖客户端状态，服务端跳过（HTTP 层认证由 server plugin 保证）
    if (import.meta.server) return

    if (!isAuthenticated()) {
      const siteUrl = useRuntimeConfig().public.siteUrl as string || ''
      const origin = siteUrl || (typeof window !== 'undefined' ? window.location.origin : '')
      const next = encodeURIComponent(buildRedirectUrl?.(to) || origin + to.fullPath)
      return navigateTo(`${loginPath}?next=${next}`, { external: true })
    }

    const routeRoles = (to.meta.roles as string[] | undefined)
    const requireRole = useRuntimeConfig().public.admin?.requireRole as string | undefined
    const requiredRoles = routeRoles ?? (requireRole ? [requireRole] : ['admin'])

    const hasAccess = requiredRoles.some(role => hasRole(role))

    if (!hasAccess) {
      throw createError({
        statusCode: 403,
        statusMessage: 'Forbidden',
        message: 'You do not have permission to access this page.',
      })
    }
  })
}
