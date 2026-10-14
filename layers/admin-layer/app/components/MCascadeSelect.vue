<script setup lang="ts" generic="T extends string | number">
import UPopover from '@nuxt/ui/components/Popover.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import { computed, ref } from 'vue'

export interface CascadeOption {
  label: string
  value: string | number
  icon?: string
  disabled?: boolean
  children?: CascadeOption[]
}

const props = withDefaults(defineProps<{
  options: CascadeOption[]
  modelValue?: T | null
  placeholder?: string
  disabled?: boolean
  clearable?: boolean
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  /** Visual variant matching Nuxt UI Select */
  variant?: 'outline' | 'soft' | 'subtle' | 'ghost' | 'none'
  /** Color variant */
  color?: 'primary' | 'secondary' | 'success' | 'info' | 'warning' | 'error' | 'neutral'
  /** Optional header text in dropdown */
  header?: string
  /** Max height of the dropdown panel. Default: 15rem */
  maxHeight?: string
}>(), {
  modelValue: null,
  placeholder: '请选择',
  disabled: false,
  clearable: false,
  size: 'md',
  variant: 'outline',
  color: 'primary',
  maxHeight: '15rem',
})

const emit = defineEmits<{
  'update:modelValue': [value: T | null]
}>()

// ========== Nuxt UI Select–compatible variant classes ==========
const variantClasses: Record<string, string> = {
  outline: 'text-highlighted bg-default ring ring-inset ring-accented hover:bg-elevated disabled:bg-default',
  soft: 'text-highlighted bg-elevated/50 hover:bg-elevated focus:bg-elevated disabled:bg-elevated/50',
  subtle: 'text-highlighted bg-elevated ring ring-inset ring-accented hover:bg-accented/75 disabled:bg-elevated',
  ghost: 'text-highlighted bg-transparent hover:bg-elevated focus:bg-elevated disabled:bg-transparent',
  none: 'text-highlighted bg-transparent',
}

const focusRingColor: Record<string, string> = {
  primary: 'focus:ring-primary',
  secondary: 'focus:ring-secondary',
  success: 'focus:ring-success',
  info: 'focus:ring-info',
  warning: 'focus:ring-warning',
  error: 'focus:ring-error',
  neutral: 'focus:ring-inverted',
}

const sizeClasses: Record<string, string> = {
  xs: 'px-2 py-1 text-xs gap-1',
  sm: 'px-2.5 py-1.5 text-xs gap-1.5',
  md: 'px-2.5 py-1.5 text-sm gap-1.5',
  lg: 'px-3 py-2 text-sm gap-2',
  xl: 'px-3 py-2 text-base gap-2',
}

const itemSizeClasses: Record<string, string> = {
  xs: 'text-xs py-1 px-2 gap-1',
  sm: 'text-xs py-1.5 px-2 gap-1.5',
  md: 'text-sm py-1.5 px-2.5 gap-1.5',
  lg: 'text-sm py-2 px-2.5 gap-2',
  xl: 'text-base py-2 px-3 gap-2',
}

const chevronSizeClasses: Record<string, string> = {
  xs: 'size-4', sm: 'size-4', md: 'size-5', lg: 'size-5', xl: 'size-6',
}

const clearSizeClasses: Record<string, string> = {
  xs: 'size-3', sm: 'size-3.5', md: 'size-4', lg: 'size-4', xl: 'size-5',
}

// Show focus ring only for outline/subtle variants
const showFocusRing = computed(() => props.variant === 'outline' || props.variant === 'subtle')

const open = ref(false)
const selectedPath = ref<CascadeOption[]>([])

function buildPath(options: CascadeOption[], target: T | null): CascadeOption[] {
  if (target === null || target === undefined) return []
  for (const opt of options) {
    if (opt.value === target) return [opt]
    if (opt.children) {
      const childPath = buildPath(opt.children, target)
      if (childPath.length) return [opt, ...childPath]
    }
  }
  return []
}

function syncFromModel() {
  selectedPath.value = buildPath(props.options, props.modelValue as T | null)
}
syncFromModel()

const displayText = computed(() => {
  if (!selectedPath.value.length) return ''
  return selectedPath.value.map(o => o.label).join(' / ')
})

const levels = computed(() => {
  const result: CascadeOption[][] = [props.options]
  for (let i = 0; i < selectedPath.value.length; i++) {
    const current = selectedPath.value[i]
    if (current?.children?.length) {
      result.push(current.children)
    }
  }
  return result
})

