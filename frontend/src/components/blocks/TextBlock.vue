<template>
  <div class="text-block markdown-body" v-html="html" />
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
.text-block { color: var(--text-primary); }
.text-block :deep(p) { margin: 6px 0; }
.text-block :deep(p:first-child) { margin-top: 0; }
.text-block :deep(p:last-child) { margin-bottom: 0; }
.text-block :deep(ul), .text-block :deep(ol) { padding-left: 22px; margin: 6px 0; }
.text-block :deep(li) { margin: 2px 0; }
.text-block :deep(blockquote) {
  border-left: 2px solid var(--border);
  padding-left: 10px;
  margin: 8px 0;
  color: var(--text-secondary);
}
.text-block :deep(a) { color: var(--accent-light); text-decoration: none; }
.text-block :deep(a:hover) { text-decoration: underline; }
.text-block :deep(h1) { font-size: 17px; margin: 10px 0 4px; }
.text-block :deep(h2) { font-size: 15px; margin: 8px 0 4px; }
.text-block :deep(h3) { font-size: 14px; margin: 8px 0 4px; }
.text-block :deep(table) {
  border-collapse: collapse; margin: 6px 0; width: 100%;
}
.text-block :deep(th), .text-block :deep(td) {
  border: 1px solid var(--border); padding: 4px 8px; font-size: 12px; text-align: left;
}
.text-block :deep(th) { background: var(--bg-input); }
.text-block :deep(pre) {
  background: var(--bg-primary); border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 10px 12px; margin: 8px 0; overflow-x: auto; font-size: 13px; line-height: 1.5;
}
.text-block :deep(code) {
  background: var(--bg-input); padding: 1px 5px; border-radius: 3px;
  font-size: 13px; font-family: var(--font-mono); color: var(--accent-light);
}
.text-block :deep(pre code) { background: transparent; padding: 0; color: var(--text-primary); }
.text-block :deep(hr) { border: none; border-top: 1px solid var(--border); margin: 10px 0; }
.text-block :deep(strong) { font-weight: 600; }
</style>
