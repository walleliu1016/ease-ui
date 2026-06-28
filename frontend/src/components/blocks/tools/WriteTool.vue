<template>
  <div class="tool-write">
    <div class="file-header write-header">
      <span class="file-icon">📝</span> Write
      <span class="file-path">{{ filename }}</span>
    </div>
    <div class="file-fullpath">{{ filePath }}</div>
    <Truncatable :max-height="320">
      <pre class="file-content"><code>{{ content }}</code></pre>
    </Truncatable>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Truncatable from '../Truncatable.vue'

const props = defineProps<{ filePath: string; content: string }>()
const filename = computed(() => {
  const p = props.filePath
  const i = Math.max(p.lastIndexOf('/'), p.lastIndexOf('\\'))
  return i >= 0 ? p.slice(i + 1) : p
})
</script>

<style scoped>
.tool-write {
  background: var(--write-bg);
  border: 1px solid var(--write-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 6px 0;
}
.file-header {
  font-weight: 600;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.write-header { color: var(--write-header); }
.file-path {
  font-family: var(--font-mono);
  background: rgba(0,0,0,0.06);
  padding: 1px 6px;
  border-radius: 3px;
  font-size: 12px;
}
.file-fullpath {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  word-break: break-all;
}
.file-content {
  background: var(--code-bg);
  color: var(--code-text);
  padding: 8px 10px;
  border-radius: 4px;
  margin: 0;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.5;
  font-family: var(--font-mono);
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
