<script setup lang="ts">
import * as z from 'zod'

definePageMeta({ layout: 'admin' })

// ==================== Schema 定义 ====================

const searchSchema = z.object({
  keyword: z.string()
    .optional()
    .meta({
      title: '关键词',
      hoverTip: '输入关键词后按回车或点击搜索按钮，支持对标题进行模糊匹配',
    }),
  category: z.enum(['全部', '文档', '教程', 'API', '博客'])
    .optional()
    .meta({ title: '分类' }),
  status: z.enum(['全部', '已发布', '草稿', '已归档'])
    .optional()
    .meta({ title: '状态' }),
  dateFrom: z.string()
    .optional()
    .meta({ title: '开始日期' }),
  dateTo: z.string()
    .optional()
    .meta({ title: '结束日期' }),
})

type SearchFormType = z.infer<typeof searchSchema>

// ==================== 搜索表单 ====================

const searchState = reactive<SearchFormType>({
  keyword: '',
  category: '全部',
  status: '全部',
  dateFrom: '',
  dateTo: '',
})



const searched = ref(false)
const results = ref<any[]>([])

const mockData = [
  { title: '快速入门指南', category: '文档', status: '已发布', icon: 'i-lucide-file-text' },
  { title: 'API 参考文档 v2.0', category: 'API', status: '已发布', icon: 'i-lucide-code' },
  { title: '高级实战教程', category: '教程', status: '草稿', icon: 'i-lucide-graduation-cap' },
  { title: '2026 年 6 月发布说明', category: '博客', status: '已发布', icon: 'i-lucide-newspaper' },
]

function doSearch(data: z.infer<typeof searchSchema>) {
  let filtered = mockData
  if (data.keyword) {
    const q = data.keyword.toLowerCase()
    filtered = filtered.filter(item => item.title.toLowerCase().includes(q))
  }
  if (data.category && data.category !== '全部') {
    filtered = filtered.filter(item => item.category === data.category)
  }
  if (data.status && data.status !== '全部') {
    filtered = filtered.filter(item => item.status === data.status)
  }
  results.value = filtered
  searched.value = true
}

function resetSearch() {
  searchState.keyword = ''
  searchState.category = '全部'
  searchState.status = '全部'
  searchState.dateFrom = ''
  searchState.dateTo = ''
  searched.value = false
  results.value = []
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">
        搜索表单
      </h1>
      <p class="text-muted-foreground mt-1">
        基于 MForm 的高级搜索表单，支持多种筛选条件。
      </p>
    </div>

    <!-- 搜索表单 -->
    <MFormSearch
      :schema="searchSchema"
      :state="searchState"
      form-class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 items-end"
      @search="doSearch"
      @reset="resetSearch"
    />

    <!-- 搜索结果 -->
    <UCard v-if="searched">
      <template #header>
        <h3 class="font-semibold">
          搜索结果（{{ results.length }} 条）
        </h3>
      </template>
      <div v-if="results.length" class="space-y-2">
        <div v-for="(item, i) in results" :key="i" class="flex items-center gap-3 p-2 rounded hover:bg-muted">
          <UIcon :name="item.icon" class="size-4 text-muted-foreground" />
          <div>
            <p class="text-sm font-medium">
              {{ item.title }}
            </p>
            <p class="text-xs text-muted-foreground">
              {{ item.category }} · {{ item.status }}
            </p>
          </div>
        </div>
      </div>
      <p v-else class="text-sm text-muted-foreground">
        未找到匹配的结果。
      </p>
    </UCard>
  </div>
</template>
