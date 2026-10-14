<template>
  <UNavigationMenu
    :collapsed="collapsed"
    :items="menuItems"
    orientation="vertical"
    popover
    highlight
    :ui="{ link: 'text-sm py-2.5', list: 'gap-0.5' }"
  />
</template>

<script setup lang="ts">
defineProps<{
  collapsed: boolean
}>()

const { menuItems } = useAdminMenu()
</script>
