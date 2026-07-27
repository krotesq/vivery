import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/utils/api'
import type { AccountMe, LoginRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AccountMe | null>(null)
  const checked = ref(false)
  const isLoggedIn = computed(() => user.value !== null)

  async function login(payload: LoginRequest) {
    user.value = await api.post<AccountMe>('/account/login', payload)
  }

  async function logout() {
    try {
      await api.post('/account/logout', {})
    } finally {
      user.value = null
    }
  }

  async function check() {
    try {
      user.value = await api.get<AccountMe>('/account/me')
    } catch {
      user.value = null
    } finally {
      checked.value = true
    }
  }

  return { user, isLoggedIn, checked, login, logout, check }
}, {
  persist: { pick: ['user'] }
})