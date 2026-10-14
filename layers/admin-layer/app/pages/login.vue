<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { useAuthStore } from '../stores/auth'

definePageMeta({ layout: false })

const authStore = useAuthStore()
const route = useRoute()
const toast = useToast()

const captchaId = ref('')
const captchaBase64 = ref('')
const captchaLoading = ref(false)
const loading = ref(false)

const schema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(1, '请输入密码'),
  code: z.string().optional(),
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
  username: '',
  password: '',
  code: '',
})

// 已登录则直接跳首页
onMounted(() => {
  if (authStore.isAuthenticated) {
    navigateTo('/', { replace: true })
    return
  }
  getCaptcha()
})

async function getCaptcha() {
  try {
    captchaLoading.value = true
    const res = await authStore.fetchCaptcha()
    captchaId.value = res.id
    captchaBase64.value = res.base64
  } catch (error: any) {
    toast.add({ title: '获取验证码失败', description: error.message, color: 'error' })
  } finally {
    captchaLoading.value = false
  }
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    loading.value = true
    await authStore.login({
      username: event.data.username,
      password: event.data.password,
      cid: captchaId.value,
      code: event.data.code,
    })

    toast.add({ title: '登录成功', color: 'success' })

    const redirect = (route.query.redirect as string) || '/'
    if (import.meta.client) {
      window.location.replace(redirect)
    }
  } catch (error: any) {
    toast.add({ title: '登录失败', description: error.message || '请检查账号密码', color: 'error' })
    getCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted p-4">
    <UCard class="w-full max-w-md">
      <div class="text-center mb-6">
        <UIcon name="i-lucide-shield" class="w-10 h-10 text-primary mx-auto mb-3" />
        <h1 class="text-xl font-semibold">{{ useAppConfig().admin?.title || 'Admin Panel' }}</h1>
        <p class="text-muted text-sm mt-1">请登录管理后台</p>
      </div>

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormField label="用户名" name="username">
          <UInput v-model="state.username" placeholder="请输入用户名" icon="i-lucide-user" />
        </UFormField>

        <UFormField label="密码" name="password">
          <UInput v-model="state.password" type="password" placeholder="请输入密码" icon="i-lucide-lock" />
        </UFormField>

        <UFormField label="验证码" name="code">
          <div class="flex gap-2">
            <UInput v-model="state.code" placeholder="请输入验证码" icon="i-lucide-shield-check" class="flex-1" />
            <button
              type="button"
              class="w-24 h-9 shrink-0 border border-default rounded-md overflow-hidden cursor-pointer hover:border-muted flex items-center justify-center text-muted"
              :disabled="captchaLoading"
              @click="getCaptcha"
            >
              <img v-if="captchaBase64" :src="captchaBase64" alt="验证码" class="w-full h-full object-cover">
              <UIcon v-else-if="captchaLoading" name="i-lucide-loader-2" class="animate-spin" />
              <span v-else class="text-xs">获取验证码</span>
            </button>
          </div>
        </UFormField>

        <UButton type="submit" block size="lg" :loading="loading">登录</UButton>
      </UForm>
    </UCard>
  </div>
</template>
