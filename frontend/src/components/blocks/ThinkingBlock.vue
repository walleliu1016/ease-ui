<template>
  <div class="thinking-block">
    <div class="thinking-label">Thinking</div>
    <div class="thinking-body markdown-body" v-html="html" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

marked.setOptions({ breaks: true, gfm: true })

const props = defineProps<{ text: string }>()
const html = computed(() => {
  try { return marked.parse(props.text) as string } catch { return props.text }
})
</script>

<style scoped>
.thinking-block {
  background: var(--thinking-bg);
  border: 1px solid var(--thinking-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 8px 0;
  font-size: 13px;
  color: var(--text-secondary);
}
.thinking-label {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.6px;
  color: var(--thinking-label);
  margin-bottom: 6px;
}
.thinking-body :deep(p) { margin: 4px 0; }
.thinking-body :deep(p:first-child) { margin-top: 0; }
.thinking-body :deep(p:last-child) { margin-bottom: 0; }
.thinking-body :deep(ul), .thinking-body :deep(ol) { padding-left: 20px; margin: 4px 0; }
.thinking-body :deep(code) {
  background: rgba(0,0,0,0.08);
  padding: 1px 4px;
  border-radius: 3px;
  font-size: 12px;
  font-family: var(--font-mono);
}
</style>
