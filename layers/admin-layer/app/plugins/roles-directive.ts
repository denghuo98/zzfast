/**
 * v-roles 指令（client plugin）
 *
 * 用法：<UButton v-roles="'admin'"> 或 <UButton v-roles="['admin','super']">
 *
 * 权限控制由后端 Casbin（AdminAuth 中间件）保证，前端角色列表尚未接入，
 * 因此该指令当前默认放行（不移除元素），保留接口供后续接入角色判断。
 */
import type { Directive, DirectiveBinding } from 'vue'

function check(_el: HTMLElement, _binding: DirectiveBinding<string | string[]>): void {
  // 预留：接入角色后在此按 OR 语义判断，无权限时 el.parentNode?.removeChild(el)
}

const rolesDirective: Directive<HTMLElement, string | string[]> = {
  mounted: check,
  updated: check,
}

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive('roles', rolesDirective)
})
