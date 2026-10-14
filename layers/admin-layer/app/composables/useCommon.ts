/**
 * useCommon — 通用页面操作（滚动、刷新等）
 */
export function useCommon() {
  const homePath = computed(() => '/')

  const scrollToTop = () => {
    if (import.meta.server) return
    const el = document.getElementById('app-main')
    if (el) el.scrollTop = 0
  }

  const smoothScrollToTop = () => {
    if (import.meta.server) return
    const el = document.getElementById('app-main')
    if (el) el.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const scrollTo = (top: number, smooth = false) => {
    if (import.meta.server) return
    const el = document.getElementById('app-main')
    if (el) el.scrollTo({ top, behavior: smooth ? 'smooth' : 'auto' })
  }

  return { homePath, scrollToTop, smoothScrollToTop, scrollTo }
}
