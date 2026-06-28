import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SessionMeta, ChatMessage, SessionState } from '../types/session'
import { ListSessions, CreateSession, SendMessage } from '../composables/useWails'

export interface PendingPerm {
  tool: string
  args: unknown
  reqId: string
}

export const useSessionsStore = defineStore('sessions', () => {
  const list = ref<SessionMeta[]>([])
  const activeId = ref<string | null>(null)
  const messages = ref<Record<string, ChatMessage[]>>({})
  const streaming = ref<Record<string, boolean>>({})
  const state = ref<Record<string, SessionState>>({})
  const pending = ref<Record<string, PendingPerm | null>>({})
  const toolBlocks = ref<Record<string, Array<{ name: string; args: unknown }>>>({})

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
    state.value[id] = 'running'
    try {
      await SendMessage(id, prompt)
    } finally {
      streaming.value[id] = false
    }
  }

  // 处理来自 Go 端 EventsEmit 的 session 事件
  function handleEvent(sid: string, line: string) {
    let evt: any
    try { evt = JSON.parse(line) } catch { return }

    switch (evt.type) {
      case 'message': {
        const d = evt.data || evt
        if (!messages.value[sid]) messages.value[sid] = []
        messages.value[sid].push({
          id: crypto.randomUUID(),
          role: d.role,
          content: d.content || '',
          ts: Date.now(),
        })
        break
      }
      case 'tool_use': {
        const d = evt.data || evt
        if (!toolBlocks.value[sid]) toolBlocks.value[sid] = []
        toolBlocks.value[sid].push({ name: d.name, args: d.args })
        break
      }
      case 'permission_request': {
        const d = evt.data || evt
        pending.value[sid] = { tool: d.tool, args: d.args, reqId: d.request_id }
        state.value[sid] = 'awaiting_permission'
        break
      }
      case 'result':
        streaming.value[sid] = false
        state.value[sid] = 'idle'
        break
      case 'done':
        streaming.value[sid] = false
        state.value[sid] = 'idle'
        break
    }
  }

  return { list, activeId, active, messages, streaming, state, pending, toolBlocks, refresh, create, select, send, handleEvent }
})
