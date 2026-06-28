<template>
  <div class="truncatable" :class="state">
    <div class="truncatable-content" ref="contentEl">
      <slot />
    </div>
    <button v-if="overflow" class="expand-btn" @click="toggle">
      {{ state === 'truncated' ? 'Show more' : 'Show less' }}
    </button>
  </div>
</template>

<script setup lang="ts">
// 复刻参考实现：内容 > maxHeight 时折叠，Show more / Show less 切换。
// 使用 ResizeObserver 响应内容动态变化（streaming / 重渲染）。
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = withDefaults(defineProps<{
  maxHeight?: number
  /** 初始就视为已折叠（默认 true：让父组件用 v-if 控制 mounted 时机） */
  initiallyCollapsed?: boolean
}>(), { maxHeight: 200, initiallyCollapsed: true })

const contentEl = ref<HTMLElement | null>(null)
const state = ref<'truncated' | 'expanded' | 'idle'>('idle')
const overflow = ref(false)
let ro: ResizeObserver | null = null

function measure() {
  const el = contentEl.value
  if (!el) return
  const h = el.scrollHeight
  overflow.value = h > props.maxHeight
  if (state.value === 'idle' && overflow.value) state.value = 'truncated'
}

function toggle() {
  state.value = state.value === 'truncated' ? 'expanded' : 'truncated'
}

onMounted(() => {
  nextTick(measure)
  if (typeof ResizeObserver !== 'undefined' && contentEl.value) {
    ro = new ResizeObserver(() => measure())
    ro.observe(contentEl.value)
  }
})

onBeforeUnmount(() => { ro?.disconnect() })
</script>

<style scoped>
.truncatable { position: relative; }
.truncatable.truncated .truncatable-content {
  max-height: v-bind('props.maxHeight + "px"');
  overflow: hidden;
}
.truncatable.truncated::after {
  content: '';
  position: absolute;
  left: 0; right: 0;
  bottom: 36px;
  height: 50px;
  background: linear-gradient(to bottom, transparent, var(--bg-panel));
  pointer-events: none;
}
.expand-btn {
  display: block;
  width: 100%;
  margin-top: 4px;
  padding: 6px 10px;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-tertiary);
  font-size: 11px;
  cursor: pointer;
}
.expand-btn:hover { color: var(--text-secondary); border-color: var(--text-tertiary); }
</style>
