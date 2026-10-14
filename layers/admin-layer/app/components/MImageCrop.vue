<!-- 图片裁剪，基于 cropperjs。左侧裁剪区，右侧实时预览 -->
<template>
  <ClientOnly>
    <div class="flex flex-col lg:flex-row gap-4">
      <!-- 左侧：裁剪操作区 -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-3">
          <UInput type="file" accept="image/*" @change="handleFile" size="sm" />
        </div>
        <div class="border border-default rounded-md overflow-hidden">
          <img
            ref="imgRef"
            :src="src"
            alt="裁剪源图"
            class="block max-w-full max-h-[55vh] object-contain"
            @load="initCropper"
          />
        </div>
      </div>

      <!-- 右侧：预览 + 下载 -->
      <div class="w-full lg:w-56 shrink-0">
        <h4 class="text-sm font-semibold mb-2">预览效果</h4>
        <div class="border border-default rounded-md overflow-hidden bg-muted mb-3 flex items-center justify-center min-h-[160px]">
          <img v-if="previewSrc" :src="previewSrc" alt="裁剪预览" class="max-w-full" />
          <span v-else class="text-sm text-muted-foreground">调整裁剪框后预览</span>
        </div>
        <UButton block size="sm" color="primary" :disabled="!src" @click="download">
          <template #leading><UIcon name="i-lucide-download" /></template>
          下载裁剪结果
        </UButton>
      </div>
    </div>
    <template #fallback>
      <div class="flex items-center justify-center h-64">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin text-muted" />
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'

defineOptions({ name: 'MImageCrop' })

const props = withDefaults(defineProps<{
  src?: string
  aspectRatio?: number
  viewMode?: number
  outputType?: 'png' | 'jpeg' | 'webp'
  quality?: number
}>(), {
  aspectRatio: NaN,
  viewMode: 1,
  outputType: 'png',
  quality: 0.92,
})

const imgRef = ref<HTMLImageElement>()
const previewSrc = ref('')
let cropper: Cropper | null = null

function initCropper() {
  if (!imgRef.value || !imgRef.value.complete || imgRef.value.naturalWidth === 0) return
  cropper?.destroy()
  cropper = new Cropper(imgRef.value, {
    viewMode: props.viewMode,
    aspectRatio: props.aspectRatio,
    dragMode: 'move',
    autoCropArea: 1,
    responsive: true,
    crop() {
      updatePreview()
    },
    ready() {
      updatePreview()
    },
  })
}

function updatePreview() {
  if (!cropper) return
  const canvas = cropper.getCroppedCanvas()
  previewSrc.value = canvas?.toDataURL(`image/${props.outputType}`, props.quality) ?? ''
}

function download() {
  if (!cropper) return
  const canvas = cropper.getCroppedCanvas()
  canvas.toBlob((blob) => {
    if (!blob) return
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = `crop.${props.outputType}`
    a.click()
    URL.revokeObjectURL(a.href)
  }, `image/${props.outputType}`, props.quality)
}

function toDataURL(): string | null {
  return cropper?.getCroppedCanvas()?.toDataURL(`image/${props.outputType}`, props.quality) ?? null
}

function toBlob(): Promise<Blob | null> {
  return new Promise((resolve) => {
    cropper?.getCroppedCanvas().toBlob(resolve, `image/${props.outputType}`, props.quality)
  })
}

function handleFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => { imgRef.value!.src = reader.result as string }
  reader.readAsDataURL(file)
}

defineExpose({ toDataURL, toBlob, download })

onBeforeUnmount(() => { cropper?.destroy(); cropper = null })
</script>
