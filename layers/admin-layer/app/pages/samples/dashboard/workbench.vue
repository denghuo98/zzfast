<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">工作台</h1>
        <p class="text-muted-foreground mt-1">Welcome back! 这是仪表盘工作台页面。</p>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UCard v-for="stat in stats" :key="stat.label">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-primary/10">
            <UIcon :name="stat.icon" class="size-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{{ stat.value }}</p>
            <p class="text-sm text-muted-foreground">{{ stat.label }}</p>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const stats = [
  { label: 'Total Users', value: '12,480', icon: 'i-lucide-users' },
  { label: 'Active Sessions', value: '342', icon: 'i-lucide-activity' },
  { label: 'Revenue', value: '$48,320', icon: 'i-lucide-dollar-sign' },
  { label: 'Open Issues', value: '23', icon: 'i-lucide-alert-circle' },
]
</script>
