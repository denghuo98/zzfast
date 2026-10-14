/**
 * useConfirmDialog — 程序化确认弹窗
 *
 * 基于 Nuxt UI useOverlay.create 实现，无需在模板中声明组件即可调用。
 *
 * @example
 * ```ts
 * const confirm = useConfirmDialog()
 * const ok = await confirm({ title: '删除', description: '确定删除吗？', variant: 'danger' })
 * if (ok) { ... }
 * ```
 */
import ConfirmDialog from '../components/ConfirmDialog.vue'

export interface ConfirmDialogOptions {
  title?: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  icon?: string
  /** confirm 普通确认，danger 危险操作（红色按钮） */
  variant?: 'confirm' | 'danger'
}

export function useConfirmDialog() {
  const overlay = useOverlay()

  return (options: ConfirmDialogOptions): Promise<boolean> => {
    const modal = overlay.create(ConfirmDialog, {
      destroyOnClose: true,
      props: options,
    })
    return modal.open()
  }
}
