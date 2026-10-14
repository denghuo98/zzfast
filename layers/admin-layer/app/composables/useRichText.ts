// 合并 nuxt-ui-templates/editor 的 5 个 composable：
// useEditorEmojis + useEditorMentions + useEditorSuggestions + useEditorDragHandle + useEditorToolbar
// 去掉了协作（y-partykit/yjs）和 AI（ai-sdk）相关部分。

import type { EditorToolbarItem, EditorSuggestionMenuItem, EditorEmojiMenuItem, EditorMentionMenuItem, DropdownMenuItem, EditorCustomHandlers } from '@nuxt/ui'
import type { Editor, JSONContent } from '@tiptap/vue-3'
import { TaskList, TaskItem } from '@tiptap/extension-list'
import { TableKit } from '@tiptap/extension-table'
import { Emoji, gitHubEmojis, shortcodeToEmoji } from '@tiptap/extension-emoji'
import { CodeBlockShiki } from 'tiptap-extension-code-block-shiki'
import { ImageUpload } from '~/components/editor/ImageUploadExtension'
import { mapEditorItems } from '@nuxt/ui/utils/editor'
import { upperFirst } from 'scule'

const CONVERTIBLE_TYPES = ['paragraph', 'heading', 'bulletList', 'orderedList', 'taskList', 'blockquote', 'codeBlock', 'listItem', 'taskItem']

// ========== Emoji ==========
function createEmojiExtension() {
  return Emoji.extend({
    markdownTokenName: 'emoji',
    markdownTokenizer: {
      name: 'emoji',
      level: 'inline',
      start: ':',
      tokenize(src: string) {
        const match = src.match(/^:([a-zA-Z0-9_+-]+):/)
        if (!match?.[1]) return undefined
        if (!shortcodeToEmoji(match[1], gitHubEmojis)) return undefined
        return { type: 'emoji', raw: match[0], name: match[1] }
      }
    },
    parseMarkdown(token: any, { createNode }: any) {
      return createNode('emoji', { name: token.name })
    },
    renderMarkdown(node: any) {
      if (!node.attrs?.name) return ''
      const emojiItem = shortcodeToEmoji(node.attrs.name, gitHubEmojis)
      return emojiItem?.emoji || `:${node.attrs.name}:`
    }
  })
}

// ========== Mentions (static fallback users) ==========
const FALLBACK_USERS: EditorMentionMenuItem[] = [
  { label: 'benjamincanac', avatar: { src: 'https://avatars.githubusercontent.com/u/739984?v=4' } },
  { label: 'atinux', avatar: { src: 'https://avatars.githubusercontent.com/u/904724?v=4' } },
  { label: 'HugoRCD', avatar: { src: 'https://avatars.githubusercontent.com/u/71938701?v=4' } },
]

// ========== Suggestions ==========
function getSuggestionItems<T extends EditorCustomHandlers>(_customHandlers?: T): EditorSuggestionMenuItem<T>[][] {
  return [[
    { type: 'label' as const, label: '样式' },
    { kind: 'paragraph' as const, label: '段落', icon: 'i-lucide-type' },
    { kind: 'heading' as const, level: 1, label: '标题 1', icon: 'i-lucide-heading-1' },
    { kind: 'heading' as const, level: 2, label: '标题 2', icon: 'i-lucide-heading-2' },
    { kind: 'heading' as const, level: 3, label: '标题 3', icon: 'i-lucide-heading-3' },
    { kind: 'bulletList' as const, label: '无序列表', icon: 'i-lucide-list' },
    { kind: 'orderedList' as const, label: '有序列表', icon: 'i-lucide-list-ordered' },
    { kind: 'taskList' as const, label: '任务列表', icon: 'i-lucide-list-check' },
    { kind: 'blockquote' as const, label: '引用', icon: 'i-lucide-text-quote' },
    { kind: 'codeBlock' as const, label: '代码块', icon: 'i-lucide-square-code' },
  ], [
    { type: 'label' as const, label: '插入' },
    { kind: 'mention' as const, label: '提及', icon: 'i-lucide-at-sign' },
    { kind: 'emoji' as const, label: '表情', icon: 'i-lucide-smile-plus' },
    { kind: 'imageUpload' as const, label: '图片', icon: 'i-lucide-image' },
    { kind: 'table' as const, label: '表格', icon: 'i-lucide-table' },
    { kind: 'horizontalRule' as const, label: '分隔线', icon: 'i-lucide-separator-horizontal' },
  ]] satisfies EditorSuggestionMenuItem<T>[][]
}

