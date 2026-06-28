<template>
  <div class="tool-generic">
    <div class="tool-header"><span class="tool-icon">⚙</span> {{ name }}</div>
    <div v-if="description" class="tool-description">{{ description }}</div>
    <Truncatable :max-height="280">
      <pre class="json" ref="preEl"><code>{{ json }}</code></pre>
    </Truncatable>
  </div>
</template>

<script setup lang="ts">
// 通用 tool_use 卡片：紫色头 + description + 语法高亮 JSON。
// 复刻参考实现的 .tool-use generic 样式，未知工具(Read/Glob/Grep/Skill/WebFetch
// 等没有专门实现的)都走这里，input 整体以 JSON 展示。
import { computed, onMounted, ref, watch } from 'vue'
import Truncatable from '../Truncatable.vue'

const props = defineProps<{ name: string; description?: string; input: Record<string, unknown> }>()
const preEl = ref<HTMLElement | null>(null)

// 移除 description 字段避免重复展示（与参考实现一致）
const displayInput = computed(() => {
  const { description: _d, ...rest } = props.input || {}
  void _d
  return Object.keys(rest).length ? rest : {}
})
const json = computed(() => JSON.stringify(displayInput.value, null, 2))

function highlight() {
  const el = preEl.value
  if (!el) return
  let text = el.textContent || ''
  text = text.replace(/"([^"]+)":/g, '<span style="color:#ce93d8">"$1"</span>:')
  text = text.replace(/: "([^"]*)"/g, ': <span style="color:#81d4fa">"$1"</span>')
  text = text.replace(/: (\d+)/g, ': <span style="color:#ffcc80">$1</span>')
  text = text.replace(/: (true|false|null)/g, ': <span style="color:#f48fb1">$1</span>')
  el.innerHTML = text
}

onMounted(highlight)
watch(json, highlight)
</script>

<style scoped>
.tool-generic {
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
.tool-icon { font-size: 13px; color: var(--accent); }
.tool-description {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  font-style: italic;
}
.json {
  background: var(--code-bg);
  color: #e0e0e0;
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
