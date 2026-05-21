import { defineStore } from "pinia"
import { ref, computed } from "vue"

export const useAuthStore = defineStore("auth", () => {
  const user = ref(null)
  const checked = ref(false)
  const isLoggedIn = computed(() => user.value !== null)

  async function login(username: string, password: string) {
    const res = await fetch("/api/account/login", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({username, password})
    })
    if (!res.ok) throw new Error("Error while logging in")
    user.value = await res.json()
  }

  function logout() {
    user.value = null
  }

  async function check() {
    try {
      const res = await fetch("/api/account/me")
      if (res.ok) {
        user.value = await res.json()
      }
      else {
        user.value = null
      }
    }
    catch {
      user.value = null
    }
    finally {
      checked.value = true
    }
  }

  return { user, isLoggedIn, checked, login, logout, check }
}, {
  persist: { // ensure that the store is still available after page reload
    pick: ["user"]
  }
})
