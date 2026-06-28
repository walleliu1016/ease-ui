<template>
  <div class="tool-ask">
    <div class="ask-header"><span class="ask-icon">❓</span> AskUserQuestion</div>
    <div v-for="(q, qi) in questions" :key="qi" class="ask-question-block">
      <div v-if="q.question" class="ask-question">{{ q.question }}</div>
      <ol v-if="q.options && q.options.length" class="ask-options">
        <li v-for="(opt, i) in q.options" :key="i" class="ask-option">
          <span class="ask-option-num">{{ numGlyph(i + 1) }}</span>
          <div class="ask-option-body">
            <div class="ask-option-label">{{ opt.label }}</div>
            <div v-if="opt.description" class="ask-option-desc">{{ opt.description }}</div>
            <details v-if="opt.preview" class="ask-option-preview">
              <summary>Preview</summary>
              <pre>{{ opt.preview }}</pre>
            </details>
          </div>
        </li>
      </ol>
      <div v-if="q.multiSelect" class="ask-hint">多选</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AskOption } from '../../../types/blocks'

export interface AskQuestion {
  question: string
  options: AskOption[]
  multiSelect?: boolean
  header?: string
}

defineProps<{ questions: AskQuestion[] }>()

const NUM_GLYPHS = ['①', '②', '③', '④']
function numGlyph(i: number) { return NUM_GLYPHS[i - 1] || `${i}.` }
</script>

<style scoped>
.tool-ask {
  background: var(--ask-bg);
  border: 1px solid var(--ask-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 6px 0;
}
.ask-header {
  font-weight: 600;
  color: var(--ask-header);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.ask-icon { font-size: 14px; }
.ask-question-block { margin-bottom: 8px; }
.ask-question-block:last-child { margin-bottom: 0; }
.ask-question {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
  margin-bottom: 6px;
  line-height: 1.5;
}
.ask-options {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.ask-option {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 6px 8px;
  background: var(--ask-option-bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}
.ask-option-num {
  flex-shrink: 0;
  font-size: 13px;
  color: var(--accent);
  font-weight: 600;
  font-family: var(--font-mono);
  line-height: 1.4;
}
.ask-option-body { flex: 1; min-width: 0; }
.ask-option-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.4;
}
.ask-option-desc {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
  line-height: 1.4;
}
.ask-option-preview { margin-top: 4px; font-size: 11px; color: var(--text-tertiary); }
.ask-option-preview summary { cursor: pointer; }
.ask-option-preview pre {
  background: var(--code-bg);
  color: var(--code-text);
  padding: 6px 8px;
  border-radius: 3px;
  font-size: 11px;
  margin: 4px 0 0;
  white-space: pre-wrap;
  font-family: var(--font-mono);
}
.ask-hint {
  margin-top: 4px;
  font-size: 10px;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
</style>
