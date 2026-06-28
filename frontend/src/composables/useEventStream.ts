import { onMounted, onBeforeUnmount, watch } from 'vue'
import { EventsOn } from './useWails'
import { useSessionsStore } from '../stores/sessions'

export function useEventStream() {
  const sessions = useSessionsStore()
  const cleanups: Array<() => void> = []
  let sessionCleanup: (() => void) | null = null
  let hookCleanup: (() => void) | null = null

  onMounted(() => {
    cleanups.push(EventsOn('app:toast', (level: string, message: string) => {
      console.log('[toast]', level, message)
    }))
    cleanups.push(EventsOn('app:fatal', (msg: string) => {
      console.error('[fatal]', msg)
    }))

    // 后端 fsnotify 监听 jsonl 变化后推送；触发列表刷新
    cleanups.push(EventsOn('sessions:list:changed', () => {
      void sessions.refresh()
    }))

    watch(
      () => sessions.activeId,
      (newId, oldId) => {
        if (oldId) {
          sessionCleanup?.()
          hookCleanup?.()
          sessionCleanup = null
          hookCleanup = null
        }
        if (newId) {
          sessionCleanup = EventsOn(`session:${newId}`, (line: string) => {
            sessions.handleEvent(newId, line)
          })
          hookCleanup = EventsOn(`hook:${newId}`, (line: string) => {
            sessions.handleHookEvent(newId, line)
          })
        }
      },
      { immediate: true }
    )
  })

  onBeforeUnmount(() => {
    sessionCleanup?.()
    hookCleanup?.()
    cleanups.forEach((fn) => fn())
  })

  return { sessions }
}
