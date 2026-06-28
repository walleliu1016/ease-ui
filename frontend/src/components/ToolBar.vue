<template>
  <div class="toolbar">
    <div class="left">
      <span class="status" :class="state"></span>
      <span class="title">{{ displayTitle }}</span>
      <span class="sep">·</span>
      <span class="meta mono">{{ projectName }}</span>
      <span class="sep">·</span>
      <span class="meta">{{ msgCount }} 条消息</span>
    </div>
    <button class="term-btn" @click="$emit('open-terminal')" title="在终端中打开会话">
      打开终端
    </button>
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
  state: 'idle' | 'running' | 'awaiting_permission'
}>()
defineEmits<{ (e: 'open-terminal'): void }>()

const displayTitle = computed(() => (props as any).aiTitle || props.title || '新会话')
const projectName = computed(() => {
  const wd = props.path
  if (!wd || wd === '/') return wd
  return wd.split('/').filter(Boolean).pop() || wd
})
</script>

<style scoped>
.toolbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 12px; border-bottom: 1px solid var(--border);
  background: var(--bg-panel); flex-shrink: 0; gap: 12px; height: 38px;
}
.left { display: flex; align-items: center; gap: 7px; min-width: 0; overflow: hidden; }
.status { width: 7px; height: 7px; border-radius: 50%; background: var(--text-tertiary); flex-shrink: 0; }
.status.running { background: var(--status-success); }
.status.awaiting_permission { background: var(--status-warn); }
.title {
  font-size: 12px; color: var(--text-primary); font-weight: 500;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.sep { color: var(--border); flex-shrink: 0; }
.meta { font-size: 11px; color: var(--text-tertiary); white-space: nowrap; }
.mono { font-family: var(--font-mono); }
.term-btn {
  flex-shrink: 0;
  padding: 4px 12px;
  background: var(--accent); color: white;
  border: none; border-radius: var(--radius-md);
  font-size: 11px; font-weight: 500; cursor: pointer;
  white-space: nowrap;
}
.term-btn:hover { background: var(--accent-deep); }
</style>
