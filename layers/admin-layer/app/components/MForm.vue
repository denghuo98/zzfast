<script setup lang="ts" generic="T extends z.ZodObject<any>">
import type { FormError, FormErrorEvent, FormSubmitEvent, InferOutput } from '@nuxt/ui'
import type * as z from 'zod'
import type { MFormConfig, MFormState } from '../types/form'
import UButton from '@nuxt/ui/components/Button.vue'
import { computed, reactive, ref, useTemplateRef } from 'vue'
import MFormPrimitive from './MFormPrimitive.vue'

const props = withDefaults(defineProps<{
  schema: T
  /** 受控模式：通过 v-model:state 传入，表单值与外部双向同步 */
  state?: MFormState<T>
  /** 非受控模式：仅作初始值（向后兼容，未传 state 时生效） */
  initialState?: MFormState<T>
  config?: MFormConfig
  /** Override the form root class. Default: 'space-y-4' */
  formClass?: string
}>(), {
  initialState: () => ({}),
})

const emit = defineEmits<{
  'update:state': [value: MFormState<T>]
  submit: [data: InferOutput<T>]
  error: [payload: FormErrorEvent]
}>()

// 非受控模式下的内部状态（仅当未传入受控 state 时使用）
const internalState = reactive({ ...props.initialState }) as MFormState<T>
// 受控优先：传了 state 就直接引用外部对象（引用共享实现双向同步），否则回退内部副本
const state = computed<MFormState<T>>(() => props.state ?? internalState)
const formRef = useTemplateRef('form')
const loading = ref(false)

const isButtonDisabled = computed(() => !props.schema.safeParse(state.value).success)

const submitButton = computed(() => {
  if (props.config?.submit !== false)
    return props.config?.submit
  return undefined
})

const submitButtonProps = computed(() => ({
  'aria-disabled': isButtonDisabled.value,
  ...submitButton.value?.props,
}))

defineExpose({
  submit,
  clear: (path?: string | RegExp) => formRef.value?.clear(path as any),
})

async function onSubmit(event: FormSubmitEvent<InferOutput<T>>) {
  event.preventDefault()
  loading.value = true
  try {
    emit('submit', event.data)
  }
  finally {
    loading.value = false
  }
}

function submit() {
  formRef.value?.submit()
}
</script>

<template>
  <MFormPrimitive
    ref="form"
    :schema="schema"
    :state="state"
    :config="config"
    :form-class="formClass"
    @submit="onSubmit"
    @error="emit('error', $event)"
    @update:state="emit('update:state', $event)"
  >
    <template v-for="(_, name) in $slots" :key="name" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>

    <template #submit>
      <slot name="submit" :disabled="isButtonDisabled">
        <div v-if="config?.submit !== false">
          <template v-if="submitButton?.component">
            <component :is="submitButton.component" v-bind="submitButtonProps" />
          </template>
          <UButton
            v-else
            type="submit"
            label="Send"
            v-bind="submitButtonProps"
          />
        </div>
      </slot>
    </template>
  </MFormPrimitive>
</template>
