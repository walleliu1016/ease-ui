<template>
  <div class="composer">
    <textarea
      ref="inputEl"
      v-model="text"
      class="input"
      :rows="1"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="autoResize"
      @keydown.enter.exact.prevent="onSend"
    />
    <button class="send" :disabled="!text.trim() || disabled" @click="onSend">发送</button>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
const props = defineProps<{ disabled?: boolean; placeholder?: string }>()
const emit = defineEmits<{ (e: 'send', text: string): void }>()

const text = ref('')
const inputEl = ref<HTMLTextAreaElement | null>(null)
const placeholder = props.placeholder ?? '输入消息…'

function autoResize() {
  nextTick(() => {
    const el = inputEl.value
    if (!el) return
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 120) + 'px'
  })
}

function onSend() {
  const t = text.value.trim()
  if (!t || props.disabled) return
  emit('send', t)
  text.value = ''
  autoResize()
}
</script>

<style scoped>
.composer {
  display: flex; align-items: center; gap: 8px; padding: 0 12px;
  background: var(--bg-panel); border-top: 1px solid var(--border);
  flex-shrink: 0; height: 46px;
}
.input {
  flex: 1; resize: none; overflow: hidden;
  background: var(--bg-input); border: 1px solid var(--border);
  border-radius: var(--radius-md); padding: 5px 10px;
  color: var(--text-primary); font-size: 12px;
  font-family: inherit; line-height: 1.5;
  min-height: 28px; max-height: 120px;
}
.input:focus { outline: none; border-color: var(--accent); }
.send {
  background: var(--accent); color: white;
  padding: 5px 14px; border-radius: var(--radius-md);
  font-size: 11px; font-weight: 500;
  flex-shrink: 0; height: 26px;
}
.send:hover:not(:disabled) { background: var(--accent-deep); }
.send:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
