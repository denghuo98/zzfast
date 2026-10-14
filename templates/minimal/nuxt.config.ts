// 最小消费工程：本地 extends admin-layer（未发版前验证编译/运行）
// 发版后业务方改为 extends ['@zzfast/admin-layer'] 并 npm 安装即可。
export default defineNuxtConfig({
  extends: ['../../layers/admin-layer'],

  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  devServer: {
    host: '0.0.0.0',
    port: 8030,
  },

  runtimeConfig: {
    public: {
      admin: {
        // 最小工程不暴露样例页面
        exposeSamples: false,
      },
    },
  },
})
