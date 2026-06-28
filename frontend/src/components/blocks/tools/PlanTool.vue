<template>
  <div class="tool-plan">
    <div class="plan-header"><span class="plan-icon">📋</span> {{ title }}</div>
    <Truncatable v-if="plan" :max-height="400">
      <div class="plan-body markdown-body" v-html="html" />
    </Truncatable>
    <div v-else-if="planFilePath" class="plan-fallback">
      <div class="plan-fallback-msg">Plan 内容在文件中：</div>
      <code class="plan-fallback-path">{{ planFilePath }}</code>
    </div>
    <div v-else class="plan-empty">无 plan 内容（EnterPlanMode 通常为空）</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'
import Truncatable from '../Truncatable.vue'

marked.setOptions({ breaks: true, gfm: true })

const props = defineProps<{ title?: string; plan: string; planFilePath?: string }>()
const title = computed(() => props.title || 'Plan')
const html = computed(() => {
  try { return marked.parse(props.plan) as string } catch { return props.plan }
})
</script>

<style scoped>
.tool-plan {
  background: var(--plan-bg);
  border: 1px solid var(--plan-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 6px 0;
}
.plan-header {
  font-weight: 600;
  color: var(--plan-header);
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.plan-icon { font-size: 14px; }
.plan-body { color: var(--text-primary); font-size: 13px; line-height: 1.5; }
.plan-body :deep(p) { margin: 4px 0; }
.plan-body :deep(p:first-child) { margin-top: 0; }
.plan-body :deep(p:last-child) { margin-bottom: 0; }
.plan-body :deep(ul), .plan-body :deep(ol) { padding-left: 20px; margin: 4px 0; }
.plan-body :deep(li) { margin: 2px 0; }
.plan-body :deep(h1), .plan-body :deep(h2), .plan-body :deep(h3) {
  margin: 8px 0 4px; font-weight: 600; font-size: 14px;
}
.plan-body :deep(code) {
  background: var(--bg-input);
  padding: 1px 4px;
  border-radius: 3px;
  font-size: 12px;
  font-family: var(--font-mono);
}
.plan-fallback {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-tertiary);
}
.plan-fallback-msg { margin-bottom: 4px; }
.plan-fallback-path {
  background: var(--code-bg);
  color: var(--code-text);
  padding: 4px 8px;
  border-radius: 3px;
  font-family: var(--font-mono);
  font-size: 11px;
  display: block;
  word-break: break-all;
}
.plan-empty {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-tertiary);
  font-style: italic;
}
</style>
