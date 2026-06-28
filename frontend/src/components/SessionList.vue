<template>
  <div class="session-list">
    <div class="header">
      <span class="label">会话</span>
      <button class="add" @click="$emit('create')" title="新建会话">+</button>
    </div>
    <div class="items">
      <SessionItem
        v-for="s in list"
        :key="s.id"
        :meta="s"
        :is-active="s.id === activeId"
        :dup="dupProjects.has(projectName(s))"
        @select="$emit('select', s.id)"
      />
      <div v-if="!list.length" class="empty">暂无会话</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SessionItem from './SessionItem.vue'
import type { SessionMeta } from '../types/session'

const props = defineProps<{ list: SessionMeta[]; activeId: string | null }>()
defineEmits<{
  (e: 'create'): void
  (e: 'select', id: string): void
}>()

// 同名 project 检测，用于 SessionItem 显示前缀区分
const dupProjects = computed(() => {
  const counts: Record<string, number> = {}
  for (const s of props.list) {
    const name = projectName(s)
    counts[name] = (counts[name] || 0) + 1
  }
  return new Set(Object.keys(counts).filter((k) => counts[k] > 1))
})

function projectName(s: SessionMeta): string {
  const wd = s.workdir
  if (!wd || wd === '/') return wd
  return wd.split('/').filter(Boolean).pop() || wd
}
</script>

<style scoped>
.session-list { display: flex; flex-direction: column; flex: 1; min-height: 0; }
.header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 12px; border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.label { font-size: 11px; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.8px; font-weight: 600; }
.add {
  width: 22px; height: 22px; border-radius: var(--radius-md);
  background: var(--bg-input); color: var(--accent-light);
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 600;
}
.add:hover { background: var(--accent); color: white; }
.items { flex: 1; overflow-y: auto; padding: 6px; min-height: 0; }
.empty { color: var(--text-tertiary); font-size: 11px; text-align: center; padding: 20px; }
</style>
