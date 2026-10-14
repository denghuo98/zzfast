import { fileURLToPath } from 'node:url'
import { resolve } from 'node:path'

const currentDir = fileURLToPath(new URL('.', import.meta.url))

// 默认值由 admin-layer 统一管理，消费者可通过 app.config.ts 或 runtimeConfig 覆盖
export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@pinia/nuxt'],
  css: [resolve(currentDir, './app/assets/css/main.css')],

  runtimeConfig: {
    public: {
      admin: {
        title: 'Admin Panel',
        requireRole: 'admin',
        /** 是否暴露样例页面，默认 false。消费者在 nuxt.config.ts 中覆盖为 true 才能访问 /samples/* */
        exposeSamples: false,
      },
    },
  },
})
