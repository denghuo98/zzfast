<script setup lang="ts">
import type { NodeViewProps } from '@tiptap/vue-3'
import { NodeViewWrapper } from '@tiptap/vue-3'

const props = defineProps<NodeViewProps>()

const { upload } = useEditorUpload()

const fileInputRef = ref<HTMLInputElement | null>(null)
const error = ref<string | null>(null)
const loading = ref(false)
const isDragging = ref(false)

async function uploadFile(file: File) {
  loading.value = true
  error.value = null

  try {
    const result = await upload(file)

    const pos = props.getPos()
    if (typeof pos !== 'number') return

    props.editor
      .chain()
      .focus()
      .deleteRange({ from: pos, to: pos + 1 })
      .setImage({ src: result.url })
      .run()
  } catch (e) {
    error.value = (e as any)?.data?.message || (e as any)?.message || '上传失败，请重试'
  } finally {
    loading.value = false
  }
}

function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) uploadFile(file)
  // reset input so same file can be re-selected
  if (fileInputRef.value) fileInputRef.value.value = ''
}

function onDrop(e: DragEvent) {
  const file = e.dataTransfer?.files?.[0]
  if (file) uploadFile(file)
  isDragging.value = false
}

function onDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
}

function onDragLeave() {
  isDragging.value = false
}

function openFilePicker() {
  fileInputRef.value?.click()
}
</script>

<template>
  <NodeViewWrapper>
    <div
      class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors"
      :class="[
        isDragging ? 'border-primary bg-primary/5' : error ? 'border-red-300 bg-red-50 dark:bg-red-950/20' : 'border-default hover:border-primary/50 hover:bg-muted/50',
        loading ? 'pointer-events-none opacity-60' : ''
      ]"
      @click="openFilePicker"
      @dragover="onDragOver"
      @dragleave="onDragLeave"
      @drop.prevent="onDrop"
    >
      <input ref="fileInputRef" type="file" accept="image/*" class="hidden" @change="onFileChange" />

      <div class="flex flex-col items-center gap-2">
        <UIcon
          :name="error ? 'i-lucide-alert-circle' : loading ? 'i-lucide-loader-2' : 'i-lucide-image-plus'"
          class="size-10"
          :class="[loading && 'animate-spin', error ? 'text-red-500' : 'text-muted-foreground']"
        />
        <p class="text-sm font-medium" :class="error ? 'text-red-500' : 'text-muted-foreground'">
          {{ error || (loading ? '上传中...' : '点击或拖拽图片到此处上传') }}
        </p>
        <p class="text-xs text-muted-foreground" v-if="!error && !loading">
          支持 SVG、PNG、JPG、GIF（最大 2MB）
        </p>
      </div>
    </div>
  </NodeViewWrapper>
</template>
