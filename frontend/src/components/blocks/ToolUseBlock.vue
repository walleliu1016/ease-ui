<template>
  <BashTool v-if="name === 'Bash'" :command="(input.command as string) || ''" :description="(input.description as string) || ''" />
  <WriteTool v-else-if="name === 'Write'" :file-path="(input.file_path as string) || 'Unknown file'" :content="(input.content as string) || ''" />
  <EditTool
    v-else-if="name === 'Edit' || name === 'MultiEdit' || name === 'update_plan'"
    :file-path="(input.file_path as string) || 'Unknown file'"
    :old-string="(input.old_string as string) || (input.OldString as string) || ''"
    :new-string="(input.new_string as string) || (input.NewString as string) || ''"
    :replace-all="!!input.replace_all"
  />
  <TodoWriteTool v-else-if="name === 'TodoWrite'" :todos="todoItems" />
  <AskUserQuestionTool
    v-else-if="name === 'AskUserQuestion'"
    :questions="askQuestions"
  />
  <PlanTool
    v-else-if="name === 'ExitPlanMode'"
    :title="'Plan'"
    :plan="(input.plan as string) || ''"
    :plan-file-path="(input.planFilePath as string) || ''"
  />
  <PlanTool
    v-else-if="name === 'EnterPlanMode'"
    :title="'Plan Mode'"
    :plan="(input.plan as string) || ''"
    :plan-file-path="(input.planFilePath as string) || ''"
  />
  <GenericTool v-else :name="name" :description="(input.description as string) || ''" :input="input" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AskOption, TodoItem } from '../../types/blocks'
import { normalizeToolInput } from '../../types/blocks'
import type { AskQuestion } from './tools/AskUserQuestionTool.vue'
import BashTool from './tools/BashTool.vue'
import WriteTool from './tools/WriteTool.vue'
import EditTool from './tools/EditTool.vue'
import TodoWriteTool from './tools/TodoWriteTool.vue'
import AskUserQuestionTool from './tools/AskUserQuestionTool.vue'
import PlanTool from './tools/PlanTool.vue'
import GenericTool from './tools/GenericTool.vue'

const props = defineProps<{ name: string; input: Record<string, unknown> }>()
const input = computed(() => normalizeToolInput(props.input))

const todoItems = computed<TodoItem[]>(() => {
  const raw = (input.value as any).todos
  if (!Array.isArray(raw)) return []
  return raw.map((t: any) => ({
    content: t?.content || t?.activeForm || '',
    status: (t?.status || 'pending') as TodoItem['status'],
    activeForm: t?.activeForm,
  }))
})

const askQuestions = computed<AskQuestion[]>(() => {
  const i = input.value as any
  if (Array.isArray(i.questions)) {
    return i.questions.map((q: any) => ({
      question: q?.question || '',
      options: Array.isArray(q?.options) ? q.options.map((o: any) => ({
        label: o?.label || '',
        description: o?.description,
        preview: o?.preview,
      })) : [],
      multiSelect: !!q?.multiSelect,
      header: q?.header,
    }))
  }
  if (i.question || i.options) {
    return [{
      question: i.question || '',
      options: Array.isArray(i.options) ? i.options.map((o: any) => ({
        label: o?.label || '',
        description: o?.description,
        preview: o?.preview,
      })) : [],
      multiSelect: !!i.multi_select,
      header: i.header,
    }]
  }
  return []
})
</script>
