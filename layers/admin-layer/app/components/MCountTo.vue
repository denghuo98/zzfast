<!-- 数字滚动动画，基于 @vueuse/core useTransition -->
<template>
  <span class="tabular-nums">
    {{ formattedValue }}
  </span>
</template>

<script setup lang="ts">
import { computed, watch, nextTick, onUnmounted, shallowRef } from 'vue'
import { useTransition, TransitionPresets } from '@vueuse/core'

const props = withDefaults(defineProps<{
  target?: number
  duration?: number
  autoStart?: boolean
  decimals?: number
  decimal?: string
  separator?: string
  prefix?: string
  suffix?: string
  easing?: keyof typeof TransitionPresets
  disabled?: boolean
}>(), {
  target: 0,
  duration: 2000,
  autoStart: true,
  decimals: 0,
  decimal: '.',
  separator: '',
  prefix: '',
  suffix: '',
  easing: 'easeOutExpo',
  disabled: false,
})

const emit = defineEmits<{
  started: [value: number]
  finished: [value: number]
  paused: [value: number]
  reset: []
}>()

const EPSILON = Number.EPSILON

function formatNumber(v: number): string {
  let r = props.decimals > 0 ? v.toFixed(props.decimals) : Math.floor(v).toString()
  if (props.separator) {
    const parts = r.split('.')
    parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, props.separator)
    r = parts.join('.')
  }
  return r
}

const safeEasing = computed(() => {
  const e = props.easing
  return (e in TransitionPresets) ? e : 'easeOutExpo'
})

const currentValue = shallowRef(0)
const targetValue = shallowRef(props.target)
const isRunning = shallowRef(false)
const isPaused = shallowRef(false)
const pausedValue = shallowRef(0)

const transitionValue = useTransition(currentValue, {
  duration: computed(() => props.duration),
  transition: computed(() => TransitionPresets[safeEasing.value]),
  onStarted() {
    isRunning.value = true
    isPaused.value = false
    emit('started', targetValue.value)
  },
  onFinished() {
    isRunning.value = false
    isPaused.value = false
    emit('finished', targetValue.value)
  },
})

const formattedValue = computed(() => {
  const v = isPaused.value ? pausedValue.value : transitionValue.value
  if (!Number.isFinite(v)) return ''
  return `${props.prefix}${formatNumber(v)}${props.suffix}`
})

function start(target?: number) {
  if (props.disabled) return
  const t = target ?? targetValue.value
  if (!Number.isFinite(t)) return
  targetValue.value = t
  const cur = isPaused.value ? pausedValue.value : transitionValue.value
  if (Math.abs(cur - t) < EPSILON) return
  if (isPaused.value) {
    currentValue.value = pausedValue.value
    isPaused.value = false
    pausedValue.value = 0
  }
  nextTick(() => { currentValue.value = t })
}

function pause() {
  if (!isRunning.value || isPaused.value) return
  isPaused.value = true
  pausedValue.value = transitionValue.value
  currentValue.value = pausedValue.value
  emit('paused', pausedValue.value)
}

function reset(newTarget = 0) {
  currentValue.value = newTarget
  targetValue.value = newTarget
  isPaused.value = false
  pausedValue.value = 0
  emit('reset')
}

function stop() {
  if (isRunning.value || isPaused.value) {
    currentValue.value = 0
    isPaused.value = false
    pausedValue.value = 0
    emit('paused', 0)
  }
}

function setTarget(target: number) {
  if (!Number.isFinite(target)) return
  targetValue.value = target
  if ((isRunning.value || props.autoStart) && !props.disabled) start(target)
}

watch(() => props.target, (t) => {
  if (props.autoStart && !props.disabled) start(t)
  else targetValue.value = t
}, { immediate: props.autoStart && !props.disabled })

watch(() => props.disabled, (d) => { if (d && isRunning.value) stop() })

onUnmounted(() => { if (isRunning.value) stop() })

defineExpose({
  start, pause, reset, stop, setTarget,
  get isRunning() { return isRunning.value },
  get isPaused() { return isPaused.value },
  get currentValue() { return isPaused.value ? pausedValue.value : transitionValue.value },
  get targetValue() { return targetValue.value },
  get progress() {
    const cur = isPaused.value ? pausedValue.value : transitionValue.value
    const tgt = targetValue.value
    return tgt === 0 ? (cur === 0 ? 1 : 0) : Math.abs(cur / tgt)
  },
})
</script>
