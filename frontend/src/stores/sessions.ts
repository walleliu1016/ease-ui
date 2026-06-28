import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SessionMeta, ChatMessage, SessionState } from '../types/session'
import { ListSessions, CreateSession, SendMessage, GetSessionMessages } from '../composables/useWails'

export interface PendingPerm {
  tool: string
  args: unknown
  reqId: string
}

// 解析 Claude API content block，提取文本拼成 markdown。
// raw 可能是 plain string，也可能是 content-block 数组。
function formatContent(raw: any): string {
  if (Array.isArray(raw)) {
    return raw
      .map((b: any) => {
        if (b.type === 'text' && b.text) return b.text
        if (b.type === 'thinking' && b.thinking) return `> ${b.thinking}`
        if (b.type === 'tool_use') return `> 🔧 **${b.name}**`
        return ''
      })
      .filter(Boolean)
      .join('\n\n')
  }
  if (typeof raw === 'string') {
    try {
      const blocks = JSON.parse(raw)
      if (Array.isArray(blocks)) return formatContent(blocks)
    } catch {}
    return raw
  }
  return String(raw || '')
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
    const meta = list.value.find((s) => s.id === id)
    if (meta && !messages.value[id]) {
      loadHistory(id, meta.workdir)
    }
  }

  async function loadHistory(sid: string, workdir: string) {
    try {
      const raw = await GetSessionMessages(sid, workdir)
      messages.value = {
        ...messages.value,
        [sid]: (raw || []).map((m: any, i: number) => ({
          id: `${sid}-${i}`,
          role: m.role || m.Role || 'assistant',
          content: formatContent(m.content || m.Content || ''),
          ts: Date.now() - ((raw || []).length - i) * 1000,
        })),
      }
    } catch (e: any) {
      console.error('[sessions] loadHistory failed:', e?.message || e)
    }
  }

  async function send(id: string, prompt: string) {
    const prev = messages.value[id] || []
    messages.value = { ...messages.value, [id]: [...prev, { id: crypto.randomUUID(), role: 'user', content: prompt, ts: Date.now() }] }
    streaming.value = { ...streaming.value, [id]: true }
    state.value = { ...state.value, [id]: 'running' }
    try {
      await SendMessage(id, prompt)
    } finally {
      streaming.value = { ...streaming.value, [id]: false }
    }
  }

  function handleEvent(sid: string, line: string) {
    let evt: any
    try { evt = JSON.parse(line) } catch { return }

    switch (evt.type) {
      case 'message': {
        const d = evt.data || evt
        const prev = messages.value[sid] || []
        messages.value = { ...messages.value, [sid]: [...prev, {
          id: crypto.randomUUID(),
          role: d.role,
          content: d.content || '',
          ts: Date.now(),
        }] }
        break
      }
      case 'tool_use': {
        const d = evt.data || evt
        const prev = toolBlocks.value[sid] || []
        toolBlocks.value = { ...toolBlocks.value, [sid]: [...prev, { name: d.name, args: d.args }] }
        break
      }
      case 'permission_request': {
        const d = evt.data || evt
        pending.value = { ...pending.value, [sid]: { tool: d.tool, args: d.args, reqId: d.request_id } }
        state.value = { ...state.value, [sid]: 'awaiting_permission' }
        break
      }
      case 'result':
        streaming.value = { ...streaming.value, [sid]: false }
        state.value = { ...state.value, [sid]: 'idle' }
        break
      case 'done':
        streaming.value = { ...streaming.value, [sid]: false }
        state.value = { ...state.value, [sid]: 'idle' }
        break
    }
  }

  return { list, activeId, active, messages, streaming, state, pending, toolBlocks, refresh, create, select, send, handleEvent }
})
