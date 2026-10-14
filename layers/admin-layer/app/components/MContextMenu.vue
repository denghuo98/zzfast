<!-- 右键菜单：基于 Nuxt UI UContextMenu 的薄封装 -->
<script setup lang="ts">
import type { ContextMenuItem } from '@nuxt/ui'

defineProps<{
  items: ContextMenuItem[][]
}>()
</script>

<template>
  <UContextMenu :items="items">
    <slot />
  </UContextMenu>
</template>
