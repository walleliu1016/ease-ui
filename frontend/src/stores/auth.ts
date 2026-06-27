import { defineStore } from 'pinia'
import { ref } from 'vue'
import { IsInitialized, Verify, LockoutState } from '../composables/useWails'

export const useAuthStore = defineStore('auth', () => {
  const initialized = ref(false)
  const attempts = ref(0)
  const lockedUntil = ref<Date | null>(null)

  async function checkInit() {
    initialized.value = await IsInitialized()
  }

  async function login(password: string): Promise<string | null> {
    try {
      await Verify(password)
      attempts.value = 0
      lockedUntil.value = null
      return null
    } catch (e: any) {
      // Reload lockout state to learn if we are now locked
      const [a, until] = await LockoutState()
      attempts.value = a
      lockedUntil.value = until && new Date(until).getTime() > Date.now() ? new Date(until) : null
      return e?.message ?? '登录失败'
    }
  }

  return { initialized, attempts, lockedUntil, checkInit, login }
})
