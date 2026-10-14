<!-- 水印组件：自研 Canvas 平铺文字水印，全屏 + 局部两种模式 -->
<template>
  <ClientOnly>
    <div :class="fullscreen ? 'fixed inset-0 pointer-events-none' : 'relative'" :style="{ zIndex: zIndex }">
      <div
        class="absolute inset-0 pointer-events-none"
        :style="{ backgroundImage: `url(${bgImage})`, backgroundRepeat: 'repeat', zIndex }"
      />
      <slot />
    </div>
    <template #fallback>
      <slot />
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'

interface WatermarkFont { fontSize?: number; color?: string; fontFamily?: string }

const props = withDefaults(defineProps<{
  content?: string | string[]
  font?: WatermarkFont
  rotate?: number
  gap?: [number, number]
  offset?: [number, number]
  opacity?: number
  zIndex?: number
  fullscreen?: boolean
}>(), {
  content: '',
  font: () => ({ fontSize: 16, color: 'rgba(128,128,128,0.15)', fontFamily: 'sans-serif' }),
  rotate: -22,
  gap: () => [200, 200] as [number, number],
  offset: () => [100, 100] as [number, number],
  opacity: 1,
  zIndex: 9999,
  fullscreen: false,
})

const bgImage = ref('')
const dpr = typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1

const texts = computed(() => {
  const c = props.content
  return Array.isArray(c) ? c : [c]
})

function drawWatermark() {
  if (typeof document === 'undefined') return
  if (!texts.value.length || !texts.value[0]) { bgImage.value = ''; return }

  const { fontSize = 16, color = 'rgba(128,128,128,0.15)', fontFamily = 'sans-serif' } = props.font
  const lineHeight = fontSize * 1.5
  const textWidth = Math.max(...texts.value.map(t => t.length * fontSize * 0.6))
  const textHeight = texts.value.length * lineHeight

  const rad = (props.rotate * Math.PI) / 180
  const absSin = Math.abs(Math.sin(rad))
  const absCos = Math.abs(Math.cos(rad))
  const cw = textWidth * absCos + textHeight * absSin + 20
  const ch = textWidth * absSin + textHeight * absCos + 20

  const canvas = document.createElement('canvas')
  canvas.width = cw * dpr
  canvas.height = ch * dpr
  const ctx = canvas.getContext('2d')!
  ctx.scale(dpr, dpr)
  ctx.globalAlpha = props.opacity
  ctx.font = `${fontSize}px ${fontFamily}`
  ctx.fillStyle = color
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'

  ctx.save()
  ctx.translate(cw / 2, ch / 2)
  ctx.rotate(rad)
  texts.value.forEach((t, i) => {
    ctx.fillText(t, 0, -(textHeight / 2) + lineHeight * (i + 0.5))
  })
  ctx.restore()

  bgImage.value = canvas.toDataURL()
}

onMounted(drawWatermark)
watch([() => props.content, () => props.font, () => props.rotate, () => props.opacity], drawWatermark)
</script>
