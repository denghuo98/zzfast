// 图片上传接口：替代 nuxt-ui-templates/editor 的 hub:blob + vercel blob
import { randomUUID } from 'uncrypto'

// 允许的图片扩展名白名单（服务端校验，不信任用户输入）
const ALLOWED_EXTENSIONS = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp'])

export default eventHandler(async (event) => {
  const data = await readMultipartFormData(event)

  if (!data || data.length === 0) {
    throw createError({ statusCode: 400, statusMessage: 'No file', data: { message: '未选择文件' } })
  }

  const file = data.find(f => f.name === 'file') ?? data[0]!

  // 扩展名：从原始文件名提取，仅作白名单校验（不信任用户输入）
  const originalExt = file.filename?.split('.').pop()?.toLowerCase() || ''
  const ext = ALLOWED_EXTENSIONS.has(originalExt) ? originalExt : 'png'

  // MIME 校验
  if (!file.type?.startsWith('image/')) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid file type', data: { message: '仅支持图片文件' } })
  }

  // 大小校验 2MB
  if (file.data.length > 2 * 1024 * 1024) {
    throw createError({ statusCode: 400, statusMessage: 'File too large', data: { message: '文件大小不能超过 2MB' } })
  }

  // 文件名完全由服务端生成（UUID + 白名单扩展名），不包含用户输入，杜绝路径穿越
  const name = `${randomUUID().slice(0, 8)}.${ext}`

  const fs = await import('node:fs/promises')
  const path = await import('node:path')
  const dest = path.resolve('public/uploads', name)

  await fs.mkdir(path.resolve('public/uploads'), { recursive: true })
  await fs.writeFile(dest, file.data)

  // ⚠️ Nuxt layers 文件在 layer 的 public/ 下需特殊处理：
  // 实际放到根 public/uploads/ 保证 CDN / dev 都能访问
  // 返回相对路径
  return {
    url: `/uploads/${name}`,
    pathname: name
  }
})
