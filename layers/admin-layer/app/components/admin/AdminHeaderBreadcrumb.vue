<template>
  <div class="flex items-center gap-2 text-sm">
    <template v-for="(crumb, index) in crumbs" :key="index">
      <span v-if="index > 0" class="text-muted-foreground">/</span>
      <ULink
        v-if="crumb.to && index < crumbs.length - 1"
        :to="crumb.to"
        class="text-muted-foreground hover:text-foreground transition-colors"
      >
        {{ crumb.label }}
      </ULink>
      <span v-else class="font-medium">
        {{ crumb.label }}
      </span>
    </template>
  </div>
</template>

<script setup lang="ts">
const { breadcrumbs } = useAdminMenu()

const crumbs = computed(() => {
  if (breadcrumbs.value.length > 0) {
    return [{ label: '首页', to: '/' }, ...breadcrumbs.value]
  }
  return [{ label: '首页' }]
})
</script>
