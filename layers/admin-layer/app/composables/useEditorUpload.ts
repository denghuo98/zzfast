// 编辑器图片上传配置（模块级共享，因为 TipTap NodeView 不在 MRichText 的 Vue 组件树内）
type UploadHandler = (file: File) => Promise<{ url: string }>

const _uploadHandler = ref<UploadHandler | null>(null)
const _uploadUrl = ref('/api/upload')

export function useEditorUpload() {
  function setUploadConfig(config: { url?: string; handler?: UploadHandler }) {
    if (config.url) _uploadUrl.value = config.url
    if (config.handler) _uploadHandler.value = config.handler
  }

  async function upload(file: File): Promise<{ url: string }> {
    // 优先使用自定义 handler
    if (_uploadHandler.value) {
      return _uploadHandler.value(file)
    }

    // 默认：走配置的 URL
    const formData = new FormData()
    formData.append('file', file)
    const res = await $fetch<{ url: string; pathname: string }>(_uploadUrl.value, {
      method: 'POST',
      body: formData,
    })
    return { url: res.url || `/uploads/${res.pathname}` }
  }

  return { setUploadConfig, upload }
}
