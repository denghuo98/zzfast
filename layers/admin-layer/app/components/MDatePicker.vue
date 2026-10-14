<script setup lang="ts">
/**
 * MDatePicker — 日期选择器（Popover + Calendar）
 *
 * - v-model 绑定 'YYYY-MM-DD' 字符串（便于直接提交后端）
 * - 内部用 @internationalized/date 的 CalendarDate 处理日期
 * - 触发器按钮宽度跟随父级（class="w-full"），与表单其他输入一致
 */
import { CalendarDate, DateFormatter, getLocalTimeZone, parseDate } from '@internationalized/date'

const props = defineProps<{
  modelValue?: string | null
  placeholder?: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const df = new DateFormatter('zh-CN', { dateStyle: 'medium' })

// 字符串 ↔ CalendarDate 双向桥接
const calendarValue = computed<CalendarDate | undefined>({
  get() {
    if (!props.modelValue) return undefined
    try { return parseDate(props.modelValue) } catch { return undefined }
  },
  set(v) {
    if (!v) { emit('update:modelValue', ''); return }
    // CalendarDate.toString() 已是 ISO 'YYYY-MM-DD'
    emit('update:modelValue', v.toString())
  },
})

const displayText = computed(() => {
  if (!calendarValue.value) return props.placeholder || '选择日期'
  return df.format(calendarValue.value.toDate(getLocalTimeZone()))
})
</script>

<template>
  <UPopover :disabled="disabled">
    <UButton
      color="neutral"
      variant="outline"
      icon="i-lucide-calendar"
      :disabled="disabled"
      class="w-full justify-start"
    >
      {{ displayText }}
    </UButton>

    <template #content>
      <UCalendar v-model="calendarValue" class="p-2" />
    </template>
  </UPopover>
</template>
