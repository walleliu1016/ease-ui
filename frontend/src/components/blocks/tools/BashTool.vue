<template>
  <div class="tool-bash">
    <div class="tool-header"><span class="tool-icon">$</span> Bash</div>
    <div v-if="description" class="tool-description">{{ description }}</div>
    <Truncatable :max-height="280">
      <pre class="bash-command"><code>{{ command }}</code></pre>
    </Truncatable>
  </div>
</template>

<script setup lang="ts">
import Truncatable from '../Truncatable.vue'
defineProps<{ command: string; description?: string }>()
</script>

<style scoped>
.tool-bash {
  background: var(--tool-bg);
  border: 1px solid var(--tool-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 6px 0;
}
.tool-header {
  font-weight: 600;
  color: var(--tool-header);
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-family: var(--font-mono);
  margin-bottom: 6px;
}
.tool-icon {
  font-size: 13px;
  color: var(--accent);
  font-weight: 700;
}
.tool-description {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  font-style: italic;
}
.bash-command {
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
