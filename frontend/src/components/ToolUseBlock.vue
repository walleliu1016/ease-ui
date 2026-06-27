<template>
  <div class="tool">
    <div class="head">
      <span class="dot" />
      <span class="name">{{ name }}</span>
    </div>
    <pre class="args"><code>{{ formattedArgs }}</code></pre>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ name: string; args: unknown }>()
const formattedArgs = computed(() => {
  try { return JSON.stringify(props.args, null, 2) } catch { return String(props.args) }
})
</script>

<style scoped>
.tool {
  background: var(--bg-input); border: 1px solid var(--border);
  border-left: 2px solid var(--accent);
  border-radius: var(--radius-md); padding: 8px 10px; margin: 6px 0;
}
.head { display: flex; align-items: center; gap: 6px; margin-bottom: 4px; }
.dot { width: 6px; height: 6px; border-radius: 50%; background: var(--status-warn); }
.name { font-size: 11px; color: var(--accent-light); font-weight: 600; font-family: var(--font-mono); }
.args { font-size: var(--font-size-code); color: var(--text-primary); overflow-x: auto; }
</style>