// ========== DragHandle ==========
function getTypeSpecificItems(editor: Editor, nodeType: string, pos?: number): DropdownMenuItem[] {
  if (CONVERTIBLE_TYPES.includes(nodeType)) {
    return [
      {
        label: '转换为',
        icon: 'i-lucide-repeat-2',
        children: [
          { kind: 'paragraph' as const, label: '段落', icon: 'i-lucide-type' },
          { kind: 'heading' as const, level: 1, label: '标题 1', icon: 'i-lucide-heading-1' },
          { kind: 'heading' as const, level: 2, label: '标题 2', icon: 'i-lucide-heading-2' },
          { kind: 'heading' as const, level: 3, label: '标题 3', icon: 'i-lucide-heading-3' },
          { kind: 'heading' as const, level: 4, label: '标题 4', icon: 'i-lucide-heading-4' },
          { kind: 'bulletList' as const, label: '无序列表', icon: 'i-lucide-list' },
          { kind: 'orderedList' as const, label: '有序列表', icon: 'i-lucide-list-ordered' },
          { kind: 'taskList' as const, label: '任务列表', icon: 'i-lucide-list-check' },
          { kind: 'blockquote' as const, label: '引用', icon: 'i-lucide-text-quote' },
          { kind: 'codeBlock' as const, label: '代码块', icon: 'i-lucide-square-code' },
        ]
      },
      { kind: 'clearFormatting' as const, pos, label: '清除格式', icon: 'i-lucide-rotate-ccw' },
    ]
  }

  if (nodeType === 'image') {
    const node = pos !== undefined ? editor.state.doc.nodeAt(pos) : null
    return [{ label: '下载图片', icon: 'i-lucide-download', to: node?.attrs?.src, download: true }]
  }

  if (nodeType === 'table') {
    return [{
      label: '清空所有内容',
      icon: 'i-lucide-square-x',
      onSelect: () => {
        if (pos === undefined) return
        const tableNode = editor.state.doc.nodeAt(pos)
        if (!tableNode) return

        const cellRanges: { from: number, to: number }[] = []
        tableNode.descendants((node: any, nodePos: number) => {
          if (node.type.name === 'tableCell' || node.type.name === 'tableHeader') {
            const cellStart = pos + 1 + nodePos + 1
            const cellEnd = cellStart + node.content.size
            if (node.content.size > 0) {
              cellRanges.push({ from: cellStart, to: cellEnd })
            }
          }
          return true
        })

        const { tr } = editor.state
        cellRanges.reverse().forEach(({ from, to }) => {
          tr.delete(from, to)
        })
        editor.view.dispatch(tr)
      }
    }]
  }

  return []
}

