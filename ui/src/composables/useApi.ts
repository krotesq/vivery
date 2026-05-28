import { ref } from 'vue'

export function useApi<T>() {
  const data    = ref<T | null>(null)
  const error   = ref<string | null>(null)
  const loading = ref(false)

  async function request(fn: () => Promise<T>) {
    loading.value = true
    error.value   = null
    try {
      data.value = await fn()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  return { data, error, loading, request }
}