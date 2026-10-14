import type { ButtonProps, InferInput } from '@nuxt/ui'
import type { ConcreteComponent } from '@vue/runtime-core'

export type MFormState<TSchema> = Partial<InferInput<TSchema>>

export type ComponentsMap = Record<string, (params: {
  key: string
  state: any
  zodType: any
  config: MFormConfig
}) => ComponentDefinition | null>

export interface ComponentDefinition {
  component: any
  componentProps?: Record<string, any>
}

export interface MFormConfig {
  /** Submit button configuration */
  submit?: {
    /** Name of the globally defined component to use as the submit button */
    component: ConcreteComponent | string
    /** Props to pass to the custom submit button */
    props?: Record<string, any>
  } | {
    component?: undefined
    /** Props to pass to the Nuxt UI `UButton` component */
    props?: ButtonProps
  } | false

  /** MFormModal default button labels */
  modal?: {
    /** Default label for the submit button. @default 'Submit' */
    submitLabel?: string
    /** Default label for the cancel/close button. @default 'Cancel' */
    closeLabel?: string
  }

  /**
   * Customize default form styles
   */
  theme?: {
    /**
     * The orientation of form fields.
     * - `vertical`: label on top, input below (default)
     * - `horizontal`: label and input side by side
     * @default 'vertical'
     */
    orientation?: 'vertical' | 'horizontal'

    /**
     * Apply `w-full` class to all input components by default.
     * When orientation is `horizontal`, this defaults to false.
     * @default true (vertical) / false (horizontal)
     */
    wFull?: boolean

    /**
     * Generate default labels from field names globally. Set to `false` to disable auto-generated labels.
     * Manually set titles via `meta().title` will still render regardless of this setting.
     * @default true
     */
    enableDefaultTitles?: boolean
  }

  /**
   * Customize the default components used for rendering form fields.
   *
   * @example
   * ```ts
   * components: {
   *   email: () => ({ component: UInput, componentProps: { type: 'email' } }),
   * }
   * ```
   */
  components?: ComponentsMap

  /**
   * Custom meta string preprocessors that transform prefixed strings in `meta()`.
   * The `$` prefix and `:` suffix are added automatically to keys.
   *
   * @example
   * ```ts
   * metaStringProcessors: {
   *   upper: value => value.toUpperCase(),
   *   slug: value => value.trim().toLowerCase().replace(/\s+/g, '-'),
   * }
   * ```
   */
  metaStringProcessors?: Record<string, (value: string, context: {
    key: string
    path: string
    meta: Record<string, any>
    config: MFormConfig
  }) => any>
}

declare module 'zod' {
  interface GlobalMeta {
    /** Auto-form specific configuration */
    theme?: {
      /**
       * Override the form field orientation for this field.
       * Falls back to config.theme.orientation.
       */
      orientation?: 'vertical' | 'horizontal'

      /**
       * Apply `w-full` class to this input component.
       * @default follows config.theme.wFull
       */
      wFull?: boolean
    }

    /**
     * Configuration for the inner Input component.
     */
    input?: {
      /**
       * Vue component or component name used to render the input.
       */
      component?: string | Component
      /**
       * Props forwarded to the input component at render time.
       * These are merged with the module's default props - user props take precedence.
       */
      props?: Record<string, any>
    }

    /**
     * Field label displayed in the form's label row.
     * Maps directly to `UFormField.label`.
     * Set to `null` to hide the label for this field.
     *
     * @default Field name in PascalCase with spaces
     */
    title?: string | null

    /**
     * Marks the field as required in the UI (asterisk, aria-required).
     * Validation is not enforced here; use `zod` refinements for that.
     */
    required?: boolean

    /**
     * Supporting text displayed between the label and the input control.
     */
    description?: string

    /**
     * Short inline hint shown next to the label.
     */
    hint?: string

    /**
     * Help text displayed below the input field.
     */
    help?: string

    /**
     * Hover tooltip text displayed as a question mark icon next to the label.
     * When set, a help icon appears; hovering it shows this text in a tooltip.
     */
    hoverTip?: string
  }
}

export {}
