import type { AdminConfig } from './types/admin'

export default defineAppConfig({
  admin: {
    title: 'Admin Panel',
    logo: 'i-lucide-shield',
    sidebar: {
      collapsed: false,
      width: 260,
      collapsedWidth: 64,
    },
    menu: [],
    // 示例菜单（由 useAdminMenu 在 runtimeConfig.public.admin.exposeSamples=true 时附加）
    samplesMenu: [
      {
        label: '仪表盘',
        icon: 'i-lucide-layout-dashboard',
        children: [
          { label: '工作台', to: '/samples/dashboard/workbench' },
          { label: '分析页', to: '/samples/dashboard/analytics' },
        ],
      },
      {
        label: '模板中心',
        icon: 'i-lucide-panels-top-left',
        children: [
          { label: '卡片', to: '/samples/template/cards' },
          { label: '横幅', to: '/samples/template/banner' },
          { label: '图表', to: '/samples/template/charts' },
          { label: '日历', to: '/samples/template/calendar' },
        ],
      },
      {
        label: '功能示例',
        icon: 'i-lucide-blocks',
        children: [
          { label: '前端权限', to: '/samples/examples/permissions' },
          { label: '基础表格', to: '/samples/examples/basic-table' },
          { label: '高级表格', to: '/samples/examples/advanced-table' },
          { label: '表单', to: '/samples/examples/form' },
          { label: '搜索表单', to: '/samples/examples/search-form' },
        ],
      },
      {
        label: '部件中心',
        icon: 'i-lucide-layout-grid',
        children: [
          { label: '图标浏览器', to: '/samples/widget/icon' },
          { label: '图片裁剪', to: '/samples/widget/image-crop' },
          { label: 'Excel', to: '/samples/widget/excel' },
          { label: '水印', to: '/samples/widget/watermark' },
          { label: '数字滚动', to: '/samples/widget/count-to' },
          { label: '富文本', to: '/samples/widget/rich-text' },
          { label: '二维码', to: '/samples/widget/qrcode' },
          { label: '右键菜单', to: '/samples/widget/context-menu' },
          { label: '拖拽排序', to: '/samples/widget/drag' },
          { label: '文字滚动', to: '/samples/widget/text-scroll' },
          { label: '烟花', to: '/samples/widget/fireworks' },
          { label: '视频播放', to: '/samples/widget/video' },
        ],
      },
    ],
    // 主题设置面板默认隐藏，消费者配置 themePanel: true 才在顶栏显示
    themePanel: false,
    homeRoute: '/samples/dashboard/workbench',
  } satisfies AdminConfig,
  ui: {
    colors: {
      neutral: 'slate',
    },
    radius: 0.5,
    tooltip: {
      slots: {
        content: 'flex items-start gap-1 bg-black/70 text-white shadow-sm rounded-sm max-w-xs h-auto ring-0 px-2.5 py-1 text-xs select-none data-[state=delayed-open]:animate-[scale-in_100ms_var(--ease-out)] data-[state=closed]:animate-[scale-out_100ms_var(--ease-out)] origin-(--reka-tooltip-content-transform-origin) pointer-events-auto',
        text: 'text-clip break-words',
        arrow: 'fill-black/70',
      },
    },
    drawer: {
      slots: {
        // 抽屉遮罩黑透（对齐 tooltip 黑透风格）；内容面板保持默认
        overlay: 'bg-black/70',
      },
    },
  },
})
