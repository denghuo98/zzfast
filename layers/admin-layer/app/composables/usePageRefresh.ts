/**
 * usePageRefresh — 软刷新当前页面内容（不刷新浏览器）。
 *
 * 通过切换 alive 标志让 admin layout 里的页面组件卸载后重新挂载，
 * setup / onMounted / 各类数据请求（含 useTable + $api 这类不走 Nuxt data 层的）全部重跑，
 * 同时刷新 useAsyncData / useFetch 缓存的数据。
 */
export function usePageRefresh() {
  const alive = useState('admin-page-alive', () => true)
  const refreshing = useState('admin-page-refreshing', () => false)

  async function reload() {
    if (refreshing.value) return
    refreshing.value = true
    // 卸载页面组件
    alive.value = false
    await nextTick()
    // 重新挂载 → 页面 setup / onMounted 重跑
    alive.value = true
    // 顺带刷新 Nuxt data 层（useAsyncData / useFetch）
    await refreshNuxtData()
    // 保留短暂动画视觉反馈
    setTimeout(() => {
      refreshing.value = false
    }, 300)
  }

  return { alive, refreshing, reload }
}
