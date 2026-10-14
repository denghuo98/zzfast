export default defineAppConfig({
  admin: {
    title: 'ZZAdmin',
    logo: 'i-lucide-shield',
    menu: [
      {
        label: '系统管理',
        icon: 'i-lucide-settings',
        children: [
          { label: '用户管理', to: '/system/user', icon: 'i-lucide-users' },
          { label: '角色管理', to: '/system/role', icon: 'i-lucide-shield' },
          { label: '菜单管理', to: '/system/menu', icon: 'i-lucide-menu' },
          { label: '登录日志', to: '/system/login-log', icon: 'i-lucide-scroll-text' },
          { label: '个人中心', to: '/system/profile', icon: 'i-lucide-user-round' },
        ],
      },
    ],
    homeRoute: '/system/user',
  },
  ui: {
    colors: {
      primary: 'violet',
      neutral: 'slate',
    },
  },
})
