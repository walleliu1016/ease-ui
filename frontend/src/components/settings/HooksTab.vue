<template>
  <div class="hooks-tab">
    <div class="path">~/.claude/settings.json</div>

    <section v-for="ev in events" :key="ev" class="section">
      <div class="section-head">
        <h3>{{ ev }}</h3>
        <button class="add" @click="addEntry(ev)">+ 添加</button>
      </div>
      <div v-if="!cfg?.[ev]?.length" class="empty">无</div>
      <div v-for="(entry, i) in cfg?.[ev] ?? []" :key="i" class="entry">
        <div class="row">
          <span class="k">matcher</span>
          <input class="v" v-model="entry.matcher" placeholder="(optional)" />
        </div>
        <div class="row">
          <span class="k">command</span>
          <input class="v mono" v-model="entry.command" />
        </div>
        <div class="row">
          <span class="k">type</span>
          <select class="v" v-model="entry.type">
            <option value="shell">shell</option>
            <option value="python">python</option>
          </select>
          <button class="del" @click="removeEntry(ev, i)">×</button>
        </div>
      </div>
    </section>

    <div class="actions">
      <button class="save" :disabled="!hooks.dirty" @click="onSave">保存</button>
      <button class="cancel" :disabled="!hooks.dirty" @click="hooks.load">取消</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useHooksStore } from '../../stores/hooks'
import type { HookEntry } from '../../types/hooks'

const hooks = useHooksStore()
const cfg = computed(() => hooks.cfg)
const events = ['PreToolUse', 'PermissionRequest', 'PostToolUse', 'Notification', 'Stop'] as const

onMounted(() => hooks.load())

function addEntry(ev: typeof events[number]) {
  if (!hooks.cfg) return
  const entry: HookEntry = { command: '', type: 'shell' }
  if (!hooks.cfg[ev]) hooks.cfg[ev] = []
  hooks.cfg[ev].push(entry)
  hooks.markDirty()
}

function removeEntry(ev: typeof events[number], i: number) {
  if (!hooks.cfg) return
  hooks.cfg[ev].splice(i, 1)
  hooks.markDirty()
}

async function onSave() {
  try {
    await hooks.save()
  } catch (e: any) {
    alert('保存失败：' + (e?.message ?? e))
  }
}
</script>

<style scoped>
.hooks-tab { padding: 20px 24px; max-width: 720px; }
.path { font-family: var(--font-mono); font-size: 11px; color: var(--text-tertiary); margin-bottom: 16px; }
.section { margin-bottom: 18px; padding: 12px 14px; background: var(--bg-panel); border: 1px solid var(--border); border-radius: var(--radius-md); }
.section-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
h3 { font-size: 12px; color: var(--accent-light); font-weight: 600; }
.add { padding: 3px 10px; background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-sm); color: var(--text-primary); font-size: 11px; }
.add:hover { background: var(--border); }
.entry { padding: 8px; background: var(--bg-input); border-radius: var(--radius-sm); margin-bottom: 6px; }
.row { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.k { font-size: 10px; color: var(--text-tertiary); width: 60px; }
.v { flex: 1; background: var(--bg-panel); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 4px 8px; color: var(--text-primary); font-size: 11px; font-family: inherit; }
.v:focus { outline: none; border-color: var(--accent); }
.mono { font-family: var(--font-mono); }
.del { width: 22px; height: 22px; color: var(--text-secondary); border-radius: var(--radius-sm); }
.del:hover { background: var(--status-error); color: white; }
.empty { color: var(--text-tertiary); font-size: 11px; padding: 8px; }
.actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 20px; }
.save { padding: 6px 16px; background: var(--accent); color: white; border-radius: var(--radius-md); }
.save:hover:not(:disabled) { background: var(--accent-light); }
.cancel { padding: 6px 16px; background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-md); }
.save:disabled, .cancel:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
