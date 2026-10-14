<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold">日历</h1>
      <p class="text-muted-foreground mt-1">自研月历组件，支持事件管理、月份切换、新增/删除事件。</p>
    </div>

    <UCard>
      <!-- Toolbar -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <UButton icon="i-lucide-chevron-left" size="sm" variant="ghost" @click="prevMonth" />
          <h3 class="text-lg font-semibold min-w-36 text-center">{{ year }}年{{ month }}月</h3>
          <UButton icon="i-lucide-chevron-right" size="sm" variant="ghost" @click="nextMonth" />
        </div>
        <div class="flex items-center gap-2">
          <UButton size="sm" variant="soft" @click="goToday">今天</UButton>
          <UBadge size="sm" variant="soft">{{ events.length }} 个事件</UBadge>
        </div>
      </div>

      <!-- 月历网格 -->
      <div class="border border-default rounded-md overflow-hidden">
        <div class="grid grid-cols-7 bg-muted/50 border-b border-default">
          <div v-for="day in weekDays" :key="day" class="py-2 text-center text-xs font-medium text-muted-foreground">
            {{ day }}
          </div>
        </div>
        <div class="grid grid-cols-7">
          <div
            v-for="(cell, i) in cells"
            :key="i"
            class="min-h-24 p-1.5 border-b border-r border-default cursor-pointer hover:bg-muted/30 transition-colors"
            :class="[
              cell.otherMonth ? 'opacity-30' : '',
              cell.isToday ? 'bg-primary/5' : '',
              i % 7 === 0 ? 'border-l-0' : '',
              Math.floor(i / 7) === weeks - 1 ? 'border-b-0' : '',
              (i % 7 === 6) ? 'border-r-0' : '',
            ]"
            @click="cell.day && !cell.otherMonth && selectDate(cell.dateKey)"
          >
            <template v-if="cell.day">
              <div class="flex items-center justify-between">
                <span
                  class="size-6 flex items-center justify-center rounded-full text-xs"
                  :class="cell.isToday ? 'bg-primary text-white font-semibold' : ''"
                >{{ cell.day }}</span>
              </div>
              <div class="mt-1 space-y-0.5">
                <div v-for="evt in cell.events" :key="evt.id" class="flex items-center gap-1 group">
                  <span class="size-1.5 rounded-full shrink-0" :class="dotColor(evt.type)" />
                  <span class="text-xs truncate group-hover:opacity-70">{{ evt.title }}</span>
                  <button class="ml-auto shrink-0 opacity-0 group-hover:opacity-100 p-0.5 rounded hover:bg-muted" @click.stop="removeEvent(evt.id)">
                    <UIcon name="i-lucide-x" class="size-3" />
                  </button>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </UCard>

    <!-- 新增事件弹窗 -->
    <MFormModal
      v-model:open="showModal"
      title="添加事件"
      :schema="eventSchema"
      :config="eventFormConfig"
      :initial-state="eventFormState"
      @submit="addEvent"
    >
      <template #description>
        <p class="text-sm text-muted-foreground">日期：{{ selectedDate }}</p>
      </template>
    </MFormModal>
  </div>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { MFormConfig } from '../../../types/form'

definePageMeta({ layout: 'admin' })

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

interface CalendarEvent { id: number; date: string; title: string; type: 'work' | 'personal' | 'important' }

const now = new Date()
const year = ref(now.getFullYear())
const month = ref(now.getMonth() + 1)

let nextId = 0
const events = ref<CalendarEvent[]>([
  { id: ++nextId, date: '2026-07-15', title: '产品评审会', type: 'work' },
  { id: ++nextId, date: '2026-07-15', title: '迭代回顾', type: 'work' },
  { id: ++nextId, date: '2026-07-20', title: '前端技术分享', type: 'important' },
  { id: ++nextId, date: '2026-07-25', title: '团建活动', type: 'personal' },
  { id: ++nextId, date: '2026-08-03', title: 'Q3 规划会', type: 'work' },
])

const showModal = ref(false)
const selectedDate = ref('')

const eventSchema = z.object({
  title: z.string().min(1, '标题不能为空').meta({ title: '事件标题', required: true }),
  type: z.enum(['work', 'personal', 'important']).meta({ title: '类型' }),
})

const eventFormConfig: MFormConfig = {
  modal: { submitLabel: '添加', closeLabel: '取消' },
  fields: {
    type: { component: 'select', props: { options: [
      { label: '工作', value: 'work' },
      { label: '个人', value: 'personal' },
      { label: '重要', value: 'important' },
    ] } },
  },
}

const eventFormState = { title: '', type: 'work' as const }

function addEvent(data: { title: string; type: string }) {
  events.value.push({
    id: ++nextId,
    date: selectedDate.value,
    title: data.title,
    type: data.type as CalendarEvent['type'],
  })
  showModal.value = false
}

const eventsByDate = computed(() => {
  const map = new Map<string, CalendarEvent[]>()
  for (const e of events.value) {
    const arr = map.get(e.date) || []
    arr.push(e)
    map.set(e.date, arr)
  }
  return map
})

interface Cell { day?: number; dateKey: string; otherMonth?: boolean; isToday?: boolean; events: CalendarEvent[] }

const weeks = computed(() => Math.ceil(cells.value.length / 7))
const cells = computed<Cell[]>(() => {
  const y = year.value; const m = month.value
  const first = new Date(y, m - 1, 1)
  const startDow = first.getDay()
  const daysInMonth = new Date(y, m, 0).getDate()
  const prevDays = new Date(y, m - 1, 0).getDate()
  const result: Cell[] = []

  for (let i = 0; i < startDow; i++) {
    const d = prevDays - startDow + i + 1
    const dk = formatDate(y, m - 1, d)
    result.push({ day: d, dateKey: dk, otherMonth: true, events: eventsByDate.value.get(dk) || [] })
  }
  for (let d = 1; d <= daysInMonth; d++) {
    const dk = formatDate(y, m, d)
    result.push({ day: d, dateKey: dk, isToday: dk === formatDate(now.getFullYear(), now.getMonth() + 1, now.getDate()), events: eventsByDate.value.get(dk) || [] })
  }
  const remaining = 7 - (result.length % 7)
  if (remaining < 7) {
    for (let d = 1; d <= remaining; d++) {
      const dk = formatDate(y, m + 1, d)
      result.push({ day: d, dateKey: dk, otherMonth: true, events: eventsByDate.value.get(dk) || [] })
    }
  }
  return result
})

function formatDate(y: number, m: number, d: number): string {
  return `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
}

function dotColor(type: string) {
  return { work: 'bg-primary', personal: 'bg-success', important: 'bg-error' }[type] || 'bg-muted'
}

function selectDate(dateKey: string) {
  selectedDate.value = dateKey
  showModal.value = true
}

function removeEvent(id: number) {
  events.value = events.value.filter(e => e.id !== id)
}

function prevMonth() {
  if (month.value === 1) { year.value--; month.value = 12 } else month.value--
}

function nextMonth() {
  if (month.value === 12) { year.value++; month.value = 1 } else month.value++
}

function goToday() {
  year.value = now.getFullYear(); month.value = now.getMonth() + 1
}
</script>
