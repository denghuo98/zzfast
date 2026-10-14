<script setup lang="ts" generic="T extends z.ZodObject<any>">
import type { FormError, FormSubmitEvent, InferOutput, ModalProps } from '@nuxt/ui'
import type * as z from 'zod'
import type { MFormConfig, MFormState } from '../types/form'
import UButton from '@nuxt/ui/components/Button.vue'
import UModal from '@nuxt/ui/components/Modal.vue'
import { computed, reactive, ref, toRaw, useTemplateRef } from 'vue'
import MFormPrimitive from './MFormPrimitive.vue'

export interface MFormModalProps<TSchema extends z.ZodObject<any>> {
  open?: boolean
  schema: TSchema
  initialState?: MFormState<TSchema>
  title?: string
  description?: string
  config?: MFormConfig
  modalProps?: Partial<ModalProps>
}

const props = withDefaults(defineProps<MFormModalProps<T>>(), {
  open: false,
  initialState: () => ({}),
  modalProps: () => ({}),
})

const emit = defineEmits<{
  'submit': [data: InferOutput<T>]
  'update:open': [value: boolean]
}>()

function mergeDeep(...sources: (Record<string, any> | undefined)[]): Record<string, any> {
  const result: Record<string, any> = {}
  for (const source of sources) {
    if (!source) continue
    for (const key of Object.keys(source)) {
      const val = source[key]
      if (val !== null && typeof val === 'object' && !Array.isArray(val) && typeof result[key] === 'object' && result[key] !== null && !Array.isArray(result[key])) {
        result[key] = mergeDeep(result[key], val)
      }
      else {
        result[key] = val
      }
    }
  }
  return result
}

const defaults: Partial<MFormConfig> = {
  modal: { submitLabel: 'Submit', closeLabel: 'Cancel' },
}

const appConfig = computed<MFormConfig>(() => {
  return mergeDeep(defaults, props.config ?? {}) as MFormConfig
})

const state = reactive<MFormState<T>>({ ...props.initialState })
const formRef = useTemplateRef('form')

// 始终允许点击提交，由 submit 事件中的 schema 校验拦截
const isButtonDisabled = ref(false)

const open = computed({
  get: () => props.open,
  set: value => emit('update:open', value),
})

defineExpose({
  submit: () => formRef.value?.submit(),
  form: formRef,
})

async function onSubmit(event: FormSubmitEvent<InferOutput<T>>) {
  event.preventDefault()
  emit('submit', toRaw(event.data) as InferOutput<T>)
}

function close() {
  open.value = false
}
</script>

<template>
  <UModal
    v-model:open="open"
    :title="title"
    :description="description"
    :ui="{ footer: 'justify-end' }"
    v-bind="modalProps"
  >
    <template #body>
      <MFormPrimitive
        ref="form"
        :schema="schema"
        :state="state as MFormState<T>"
        :config="config"
        @submit="onSubmit"
      >
        <template v-for="(_, name) in $slots" :key="name" #[name]="slotData">
          <slot :name="name" v-bind="slotData" />
        </template>
      </MFormPrimitive>
    </template>

    <template #footer>
      <slot name="footer" :disabled="isButtonDisabled" :submit="() => formRef?.submit()" :close="close">
        <UButton
          :label="appConfig.modal?.closeLabel"
          color="neutral"
          variant="outline"
          @click="close"
        />
        <UButton
          :label="appConfig.modal?.submitLabel"
          :disabled="isButtonDisabled"
          @click="formRef?.submit()"
        />
      </slot>
    </template>
  </UModal>
</template>
