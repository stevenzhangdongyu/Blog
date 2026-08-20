<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const errorMessage = ref('')
const submitting = ref(false)

async function login() {
  submitting.value = true
  errorMessage.value = ''
  try {
    const response = await fetch('/api/admin/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: username.value, password: password.value }),
    })
    if (!response.ok) {
      const result = await response.json().catch(() => ({}))
      throw new Error(result.error || '登录失败')
    }
    await router.push('/admin/articles')
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '登录失败'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="admin-shell">
    <form class="admin-form" @submit.prevent="login">
      <p class="eyebrow">ADMIN</p>
      <h1>登录管理端</h1>
      <label>用户名 <input v-model="username" autocomplete="username" required /></label>
      <label>密码 <input v-model="password" type="password" autocomplete="current-password" required /></label>
      <p v-if="errorMessage" class="form-error">{{ errorMessage }}</p>
      <button type="submit" :disabled="submitting">{{ submitting ? '登录中...' : '登录' }}</button>
    </form>
  </main>
</template>
