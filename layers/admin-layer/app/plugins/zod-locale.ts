/**
 * Zod 中文错误提示 — Zod 4 内置 locale，无需额外依赖
 */
import { zhCN } from 'zod/locales'
import * as z from 'zod'

export default defineNuxtPlugin(() => {
  z.config({ localeError: zhCN().localeError })
})