// ========== Toolbar ==========
function getToolbarItems<T extends EditorCustomHandlers>(customHandlers?: T) {
  const toolbarItems: EditorToolbarItem<T>[][] = [[
    { kind: 'undo' as const, icon: 'i-lucide-undo', tooltip: { text: '撤销' } },
    { kind: 'redo' as const, icon: 'i-lucide-redo', tooltip: { text: '重做' } },
  ], [
    { kind: 'imageUpload' as const, label: '添加', icon: 'i-lucide-image', tooltip: { text: '添加图片' } },
  ]]

  // Bubble toolbar — no AI sparkles group (removed), only formatting
  const bubbleToolbarItems = computed<EditorToolbarItem<T>[][]>(() => [[
    {
      label: '转换为',
      trailingIcon: 'i-lucide-chevron-down',
      activeColor: 'neutral' as const,
      activeVariant: 'ghost' as const,
      tooltip: { text: '转换为' },
      content: { align: 'start' as const },
      ui: { label: 'text-xs' },
      items: [
        { type: 'label' as const, label: '转换为' },
        { kind: 'paragraph' as const, label: '段落', icon: 'i-lucide-type' },
        { kind: 'heading' as const, level: 1, label: '标题 1', icon: 'i-lucide-heading-1' },
        { kind: 'heading' as const, level: 2, label: '标题 2', icon: 'i-lucide-heading-2' },
        { kind: 'heading' as const, level: 3, label: '标题 3', icon: 'i-lucide-heading-3' },
        { kind: 'heading' as const, level: 4, label: '标题 4', icon: 'i-lucide-heading-4' },
        { kind: 'bulletList' as const, label: '无序列表', icon: 'i-lucide-list' },
        { kind: 'orderedList' as const, label: '有序列表', icon: 'i-lucide-list-ordered' },
        { kind: 'taskList' as const, label: '任务列表', icon: 'i-lucide-list-check' },
        { kind: 'blockquote' as const, label: '引用', icon: 'i-lucide-text-quote' },
        { kind: 'codeBlock' as const, label: '代码块', icon: 'i-lucide-square-code' },
      ]
    },
  ], [
    { kind: 'mark' as const, mark: 'bold' as const, icon: 'i-lucide-bold', tooltip: { text: '加粗' } },
    { kind: 'mark' as const, mark: 'italic' as const, icon: 'i-lucide-italic', tooltip: { text: '斜体' } },
    { kind: 'mark' as const, mark: 'underline' as const, icon: 'i-lucide-underline', tooltip: { text: '下划线' } },
    { kind: 'mark' as const, mark: 'strike' as const, icon: 'i-lucide-strikethrough', tooltip: { text: '删除线' } },
    { kind: 'mark' as const, mark: 'code' as const, icon: 'i-lucide-code', tooltip: { text: '行内代码' } },
  ], [
    { slot: 'link' as const, icon: 'i-lucide-link' },
    { kind: 'imageUpload' as const, icon: 'i-lucide-image', tooltip: { text: '图片' } },
  ]] satisfies EditorToolbarItem<T>[][])

  return { toolbarItems, bubbleToolbarItems }
}

function getImageToolbarItems<T extends EditorCustomHandlers>(editor: Editor): EditorToolbarItem<T>[][] {
  const node = editor.state.doc.nodeAt(editor.state.selection.from)

  return [[
    { icon: 'i-lucide-download', to: node?.attrs?.src, download: true, tooltip: { text: '下载' } },
    {
      icon: 'i-lucide-refresh-cw',
      tooltip: { text: '替换' },
      onClick: () => {
        const { state } = editor
        const { selection } = state
        const pos = selection.from
        const node = state.doc.nodeAt(pos)
        if (node && node.type.name === 'image') {
          editor.chain().focus().deleteRange({ from: pos, to: pos + node.nodeSize }).insertContentAt(pos, { type: 'imageUpload' }).run()
        }
      }
    },
  ], [
    {
      icon: 'i-lucide-trash',
      tooltip: { text: '删除' },
      onClick: () => {
        const { state } = editor
        const { selection } = state
        const pos = selection.from
        const node = state.doc.nodeAt(pos)
        if (node && node.type.name === 'image') {
          editor.chain().focus().deleteRange({ from: pos, to: pos + node.nodeSize }).run()
        }
      }
    }
  ]]
}

