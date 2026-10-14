<script setup lang="ts" generic="T extends z.ZodObject<any>">
import type { FormError, FormErrorEvent, FormSubmitEvent, InferOutput } from '@nuxt/ui'
import type * as z from 'zod'
import type { ComponentDefinition } from '../utils/form/components_map'
import type { MFormConfig, MFormState } from '../types/form'
import UForm from '@nuxt/ui/components/Form.vue'
import UFormField from '@nuxt/ui/components/FormField.vue'
import { splitByCase, upperFirst } from 'scule'
import { computed, useSlots, useTemplateRef } from 'vue'
import { COMPONENTS_MAP, mapZodTypeToComponent } from '../utils/form/components_map'
import { useMetaProcessor } from '../utils/form/useMetaProcessor'

/** Simple deep merge — later keys take priority */
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

type FormInstance = InstanceType<typeof UForm>
type FieldKey = keyof MFormState<T> & string

interface FieldDefinition {
  key: FieldKey
  formField: ReturnType<typeof parseMeta> & {
    name: FieldKey
    slots: string[]
    orientation?: 'vertical' | 'horizontal'
  }
  component: ComponentDefinition['component']
  props: NonNullable<ComponentDefinition['componentProps']>
}

const props = defineProps<{
  schema: T
  state: MFormState<T>
  config?: MFormConfig
  id?: string | number
  validate?: (state: MFormState<T>) => FormError[] | Promise<FormError[]>
  validateOn?: ('input' | 'change' | 'blur')[]
  disabled?: boolean
  validateOnInputDelay?: number
  transform?: boolean
  loadingAuto?: boolean
  /** Override the form root class. Default: 'space-y-4' */
  formClass?: string
}>()

const emits = defineEmits<{
  submit: [payload: FormSubmitEvent<InferOutput<T>>]
  error: [payload: FormErrorEvent]
  'update:state': [value: MFormState<T>]
}>()

const slots = useSlots()
const formRef = useTemplateRef<FormInstance>('form')
const shape = (props.schema as z.ZodObject<any>).shape

const defaults: Partial<MFormConfig> = {
  components: COMPONENTS_MAP,
  theme: {
    orientation: 'vertical',
    wFull: true,
  },
}

const appConfig = computed<MFormConfig>(() => {
  return mergeDeep(defaults, props.config ?? {}) as MFormConfig
})

const { processFieldMeta } = useMetaProcessor(appConfig)

const fields = computed<FieldDefinition[]>(() => {
  const resolvedFields: FieldDefinition[] = []

  for (const key of Object.keys(shape) as FieldKey[]) {
    const zodType = shape[key]
    const result = mapZodTypeToComponent(key, zodType, appConfig.value, props.state)
    if (!result)
      continue

    const rawMeta = typeof zodType.meta === 'function' ? zodType.meta() || {} : {}
    const meta = processFieldMeta(rawMeta, key)

    // Resolve orientation: field-level > global config > default ('vertical')
    const fieldOrientation = meta?.theme?.orientation ?? appConfig.value.theme?.orientation ?? 'vertical'

    // wFull: field-level > global config; auto-adjust based on orientation
    const resolvedWFull = meta?.theme?.wFull ?? appConfig.value.theme?.wFull
    const shouldApplyWFull = resolvedWFull !== undefined
      ? resolvedWFull
      : appConfig.value.theme?.orientation !== 'horizontal'

    const defaultProps: Record<string, any> = {}
    if (shouldApplyWFull) {
      defaultProps.class = 'w-full'
    }

    // Merge: default → component result → meta input props
    const mergedProps = mergeDeep(defaultProps, result.componentProps ?? {}, meta?.input?.props ?? {})

    // Concatenate class strings instead of replacing
    if (shouldApplyWFull || result.componentProps?.class || meta?.input?.props?.class) {
      const classes: string[] = []
      if (shouldApplyWFull) classes.push('w-full')
      if (result.componentProps?.class) classes.push(result.componentProps.class)
      if (meta?.input?.props?.class) classes.push(meta.input.props.class)
      mergedProps.class = classes.join(' ')
    }

    resolvedFields.push({
      key,
      formField: {
        name: key,
        slots: findSlots(key),
        orientation: fieldOrientation,
        ...parseMeta(meta, key, appConfig.value),
      },
      component: meta?.input?.component ?? result.component,
      props: mergedProps,
    })
  }

  return resolvedFields
})

