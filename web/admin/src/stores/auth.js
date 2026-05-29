import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi } from '@/api/auth'

const TOKEN_KEY = 'zl_token'
const ADMIN_KEY = 'zl_admin'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) ?? null)
  const admin = ref(JSON.parse(localStorage.getItem(ADMIN_KEY) ?? 'null'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const data = await loginApi(username, password)
    token.value = data.token
    admin.value = data.admin
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(ADMIN_KEY, JSON.stringify(data.admin))
  }

  function logout() {
    token.value = null
    admin.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(ADMIN_KEY)
  }

  return { token, admin, isLoggedIn, login, logout }
})
