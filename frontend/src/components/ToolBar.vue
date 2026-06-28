<template>
  <div class="toolbar">
    <div class="main">
      <span class="status" :class="state"></span>
      <span class="title">{{ displayTitle }}</span>
    </div>
    <div class="meta-row">
      <span class="item mono session-id" :title="sessionId">#{{ shortId }}</span>
      <span class="sep">·</span>
      <span class="item mono">{{ projectName }}</span>
      <span class="sep">·</span>
      <span class="item">{{ msgCount }} 条消息</span>
      <span class="sep">·</span>
      <span class="item">{{ formatSize }}</span>
      <button class="term-btn" @click="$emit('open-terminal')" title="在终端中打开">
        ›_
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  aiTitle?: string
  path: string
  sessionId: string
  msgCount: number
  size: number
  state: 'idle' | 'running' | 'awaiting_permission'
}>()
defineEmits<{ (e: 'open-terminal'): void }>()

const displayTitle = computed(() => (props as any).aiTitle || props.title || '新会话')
const shortId = computed(() => props.sessionId?.slice(0, 8) || '')
const projectName = computed(() => {
  const wd = props.path
  if (!wd || wd === '/') return wd
  return wd.split('/').filter(Boolean).pop() || wd
})
const formatSize = computed(() => {
  const s = props.size
  if (s < 1024) return `${s} B`
  if (s < 1024 * 1024) return `${(s / 1024).toFixed(1)} KB`
  return `${(s / 1024 / 1024).toFixed(2)} MB`
})
</script>

<style scoped>
.toolbar {
  padding: 8px 14px; border-bottom: 1px solid var(--border);
  background: var(--bg-panel); flex-shrink: 0;
  max-height: 52px;
}
.main { display: flex; align-items: center; gap: 8px; max-width: 100%; }
.status { width: 6px; height: 6px; border-radius: 50%; background: var(--text-tertiary); flex-shrink: 0; }
.status.running { background: var(--status-success); }
.status.awaiting_permission { background: var(--status-warn); }
.title {
  font-size: 12px; color: var(--text-primary); font-weight: 500;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.meta-row {
  display: flex; align-items: center; gap: 6px; margin-top: 3px;
}
.item { font-size: 10px; color: var(--text-tertiary); }
.mono { font-family: var(--font-mono); }
.session-id { color: var(--accent-light); }
.sep { color: var(--text-tertiary); margin: 0 -2px; }
.term-btn {
  font-size: 10px; padding: 1px 6px; margin-left: auto;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-tertiary);
  font-family: var(--font-mono); cursor: pointer;
}
.term-btn:hover { background: var(--bg-input); color: var(--text-secondary); }
</style>