function findSlots(key: string): string[] {
  return Object.keys(slots)
    .filter(name => name.startsWith(`${key}-`))
    .map(name => name.slice(key.length + 1))
}

function parseMeta(meta: any, key: string, config: MFormConfig) {
  const enableDefaultTitles = config?.theme?.enableDefaultTitles !== false
  const label = typeof meta.title === 'string'
    ? meta.title
    : meta.title === undefined && enableDefaultTitles
      ? upperFirst(splitByCase(key).join(' ').toLowerCase())
      : undefined
  return {
    label,
    required: meta.required,
    description: meta.description,
    hint: meta.hint,
    help: meta.help,
    hoverTip: meta.hoverTip,
  }
}

function getFieldValue(key: string) {
  return props.state[key as keyof MFormState<T>]
}

function updateFieldValue(key: string, value: unknown) {
  ;(props.state as Record<string, unknown>)[key] = value
  emits('update:state', props.state)
}

const formProps = computed(() => ({
  schema: props.schema,
  state: props.state,
  id: props.id,
  validate: props.validate,
  validateOn: props.validateOn,
  disabled: props.disabled,
  validateOnInputDelay: props.validateOnInputDelay,
  transform: props.transform,
  loadingAuto: props.loadingAuto,
}))

defineExpose({
  form: formRef,
  submit: () => formRef.value?.submit(),
  validate: (opts?: any) => formRef.value?.validate(opts),
  clear: (path?: string | RegExp) => formRef.value?.clear(path as any),
  getErrors: (path?: string | RegExp) => formRef.value?.getErrors(path as any),
  setErrors: (errors: FormError[], name?: string | RegExp) => formRef.value?.setErrors(errors, name as any),
  errors: computed(() => formRef.value?.errors),
  disabled: computed(() => formRef.value?.disabled),
  dirty: computed(() => formRef.value?.dirty),
  dirtyFields: computed(() => formRef.value?.dirtyFields),
  touchedFields: computed(() => formRef.value?.touchedFields),
  blurredFields: computed(() => formRef.value?.blurredFields),
})
</script>

<template>
  <UForm
    :id="formProps.id"
    ref="form"
    :schema="formProps.schema"
    :state="formProps.state"
    :validate="formProps.validate"
    :validate-on="formProps.validateOn"
    :disabled="formProps.disabled"
    :validate-on-input-delay="formProps.validateOnInputDelay"
    :transform="formProps.transform"
    :loading-auto="formProps.loadingAuto"
    :class="formClass ?? 'space-y-4'"
    @submit="emits('submit', $event)"
    @error="emits('error', $event)"
  >
    <slot name="before-fields" />

    <UFormField
      v-for="field in fields"
      :key="field.key"
      :name="field.formField.name"
      :label="field.formField.label"
      :description="field.formField.description"
      :hint="field.formField.hint"
      :help="field.formField.help"
      :required="field.formField.required"
      :orientation="field.formField.orientation"
    >
      <!-- hoverTip: show question mark icon next to label with tooltip -->
      <template v-if="field.formField.hoverTip" #label="{ label }">
        <span class="inline-flex items-center gap-1">
          <span>{{ label }}</span>
          <UTooltip :text="field.formField.hoverTip" :content="{ side: 'top', sideOffset: 4 }">
            <UIcon name="i-lucide-circle-help" class="size-3.5 text-muted-foreground cursor-help" />
          </UTooltip>
        </span>
      </template>
      <template v-for="slot in field.formField.slots" :key="slot" #[slot]>
        <slot :name="`${field.key}-${slot}`" />
      </template>

      <slot
        :name="field.key"
        :field="field.key"
        :state="props.state"
      >
        <component
          :is="field.component"
          v-bind="field.props"
          :model-value="getFieldValue(field.key)"
          @update:model-value="(value: unknown) => updateFieldValue(field.key, value)"
        >
          <slot :name="`${field.key}-content`" />
        </component>
      </slot>
    </UFormField>

    <slot name="after-fields" />

    <slot name="submit" />
  </UForm>
</template>
