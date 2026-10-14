/**
 * useAdminMenu — 管理后台菜单状态与操作
 *
 * 基于 app.config.ts 中的 admin.menu 配置，
 * 提供角色过滤、活跃路径解析、面包屑构建等功能。
 */
import type { AdminMenuItem } from '../types/admin'

export const useAdminMenu = () => {
  const appConfig = useAppConfig()
  const runtimeConfig = useRuntimeConfig()
  const route = useRoute()

  // 原始菜单配置（业务菜单 + 在 exposeSamples=true 时附加 layer 提供的示例菜单）
  const rawMenu = computed<AdminMenuItem[]>(() => {
    const admin = appConfig.admin as Record<string, unknown> | undefined
    const business = (admin?.menu as AdminMenuItem[]) ?? []
    const exposeSamples = (runtimeConfig.public.admin as Record<string, unknown> | undefined)?.exposeSamples
    if (!exposeSamples) return business
    const samples = (admin?.samplesMenu as AdminMenuItem[]) ?? []
    return [...business, ...samples]
  })

  // 角色过滤后的菜单项
  const menuItems = computed(() => filterByRoles(rawMenu.value))

  // 当前活跃的菜单路径
  const activePath = computed(() => {
    return findActivePath(rawMenu.value, route.path)
  })

  // 当前页面的面包屑路径
  const breadcrumbs = computed(() => {
    const metaCrumb = (route.meta?.breadcrumb as string | undefined) || ''
    return buildBreadcrumbs(rawMenu.value, route.path, [], metaCrumb)
  })

  // 扁平化菜单（用于搜索/CommandPalette）
  const flatMenu = computed(() => {
    return flattenMenu(menuItems.value)
  })

  // 检查菜单项是否活跃
  function isMenuItemActive(item: AdminMenuItem): boolean {
    if (!item.to) return false
    if (item.external) return false
    return route.path === item.to || route.path.startsWith(item.to + '/')
  }

  return {
    menuItems,
    activePath,
    breadcrumbs,
    flatMenu,
    isMenuItemActive,
  }
}

// === 内部工具函数 ===

function filterByRoles(items: AdminMenuItem[]): AdminMenuItem[] {
  // admin-layer 不依赖 auth-module；通过注入的 auth adapter 做角色过滤
  const { hasRole, hasPermission } = useAdminAuth()
  return doFilter(items, hasRole, hasPermission)
}

function doFilter(
  items: AdminMenuItem[],
  hasRole: (role: string) => boolean,
  hasPermission: (perm: string) => boolean,
): AdminMenuItem[] {
  return items
    .filter(item => !item.hidden)
    .filter((item) => {
      if (item.roles?.length && !item.roles.some(r => hasRole(r))) return false
      if (item.permissions?.length && !item.permissions.some(p => hasPermission(p))) return false
      return true
    })
    .map((item) => {
      if (item.children && item.children.length > 0) {
        const filtered = doFilter(item.children, hasRole, hasPermission)
        return filtered.length > 0 ? { ...item, children: filtered } : null
      }
      return item
    })
    .filter(Boolean) as AdminMenuItem[]
}

function findActivePath(items: AdminMenuItem[], currentPath: string): string | null {
  for (const item of items) {
    if (item.to && !item.external && (currentPath === item.to || currentPath.startsWith(item.to + '/'))) {
      return item.to
    }
    if (item.children) {
      const found = findActivePath(item.children, currentPath)
      if (found) return found
    }
  }
  return null
}

function buildBreadcrumbs(
  items: AdminMenuItem[],
  currentPath: string,
  ancestors: { label: string; to?: string; icon?: string }[] = [],
  metaCrumb = '',
): { label: string; to?: string; icon?: string }[] {
  // 1. 精确匹配：当前菜单项 to === currentPath → 直接返回完整 trail
  for (const item of items) {
    const crumb = { label: item.label, to: item.to, icon: item.icon }
    const trail = [...ancestors, crumb]

    if (item.to === currentPath) return trail
    if (item.children) {
      const found = buildBreadcrumbs(item.children, currentPath, trail, metaCrumb)
      if (found.length > 0) return found
    }
  }

  // 2. 前缀匹配（当前层级）：currentPath 以某个 item.to 开头 → 继承父链并追加末段
  //    例如 /procurement/requisitions/123 → [采购管理 / 采购需求单 / 需求单详情]
  for (const item of items) {
    if (item.to && item.to !== '/' && currentPath.startsWith(item.to + '/')) {
      const crumb = { label: item.label, to: item.to, icon: item.icon }
      const trail = [...ancestors, crumb]
      const tailLabel = metaCrumb || prettifySegment(currentPath.slice(item.to.length + 1))
      if (tailLabel) trail.push({ label: tailLabel })
      return trail
    }
  }

  // 3. 完全找不到：从 URL segments 兜底（英文化）；优先使用 metaCrumb 作为末段
  if (ancestors.length === 0) {
    const segments = currentPath.split('/').filter(Boolean)
    const crumbs: { label: string; to?: string }[] = []
    let accumulated = ''
    for (let i = 0; i < segments.length; i++) {
      accumulated += '/' + segments[i]
      const isLast = i === segments.length - 1
      crumbs.push({
        label: isLast && metaCrumb ? metaCrumb : prettifySegment(segments[i]!),
        to: accumulated,
      })
    }
    return crumbs
  }

  return []
}

function prettifySegment(seg: string): string {
  if (!seg) return ''
  // 把第一段（可能含 / 之后的尾段）转成"首字母大写 + 短横替空格"
  const first = seg.split('/')[0]!
  return first.charAt(0).toUpperCase() + first.slice(1).replace(/-/g, ' ')
}

function flattenMenu(items: AdminMenuItem[]): AdminMenuItem[] {
  return items.flatMap((item) => {
    const children = item.children ? flattenMenu(item.children) : []
    return [item, ...children]
  })
}
