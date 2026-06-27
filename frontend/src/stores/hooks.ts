import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HooksConfig } from '../types/hooks'
import { GetHooksConfig, SaveHooksConfig } from '../composables/useWails'

export const useHooksStore = defineStore('hooks', () => {
  const cfg = ref<HooksConfig | null>(null)
  const dirty = ref(false)

  async function load() {
    cfg.value = await GetHooksConfig()
    dirty.value = false
  }

  async function save() {
    if (!cfg.value) return
    await SaveHooksConfig(cfg.value)
    dirty.value = false
  }

  function markDirty() { dirty.value = true }

  return { cfg, dirty, load, save, markDirty }
})
