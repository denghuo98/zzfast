<script setup lang="ts">
const { refreshing, reload } = usePageRefresh()
</script>

<template>
  <UButton
    icon="lucide:rotate-cw"
    :class="{ 'animate-spin': refreshing, 'text-muted': true }"
    :disabled="refreshing"
    variant="ghost"
    color="neutral"
    square
    aria-label="刷新页面"
    @click="reload"
  />
</template>
