<script setup lang="ts">
withDefaults(defineProps<{ lines?: number }>(), { lines: 3 })
</script>

<template>
  <div class="animate-pulse space-y-3" aria-hidden="true">
    <div v-for="n in lines" :key="n" class="h-4 rounded bg-surface-sunken" :style="{ width: `${100 - ((n * 17) % 45)}%` }"></div>
  </div>
</template>
