<template>
  <UDashboardGroup storage="cookie" storage-key="admin-layout">
    <!-- 侧边栏 -->
    <UDashboardSidebar
      class="dark bg-slate-900"
      collapsible
      resizable
      :default-size="15"
      :min-size="10"
      :max-size="25"
      :collapsed-size="4"
    >
      <template #header="{ collapsed }">
        <AdminSidebarHeader :collapsed="collapsed" />
      </template>

      <template #default="{ collapsed }">
        <AdminSidebarMenu :collapsed="collapsed" />
      </template>

      <template #footer="{ collapsed }">
        <AdminUserDropdown
          :collapsed="collapsed"
          :user-info="userInfo"
          @logout="handleLogout"
        />
      </template>
    </UDashboardSidebar>

    <!-- 主内容区 -->
    <UDashboardPanel :ui="{ body: 'bg-muted' }">
      <template #header>
        <UDashboardNavbar :ui="{ root: 'border-b border-default' }">
          <template #leading>
            <AdminHeaderCollapseButton />
            <AdminHeaderRefreshButton />
            <AdminProjectSelect :token="token" />
            <AdminHeaderBreadcrumb />
          </template>

          <template #right>
            <div class="flex items-center gap-2">
              <UDashboardSearchButton class="text-muted min-w-48" />
              <MThemePanel v-if="showThemePanel" class="text-muted" />
              <MFullscreenToggle />
              <AdminColorToggle class="text-muted" />
            </div>
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <slot v-if="alive" />
      </template>
    </UDashboardPanel>

    <!-- 命令面板（搜索弹窗） -->
    <UDashboardSearch :groups="searchGroups" />
  </UDashboardGroup>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Admin layout — 所有管理后台页面使用 definePageMeta({ layout: 'admin' })
// header 的认证数据（token / userInfo / logout）由上层通过 props 传入，
// admin-layer 不再依赖 auth-module 的注入。
const props = defineProps<{
  /** Nuxt 内部传入的 layout 名称 */
  name?: string
  /** 访问令牌，用于项目切换下拉框请求项目列表；为空时不请求 */
  token?: string | null
  /** 当前登录用户信息 */
  userInfo?: { fullname?: string; name?: string } | null
  /** 登出回调 */
  logout?: () => void | Promise<void>
}>()

const { flatMenu } = useAdminMenu()

function handleLogout() {
  props.logout?.()
}

// 软刷新：控制页面内容组件卸载/重挂
const { alive } = usePageRefresh()

const appConfig = useAppConfig()
// 主题面板默认隐藏，仅当 admin.themePanel 显式为 true 时在顶栏显示
const showThemePanel = computed(() => (appConfig.admin as Record<string, unknown> | undefined)?.themePanel === true)

const searchGroups = computed(() => [
  {
    id: 'pages',
    label: 'Pages',
    items: flatMenu.value
      .filter(item => !item.external && item.to)
      .map(item => ({
        id: item.to,
        label: item.label,
        icon: item.icon,
        to: item.to,
      })),
  },
])
</script>
