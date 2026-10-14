<script setup lang="ts" generic="T extends z.ZodObject<any>">
import * as z from 'zod'

const props = defineProps<{
  schema: T
  state: Record<string, any>
  config?: Record<string, any>
  formClass?: string
  /** 搜索按钮文字 */
  searchLabel?: string
  /** 重置按钮文字 */
  resetLabel?: string
}>()

const emit = defineEmits<{
  search: [data: z.infer<T>]
  reset: []
  'update:state': [value: Record<string, any>]
}>()

const searchLabel = computed(() => props.searchLabel ?? '搜索')
const resetLabel = computed(() => props.resetLabel ?? '重置')

function onSearch(data: z.infer<T>) {
  emit('search', data)
}

function onReset() {
  emit('reset')
}
</script>

<template>
  <div id="m-search-form" class="shrink-0 bg-default rounded-md ring ring-default p-4">
    <MForm
      :schema="schema"
      :state="state"
      :form-class="formClass ?? 'flex flex-wrap items-end gap-3'"
      :config="config ?? { theme: { orientation: 'vertical' } }"
      @submit="onSearch"
      @update:state="emit('update:state', $event)"
    >
      <template v-for="(_, name) in $slots" :key="name" #[name]="slotData">
        <slot :name="name" v-bind="slotData ?? {}" />
      </template>

      <template #submit>
        <div class="flex items-end gap-2">
          <UButton type="submit" icon="i-lucide-search">
            {{ searchLabel }}
          </UButton>
          <UButton variant="outline" color="neutral" icon="i-lucide-rotate-ccw" @click="onReset">
            {{ resetLabel }}
          </UButton>
        </div>
      </template>
    </MForm>
  </div>
</template>
