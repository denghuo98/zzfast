<script setup lang="ts">
import * as z from 'zod'
import type { MFormConfig } from '../../../types/form'
import type { CascadeOption } from '#components/MCascadeSelect.vue'

definePageMeta({ layout: 'admin' })

// ==================== Schema 定义 ====================

const schema = z.object({
  name: z.string()
    .min(1, '姓名不能为空')
    .meta({
      title: '姓名',
      description: '请输入您的真实姓名',
      required: true,
    }),
  email: z.email('邮箱格式不正确')
    .meta({
      title: '邮箱',
      description: '我们不会泄露您的邮箱',
      hoverTip: '请使用常用邮箱，用于接收通知和找回密码',
    }),
  role: z.enum(['管理员', '编辑', '普通用户'])
    .meta({
      title: '角色',
      hint: '请选择您的角色',
    }),
  bio: z.string()
    .optional()
    .meta({
      title: '个人简介',
      description: '简单介绍一下自己',
    }),
  agree: z.boolean()
    .default(false)
    .meta({
      title: '同意服务条款',
    }),
})

const schema2 = z.object({
  username: z.string()
    .min(3, '至少 3 个字符')
    .meta({ title: '用户名', required: true }),
  password: z.string()
    .min(6, '至少 6 个字符')
    .meta({ title: '密码', required: true }),
  age: z.number()
    .min(18, '必须大于 18 岁')
    .optional()
    .meta({ title: '年龄' }),
})

type SchemaType = z.infer<typeof schema>

// ==================== 表单 1：基础 MForm ====================

const formData = reactive<SchemaType>({
  name: '',
  email: '',
  role: '普通用户',
  bio: '',
  agree: false,
})

const submitted = ref(false)
const submittedData = ref<any>(null)
const orientation = ref<'vertical' | 'horizontal'>('vertical')

const formConfig = computed<MFormConfig>(() => ({
  theme: {
    orientation: orientation.value,
  },
  submit: {
    props: { label: '提交' },
  },
}))

function onSubmit(data: z.infer<typeof schema>) {
  submitted.value = true
  submittedData.value = data
}

function resetForm() {
  formData.name = ''
  formData.email = ''
  formData.role = '普通用户'
  formData.bio = ''
  formData.agree = false
  submitted.value = false
  submittedData.value = null
}

// ==================== 表单 2：Modal ====================
const modalOpen = ref(false)
const modalData = reactive({
  username: '',
  password: '',
  age: undefined as number | undefined,
})

function onModalSubmit(data: z.infer<typeof schema2>) {
  useToast().add({ title: '提交成功', description: JSON.stringify(data, null, 2) })
  modalOpen.value = false
}

// ==================== 级联选择 ====================
const cascadeOptions: CascadeOption[] = [
  {
    label: '电子产品',
    value: 'electronics',
    icon: 'i-lucide-laptop',
    children: [
      {
        label: '手机',
        value: 'phones',
        icon: 'i-lucide-smartphone',
        children: [
          { label: 'iPhone 16', value: 'iphone16' },
          { label: 'iPhone 16 Pro', value: 'iphone16pro' },
          { label: '华为 Mate 70', value: 'mate70' },
          { label: '小米 15', value: 'mi15' },
        ],
      },
      {
        label: '电脑',
        value: 'computers',
        icon: 'i-lucide-monitor',
        children: [
          { label: 'MacBook Pro', value: 'macbook-pro' },
          { label: 'MacBook Air', value: 'macbook-air' },
          { label: 'ThinkPad X1', value: 'thinkpad-x1' },
        ],
      },
      {
        label: '平板',
        value: 'tablets',
        icon: 'i-lucide-tablet',
        children: [
          { label: 'iPad Pro', value: 'ipad-pro' },
          { label: 'iPad Air', value: 'ipad-air' },
          { label: '华为 MatePad', value: 'matepad' },
        ],
      },
    ],
  },
  {
    label: '服装',
    value: 'clothing',
    icon: 'i-lucide-shirt',
    children: [
      {
        label: '男装',
        value: 'mens',
        children: [
          { label: 'T 恤', value: 't-shirt' },
          { label: '衬衫', value: 'shirt' },
          { label: '外套', value: 'jacket' },
        ],
      },
      {
        label: '女装',
        value: 'womens',
        children: [
          { label: '连衣裙', value: 'dress' },
          { label: '半身裙', value: 'skirt' },
          { label: '针织衫', value: 'knit' },
        ],
      },
    ],
  },
  {
    label: '图书',
    value: 'books',
    icon: 'i-lucide-book-open',
    children: [
      { label: '技术', value: 'tech', children: [
        { label: '前端开发', value: 'frontend' },
        { label: '后端开发', value: 'backend' },
        { label: '数据分析', value: 'data-analysis' },
        { label: '人工智能', value: 'ai' },
        { label: '机器学习', value: 'ml' },
        { label: '区块链', value: 'blockchain' },
        { label: '云计算', value: 'cloud' },
        { label: '网络安全', value: 'security' },
        { label: '嵌入式', value: 'embedded' },
        { label: '游戏开发', value: 'gamedev' },
        { label: '移动开发', value: 'mobile' },
        { label: '测试工程', value: 'testing' },
        { label: '运维工程', value: 'devops' },
        { label: '架构设计', value: 'architecture' },
        { label: '产品设计', value: 'product-design' },
        { label: 'UI 设计', value: 'ui-design' },
        { label: 'UX 研究', value: 'ux-research' },
      ] },
      { label: '文学', value: 'literature', children: [{ label: '小说', value: 'novel' }, { label: '散文', value: 'prose' }] },
    ],
  },
  { label: '其他', value: 'other' },
]

