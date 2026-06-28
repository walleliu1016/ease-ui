<template>
  <div class="tool-edit">
    <div class="file-header edit-header">
      <span class="file-icon">✏️</span> Edit
      <span class="file-path">{{ filename }}</span>
      <span v-if="replaceAll" class="replace-all">(replace all)</span>
    </div>
    <div class="file-fullpath">{{ filePath }}</div>
    <Truncatable :max-height="320">
      <div class="edit-section edit-old">
        <div class="edit-label">−</div>
        <pre class="edit-content"><code>{{ oldString }}</code></pre>
      </div>
      <div class="edit-section edit-new">
        <div class="edit-label">+</div>
        <pre class="edit-content"><code>{{ newString }}</code></pre>
      </div>
    </Truncatable>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Truncatable from '../Truncatable.vue'

const props = defineProps<{ filePath: string; oldString: string; newString: string; replaceAll?: boolean }>()
const filename = computed(() => {
  const p = props.filePath
  const i = Math.max(p.lastIndexOf('/'), p.lastIndexOf('\\'))
  return i >= 0 ? p.slice(i + 1) : p
})
</script>

<style scoped>
.tool-edit {
  background: var(--edit-bg);
  border: 1px solid var(--edit-border);
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
.edit-header { color: var(--edit-header); }
.file-path {
  font-family: var(--font-mono);
  background: rgba(0,0,0,0.06);
  padding: 1px 6px;
  border-radius: 3px;
  font-size: 12px;
}
.replace-all {
  font-size: 11px;
  font-weight: normal;
  color: var(--text-tertiary);
}
.file-fullpath {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  word-break: break-all;
}
.edit-section {
  display: flex;
  margin: 4px 0;
  border-radius: 4px;
  overflow: hidden;
  background: var(--code-bg);
}
.edit-label {
  padding: 6px 10px;
  font-weight: bold;
  font-family: var(--font-mono);
  display: flex;
  align-items: flex-start;
  background: rgba(0,0,0,0.2);
  color: var(--text-tertiary);
}
.edit-old .edit-label { color: var(--edit-old-label); }
.edit-new .edit-label { color: var(--edit-new-label); }
.edit-content {
  margin: 0;
  flex: 1;
  background: transparent;
  font-size: 12px;
  font-family: var(--font-mono);
  padding: 6px 10px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--text-primary);
}
</style>
