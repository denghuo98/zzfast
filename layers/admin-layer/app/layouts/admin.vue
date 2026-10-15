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
          :user-info="adminAuth.userInfo.value"
          @logout="adminAuth.logout()"
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
import { useAuthStore } from '../stores/auth'

// Admin layout — 所有管理后台页面使用 definePageMeta({ layout: 'admin' })
// 认证数据（token / userInfo / logout）从内置 auth store 读取。
const adminAuth = useAdminAuth()
const authStore = useAuthStore()

// 刷新页面后 userInfo 会丢失（仅 token 持久化在 cookie），客户端重新拉取
onMounted(async () => {
  if (authStore.isAuthenticated && !authStore.userInfo) {
    await authStore.fetchUserInfo().catch(() => {})
  }
})

const { flatMenu } = useAdminMenu()

// 软刷新：控制页面内容组件卸载/重挂
const { alive } = usePageRefresh()

const appConfig = useAppConfig()
// 主题面板默认隐藏，仅当 admin.themePanel 显式为 true 时在顶栏显示
const showThemePanel = computed(() => (appConfig.admin as Record<string, unknown> | undefined)?.themePanel === true)

// 浏览器标题使用品牌名（app.config 的 admin.title，可覆盖）
const brandTitle = computed(() => (appConfig.admin as Record<string, unknown> | undefined)?.title as string || 'Admin Panel')
useHead({ title: brandTitle })

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
