/**
 * 样例页面全局守卫：拦截所有 /samples/* 路由。
 * 消费者需在 nuxt.config.ts 的 runtimeConfig.public.admin.exposeSamples 设为 true 才能启用。
 */
export default defineNuxtRouteMiddleware((to) => {
  if (!to.path.startsWith('/samples/') && to.path !== '/samples') return

  const config = useRuntimeConfig()
  const enabled = (config.public.admin as Record<string, unknown>)?.exposeSamples as boolean ?? false

  if (!enabled) {
    if (import.meta.server) {
      throw createError({ statusCode: 404, message: 'Page not found' })
    }
    return navigateTo('/', { replace: true })
  }
})
