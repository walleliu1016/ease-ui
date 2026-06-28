<template>
  <div class="image-block">
    <img :src="`data:${mediaType};base64,${data}`" :alt="alt" />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ mediaType: string; data: string; alt?: string }>(), { alt: '' })
</script>

<style scoped>
.image-block { margin: 4px 0; }
.image-block img { max-width: 100%; border-radius: var(--radius-sm); display: block; }
</style>