function isSelected(option: CascadeOption, levelIndex: number): boolean {
  return selectedPath.value[levelIndex]?.value === option.value
}

function selectOption(option: CascadeOption, levelIndex: number) {
  if (option.disabled) return
  selectedPath.value = [...selectedPath.value.slice(0, levelIndex), option]
  if (!option.children?.length) {
    emit('update:modelValue', option.value as T)
    open.value = false
  }
}

function clearSelection() {
  selectedPath.value = []
  emit('update:modelValue', null)
  open.value = false
}

function onOpenChange(val: boolean) {
  open.value = val
  if (val) {
    syncFromModel()
  }
  else {
    // 关闭时如果没有选中叶子节点，清空中间选择
    const last = selectedPath.value[selectedPath.value.length - 1]
    if (last?.children?.length) {
      selectedPath.value = []
    }
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') open.value = false
}
</script>

<template>
  <UPopover v-model:open="open" :content="{ side: 'bottom', sideOffset: 4, align: 'start', collisionPadding: 8 }" @update:open="onOpenChange">
    <!-- Trigger — matches USelect styling exactly -->
    <button
      type="button"
      :class="[
        'relative group rounded-md inline-flex items-center focus:outline-none disabled:cursor-not-allowed disabled:opacity-75 transition-colors w-full cursor-pointer text-left',
        variantClasses[variant],
        sizeClasses[size],
        showFocusRing ? ['focus:ring-2 focus:ring-inset', focusRingColor[color]] : '',
      ]"
      role="combobox"
      :aria-expanded="open"
      :aria-disabled="disabled"
      :disabled="disabled"
      @keydown="onKeydown"
    >
      <span v-if="displayText" class="flex-1 truncate text-highlighted">
        {{ displayText }}
      </span>
      <span v-else class="flex-1 truncate text-dimmed">
        {{ placeholder }}
      </span>

      <!-- Clear -->
      <button
        v-if="clearable && displayText && !disabled"
        type="button"
        class="shrink-0 inline-flex items-center justify-center rounded hover:bg-elevated/50"
        aria-label="清除"
        @click.stop="clearSelection"
      >
        <UIcon name="i-lucide-x" :class="[clearSizeClasses[size], 'text-dimmed shrink-0']" />
      </button>

      <!-- Chevron -->
      <UIcon
        :name="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
        :class="[chevronSizeClasses[size], 'text-dimmed shrink-0']"
      />
    </button>

    <!-- Dropdown Panel — matches USelect content + custom columns -->
    <template #content>
      <div class="bg-default shadow-lg rounded-md ring ring-default overflow-hidden flex" :style="{ maxHeight }" @keydown="onKeydown">
        <!-- Columns -->
        <div class="flex divide-x divide-default min-h-28" style="flex: 1; min-height: 0;">
          <div
            v-for="(levelOptions, levelIndex) in levels"
            :key="levelIndex"
            class="flex flex-col min-w-36 max-w-52 py-1 overflow-y-auto"
            style="flex: 1; min-height: 0;"
          >
            <button
              v-for="option in levelOptions"
              :key="option.value"
              :class="[
                'relative w-full flex items-center outline-none transition-colors text-left',
                'data-disabled:cursor-not-allowed data-disabled:opacity-75',
                itemSizeClasses[size],
                isSelected(option, levelIndex)
                  ? 'text-highlighted font-medium before:absolute before:z-[-1] before:inset-px before:rounded-md before:bg-elevated/50'
                  : 'text-default hover:text-highlighted hover:before:absolute hover:before:z-[-1] hover:before:inset-px hover:before:rounded-md hover:before:bg-elevated/50',
                option.disabled ? 'opacity-40 cursor-not-allowed' : 'cursor-pointer',
              ]"
              :disabled="option.disabled"
              @click="selectOption(option, levelIndex)"
            >
              <UIcon v-if="option.icon" :name="option.icon" class="size-4 shrink-0 text-dimmed" />
              <span class="flex-1 truncate">{{ option.label }}</span>
              <UIcon
                v-if="option.children?.length"
                name="i-lucide-chevron-right"
                class="size-3.5 shrink-0 text-dimmed"
              />
            </button>

            <div v-if="!levelOptions.length" class="flex items-center justify-center flex-1 p-4">
              <span class="text-xs text-dimmed">无可用选项</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </UPopover>
</template>
