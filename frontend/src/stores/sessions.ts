import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SessionMeta, ChatMessage } from '../types/session'
import { ListSessions, CreateSession, SendMessage } from '../composables/useWails'

export const useSessionsStore = defineStore('sessions', () => {
  const list = ref<SessionMeta[]>([])
  const activeId = ref<string | null>(null)
  const messages = ref<Record<string, ChatMessage[]>>({})
  const streaming = ref<Record<string, boolean>>({})

  const active = computed(() => list.value.find((s) => s.id === activeId.value) ?? null)

  async function refresh() {
    list.value = await ListSessions()
  }

  async function create(workdir: string, prompt: string) {
    const id = await CreateSession(workdir, prompt)
    await refresh()
    activeId.value = id
    return id
  }

  function select(id: string) {
    activeId.value = id
  }

  async function send(id: string, prompt: string) {
    if (!messages.value[id]) messages.value[id] = []
    messages.value[id].push({ id: crypto.randomUUID(), role: 'user', content: prompt, ts: Date.now() })
    streaming.value[id] = true
    try {
      await SendMessage(id, prompt)
    } finally {
      streaming.value[id] = false
    }
  }

  return { list, activeId, active, messages, streaming, refresh, create, select, send }
})
