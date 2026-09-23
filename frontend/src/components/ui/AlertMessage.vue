<script setup lang="ts">
import { computed } from 'vue'

import AppIcon from './AppIcon.vue'

const props = withDefaults(defineProps<{ tone?: 'error' | 'success' | 'warning' | 'info' }>(), { tone: 'error' })

const styles = {
  error: { box: 'border-rose-200 bg-rose-50 text-rose-800 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-200', icon: 'warning' },
  success: { box: 'border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-200', icon: 'check' },
  warning: { box: 'border-amber-200 bg-amber-50 text-amber-900 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-200', icon: 'warning' },
  info: { box: 'border-brand-200 bg-brand-50 text-brand-900 dark:border-brand-500/30 dark:bg-brand-500/10 dark:text-brand-100', icon: 'info' },
} as const

const style = computed(() => styles[props.tone])
</script>

<!-- The caller sets role="alert" or role="status" on the component, as fits. -->
<template>
  <div class="flex items-start gap-3 rounded-lg border px-3.5 py-3 text-sm" :class="style.box">
    <AppIcon :name="style.icon" :size="18" class="mt-px" />
    <div class="min-w-0 flex-1"><slot /></div>
  </div>
</template>