function getTableToolbarItems<T extends EditorCustomHandlers>(editor: Editor): EditorToolbarItem<T>[][] {
  return [[
    { icon: 'i-lucide-between-vertical-start', tooltip: { text: '上方插入行' }, onClick: () => { editor.chain().focus().addRowBefore().run() } },
    { icon: 'i-lucide-between-vertical-end', tooltip: { text: '下方插入行' }, onClick: () => { editor.chain().focus().addRowAfter().run() } },
    { icon: 'i-lucide-between-horizontal-start', tooltip: { text: '左侧插入列' }, onClick: () => { editor.chain().focus().addColumnBefore().run() } },
    { icon: 'i-lucide-between-horizontal-end', tooltip: { text: '右侧插入列' }, onClick: () => { editor.chain().focus().addColumnAfter().run() } },
  ], [
    { icon: 'i-lucide-rows-3', tooltip: { text: '删除行' }, onClick: () => { editor.chain().focus().deleteRow().run() } },
    { icon: 'i-lucide-columns-3', tooltip: { text: '删除列' }, onClick: () => { editor.chain().focus().deleteColumn().run() } },
  ], [
    { icon: 'i-lucide-trash', tooltip: { text: '删除表格' }, onClick: () => { editor.chain().focus().deleteTable().run() } },
  ]]
}

// ========== Main composable ==========
export function useRichText<T extends EditorCustomHandlers>(customHandlers?: T) {
  const EmojiExtension = createEmojiExtension()

  const extensions = computed(() => [
    CodeBlockShiki.configure({
      defaultTheme: 'material-theme',
      themes: {
        light: 'material-theme-lighter',
        dark: 'material-theme-palenight'
      }
    }),
    EmojiExtension,
    ImageUpload,
    TableKit,
    TaskList,
    TaskItem,
  ])

  const emojiItems: EditorEmojiMenuItem[] = gitHubEmojis.filter(
    emoji => !emoji.name.startsWith('regional_indicator_')
  )

  const mentionItems: EditorMentionMenuItem[] = FALLBACK_USERS

  const suggestionItems = computed(() => getSuggestionItems(customHandlers))

  const selectedNode = ref<{ node: JSONContent | null, pos: number }>()

  const getDragHandleItems = (editor: Editor): DropdownMenuItem[][] => {
    if (!selectedNode.value?.node?.type) return []

    const nodeType = selectedNode.value.node.type
    const pos = selectedNode.value.pos
    const typeSpecificItems = getTypeSpecificItems(editor, nodeType, pos)

    return mapEditorItems(editor, [
      [
        { type: 'label' as const, label: upperFirst(nodeType) },
        ...typeSpecificItems
      ],
      [
        { kind: 'duplicate' as const, pos, label: '复制', icon: 'i-lucide-copy' },
        {
          label: '复制到剪贴板',
          icon: 'i-lucide-clipboard',
          onSelect: async () => {
            if (!selectedNode.value) return
            const node = editor.state.doc.nodeAt(pos)
            if (node) {
              await navigator.clipboard.writeText(node.textContent)
            }
          }
        }
      ],
      [
        { kind: 'moveUp' as const, pos, label: '上移', icon: 'i-lucide-arrow-up' },
        { kind: 'moveDown' as const, pos, label: '下移', icon: 'i-lucide-arrow-down' },
      ],
      [
        { kind: 'delete' as const, pos, label: '删除', icon: 'i-lucide-trash' },
      ],
    ]) as DropdownMenuItem[][]
  }

  const onNodeChange = (event: { node: JSONContent | null, pos: number }) => {
    selectedNode.value = event
  }

  const { toolbarItems, bubbleToolbarItems } = getToolbarItems(customHandlers)

  return {
    extensions,
    emojiItems,
    mentionItems,
    suggestionItems,
    toolbarItems,
    bubbleToolbarItems,
    selectedNode,
    getDragHandleItems,
    onNodeChange,
    getImageToolbarItems,
    getTableToolbarItems,
  }
}
