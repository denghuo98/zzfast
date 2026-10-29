<!--
  MImageUpload —— 通用「图片上传 + 裁剪」组件
  · 选择本地图片 → 弹窗用 cropperjs 裁剪（固定裁剪框，可拖动/缩放）→ 导出上传 → 回填 URL
  · 上传接口默认 /admin/upload/file（form-data：file + uploadType 头），返回 data.fileUrl
  用法：<MImageUpload v-model="form.avatar" :aspect="3/4" />
-->
<template>
  <div class="flex items-center gap-4">
    <div
      class="flex size-16 shrink-0 items-center justify-center overflow-hidden bg-muted ring ring-default"
      :class="previewShape === 'circle' ? 'rounded-full' : 'rounded-md'"
    >
      <img v-if="modelValue" :src="modelValue" alt="图片" class="size-full object-cover">
      <UIcon v-else name="i-lucide-image" class="size-6 text-muted" />
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <UButton size="sm" variant="soft" color="primary" icon="i-lucide-upload" @click="pick">
        {{ modelValue ? '更换图片' : '上传图片' }}
      </UButton>
      <UButton v-if="modelValue" size="sm" variant="ghost" color="error" @click="emit('update:modelValue', '')">
        移除
      </UButton>
    </div>

    <input ref="fileInput" type="file" :accept="accept" class="hidden" @change="onPick">

    <UModal v-model:open="cropOpen" title="裁剪图片">
      <template #body>
        <div class="space-y-3">
          <div class="flex h-80 w-full items-center justify-center overflow-hidden bg-muted">
            <img v-if="srcUrl" ref="imgRef" :src="srcUrl" alt="待裁剪图片" class="block">
          </div>
          <div class="flex items-center gap-3">
            <span class="shrink-0 text-sm text-muted">缩放</span>
            <input v-model="zoom" type="range" min="0.5" max="3" step="0.01" class="flex-1" @input="onZoomInput">
            <UButton size="xs" variant="ghost" color="neutral" icon="i-lucide-rotate-ccw" @click="resetCrop">重置</UButton>
          </div>
        </div>
      </template>
      <template #footer>
        <div class="flex w-full justify-end gap-2">
          <UButton color="neutral" variant="ghost" @click="cropOpen = false">取消</UButton>
          <UButton color="primary" :loading="saving" @click="confirmCrop">确定</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'

defineOptions({ name: 'MImageUpload' })

const props = withDefaults(defineProps<{
  modelValue?: string
  /** 裁剪比例（宽/高），默认 1:1；竖版头像可用 3/4 */
  aspect?: number
  /** 导出宽度（px），高度按 aspect 计算 */
  outputWidth?: number
  /** 输出格式 */
  outputType?: 'png' | 'jpeg' | 'webp'
  /** 可选择的文件类型 */
  accept?: string
  /** 预览图形状 */
  previewShape?: 'square' | 'circle'
  /** 上传接口 */
  uploadUrl?: string
  /** 上传类型（作为 uploadType 请求头发给后端） */
  uploadType?: string
}>(), {
  modelValue: '',
  aspect: 1,
  outputWidth: 600,
  outputType: 'png',
  accept: 'image/*',
  previewShape: 'square',
  uploadUrl: '/admin/upload/file',
  uploadType: 'image',
})

const emit = defineEmits<{ 'update:modelValue': [v: string] }>()

const api = useApi()
const toast = useToast()

const fileInput = ref<HTMLInputElement | null>(null)
const imgRef = ref<HTMLImageElement | null>(null)
const cropOpen = ref(false)
const srcUrl = ref('')
const zoom = ref(1)
const saving = ref(false)
let cropper: Cropper | null = null

function pick() {
  fileInput.value?.click()
}

function onPick(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (!file.type.startsWith('image/')) {
    toast.add({ title: '请选择图片文件', color: 'error' })
    return
  }
  if (srcUrl.value) URL.revokeObjectURL(srcUrl.value)
  srcUrl.value = URL.createObjectURL(file)
  cropOpen.value = true
}

function initCropper() {
  const img = imgRef.value
  if (!img) return
  cropper?.destroy()
  cropper = new Cropper(img, {
    viewMode: 1,
    dragMode: 'move',
    aspectRatio: props.aspect,
    autoCropArea: 0.9,
    cropBoxMovable: false,
    cropBoxResizable: false,
    background: false,
    guides: false,
    center: false,
    highlight: false,
    responsive: true,
    zoomable: true,
    zoomOnWheel: true,
  })
  zoom.value = 1
}

watch(cropOpen, async (open) => {
  if (open) {
    await nextTick()
    initCropper()
  } else {
    cropper?.destroy()
    cropper = null
  }
})

function onZoomInput() {
  cropper?.zoomTo(Number(zoom.value))
}

function resetCrop() {
  cropper?.reset()
  zoom.value = 1
}

async function confirmCrop() {
  if (!cropper) return
  const width = props.outputWidth
  const height = Math.round(width / props.aspect)
  const canvas = cropper.getCroppedCanvas({ width, height, fillColor: '#ffffff', imageSmoothingQuality: 'high' })
  if (!canvas) return

  saving.value = true
  try {
    const mime = `image/${props.outputType}`
    const blob = await new Promise<Blob | null>(resolve => canvas.toBlob(resolve, mime, 0.92))
    if (!blob) throw new Error('裁剪失败')

    const formData = new FormData()
    formData.append('file', new File([blob], `image-${Date.now()}.${props.outputType}`, { type: mime }))

    const data = await api.request<{ fileUrl: string }>(props.uploadUrl, {
      method: 'POST',
      body: formData,
      headers: { uploadType: props.uploadType },
    })
    if (!data?.fileUrl) throw new Error('上传失败')

    emit('update:modelValue', data.fileUrl)
    toast.add({ title: '图片已上传', color: 'success' })
    cropOpen.value = false
  } catch (e: any) {
    toast.add({ title: '上传失败', description: e.message, color: 'error' })
  } finally {
    saving.value = false
  }
}

onBeforeUnmount(() => {
  cropper?.destroy()
  if (srcUrl.value) URL.revokeObjectURL(srcUrl.value)
})
</script>
