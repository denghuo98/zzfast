/**
 * 管理后台菜单项配置
 */
export interface AdminMenuItem {
  /** 显示文本 */
  label: string
  /** 图标名称（支持 @nuxt/icon 格式，如 'i-lucide-home'） */
  icon?: string
  /** 路由路径 */
  to?: string
  /** 限制可见性：需要用户具备的角色列表（任一匹配即显示） */
  roles?: string[]
  /** 限制可见性：需要用户具备的权限列表 */
  permissions?: string[]
  /** 子菜单项 */
  children?: AdminMenuItem[]
  /** 徽标（数字或文本） */
  badge?: string | number
  /** 是否为外部链接 */
  external?: boolean
  /** 链接打开方式（external 为 true 时生效） */
  target?: string
  /** 在此项前添加分隔线 */
  divider?: boolean
  /** 强制隐藏（忽略角色/权限检查） */
  hidden?: boolean
}

/**
 * 管理后台主题配置
 */
export interface AdminTheme {
  /** 主色调（Tailwind 颜色名，如 'blue', 'indigo', 'violet'） */
  primary: string
  /** 圆角比例 */
  radius: number
  /** 灰色调（如 'slate', 'neutral', 'zinc'） */
  gray: string
}

/**
 * 管理后台全局配置
 */
export interface AdminConfig {
  /** 应用标题 */
  title: string
  /** Logo：图片路径或图标名称 */
  logo: string
  /** 侧边栏配置 */
  sidebar: {
    /** 默认是否折叠 */
    collapsed: boolean
    /** 展开宽度（px） */
    width: number
    /** 折叠宽度（px） */
    collapsedWidth: number
  }
  /** 菜单项列表 */
  menu: AdminMenuItem[]
  /** 主题配置 */
  theme: AdminTheme
  /** 是否显示主题设置面板（默认隐藏，配置为 true 才在顶栏显示） */
  themePanel?: boolean
  /** 首页路由（访问 /admin 时重定向） */
  homeRoute: string
}
