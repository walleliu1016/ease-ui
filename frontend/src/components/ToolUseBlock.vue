<template>
  <div class="tool" :class="toolClass">
    <div class="head">
      <span class="icon">{{ toolIcon }}</span>
      <span class="name">{{ name }}</span>
    </div>
    <pre class="args" ref="argsEl"><code>{{ formattedArgs }}</code></pre>
    <button v-if="truncated" class="expand" @click="expanded = !expanded">
      {{ expanded ? '收起' : '展开' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
const props = defineProps<{ name: string; args: unknown }>()

const argsEl = ref<HTMLElement | null>(null)
const truncated = ref(false)
const expanded = ref(false)

onMounted(() => {
  if (argsEl.value && argsEl.value.scrollHeight > 200) {
    truncated.value = true
  }
})

const formattedArgs = computed(() => {
  try { return JSON.stringify(props.args, null, 2) } catch { return String(props.args) }
})

const toolIcon = computed(() => {
  switch (props.name?.toLowerCase()) {
    case 'bash': return '>_'
    case 'write': return '✎'
    case 'edit': return '✐'
    case 'read': return '☰'
    case 'grep': return '⌕'
    case 'glob': return '◉'
    case 'task': return '☑'
    case 'todowrite': return '☑'
    default: return '⚙'
  }
})

const toolClass = computed(() => props.name?.toLowerCase() || '')
</script>

<style scoped>
.tool {
  background: var(--bg-input); border: 1px solid var(--border);
  border-left: 2px solid var(--accent);
  border-radius: var(--radius-md); padding: 8px 10px; margin: 6px 0;
}
.tool :deep(pre) { max-height: 200px; overflow: hidden; }
.tool :deep(pre.expanded) { max-height: none; }
.head { display: flex; align-items: center; gap: 6px; margin-bottom: 4px; }
.icon { font-size: 11px; color: var(--accent-light); font-family: var(--font-mono); font-weight: 700; }
.name { font-size: 11px; color: var(--accent-light); font-weight: 600; font-family: var(--font-mono); }
.args {
  font-size: var(--font-size-code); color: var(--text-primary); overflow-x: auto;
  background: var(--bg-primary); padding: 6px 8px; border-radius: 4px;
  max-height: 200px; overflow-y: auto;
}
.args.expanded { max-height: none; }
.expand {
  margin-top: 4px; font-size: 10px; color: var(--accent-light);
  background: transparent; border: none; cursor: pointer;
}
.expand:hover { color: var(--accent); }

/* 工具类型样式 */
.tool.bash { border-left-color: #9c27b0; }
.tool.write { border-left-color: #2196f3; }
.tool.edit { border-left-color: #ff9800; }
.tool.read { border-left-color: #4caf50; }
.tool.grep { border-left-color: #00bcd4; }
.tool.glob { border-left-color: #795548; }
.tool.task, .tool.todowrite { border-left-color: #8bc34a; }
</style>
