<!-- 图标浏览器 — 网格布局，可搜索 Lucide 图标集，点击复制名称 -->
<template>
  <ClientOnly>
    <div class="space-y-3">
      <UInput v-model="query" placeholder="搜索图标（如 home、arrow、settings）..." icon="i-lucide-search" clearable size="sm" />
      <p class="text-xs text-muted-foreground">
        共 {{ filtered.length }} 个图标{{ query ? '（搜索结果）' : '' }}
        <span v-if="copiedName" class="ml-2 text-primary">✓ 已复制 {{ copiedName }}</span>
      </p>
      <div class="h-[420px] overflow-auto border border-default rounded-md p-3">
        <div class="grid grid-cols-[repeat(auto-fill,minmax(88px,1fr))] gap-1.5">
          <div
            v-for="name in filtered" :key="name"
            class="flex flex-col items-center gap-1 p-2 rounded-md hover:bg-muted cursor-pointer group text-center"
            @click="copy(name)"
          >
            <UIcon :name="'i-lucide-' + name" class="size-5 shrink-0 text-muted group-hover:text-foreground" />
            <span class="text-[10px] leading-tight text-muted-foreground group-hover:text-foreground truncate w-full">{{ name }}</span>
          </div>
        </div>
      </div>
    </div>
    <template #fallback>
      <div class="flex items-center justify-center h-64">
        <UIcon name="i-lucide-loader-2" class="size-6 animate-spin text-muted" />
      </div>
    </template>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

defineOptions({ name: 'MIconBrowser' })

const query = ref('')
const icons = ref<string[]>([])
const copiedName = ref('')

const filtered = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return icons.value
  return icons.value.filter(n => n.toLowerCase().includes(q))
})

function copy(name: string) {
  const text = `i-lucide-${name}`
  navigator.clipboard.writeText(text).catch(() => {})
  copiedName.value = text
  setTimeout(() => { copiedName.value = '' }, 1500)
}

onMounted(async () => {
  try {
    const mod = await import('@iconify-json/lucide/icons.json')
    icons.value = Object.keys((mod as any).default?.icons ?? (mod as any).icons ?? {})
  } catch {
    // 图标集未安装或无法加载
  }
})
</script>
