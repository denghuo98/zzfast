<template>
  <!-- 展开时：logo + title -->
  <NuxtLink v-if="!collapsed" to="/" class="flex items-center gap-2 px-3 py-2 hover:opacity-80 transition-opacity">
    <img v-if="isImageLogo" :src="logo" alt="Logo" class="size-8 shrink-0" />
    <UIcon v-else :name="logo" class="size-8 shrink-0 text-white" />
    <span class="truncate text-lg font-semibold text-white">{{ title }}</span>
  </NuxtLink>
  <!-- 折叠时：仅 logo -->
  <NuxtLink v-else to="/" class="flex justify-center py-2 hover:opacity-80 transition-opacity">
    <img v-if="isImageLogo" :src="logo" alt="Logo" class="size-8 shrink-0" />
    <UIcon v-else :name="logo" class="size-8 shrink-0 text-white" />
  </NuxtLink>
</template>

<script setup lang="ts">
// 品牌名与 logo 通过 app.config 的 admin.title / admin.logo 配置，各项目可覆盖。
// admin.logo 支持图标名（如 i-lucide-shield）或图片路径/URL。
const appConfig = useAppConfig()
const admin = computed(() => appConfig.admin as { logo?: string; title?: string } | undefined)

const title = computed(() => admin.value?.title || 'Admin Panel')
const logo = computed(() => admin.value?.logo || 'i-lucide-shield')
const isImageLogo = computed(() => logo.value.startsWith('/') || logo.value.startsWith('http'))

defineProps<{
  collapsed: boolean
}>()
</script>
