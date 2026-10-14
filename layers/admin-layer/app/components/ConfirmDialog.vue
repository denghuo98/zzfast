<script setup lang="ts">
interface ConfirmDialogProps {
  title?: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  confirmColor?: string
  icon?: string
  variant?: 'confirm' | 'danger'
}

const props = withDefaults(defineProps<ConfirmDialogProps>(), {
  title: '确认操作',
  confirmLabel: '确认',
  cancelLabel: '取消',
  confirmColor: 'primary',
  variant: 'confirm',
})

const emit = defineEmits<{
  close: [value: boolean]
}>()

const confirmColor = computed(() =>
  props.variant === 'danger' ? 'error' : props.confirmColor,
)
</script>

<template>
  <UModal
    :title="title"
    :dismissible="false"
    :ui="{ footer: 'justify-end' }"
  >
    <template v-if="icon || description" #body>
      <div class="flex items-start gap-3">
        <UIcon
          v-if="icon"
          :name="icon"
          class="size-5 shrink-0 mt-0.5"
          :class="variant === 'danger' ? 'text-error' : 'text-primary'"
        />
        <p v-if="description" class="text-sm text-muted">
          {{ description }}
        </p>
      </div>
    </template>

    <template #footer>
      <UButton :label="cancelLabel" color="neutral" variant="outline" @click="emit('close', false)" />
      <UButton :label="confirmLabel" :color="confirmColor" @click="emit('close', true)" />
    </template>
  </UModal>
</template>
