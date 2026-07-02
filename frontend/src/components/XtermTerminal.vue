<template>
  <div ref="terminalEl" class="xterm-container" @click="focusTerminal" />
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'
import { EventsOn, ResizeTerminal } from '../composables/useWails'

const props = defineProps<{
  sessionId: string
}>()

const emit = defineEmits<{
  (e: 'data', data: string): void
}>()

const terminalEl = ref<HTMLElement | null>(null)

let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let cleanupEvents: (() => void) | null = null
let resizeObserver: ResizeObserver | null = null

function focusTerminal() {
  term?.focus()
}

function calcCols(width: number): number {
  return Math.max(40, Math.floor((width - 16) / 8.4))
}
function calcRows(height: number): number {
  return Math.max(10, Math.floor((height - 8) / 17))
}

onMounted(() => {
  if (!terminalEl.value) return

  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Consolas, "Courier New", monospace',
    theme: { background: '#1e1e1e', foreground: '#d4d4d4' },
    allowProposedApi: true,
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())

  term.open(terminalEl.value)

  nextTick(() => {
    fitAddon?.fit()
    emitResize()
  })

  resizeObserver = new ResizeObserver(() => {
    fitAddon?.fit()
    emitResize()
  })
  resizeObserver.observe(terminalEl.value)

  term.onData((data: string) => {
    emit('data', data)
  })

  const topic = `session:${props.sessionId}`
  cleanupEvents = EventsOn(topic, (line: string) => {
    if (line === '{"type":"done"}') return
    term?.write(line)
  })
})

function emitResize() {
  if (!term || !terminalEl.value) return
  const cols = calcCols(terminalEl.value.clientWidth)
  const rows = calcRows(terminalEl.value.clientHeight)
  term.resize(cols, rows)
  ResizeTerminal(props.sessionId, cols, rows).catch(() => {})
}

watch(() => props.sessionId, () => {
  cleanupEvents?.()
  term?.dispose()
})

onBeforeUnmount(() => {
  cleanupEvents?.()
  resizeObserver?.disconnect()
  term?.dispose()
})
</script>

<style scoped>
.xterm-container {
  width: 100%;
  height: 100%;
  min-height: 0;
}
</style>
