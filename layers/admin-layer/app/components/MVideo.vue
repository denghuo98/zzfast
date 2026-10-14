<!-- 视频播放器：原生 video + @vueuse/IntersectionObserver 视口自动播放 -->
<template>
  <div class="relative" ref="containerRef">
    <video
      ref="videoRef"
      :src="src"
      :poster="poster"
      :controls="controls"
      :loop="loop"
      :muted="muted || autoplay"
      :autoplay="false"
      class="w-full rounded-md"
      preload="metadata"
    />
    <slot />
  </div>
</template>

<script setup lang="ts">
import { useIntersectionObserver } from '@vueuse/core'

const props = withDefaults(defineProps<{
  src: string
  poster?: string
  autoplay?: boolean
  loop?: boolean
  muted?: boolean
  controls?: boolean
}>(), {
  controls: true,
  loop: false,
  muted: false,
  autoplay: false,
})

const videoRef = ref<HTMLVideoElement | null>(null)
const containerRef = ref<HTMLElement | null>(null)

// 视口内自动播放，出视口暂停
const { stop } = useIntersectionObserver(containerRef, ([{ isIntersecting }]) => {
  if (!props.autoplay || !videoRef.value) return
  if (isIntersecting) {
    videoRef.value.play().catch(() => {})
  } else {
    videoRef.value.pause()
  }
})

onBeforeUnmount(() => {
  stop()
})

function play() { videoRef.value?.play() }
function pause() { videoRef.value?.pause() }
function getCurrentTime() { return videoRef.value?.currentTime || 0 }

defineExpose({ play, pause, getCurrentTime, videoRef })
</script>
