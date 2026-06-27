<template>
  <div class="tip">
    <div class="row"><span class="k">Session ID</span><span class="v mono">{{ meta.id }}</span></div>
    <div class="row"><span class="k">工作目录</span><span class="v">{{ meta.workdir }}</span></div>
    <div class="row"><span class="k">修改时间</span><span class="v">{{ formatTime }}</span></div>
    <div class="row"><span class="k">消息数</span><span class="v">{{ meta.msg_count }}</span></div>
    <div class="row"><span class="k">文件大小</span><span class="v">{{ formatSize }}</span></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SessionMeta } from '../types/session'

const props = defineProps<{ meta: SessionMeta }>()

const formatTime = computed(() => {
  return new Date(props.meta.mtime * 1000).toLocaleString('zh-CN', { hour12: false })
})
const formatSize = computed(() => {
  const s = props.meta.size
  if (s < 1024) return `${s} B`
  if (s < 1024 * 1024) return `${(s / 1024).toFixed(1)} KB`
  return `${(s / 1024 / 1024).toFixed(2)} MB`
})
</script>

<style scoped>
.tip {
  position: absolute; left: calc(100% + 8px); top: 0;
  background: var(--bg-panel); border: 1px solid var(--border);
  border-radius: var(--radius-md); padding: 10px 12px;
  min-width: 280px; z-index: 100;
  box-shadow: var(--shadow-window);
}
.row { display: flex; justify-content: space-between; gap: 12px; font-size: 11px; padding: 2px 0; }
.k { color: var(--text-tertiary); }
.v { color: var(--text-primary); text-align: right; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mono { font-family: var(--font-mono); font-size: 10px; }
</style>
