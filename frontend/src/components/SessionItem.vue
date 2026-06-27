<template>
  <div
    class="session-item"
    :class="{ active: isActive }"
    @click="$emit('select')"
    @mouseenter="showTip = true"
    @mouseleave="showTip = false"
  >
    <div class="status-dot" :class="state"></div>
    <div class="title">
      <div class="prompt">{{ displayName }}</div>
      <div class="meta">{{ displayTime }}</div>
    </div>
    <SessionTooltip v-if="showTip" :meta="meta" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import SessionTooltip from './SessionTooltip.vue'
import type { SessionMeta } from '../types/session'

const props = defineProps<{ meta: SessionMeta; isActive: boolean }>()
defineEmits<{ (e: 'select'): void }>()

const showTip = ref(false)

const displayName = computed(() => props.meta.first_prompt || '新会话')
const displayTime = computed(() => {
  const d = new Date(props.meta.mtime * 1000)
  return d.toLocaleString('zh-CN', { hour12: false })
})
const state = computed(() => 'idle')
</script>

<style scoped>
.session-item {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 10px; border-radius: var(--radius-md);
  cursor: pointer; position: relative;
}
.session-item:hover { background: rgba(255,255,255,0.04); }
.session-item.active { background: rgba(124,58,237,0.15); }
.status-dot {
  width: 6px; height: 6px; border-radius: 50%;
  background: var(--text-tertiary); flex-shrink: 0;
}
.status-dot.running { background: var(--status-success); }
.status-dot.awaiting_permission { background: var(--status-warn); }
.title { flex: 1; min-width: 0; }
.prompt {
  font-size: 12px; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.meta { font-size: 10px; color: var(--text-tertiary); margin-top: 2px; }
</style>
