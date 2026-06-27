<template>
  <div class="bubble" :class="role">
    <div class="role-label">{{ roleLabel }}</div>
    <div class="content">
      <pre v-if="role === 'tool'" class="code"><code>{{ content }}</code></pre>
      <div v-else class="text">{{ content }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ role: 'user' | 'assistant' | 'tool'; content: string }>()
const roleLabel = computed(() => ({
  user: '你',
  assistant: 'Claude',
  tool: '工具',
}[props.role]))
</script>

<style scoped>
.bubble { padding: 8px 0; }
.role-label { font-size: 10px; color: var(--text-tertiary); text-transform: uppercase; letter-spacing: 0.6px; margin-bottom: 4px; }
.content { font-size: 12px; line-height: 1.6; }
.text { white-space: pre-wrap; word-wrap: break-word; }
.code {
  background: var(--bg-input); border: 1px solid var(--border);
  border-radius: var(--radius-md); padding: 8px 10px;
  font-size: var(--font-size-code);
  overflow-x: auto;
}
.user .text { color: var(--text-primary); }
.assistant .text { color: var(--text-primary); }
.tool .code { color: var(--accent-light); }
</style>
