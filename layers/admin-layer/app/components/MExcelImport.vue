<!-- Excel 导入，基于 xlsx。支持多 sheet，emit 含 sheetNames -->
<template>
  <ClientOnly>
    <div class="inline-block">
      <UButton :color="color" :size="size" :disabled="disabled" as="label" class="cursor-pointer">
        <slot>导入 Excel</slot>
        <input type="file" accept=".xlsx,.xls" class="hidden" @change="handleFile" />
      </UButton>
    </div>
    <template #fallback>
      <UButton :color="color" :size="size" :disabled="disabled">
        <slot>导入 Excel</slot>
      </UButton>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  color?: 'primary' | 'neutral' | 'success' | 'error' | 'warning' | 'info'
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  disabled?: boolean
}>(), {
  color: 'primary',
  size: 'sm',
  disabled: false,
})

const emit = defineEmits<{
  success: [result: { sheets: Record<string, Record<string, unknown>[]>; sheetNames: string[] }]
  error: [error: Error]
}>()

async function handleFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const XLSX = await import('xlsx')
    const buf = await file.arrayBuffer()
    const wb = XLSX.read(buf, { type: 'array' })
    const sheets: Record<string, Record<string, unknown>[]> = {}
    for (const name of wb.SheetNames) {
      const ws = wb.Sheets[name]!
      sheets[name] = XLSX.utils.sheet_to_json(ws) as Record<string, unknown>[]
    }
    emit('success', { sheets, sheetNames: wb.SheetNames })
  } catch (err) {
    emit('error', err as Error)
  }
  input.value = ''
}
</script>
