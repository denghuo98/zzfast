export default defineAppConfig({
  admin: {
    title: 'Minimal Admin',
    logo: 'i-lucide-layout-dashboard',
    menu: [
      {
        label: '首页',
        icon: 'i-lucide-home',
        to: '/',
      },
    ],
    homeRoute: '/',
  },
})
