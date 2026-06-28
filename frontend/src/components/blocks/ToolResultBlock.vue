<template>
  <div class="tool-result" :class="{ 'tool-error': isError }">
    <Truncatable>
      <template v-for="(part, i) in renderedParts" :key="i">
        <div v-if="part.kind === 'commit'" class="commit-card">
          <a v-if="githubRepo" :href="`https://github.com/${githubRepo}/commit/${part.hash}`" target="_blank" rel="noopener">
            <span class="commit-card-hash">{{ part.hash.slice(0, 7) }}</span> {{ part.msg }}
          </a>
          <template v-else>
            <span class="commit-card-hash">{{ part.hash.slice(0, 7) }}</span> {{ part.msg }}
          </template>
        </div>
        <pre v-else-if="part.kind === 'text'">{{ part.text }}</pre>
        <ImageBlock v-else :media-type="part.mediaType" :data="part.data" />
      </template>
    </Truncatable>
  </div>
</template>

<script setup lang="ts">
// 复刻参考实现：tool_result 块识别 commit pattern + image + 文本。
// 整块用 Truncatable 包裹。
import { computed } from 'vue'
import Truncatable from './Truncatable.vue'
import ImageBlock from './ImageBlock.vue'
import type { ToolResultBlock as Block } from '../../types/blocks'

const props = withDefaults(defineProps<{
  content: Block[]
  isError: boolean
  githubRepo?: string
}>(), { githubRepo: '' })

// 跟参考实现保持一致的 commit 识别
const COMMIT_PATTERN = /\[[\w\-/]+ ([a-f0-9]{7,})\] (.+?)(?:\n|$)/g

type Part = { kind: 'text'; text: string } | { kind: 'commit'; hash: string; msg: string } | { kind: 'image'; mediaType: string; data: string }

const renderedParts = computed<Part[]>(() => {
  const out: Part[] = []
  for (const c of props.content) {
    if (c.type === 'image') {
      out.push({ kind: 'image', mediaType: c.mediaType, data: c.data })
    } else if (c.type === 'text') {
      // 在 text 块内尝试识别 commit
      const matches: { idx: number; len: number; hash: string; msg: string }[] = []
      let m: RegExpExecArray | null
      COMMIT_PATTERN.lastIndex = 0
      while ((m = COMMIT_PATTERN.exec(c.text)) !== null) {
        matches.push({ idx: m.index, len: m[0].length, hash: m[1], msg: m[2] })
      }
      if (matches.length === 0) {
        out.push({ kind: 'text', text: c.text })
      } else {
        let cursor = 0
        for (const cm of matches) {
          if (cm.idx > cursor) {
            const before = c.text.slice(cursor, cm.idx).trim()
            if (before) out.push({ kind: 'text', text: before })
          }
          out.push({ kind: 'commit', hash: cm.hash, msg: cm.msg })
          cursor = cm.idx + cm.len
        }
        if (cursor < c.text.length) {
          const after = c.text.slice(cursor).trim()
          if (after) out.push({ kind: 'text', text: after })
        }
      }
    }
  }
  return out
})
</script>

<style scoped>
.tool-result {
  background: var(--tool-result-bg);
  border-radius: var(--radius-md);
  padding: 8px 10px;
  margin: 6px 0;
  font-size: 13px;
}
.tool-result.tool-error { background: var(--tool-error-bg); }
.tool-result pre {
  background: var(--code-bg);
  color: var(--code-text);
  padding: 8px 10px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.5;
  margin: 4px 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: var(--font-mono);
}
.commit-card {
  margin: 6px 0;
  padding: 8px 10px;
  background: var(--commit-card-bg);
  border-left: 3px solid var(--commit-card-border);
  border-radius: 4px;
  font-size: 12px;
}
.commit-card a {
  text-decoration: none;
  color: var(--text-primary);
  display: block;
}
.commit-card a:hover { color: var(--accent); }
.commit-card-hash {
  font-family: var(--font-mono);
  color: var(--accent);
  font-weight: 600;
  margin-right: 6px;
}
</style>
