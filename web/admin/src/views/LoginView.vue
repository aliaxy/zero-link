<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <span class="login-logo">zero·link</span>
        <p class="login-subtitle">Administrator sign in</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="Username"
            size="large"
            autocomplete="username"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="Password"
            size="large"
            autocomplete="current-password"
            :prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-button
          type="primary"
          size="large"
          :loading="loading"
          style="width: 100%; margin-top: 4px"
          native-type="submit"
          @click="handleLogin"
        >
          Sign In
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const formRef = ref(null)
const loading = ref(false)
const form = ref({ username: '', password: '' })

const rules = {
  username: [{ required: true, message: 'Enter your username', trigger: 'blur' }],
  password: [{ required: true, message: 'Enter your password', trigger: 'blur' }],
}

async function handleLogin() {
  await formRef.value.validate()
  loading.value = true
  try {
    await auth.login(form.value.username, form.value.password)
    const redirect = route.query.redirect ?? '/links'
    router.push(redirect)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 360px;
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-overlay);
  border: 1px solid var(--color-separator);
  padding: 40px 36px 36px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  display: block;
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-accent);
  letter-spacing: -0.02em;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 13px;
  color: var(--color-label-secondary);
}
</style>
