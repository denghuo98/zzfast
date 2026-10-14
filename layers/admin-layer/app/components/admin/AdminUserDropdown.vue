<template>
  <UDropdownMenu :items="dropdownItems" :content="{ align: 'start' }">
    <UButton
      v-if="!collapsed"
      :avatar="avatarProps"
      :label="displayName"
      color="neutral"
      variant="ghost"
      class="w-full justify-start"
      block
    />
    <UButton
      v-else
      :avatar="avatarProps"
      color="neutral"
      variant="ghost"
      class="w-full"
      square
    />
  </UDropdownMenu>
</template>

<script setup lang="ts">
import type { AdminUserInfo } from '../../types/admin-api'

const props = defineProps<{
  collapsed?: boolean
  userInfo?: AdminUserInfo | null
}>()

const emit = defineEmits<{
  logout: []
}>()

const displayName = computed(() => {
  return props.userInfo?.realName || props.userInfo?.username || 'User'
})

const avatarProps = computed(() => ({
  alt: displayName.value,
  text: displayName.value.charAt(0).toUpperCase(),
  size: 'sm' as const,
}))

const dropdownItems = computed(() => [
  [
    {
      label: displayName.value,
      type: 'label' as const,
    },
  ],
  [
    {
      label: '退出登录',
      icon: 'i-lucide-log-out',
      onSelect: () => emit('logout'),
    },
  ],
])
</script>
