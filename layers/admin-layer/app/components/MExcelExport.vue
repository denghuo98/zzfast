<!-- Excel 导出，基于 xlsx + file-saver -->
<template>
  <ClientOnly>
    <UButton :color="color" :size="size" :disabled="disabled || !data.length" :loading="exporting" @click="handleExport">
      <template #leading>
        <UIcon name="i-lucide-download" v-if="!exporting" />
      </template>
      <slot>导出 Excel</slot>
    </UButton>
    <template #fallback>
      <UButton :color="color" :size="size" disabled>
        <slot>导出 Excel</slot>
      </UButton>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface ColumnConfig { title?: string; formatter?: (v: unknown) => string }

const props = withDefaults(defineProps<{
  data?: Record<string, unknown>[]
  filename?: string
  sheetName?: string
  color?: 'primary' | 'neutral' | 'success' | 'error' | 'warning' | 'info'
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  disabled?: boolean
  columns?: Record<string, ColumnConfig>
  headers?: Record<string, string>
  autoIndex?: boolean
  indexColumnTitle?: string
}>(), {
  data: () => [],
  filename: 'export',
  sheetName: 'Sheet1',
  color: 'success',
  size: 'sm',
  disabled: false,
  columns: () => ({}),
  headers: () => ({}),
  autoIndex: false,
  indexColumnTitle: '序号',
})

const emit = defineEmits<{
  success: [filename: string, rowCount: number]
  error: [error: Error]
}>()

const exporting = ref(false)

async function handleExport() {
  exporting.value = true
  try {
    const [XLSX, { saveAs }] = await Promise.all([
      import('xlsx'),
      import('file-saver'),
    ])

    const rows = props.data.map((row, idx) => {
      const out: Record<string, string> = {}
      if (props.autoIndex) out[props.indexColumnTitle] = String(idx + 1)
      for (const [k, v] of Object.entries(row)) {
        const title = props.columns[k]?.title ?? props.headers[k] ?? k
        const fmt = props.columns[k]?.formatter
        out[title] = fmt ? fmt(v) : (v == null ? '' : String(v))
      }
      return out
    })

    const ws = XLSX.utils.json_to_sheet(rows)
    // 自动列宽
    if (rows.length) {
      ws['!cols'] = Object.keys(rows[0]!).map(key => {
        const maxLen = Math.max(
          key.length,
          ...rows.slice(0, 100).map(r => (r[key] ?? '').length),
        )
        return { wch: Math.min(Math.max(maxLen + 2, 8), 50) }
      })
    }

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, props.sheetName)
    const buf = XLSX.write(wb, { bookType: 'xlsx', type: 'array', compression: true })
    const blob = new Blob([buf], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    saveAs(blob, `${props.filename}.xlsx`)
    emit('success', props.filename, rows.length)
  } catch (err) {
    emit('error', err as Error)
  } finally {
    exporting.value = false
  }
}
</script>
