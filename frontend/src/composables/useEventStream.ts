import { onMounted, onBeforeUnmount, watch } from 'vue'
import { EventsOn } from './useWails'
import { useSessionsStore } from '../stores/sessions'

export function useEventStream() {
  const sessions = useSessionsStore()
  const cleanups: Array<() => void> = []
  let sessionCleanup: (() => void) | null = null

  onMounted(() => {
    // 应用级事件
    cleanups.push(EventsOn('app:toast', (level: string, message: string) => {
      console.log('[toast]', level, message)
    }))
    cleanups.push(EventsOn('app:fatal', (msg: string) => {
      console.error('[fatal]', msg)
    }))

    // 监听活跃会话切换，订阅对应 session 流事件
    watch(
      () => sessions.activeId,
      (newId, oldId) => {
        if (oldId) {
          sessionCleanup?.()
          sessionCleanup = null
        }
        if (newId) {
          sessionCleanup = EventsOn(`session:${newId}`, (line: string) => {
            sessions.handleEvent(newId, line)
          })
        }
      },
      { immediate: true }
    )
  })

  onBeforeUnmount(() => {
    sessionCleanup?.()
    cleanups.forEach((fn) => fn())
  })

  return { sessions }
}
