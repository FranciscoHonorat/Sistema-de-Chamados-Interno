<script setup lang="ts">
import { formatDuration } from '../tickets/format'
import AppIcon from './ui/AppIcon.vue'

defineProps<{ resolution?: string; openedAt: string; closedAt: string }>()
</script>

<template>
  <section
    aria-label="Relatório de fechamento"
    class="rounded-xl border border-emerald-200 bg-emerald-50 p-5 text-sm dark:border-emerald-500/30 dark:bg-emerald-500/10"
  >
    <h3 class="mb-2 flex items-center gap-2 font-semibold text-emerald-900 dark:text-emerald-200">
      <AppIcon name="check" :size="18" />Relatório de fechamento
    </h3>
    <p class="whitespace-pre-line text-ink break-words">{{ resolution }}</p>
    <p class="mt-3 text-emerald-800/80 dark:text-emerald-300/80">Tempo até fechar: {{ formatDuration(openedAt, closedAt) }}</p>
  </section>
</template>
