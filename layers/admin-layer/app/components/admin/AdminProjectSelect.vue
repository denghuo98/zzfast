<script setup lang="ts">
interface Project {
  id: number
  name: string
  code: string
  cost_code?: string
  cost_name?: string
}

const props = withDefaults(defineProps<{
  /** 访问令牌；为空时不请求项目列表 */
  token?: string | null
  /** 项目列表 API 基础路径 */
  apiBase?: string
  /** 缓存有效期（毫秒），默认 30 分钟 */
  cacheTTL?: number
}>(), {
  apiBase: '/proxy/authv2',
  cacheTTL: 30 * 60 * 1000,
})

const { selectedProject, selectProject } = useProject()

const projects = ref<Project[]>([])
const pending = ref(false)
const modelValue = ref<string | null>(selectedProject.value?.code ?? null)

const CACHE_KEY = 'admin-project-list-cache'

interface ProjectCache {
  data: Project[]
  timestamp: number
  tokenHash: string
}

/** 简单 hash token 避免不同用户共用缓存 */
function hashToken(token: string): string {
  let h = 0
  for (let i = 0; i < token.length; i++) {
    h = (Math.imul(31, h) + token.charCodeAt(i)) | 0
  }
  return h.toString(36)
}

function readCache(): ProjectCache | null {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (!raw) return null
    return JSON.parse(raw) as ProjectCache
  } catch {
    return null
  }
}

function writeCache(data: Project[], token: string) {
  try {
    const cache: ProjectCache = { data, timestamp: Date.now(), tokenHash: hashToken(token) }
    localStorage.setItem(CACHE_KEY, JSON.stringify(cache))
  } catch { /* quota exceeded, ignore */ }
}

function isCacheFresh(cache: ProjectCache | null, token: string): cache is ProjectCache {
  if (!cache) return false
  if (cache.tokenHash !== hashToken(token)) return false
  return Date.now() - cache.timestamp < props.cacheTTL
}

// 仅客户端初始化
onMounted(async () => {
  if (!props.token) return
  pending.value = true

  // 1. 尝试从 localStorage 恢复缓存（立即展示，无闪烁）
  const cache = readCache()
  if (isCacheFresh(cache, props.token)) {
    projects.value = cache.data
    pending.value = false
    if (selectedProject.value?.code) {
      modelValue.value = selectedProject.value.code
    } else if (projects.value.length > 0) {
      onSelect(projects.value[0].code)
    }
    return
  }

  // 2. 缓存过期或不存在，从 API 拉取
  try {
    const res = await $fetch<{ items: Project[] }>(`${props.apiBase}/projects`, {
      headers: { 'X-Access-Token': props.token },
    })
    projects.value = [...(res.items ?? [])]
      .filter(p => p.code && p.code.length > 1)
      .sort((a, b) => a.code.localeCompare(b.code, 'zh'))

    writeCache(projects.value, props.token)

    // 恢复之前选中的项目；若无则默认选中第一个
    if (selectedProject.value?.code) {
      modelValue.value = selectedProject.value.code
    } else if (projects.value.length > 0) {
      onSelect(projects.value[0].code)
    }
  } finally {
    pending.value = false
  }
})

function onSelect(code: string) {
  const oldCode = selectedProject.value?.code
  modelValue.value = code
  const found = projects.value.find(p => p.code === code)
  selectProject(found
    ? {
        code: found.code,
        name: found.name,
        costCode: found.cost_code,
        costName: found.cost_name,
      }
    : null,
  )
  // 切换到不同项目时刷新页面，确保所有项目敏感数据重新加载
  // 首次选择（oldCode 为空）不刷新；选中同一项目也不刷新
  if (oldCode && oldCode !== code && import.meta.client) {
    window.location.reload()
  }
}

// 格式：code (name)，按代号首字母分组
const groupedOptions = computed(() => {
  const items = projects.value
  const groups = new Map<string, typeof items>()

  for (const p of items) {
    const key = p.code.charAt(0).toUpperCase()
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key)!.push(p)
  }

  return [...groups.entries()].map(([letter, groupItems]) => [
    { type: 'label' as const, label: letter },
    ...groupItems.map(p => ({
      label: `${p.code} (${p.name})`,
      value: p.code,
      code: p.code,
      name: p.name,
    })),
  ])
})
</script>

<template>
  <ClientOnly>
    <USelectMenu
      :items="groupedOptions"
      :model-value="modelValue"
      :loading="pending"
      placeholder="选择项目..."
      :search-input="{ placeholder: '搜索项目代号/名称...' }"
      variant="soft"
      value-key="value"
      label-key="label"
      :filter-fields="['code', 'name', 'label']"
      virtualize
      class="w-56 border-none"
      @update:model-value="onSelect"
    />
  </ClientOnly>
</template>
