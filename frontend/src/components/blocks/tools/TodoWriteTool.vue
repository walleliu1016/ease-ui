<template>
  <div class="tool-todo">
    <div class="todo-header"><span class="todo-icon">☰</span> Task List</div>
    <ul class="todo-items">
      <li v-for="(todo, i) in todos" :key="i" class="todo-item" :class="`todo-${todo.status}`">
        <span class="todo-glyph">{{ glyph(todo.status) }}</span>
        <span class="todo-content">{{ todo.content }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import type { TodoItem } from '../../../types/blocks'
defineProps<{ todos: TodoItem[] }>()
const glyph = (s: TodoItem['status']) => s === 'completed' ? '✓' : s === 'in_progress' ? '→' : '○'
</script>

<style scoped>
.tool-todo {
  background: var(--todo-bg);
  border: 1px solid var(--todo-border);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  margin: 6px 0;
}
.todo-header {
  font-weight: 600;
  color: var(--todo-header);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.todo-icon { font-size: 13px; }
.todo-items { list-style: none; margin: 0; padding: 0; }
.todo-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px solid rgba(0,0,0,0.06);
  font-size: 13px;
}
.todo-item:last-child { border-bottom: none; }
.todo-glyph {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  border-radius: 50%;
  font-size: 11px;
  font-family: var(--font-mono);
}
.todo-completed .todo-glyph { color: var(--status-success); background: rgba(52,211,153,0.15); }
.todo-completed .todo-content { color: var(--text-tertiary); text-decoration: line-through; }
.todo-in_progress .todo-glyph { color: var(--status-warn); background: rgba(251,191,36,0.15); }
.todo-in_progress .todo-content { color: var(--text-primary); font-weight: 500; }
.todo-pending .todo-glyph { color: var(--text-tertiary); background: rgba(0,0,0,0.05); }
.todo-pending .todo-content { color: var(--text-secondary); }
</style>