const cascadeValue = ref<string | number | null>(null)
const cascadeValue2 = ref<string | number | null>(null)
const cascadeValue3 = ref<string | number | null>(null)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">
        MForm — Schema 驱动表单
      </h1>
      <p class="text-muted-foreground mt-1">
        基于 Zod Schema 自动生成 Nuxt UI 表单，支持垂直/水平布局切换。
      </p>
    </div>

    <!-- 表单 1: MForm -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="font-semibold">
            MForm 基础用法
          </h3>
          <div class="flex items-center gap-2">
            <span class="text-xs text-muted-foreground">布局方向：</span>
            <USelect
              v-model="orientation"
              :items="['vertical', 'horizontal']"
              size="xs"
              class="w-28"
            />
          </div>
        </div>
      </template>

      <div class="max-w-lg">
        <MForm
          :schema="schema"
          :initial-state="formData"
          :config="formConfig"
          @submit="onSubmit"
        >
          <!-- 自定义 bio 字段用 textarea -->
          <template #bio="{ field }">
            <UTextarea
              :model-value="formData.bio"
              placeholder="简单介绍一下自己..."
              :rows="3"
              class="w-full"
              @update:model-value="formData.bio = $event"
            />
          </template>
        </MForm>

        <div class="flex gap-3 mt-4">
          <UButton variant="outline" color="neutral" @click="resetForm">
            重置
          </UButton>
        </div>

        <div v-if="submitted" class="mt-6 p-4 bg-muted rounded-lg">
          <p class="text-sm font-semibold mb-2">
            提交数据：
          </p>
          <pre class="text-xs">{{ JSON.stringify(submittedData, null, 2) }}</pre>
        </div>
      </div>
    </UCard>

    <!-- 表单 2: MFormModal -->
    <UCard>
      <template #header>
        <h3 class="font-semibold">
          MFormModal 弹窗表单
        </h3>
      </template>

      <div class="flex gap-3">
        <UButton label="打开弹窗表单" @click="modalOpen = true" />
      </div>

      <MFormModal
        v-model:open="modalOpen"
        title="编辑用户信息"
        description="请填写您的个人资料"
        :schema="schema2"
        :initial-state="modalData"
        :config="{ modal: { submitLabel: '保存', closeLabel: '取消' }, theme: { orientation: 'vertical' } }"
        @submit="onModalSubmit"
      />
    </UCard>

    <!-- MCascadeSelect 级联选择 -->
    <UCard>
      <template #header>
        <h3 class="font-semibold">
          MCascadeSelect 级联选择
        </h3>
      </template>

      <div class="max-w-sm space-y-4">
        <div>
          <span class="text-xs text-muted mb-1 block">商品分类（outline · 默认）</span>
          <MCascadeSelect
            v-model="cascadeValue"
            :options="cascadeOptions"
            placeholder="请选择分类..."
            clearable
          />
        </div>

        <div>
          <span class="text-xs text-muted mb-1 block">商品分类（subtle 变体）</span>
          <MCascadeSelect
            v-model="cascadeValue2"
            :options="cascadeOptions"
            placeholder="请选择分类..."
            variant="subtle"
            clearable
          />
        </div>

        <div>
          <span class="text-xs text-muted mb-1 block">商品分类（ghost 变体）</span>
          <MCascadeSelect
            v-model="cascadeValue3"
            :options="cascadeOptions"
            placeholder="请选择分类..."
            variant="ghost"
            clearable
          />
        </div>

        <div v-if="cascadeValue" class="p-3 bg-muted rounded text-xs">
          已选值：<code class="font-mono bg-primary/10 px-1 rounded">{{ cascadeValue }}</code>
        </div>
      </div>
    </UCard>
  </div>
</template>
