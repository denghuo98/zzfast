<!-- 全屏切换按钮 -->
<script setup lang="ts">
import { useFullscreen } from '@vueuse/core'

const { isFullscreen, toggle } = useFullscreen()
</script>

<template>
  <UTooltip :text="isFullscreen ? '退出全屏' : '全屏'">
    <UButton
      :icon="isFullscreen ? 'i-lucide-minimize-2' : 'i-lucide-maximize-2'"
      color="neutral"
      variant="ghost"
      size="sm"
      class="text-muted"
      @click="toggle()"
    />
  </UTooltip>
</template>
