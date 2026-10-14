<!-- 文字滚动：自研 CSS @keyframes 动画，支持 4 方向 + 速度 + 关闭按钮 -->
<template>
  <ClientOnly>
    <div v-if="visible" class="relative overflow-hidden rounded-md" :class="typeClass" ref="containerRef">
      <div
        ref="contentRef"
        class="whitespace-nowrap inline-block px-3 py-1.5 text-sm font-medium"
        :class="animateClass"
        :style="animateStyle"
      >
        {{ text }}
      </div>
      <button
        v-if="showClose"
        class="absolute right-1 top-1/2 -translate-y-1/2 p-0.5 rounded hover:bg-black/10 transition-colors"
        @click="handleClose"
      >
        <UIcon name="i-lucide-x" class="size-3.5" />
      </button>
    </div>
    <template #fallback>
      <div class="overflow-hidden rounded-md px-3 py-1.5 text-sm font-medium bg-muted text-muted-foreground">
        {{ text }}
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  text: string
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  speed?: number  // px/s
  direction?: 'left' | 'right' | 'up' | 'down'
  showClose?: boolean
  alwaysScroll?: boolean
}>(), {
  type: 'info',
  speed: 50,
  direction: 'left',
  showClose: false,
  alwaysScroll: true,
})

const emit = defineEmits<{ close: [] }>()

const visible = ref(true)
const shouldScroll = ref(props.alwaysScroll)
const containerRef = ref<HTMLElement | null>(null)
const contentRef = ref<HTMLElement | null>(null)

const typeClass = computed(() => ({
  primary: 'bg-primary text-primary-foreground',
  success: 'bg-green-500 text-white',
  warning: 'bg-amber-500 text-white',
  danger: 'bg-red-500 text-white',
  info: 'bg-blue-500 text-white',
}[props.type]))

const isHorizontal = computed(() => props.direction === 'left' || props.direction === 'right')

const animateClass = computed(() => {
  if (!shouldScroll.value) return ''
  return isHorizontal.value ? 'animate-marquee-x' : 'animate-marquee-y'
})

const duration = computed(() => {
  if (!shouldScroll.value) return 0
  return Math.max(3, props.text.length / (props.speed / 20))
})

const animateStyle = computed(() => {
  if (!shouldScroll.value) return {}
  const dur = `${duration.value}s`
  return {
    animationDuration: dur,
    animationTimingFunction: 'linear',
    animationIterationCount: 'infinite',
    '--marquee-dir': props.direction === 'right' || props.direction === 'down' ? 'reverse' : 'normal',
  }
})

// 检测内容宽度 → 决定是否滚动
if (!props.alwaysScroll) {
  onMounted(() => {
    if (containerRef.value && contentRef.value) {
      shouldScroll.value = contentRef.value.scrollWidth > containerRef.value.clientWidth
    }
  })
}

function handleClose() {
  visible.value = false
  emit('close')
}
</script>

<style scoped>
@keyframes marquee-x {
  0% { transform: translateX(100%); }
  100% { transform: translateX(-100%); }
}
@keyframes marquee-y {
  0% { transform: translateY(100%); }
  100% { transform: translateY(-100%); }
}
.animate-marquee-x {
  animation-name: marquee-x;
  animation-direction: var(--marquee-dir, normal);
}
.animate-marquee-y {
  animation-name: marquee-y;
  animation-direction: var(--marquee-dir, normal);
}
</style>
