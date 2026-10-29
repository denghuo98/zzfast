import { fileURLToPath } from 'node:url'
import { resolve } from 'node:path'

const currentDir = fileURLToPath(new URL('.', import.meta.url))

// 默认值由 admin-layer 统一管理，消费者可通过 app.config.ts 或 runtimeConfig 覆盖
export default defineNuxtConfig({
  modules: ['@nuxt/ui', '@pinia/nuxt'],
  css: [resolve(currentDir, './app/assets/css/main.css')],

  // 运行时图标接口默认挂在 /api/_nuxt_icon，生产环境 /api 会反代到 Go 后端导致图标 404，
  // 换到 /_nuxt_icon，交给前端自己处理。
  icon: {
    localApiEndpoint: '/_nuxt_icon',
  },

  vite: {
    resolve: {
      dedupe: ['vue', 'vue-router', '@vueuse/core'],
    },
    optimizeDeps: {
      // TipTap 依赖链会从多个入口重复预构建 prosemirror-state，
      // 导致同一个 PluginKey 被注册两次（Adding different instances of a keyed plugin）。
      // 固定从 @nuxt/ui 侧预构建，保证整个应用只有一份 prosemirror-state。
      include: ['@nuxt/ui > prosemirror-state'],
    },
    ssr: {
      // pinia 3.x 的 exports 在 node+production 条件下把入口指到 CJS（pinia.prod.cjs），
      // Rollup 会把其中的 require('vue') 转成 `import vue__default from 'vue'`，
      // 而 vue 只有命名导出 → SSR 运行时 "does not provide an export named 'default'" 直接崩。
      // 把这些包一起内联进 SSR bundle（noExternal），避免外部 ESM 默认导入问题。
      noExternal: ['vue', 'vue-router', 'pinia', '@pinia/nuxt', '@vueuse/core', 'vue-demi'],
    },
  },

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
