<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">横幅</h1>
      <p class="text-muted-foreground mt-1">Banner / Alert 组件示例。</p>
    </div>

    <UCard>
      <div class="space-y-4">
        <div v-for="banner in banners" :key="banner.title"
          class="p-4 rounded-lg border"
          :class="{
            'bg-primary/5 border-primary/20': banner.type === 'info',
            'bg-success/5 border-success/20': banner.type === 'success',
            'bg-warning/5 border-warning/20': banner.type === 'warning',
            'bg-error/5 border-error/20': banner.type === 'error',
          }">
          <div class="flex items-start gap-3">
            <UIcon :name="banner.icon" class="size-5 mt-0.5 shrink-0"
              :class="{
                'text-primary': banner.type === 'info',
                'text-success': banner.type === 'success',
                'text-warning': banner.type === 'warning',
                'text-error': banner.type === 'error',
              }" />
            <div>
              <h4 class="font-semibold text-sm">{{ banner.title }}</h4>
              <p class="text-sm text-muted-foreground mt-1">{{ banner.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const banners = [
  { type: 'info', title: 'Information', description: 'This is an informational banner with general tips.', icon: 'i-lucide-info' },
  { type: 'success', title: 'Success', description: 'Operation completed successfully!', icon: 'i-lucide-circle-check' },
  { type: 'warning', title: 'Warning', description: 'Please review your settings before proceeding.', icon: 'i-lucide-triangle-alert' },
  { type: 'error', title: 'Error', description: 'Something went wrong. Please try again.', icon: 'i-lucide-octagon-x' },
]
</script>
