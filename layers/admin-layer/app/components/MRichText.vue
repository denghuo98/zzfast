<!-- 富文本编辑器：基于 Nuxt UI UEditor (TipTap v3) + 官方 editor 模板复刻。
     去掉协作 + AI，图片上传走自建 /api/upload 后端接口。 -->
<script setup lang="ts">
import type { EditorCustomHandlers } from '@nuxt/ui'
import type { Editor } from '@tiptap/core'
import { CellSelection } from '@tiptap/pm/tables'

const props = withDefaults(defineProps<{
  modelValue?: string
  contentType?: 'markdown' | 'html' | 'json'
  placeholder?: string
  editable?: boolean
  autofocus?: boolean
  /** 自定义上传接口路径，默认 /api/upload */
  uploadUrl?: string
  /** 自定义上传函数，优先级高于 uploadUrl。签名: (file: File) => Promise<{ url: string }> */
  uploadHandler?: (file: File) => Promise<{ url: string }>
}>(), {
  contentType: 'markdown',
  placeholder: "输入内容，输入 '/' 打开命令菜单...",
  editable: true,
  autofocus: true,
  uploadUrl: '/api/upload',
})

// 注入上传配置到模块级 composable（TipTap NodeView 不在当前 Vue 组件树内）
const { setUploadConfig } = useEditorUpload()
setUploadConfig({ url: props.uploadUrl, handler: props.uploadHandler })

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editorRef = useTemplateRef('editorRef')

// Custom handlers
const customHandlers = {
  imageUpload: {
    canExecute: () => true,
    execute: (editor: Editor) => editor.chain().focus().insertContent({ type: 'imageUpload' }).run(),
    isActive: (editor: Editor) => editor.isActive('imageUpload'),
    isDisabled: undefined,
  },
  table: {
    canExecute: (editor: Editor) => editor.can().insertTable({ rows: 3, cols: 3, withHeaderRow: true }),
    execute: (editor: Editor) => editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }),
    isActive: (editor: Editor) => editor.isActive('table'),
    isDisabled: undefined,
  },
} satisfies EditorCustomHandlers

const {
  extensions,
  emojiItems,
  mentionItems,
  suggestionItems,
  toolbarItems,
  bubbleToolbarItems,
  getDragHandleItems,
  onNodeChange,
  getImageToolbarItems,
  getTableToolbarItems,
} = useRichText(customHandlers)

const content = ref(props.modelValue || '')

watch(() => props.modelValue, (v) => {
  if (v !== undefined && v !== content.value) {
    content.value = v
  }
})

function onUpdate(value: string) {
  content.value = value
  emit('update:modelValue', value)
}

// Expose for parent access
const editorInstance = ref<Editor | null>(null)
function onCreate({ editor }: { editor: Editor }) {
  editorInstance.value = editor
}

function getHtml() { return editorInstance.value?.getHTML() || '' }
function getMarkdown() { return content.value }
function getText() { return editorInstance.value?.getText() || '' }
function clear() { editorInstance.value?.commands.clearContent() }

defineExpose({ editor: editorInstance, getHtml, getMarkdown, getText, clear })
</script>

<template>
  <ClientOnly>
    <div class="border border-default rounded-lg">
      <UEditor
        ref="editorRef"
        v-slot="{ editor, handlers }"
        :model-value="content"
        :content-type="contentType"
        :extensions="extensions"
        :handlers="customHandlers"
        :autofocus="autofocus"
        :editable="editable"
        :placeholder="placeholder"
        class="min-h-64"
        :ui="{
          base: 'p-4 sm:p-8',
          content: 'max-w-4xl mx-auto'
        }"
        @update:model-value="onUpdate"
        @create="onCreate"
      >
        <!-- Fixed toolbar -->
        <div class="flex items-center justify-center mb-4">
          <UEditorToolbar
            :editor="editor"
            :items="toolbarItems"
          />
        </div>

        <!-- Bubble toolbar (text selection) -->
        <UEditorToolbar
          :editor="editor"
          :items="bubbleToolbarItems"
          layout="bubble"
          :should-show="({ editor, view, state }: any) => {
            if (editor.isActive('imageUpload') || editor.isActive('image') || state.selection instanceof CellSelection) {
              return false
            }
            const { selection } = state
            return view.hasFocus() && !selection.empty
          }"
        >
          <template #link>
            <EditorLinkPopover :editor="editor" />
          </template>
        </UEditorToolbar>

        <!-- Image bubble toolbar -->
        <UEditorToolbar
          :editor="editor"
          :items="getImageToolbarItems(editor)"
          layout="bubble"
          :should-show="({ editor, view }: any) => {
            return editor.isActive('image') && view.hasFocus()
          }"
        />

        <!-- Table bubble toolbar -->
        <UEditorToolbar
          :editor="editor"
          :items="getTableToolbarItems(editor)"
          layout="bubble"
          :should-show="({ editor, view }: any) => {
            return editor.state.selection instanceof CellSelection && view.hasFocus()
          }"
        />

        <!-- Menus -->
        <UEditorEmojiMenu :editor="editor" :items="emojiItems" />
        <UEditorMentionMenu :editor="editor" :items="mentionItems" />
        <UEditorSuggestionMenu :editor="editor" :items="suggestionItems" />

        <!-- Drag handle -->
        <UEditorDragHandle
          v-slot="{ ui, onClick }"
          :editor="editor"
          @node-change="onNodeChange"
        >
          <UButton
            icon="i-lucide-plus"
            color="neutral"
            variant="ghost"
            size="sm"
            :class="ui.handle()"
            @click="(e: MouseEvent) => {
              e.stopPropagation()
              const node = onClick()
              handlers.suggestion?.execute(editor, { pos: node?.pos }).run()
            }"
          />

          <UDropdownMenu
            v-slot="{ open }"
            :modal="false"
            :items="getDragHandleItems(editor)"
            :content="{ side: 'left' }"
            :ui="{ content: 'w-48', label: 'text-xs' }"
            @update:open="editor.chain().setMeta('lockDragHandle', $event).run()"
          >
            <UButton
              color="neutral"
              variant="ghost"
              active-variant="soft"
              size="sm"
              icon="i-lucide-grip-vertical"
              :active="open"
              :class="ui.handle()"
            />
          </UDropdownMenu>
        </UEditorDragHandle>
      </UEditor>
    </div>

    <template #fallback>
      <div class="border border-default rounded-lg p-4 sm:p-8 min-h-64 flex items-center justify-center">
        <UIcon name="i-lucide-loader-2" class="animate-spin size-6 text-muted-foreground" />
      </div>
    </template>
  </ClientOnly>
</template>
