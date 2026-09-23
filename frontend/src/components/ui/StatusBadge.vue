<script setup lang="ts">
import { computed } from 'vue'

import { statusLabel, TicketStatus } from '../../tickets/status'

const props = defineProps<{ status: TicketStatus }>()

const tones: Record<TicketStatus, string> = {
  [TicketStatus.Open]: 'bg-sky-50 text-sky-700 ring-sky-600/20 dark:bg-sky-400/10 dark:text-sky-300 dark:ring-sky-400/30',
  [TicketStatus.InProgress]:
    'bg-amber-50 text-amber-800 ring-amber-600/20 dark:bg-amber-400/10 dark:text-amber-300 dark:ring-amber-400/30',
  [TicketStatus.Closed]:
    'bg-emerald-50 text-emerald-700 ring-emerald-600/20 dark:bg-emerald-400/10 dark:text-emerald-300 dark:ring-emerald-400/30',
}
const dots: Record<TicketStatus, string> = {
  [TicketStatus.Open]: 'bg-sky-500',
  [TicketStatus.InProgress]: 'bg-amber-500',
  [TicketStatus.Closed]: 'bg-emerald-500',
}

const tone = computed(() => tones[props.status] ?? 'bg-surface-sunken text-ink-muted ring-line')
const dot = computed(() => dots[props.status] ?? 'bg-ink-subtle')
</script>

<template>
  <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset" :class="tone"
    ><span aria-hidden="true" class="size-1.5 rounded-full" :class="dot"></span>{{ statusLabel(status) }}</span
  >
</template>
