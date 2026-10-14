<script setup lang="ts">
import { useAuthStore } from '../../../stores/auth'
import { ADMIN_TOKEN_COOKIE } from '../../../composables/useApi'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const authStore = useAuthStore()
const api = useApi()
const toast = useToast()

const profileForm = reactive({
  realName: '',
  email: '',
  mobile: '',
  sex: 1,
})

const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

async function loadProfile() {
  const info = authStore.userInfo
  if (info) {
    profileForm.realName = info.realName
    profileForm.email = info.email
    profileForm.mobile = info.mobile
    profileForm.sex = info.sex || 1
  }
}

onMounted(async () => {
  if (!authStore.userInfo) {
    await authStore.fetchUserInfo().catch(() => {})
  }
  loadProfile()
})

const savingProfile = ref(false)

async function saveProfile() {
  if (!authStore.userInfo) return
  savingProfile.value = true
  try {
    await api.post('/admin/member/update-profile', {
      id: authStore.userInfo.id,
      realName: profileForm.realName,
      email: profileForm.email,
      mobile: profileForm.mobile,
      sex: profileForm.sex,
    })
    toast.add({ title: '资料已更新', color: 'success' })
    await authStore.fetchUserInfo().catch(() => {})
  } catch (e: any) {
    toast.add({ title: '保存失败', description: e.message, color: 'error' })
  } finally {
    savingProfile.value = false
  }
}

// ==================== 头像上传 ====================

const uploading = ref(false)

async function uploadAvatar(file: File) {
  const token = useCookie<string>(ADMIN_TOKEN_COOKIE)
  const formData = new FormData()
  formData.append('file', file)

  uploading.value = true
  try {
    const res = await $fetch<{ code: number; data: { fileUrl: string } }>('/admin/upload/file', {
      method: 'POST',
      body: formData,
      headers: {
        ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
        uploadType: 'image',
      },
    })
    if (res.code !== 0) throw new Error('上传失败')
    return res.data.fileUrl
  } finally {
    uploading.value = false
  }
}

async function onAvatarChange(file: File) {
  try {
    const url = await uploadAvatar(file)
    if (!authStore.userInfo) return
    await api.post('/admin/member/update-profile', {
      id: authStore.userInfo.id,
      avatar: url,
      realName: profileForm.realName,
      email: profileForm.email,
      mobile: profileForm.mobile,
      sex: profileForm.sex,
    })
    toast.add({ title: '头像已更新', color: 'success' })
    await authStore.fetchUserInfo().catch(() => {})
  } catch (e: any) {
    toast.add({ title: '头像更新失败', description: e.message, color: 'error' })
  }
}

// ==================== 修改密码 ====================

const savingPwd = ref(false)

async function savePassword() {
  if (pwdForm.newPassword !== pwdForm.confirmPassword) {
    toast.add({ title: '两次输入的新密码不一致', color: 'error' })
    return
  }
  if (!authStore.userInfo) return
  savingPwd.value = true
  try {
    await api.post('/admin/member/update-pwd', {
      id: authStore.userInfo.id,
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword,
    })
    toast.add({ title: '密码已修改', color: 'success' })
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
  } catch (e: any) {
    toast.add({ title: '修改失败', description: e.message, color: 'error' })
  } finally {
    savingPwd.value = false
  }
}
</script>

<template>
  <div class="space-y-4">
    <UCard>
      <template #header>
        <h3 class="font-semibold">个人资料</h3>
      </template>
      <div class="flex gap-6 items-start">
        <div class="flex flex-col items-center gap-2">
          <UAvatar :src="authStore.userInfo?.avatar" :text="(authStore.userInfo?.realName || authStore.userInfo?.username || 'U').charAt(0).toUpperCase()" size="3xl" />
          <UButton size="xs" variant="outline" :loading="uploading" @click="() => {}">更换头像</UButton>
          <input type="file" accept="image/*" class="hidden" @change="(e) => { const f = (e.target as HTMLInputElement).files?.[0]; if (f) onAvatarChange(f) }">
        </div>

        <div class="flex-1 max-w-lg space-y-4">
          <UFormField label="账号">
            <UInput :model-value="authStore.userInfo?.username" disabled />
          </UFormField>
          <UFormField label="姓名">
            <UInput v-model="profileForm.realName" />
          </UFormField>
          <UFormField label="性别">
            <USelect v-model="profileForm.sex" :items="[{ label: '男', value: 1 }, { label: '女', value: 2 }]" />
          </UFormField>
          <UFormField label="邮箱">
            <UInput v-model="profileForm.email" />
          </UFormField>
          <UFormField label="手机号">
            <UInput v-model="profileForm.mobile" />
          </UFormField>
          <UButton color="primary" :loading="savingProfile" @click="saveProfile">保存资料</UButton>
        </div>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <h3 class="font-semibold">修改密码</h3>
      </template>
      <div class="max-w-lg space-y-4">
        <UFormField label="旧密码" required>
          <UInput v-model="pwdForm.oldPassword" type="password" />
        </UFormField>
        <UFormField label="新密码" required>
          <UInput v-model="pwdForm.newPassword" type="password" />
        </UFormField>
        <UFormField label="确认新密码" required>
          <UInput v-model="pwdForm.confirmPassword" type="password" />
        </UFormField>
        <UButton color="primary" :loading="savingPwd" @click="savePassword">修改密码</UButton>
      </div>
    </UCard>
  </div>
</template>
