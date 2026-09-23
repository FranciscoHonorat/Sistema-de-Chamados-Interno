<script setup lang="ts">
import AppIcon, { type IconName } from './AppIcon.vue'

withDefaults(defineProps<{ icon?: IconName; title: string; description?: string }>(), { icon: 'inbox' })
</script>

<template>
  <div class="flex flex-col items-center justify-center px-6 py-12 text-center">
    <span class="mb-4 grid size-12 place-items-center rounded-full bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-300">
      <AppIcon :name="icon" :size="24" />
    </span>
    <p class="text-base font-semibold text-ink">{{ title }}</p>
    <p v-if="description" class="mt-1 max-w-sm text-sm text-ink-muted">{{ description }}</p>
    <div v-if="$slots.default" class="mt-5">
      <slot />
    </div>
  </div>
</template>
