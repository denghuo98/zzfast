<!-- 拖拽排序组件，基于 vue-draggable-plus (SortableJS) -->
<template>
  <ClientOnly>
    <VueDraggable
      :model-value="modelValue"
      :animation="animation"
      :handle="handle"
      :group="group"
      :target="target"
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <slot />
    </VueDraggable>
    <template #fallback>
      <div class="flex items-center justify-center h-32">
        <UIcon name="i-lucide-loader-2" class="animate-spin size-6 text-muted-foreground" />
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'

defineProps<{
  modelValue: any[]
  animation?: number
  handle?: string
  group?: string | { name: string; pull?: boolean | 'clone'; put?: boolean | string[] }
  target?: string
}>()

defineEmits<{
  'update:modelValue': [value: any[]]
}>()
</script>
