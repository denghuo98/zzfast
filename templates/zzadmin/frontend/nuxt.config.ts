// zzadmin 前端案例：extends 本地 admin-layer，配置静态菜单与后端代理。
// 发布后业务方改为 extends ['@zzfast/admin-layer'] 并 npm 安装即可。
export default defineNuxtConfig({
  extends: ['../../../layers/admin-layer'],

  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  devServer: {
    host: '0.0.0.0',
    port: 3000,
  },

  runtimeConfig: {
    public: {
      admin: {
        title: 'ZZAdmin',
        // 与后端 templates/zzadmin/backend/config.yaml 的 system.requestEncryptKey 保持一致
        requestEncryptKey: '1234567890abcdef',
      },
    },
  },

  // 开发环境将 /admin 代理到后端（生产由 nginx 反代）
  nitro: {
    devProxy: {
      '/admin': {
        target: 'http://127.0.0.1:8000/admin',
        changeOrigin: true,
      },
    },
  },
})
