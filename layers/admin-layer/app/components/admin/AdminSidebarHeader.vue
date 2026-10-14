<template>
  <!-- 展开时：logo + title -->
  <NuxtLink v-if="!collapsed" to="/" class="flex items-center gap-2 px-3 py-2 hover:opacity-80 transition-opacity">
    <!-- light mode -->
    <img :src="logoLight" alt="Logo" class="size-8 shrink-0 dark:hidden" />
    <img :src="titleLight" alt="Title" class="h-4 shrink-0 dark:hidden" />
    <!-- dark mode -->
    <img :src="logoDark" alt="Logo" class="size-8 shrink-0 hidden dark:block" />
    <img :src="titleDark" alt="Title" class="h-4 shrink-0 hidden dark:block" />
  </NuxtLink>
  <!-- 折叠时：仅 logo -->
  <NuxtLink v-else to="/" class="flex justify-center py-2 hover:opacity-80 transition-opacity">
    <!-- light mode -->
    <img :src="logoLight" alt="Logo" class="size-8 shrink-0 dark:hidden" />
    <!-- dark mode -->
    <img :src="logoDark" alt="Logo" class="size-8 shrink-0 hidden dark:block" />
  </NuxtLink>
</template>

<script setup lang="ts">
import logoLight from '../../assets/image/matrix-logo.png'
import logoDark from '../../assets/image/matrix-logo-dark.png'
import titleLight from '../../assets/image/matrix-title.svg'
import titleDark from '../../assets/image/matrix-title-dark.svg'

defineProps<{
  collapsed: boolean
}>()
</script>
