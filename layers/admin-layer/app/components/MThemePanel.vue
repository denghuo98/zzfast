<!-- 主题设置面板：主色 / 圆角 / 导航风格 / 内容宽度 -->
<script setup lang="ts">
const appConfig = useAppConfig()
const toast = useToast()

const open = ref(false)

// 预设主题色
const primaryColors = [
  { label: 'Green', value: 'green', color: 'bg-green-500' },
  { label: 'Blue', value: 'blue', color: 'bg-blue-500' },
  { label: 'Violet', value: 'violet', color: 'bg-violet-500' },
  { label: 'Rose', value: 'rose', color: 'bg-rose-500' },
  { label: 'Amber', value: 'amber', color: 'bg-amber-500' },
  { label: 'Teal', value: 'teal', color: 'bg-teal-500' },
  { label: 'Slate', value: 'slate', color: 'bg-slate-500' },
  { label: 'Zinc', value: 'zinc', color: 'bg-zinc-500' },
]

const currentPrimary = ref(appConfig.ui.colors?.primary || 'green')
const currentRadius = ref(appConfig.ui.radius ?? 0.5)

const radiusOptions = [
  { label: '无圆角', value: 0 },
  { label: '小', value: 0.25 },
  { label: '中', value: 0.5 },
  { label: '大', value: 0.75 },
  { label: '全圆角', value: 1 },
]

function setPrimary(color: string) {
  currentPrimary.value = color
  appConfig.ui.colors = { ...appConfig.ui.colors, primary: color as any }
}

function setRadius(value: number) {
  currentRadius.value = value
  appConfig.ui.radius = value
}

function reset() {
  setPrimary('green')
  setRadius(0.5)
  toast.add({ title: '已重置', description: '主题设置已恢复默认' })
}
</script>

<template>
  <div>
    <UTooltip text="主题设置">
      <UButton
        icon="i-lucide-palette"
        color="neutral"
        variant="ghost"
        size="sm"
        class="text-muted"
        @click="open = true"
      />
    </UTooltip>

    <USlideover v-model:open="open" title="主题设置" :ui="{ content: 'w-80' }">
      <template #body>
        <div class="space-y-6">
          <!-- 主题色 -->
          <div>
            <h4 class="text-sm font-medium mb-3">主题色</h4>
            <div class="grid grid-cols-4 gap-2">
              <button
                v-for="c in primaryColors"
                :key="c.value"
                class="flex flex-col items-center gap-1 p-2 rounded-md border-2 transition-colors"
                :class="currentPrimary === c.value ? 'border-primary' : 'border-transparent hover:border-default'"
                @click="setPrimary(c.value)"
              >
                <span class="size-6 rounded-full" :class="c.color" />
                <span class="text-xs text-muted-foreground">{{ c.label }}</span>
              </button>
            </div>
          </div>

          <!-- 圆角 -->
          <div>
            <h4 class="text-sm font-medium mb-3">圆角半径</h4>
            <div class="flex gap-2">
              <UButton
                v-for="r in radiusOptions"
                :key="r.value"
                size="sm"
                :variant="currentRadius === r.value ? 'solid' : 'soft'"
                @click="setRadius(r.value)"
              >
                {{ r.label }}
              </UButton>
            </div>
          </div>

          <!-- 重置 -->
          <div>
            <UButton variant="soft" color="error" size="sm" block @click="reset">
              重置默认设置
            </UButton>
          </div>
        </div>
      </template>
    </USlideover>
  </div>
</template>
