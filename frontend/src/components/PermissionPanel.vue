<template>
  <div class="panel">
    <div class="head">
      <span class="warn">⚠</span>
      <span class="title">权限请求</span>
    </div>
    <div class="body">
      <div class="row"><span class="k">工具</span><span class="v mono">{{ tool }}</span></div>
      <pre class="args"><code>{{ formattedArgs }}</code></pre>
    </div>
    <div class="actions">
      <button class="deny" @click="$emit('respond', false)">拒绝</button>
      <button class="allow" @click="$emit('respond', true)">允许</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ tool: string; args: unknown }>()
defineEmits<{ (e: 'respond', allow: boolean): void }>()
const formattedArgs = computed(() => {
  try { return JSON.stringify(props.args, null, 2) } catch { return String(props.args) }
})
</script>

<style scoped>
.panel {
  background: var(--bg-panel);
  border: 1px solid var(--status-warn);
  border-radius: var(--radius-md); padding: 10px 12px; margin: 8px 0;
}
.head { display: flex; align-items: center; gap: 6px; margin-bottom: 6px; }
.warn { color: var(--status-warn); }
.title { font-size: 11px; color: var(--status-warn); font-weight: 600; text-transform: uppercase; letter-spacing: 0.6px; }
.body { margin-bottom: 8px; }
.row { display: flex; gap: 8px; font-size: 11px; }
.k { color: var(--text-tertiary); }
.v { color: var(--text-primary); }
.mono { font-family: var(--font-mono); }
.args { background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 6px 8px; font-size: var(--font-size-code); margin-top: 4px; overflow-x: auto; }
.actions { display: flex; gap: 8px; justify-content: flex-end; }
.deny, .allow {
  padding: 5px 12px; border-radius: var(--radius-md); font-size: 11px; font-weight: 500;
}
.deny { background: var(--bg-input); color: var(--text-primary); border: 1px solid var(--border); }
.deny:hover { background: var(--status-error); color: white; border-color: var(--status-error); }
.allow { background: var(--accent); color: white; }
.allow:hover { background: var(--accent-light); }
</style>
