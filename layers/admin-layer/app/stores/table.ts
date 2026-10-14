/**
 * 表格外观配置 Store — 边框/斑马纹/尺寸/表头背景
 */
import { defineStore } from 'pinia'
import { TableSizeEnum } from '../enums/formEnum'

export const useTableStore = defineStore('tableStore', () => {
  // 仅在客户端使用 localStorage
  const storage = (key: string, fallback: any) => {
    if (import.meta.server) return ref(fallback)
    try {
      const stored = localStorage.getItem(key)
      return ref(stored ? JSON.parse(stored) : fallback)
    } catch { return ref(fallback) }
  }

  const tableSize = ref<TableSizeEnum>(TableSizeEnum.DEFAULT)
  const isZebra = ref(true)
  const isBorder = ref(true)
  const isHeaderBackground = ref(true)
  const isFullScreen = ref(false)

  // 客户端挂载后从 localStorage 恢复
  onMounted(() => {
    try {
      const saved = localStorage.getItem('table.tableSize')
      if (saved) tableSize.value = JSON.parse(saved) as TableSizeEnum
    } catch { /* ignore */ }
    try {
      const saved = localStorage.getItem('table.isZebra')
      if (saved !== null) isZebra.value = JSON.parse(saved)
    } catch { /* ignore */ }
    try {
      const saved = localStorage.getItem('table.isBorder')
      if (saved !== null) isBorder.value = JSON.parse(saved)
    } catch { /* ignore */ }
    try {
      const saved = localStorage.getItem('table.isHeaderBackground')
      if (saved !== null) isHeaderBackground.value = JSON.parse(saved)
    } catch { /* ignore */ }
  })

  // 持久化的 setter
  const persist = (key: string, val: any) => {
    if (import.meta.client) localStorage.setItem(key, JSON.stringify(val))
  }

  return {
    tableSize, isZebra, isBorder, isHeaderBackground, isFullScreen,
    setTableSize: (s: TableSizeEnum) => { tableSize.value = s; persist('table.tableSize', s) },
    setIsZebra: (v: boolean) => { isZebra.value = v; persist('table.isZebra', v) },
    setIsBorder: (v: boolean) => { isBorder.value = v; persist('table.isBorder', v) },
    setIsHeaderBackground: (v: boolean) => { isHeaderBackground.value = v; persist('table.isHeaderBackground', v) },
    setIsFullScreen: (v: boolean) => { isFullScreen.value = v; persist('table.isFullScreen', v) },
  }
})
