<script setup lang="ts">
import { computed } from 'vue'

import { priorityLabel } from '../../tickets/priority'

const props = defineProps<{ priority?: string }>()

const tones: Record<string, string> = {
  High: 'text-rose-700 dark:text-rose-300',
  Medium: 'text-amber-700 dark:text-amber-300',
  Low: 'text-ink-muted',
}
const bars: Record<string, number> = { High: 3, Medium: 2, Low: 1 }

const tone = computed(() => tones[props.priority ?? ''] ?? 'text-ink-subtle')
const level = computed(() => bars[props.priority ?? ''] ?? 0)
</script>

<template>
  <span class="inline-flex items-center gap-1.5 text-sm font-medium" :class="tone"
    ><span v-if="level" aria-hidden="true" class="inline-flex items-end gap-0.5"
      ><span v-for="n in 3" :key="n" class="w-1 rounded-sm" :class="[n <= level ? 'bg-current' : 'bg-line-strong', ['h-1.5', 'h-2.5', 'h-3.5'][n - 1]]"></span></span
    >{{ priorityLabel(priority) }}</span
  >
</template>
