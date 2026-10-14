/**
 * useTableColumns — 表格列配置管理
 */
import { ref, computed, watch } from 'vue'
import type { ColumnOption } from '../../types/table'

const SPECIAL_COLUMNS: Record<string, { accessorKey: string; header: string }> = {
  __selection__: { accessorKey: '__selection__', header: '选择' },
  __expand__: { accessorKey: '__expand__', header: '展开' },
  __index__: { accessorKey: '__index__', header: '序号' },
}

const getColumnKey = <T>(col: ColumnOption<T>) =>
  SPECIAL_COLUMNS[col.accessorKey as keyof typeof SPECIAL_COLUMNS]?.accessorKey ?? (col.accessorKey as string)

const getColumnVisibility = <T>(col: ColumnOption<T>): boolean => {
  if (col.checked !== undefined) return col.checked
  return col.visible ?? true
}

const getColumnChecks = <T>(columns: ColumnOption<T>[]) =>
  columns.map(col => {
    const special = col.accessorKey && SPECIAL_COLUMNS[col.accessorKey as keyof typeof SPECIAL_COLUMNS]
    const visibility = getColumnVisibility(col)
    if (special) return { ...col, accessorKey: special.accessorKey, header: special.header, checked: true, visible: true }
    return { ...col, checked: visibility, visible: visibility }
  })

export function useTableColumns<T = any>(
  columnsFactory: () => ColumnOption<T>[],
) {
  const dynamicColumns = ref<ColumnOption<T>[]>(columnsFactory())
  const columnChecks = ref<ColumnOption<T>[]>(getColumnChecks(dynamicColumns.value))
  const updateTrigger = ref(0)

  watch(dynamicColumns, (newCols) => {
    const visibilityMap = new Map(columnChecks.value.map(c => [getColumnKey(c), getColumnVisibility(c)]))
    columnChecks.value = getColumnChecks(newCols).map(c => {
      const key = getColumnKey(c)
      if (visibilityMap.has(key)) {
        const v = visibilityMap.get(key)!
        return { ...c, checked: v, visible: v }
      }
      return c
    })
  }, { deep: true })

  watch(columnChecks, () => { updateTrigger.value++ }, { deep: true })

  const columns = computed(() => {
    void updateTrigger.value
    const colMap = new Map(dynamicColumns.value.map(c => [getColumnKey(c), c]))
    return columnChecks.value
      .filter(c => getColumnVisibility(c))
      .map(c => colMap.get(getColumnKey(c)))
      .filter(Boolean) as ColumnOption<T>[]
  })

  const setDynamicColumns = (updater: (cols: ColumnOption<T>[]) => void | ColumnOption<T>[]) => {
    const copy = [...dynamicColumns.value]
    const result = updater(copy)
    dynamicColumns.value = Array.isArray(result) ? result : copy
  }

  return {
    columns,
    columnChecks,
    addColumn: (col: ColumnOption<T>, index?: number) =>
      setDynamicColumns((cols) => {
        if (typeof index === 'number' && index >= 0 && index <= cols.length)
          cols.splice(index, 0, col)
        else cols.push(col)
      }),
    removeColumn: (prop: string | string[]) =>
      setDynamicColumns(cols => cols.filter(c => !(Array.isArray(prop) ? prop : [prop]).includes(getColumnKey(c)))),
    toggleColumn: (prop: string, visible?: boolean) => {
      const i = columnChecks.value.findIndex(c => getColumnKey(c) === prop)
      if (i > -1) {
        const next = [...columnChecks.value]
        next[i] = { ...next[i], checked: visible ?? !getColumnVisibility(next[i]!), visible: visible ?? !getColumnVisibility(next[i]!) } as ColumnOption<T>
        columnChecks.value = next
      }
    },
    updateColumn: (prop: string, updates: Partial<ColumnOption<T>>) =>
      setDynamicColumns(cols => cols.map(c => getColumnKey(c) === prop ? { ...c, ...updates } : c)),
    batchUpdateColumns: (updates: Array<{ prop: string; updates: Partial<ColumnOption<T>> }>) =>
      setDynamicColumns((cols) => {
        const map = new Map(updates.map(u => [u.prop, u.updates]))
        return cols.map(c => { const upd = map.get(getColumnKey(c)); return upd ? { ...c, ...upd } as ColumnOption<T> : c })
      }),
    reorderColumns: (fromIdx: number, toIdx: number) =>
      setDynamicColumns((cols) => {
        if (fromIdx < 0 || fromIdx >= cols.length || toIdx < 0 || toIdx >= cols.length || fromIdx === toIdx) return cols
        const next = [...cols]; const [moved] = next.splice(fromIdx, 1)
        if (moved) next.splice(toIdx, 0, moved)
        return next
      }),
    getColumnConfig: (prop: string) => dynamicColumns.value.find(c => getColumnKey(c) === prop),
    getAllColumns: () => [...dynamicColumns.value],
    resetColumns: () => { dynamicColumns.value = columnsFactory() },
  }
}
