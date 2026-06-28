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

const PAGE_SIZE = 100

export const useSessionsStore = defineStore('sessions', () => {
  const list = ref<SessionMeta[]>([])
  const activeId = ref<string | null>(null)
  const messages = ref<Record<string, ChatMessage[]>>({})
  const streaming = ref<Record<string, boolean>>({})
  const state = ref<Record<string, SessionState>>({})
  const pending = ref<Record<string, PendingPerm | null>>({})
  const toolBlocks = ref<Record<string, Array<{ name: string; args: unknown }>>>({})
  const historyOffset = ref<Record<string, number>>({})
  const hasMore = ref<Record<string, boolean>>({})

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
      // 默认加载最后 PAGE_SIZE 条
      const total = meta.msg_count
      loadHistory(id, meta.workdir, Math.max(0, total - PAGE_SIZE), PAGE_SIZE, true)
    }
  }

  async function loadHistory(sid: string, workdir: string, offset: number, limit: number, isFirst: boolean) {
    try {
      const raw = await GetSessionMessages(sid, workdir, offset, limit)
      const msgs = (raw || []).map((m: any, i: number) => ({
        id: `${sid}-${offset + i}`,
        role: m.role || m.Role || 'assistant',
        content: formatContent(m.content || m.Content || ''),
        ts: Date.now() - ((raw?.length || 0) - i) * 1000,
      }))
      if (isFirst) {
        messages.value = { ...messages.value, [sid]: msgs }
      } else {
        const prev = messages.value[sid] || []
        messages.value = { ...messages.value, [sid]: [...msgs, ...prev] }
      }
      historyOffset.value = { ...historyOffset.value, [sid]: offset + (raw?.length || 0) }
      hasMore.value = { ...hasMore.value, [sid]: offset > 0 }
    } catch (e: any) {
      console.error('[sessions] loadHistory failed:', e?.message || e)
    }
  }

  async function loadMore() {
    const id = activeId.value
    if (!id) return
    const meta = list.value.find((s) => s.id === id)
    if (!meta) return
    const loaded = historyOffset.value[id] || meta.msg_count
    const nextStart = Math.max(0, loaded - PAGE_SIZE)
    if (nextStart >= loaded) return
    await loadHistory(id, meta.workdir, nextStart, loaded - nextStart, false)
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

  return { list, activeId, active, messages, streaming, state, pending, toolBlocks,
    hasMore, refresh, create, select, send, handleEvent, loadMore }
})
