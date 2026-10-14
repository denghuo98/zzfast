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
        /** 登录密码 AES-GCM 加密密钥，必须与后端 system.requestEncryptKey 一致（16/24/32 字节） */
        requestEncryptKey: '1234567890abcdef',
      },
    },
  },
})
