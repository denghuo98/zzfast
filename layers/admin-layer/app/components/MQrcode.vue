<!-- 二维码组件：基于 qrcode.vue，支持 svg/canvas + 下载 -->
<template>
  <ClientOnly>
    <div ref="containerRef">
      <QrcodeVue
        :value="value"
        :size="size"
        :level="level"
        :render-as="renderAs"
        :margin="margin"
        :background="background"
        :foreground="foreground"
        :image-settings="logoSettings"
      />
    </div>
    <template #fallback>
      <div class="flex items-center justify-center" :style="{ width: `${size}px`, height: `${size}px` }">
        <UIcon name="i-lucide-loader-2" class="animate-spin size-8 text-muted-foreground" />
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import QrcodeVue from 'qrcode.vue'
import { saveAs } from 'file-saver'

interface LogoConfig { src?: string; width?: number; height?: number }

const props = withDefaults(defineProps<{
  value: string
  size?: number
  level?: 'L' | 'M' | 'Q' | 'H'
  renderAs?: 'svg' | 'canvas'
  margin?: number
  background?: string
  foreground?: string
  logo?: LogoConfig | null
}>(), {
  size: 200,
  level: 'M',
  renderAs: 'svg',
  margin: 2,
  background: '#ffffff',
  foreground: '#000000',
  logo: null,
})

const containerRef = ref<HTMLElement | null>(null)

const logoSettings = computed(() => {
  if (!props.logo) return undefined
  return {
    src: props.logo.src,
    ...(props.logo.width && props.logo.height ? { width: props.logo.width, height: props.logo.height } : {}),
  }
})

// 直接从容器 DOM 找 canvas/svg，不依赖 qrcode.vue 内部组件实例 $el
function findCanvas(): HTMLCanvasElement | null {
  return containerRef.value?.querySelector?.('canvas') ?? null
}
function findSvg(): SVGSVGElement | null {
  return containerRef.value?.querySelector?.('svg') ?? null
}

function toDataURL(): string | undefined {
  const canvas = findCanvas()
  if (canvas) return canvas.toDataURL('image/png')

  // SVG 模式
  const svg = findSvg()
  if (svg && typeof window !== 'undefined') {
    const serializer = new XMLSerializer()
    return 'data:image/svg+xml;base64,' + window.btoa(unescape(encodeURIComponent(serializer.serializeToString(svg))))
  }
  return undefined
}

function download(filename = 'qrcode.png') {
  const canvas = findCanvas()
  if (canvas) {
    canvas.toBlob((blob) => {
      if (blob) saveAs(blob, filename)
    }, 'image/png')
    return
  }

  // SVG fallback: render to canvas → blob
  const svg = findSvg()
  if (!svg || typeof window === 'undefined') return

  const c = document.createElement('canvas')
  c.width = props.size
  c.height = props.size
  const ctx = c.getContext('2d')!
  const img = new Image()
  const svgData = new XMLSerializer().serializeToString(svg)
  const blob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  img.onload = () => {
    ctx.drawImage(img, 0, 0)
    c.toBlob((b) => {
      if (b) saveAs(b, filename)
      URL.revokeObjectURL(url)
    }, 'image/png')
  }
  img.src = url
}

defineExpose({ toDataURL, download })
</script>
