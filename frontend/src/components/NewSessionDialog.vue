<template>
  <div v-if="open" class="overlay" @click.self="$emit('close')">
    <div class="dialog">
      <div class="head">
        <h2>新建会话</h2>
        <button class="close" @click="$emit('close')">✕</button>
      </div>
      <form @submit.prevent="onSubmit" class="form">
        <div class="form-group">
          <label class="form-label">工作目录</label>
          <div class="dir-row">
            <input class="form-input" v-model="workdir" placeholder="点击右侧按钮选择目录" readonly @click="onPick" />
            <button type="button" class="pick-btn" @click="onPick">选择…</button>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">提示词</label>
          <textarea class="form-input area" v-model="prompt" rows="4" placeholder="你想让 Claude 做什么？"></textarea>
        </div>
        <div class="form-actions">
          <button type="button" class="cancel" @click="$emit('close')">取消</button>
          <button type="submit" class="primary" :disabled="!workdir.trim() || !prompt.trim()">创建</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { PickDirectory } from '../composables/useWails'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', workdir: string, prompt: string): void
}>()

const workdir = ref('')
const prompt = ref('')

async function onPick() {
  try {
    const dir = await PickDirectory()
    if (dir) workdir.value = dir
  } catch {}
}

function onSubmit() {
  if (!workdir.value.trim() || !prompt.value.trim()) return
  emit('create', workdir.value.trim(), prompt.value.trim())
  workdir.value = ''
  prompt.value = ''
  emit('close')
}
</script>

<style scoped>
.overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000;
}
.dialog {
  width: 480px; background: var(--bg-panel);
  border: 1px solid var(--border); border-radius: var(--radius-lg);
  box-shadow: var(--shadow-window);
  padding: 20px 24px;
}
.head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
h2 { font-size: 14px; color: var(--text-primary); }
.close { color: var(--text-secondary); font-size: 14px; padding: 2px 6px; border-radius: var(--radius-sm); }
.close:hover { background: var(--bg-input); color: var(--text-primary); }
.form-group { margin-bottom: 12px; }
.form-label { display: block; font-size: 10px; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.6px; margin-bottom: 4px; }
.dir-row { display: flex; gap: 8px; }
.dir-row .form-input { flex: 1; cursor: pointer; }
.form-input {
  width: 100%; background: var(--bg-input); border: 1px solid var(--border);
  border-radius: var(--radius-md); padding: 8px 10px;
  color: var(--text-primary); font-size: 12px; font-family: inherit;
}
.form-input:focus { outline: none; border-color: var(--accent); }
.pick-btn {
  padding: 8px 14px; background: var(--bg-input); border: 1px solid var(--border);
  border-radius: var(--radius-md); color: var(--text-primary); font-size: 12px;
  white-space: nowrap; cursor: pointer;
}
.pick-btn:hover { background: var(--border); }
.area { resize: vertical; min-height: 80px; }
.form-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 18px; }
.cancel { padding: 7px 16px; background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-md); font-size: 12px; color: var(--text-primary); }
.cancel:hover { background: var(--border); }
.primary { padding: 7px 18px; background: var(--accent); color: white; border-radius: var(--radius-md); font-size: 12px; font-weight: 500; }
.primary:hover:not(:disabled) { background: var(--accent-deep); }
.primary:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
