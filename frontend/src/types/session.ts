export type SessionState = 'idle' | 'running' | 'awaiting_permission'

export interface SessionMeta {
  id: string
  workdir: string
  mtime: number
  msg_count: number
  first_prompt: string
  ai_title: string
  size: number
}

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant' | 'tool'
  content: string
  tool?: { name: string; args: unknown }
  ts: number
}
