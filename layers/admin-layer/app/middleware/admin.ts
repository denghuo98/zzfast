/**
 * Admin 页面中间件 —— 默认实现（无 auth 依赖）。
 *
 * 当 auth-module 存在时，auth-module 会注册同名的 `admin` 中间件覆盖此 fallback，
 * 注入真实的认证/鉴权能力；单独使用 admin-layer 时此 fallback 会放行所有请求。
 */
export default createAdminMiddleware({
  isAuthenticated: () => true,
  hasRole: () => true,
})
