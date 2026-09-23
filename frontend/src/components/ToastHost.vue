<script setup lang="ts">
import { useToast } from '../composables/useToast'
import AppIcon from './ui/AppIcon.vue'

const { toasts, dismiss } = useToast()

const tones = {
  success: 'text-emerald-600 dark:text-emerald-400',
  error: 'text-rose-600 dark:text-rose-400',
  info: 'text-brand-600 dark:text-brand-300',
} as const
const icons = { success: 'check', error: 'warning', info: 'info' } as const
</script>

<template>
  <Teleport to="body">
    <div
      aria-live="polite"
      class="pointer-events-none fixed inset-x-0 bottom-0 z-[60] flex flex-col items-center gap-2 p-4 sm:items-end sm:p-6"
    >
      <TransitionGroup
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="translate-y-2 opacity-0"
        leave-active-class="transition duration-150 ease-in"
        leave-to-class="opacity-0"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="card pointer-events-auto flex w-full max-w-sm items-center gap-3 px-4 py-3 text-sm shadow-overlay"
        >
          <AppIcon :name="icons[toast.tone]" :class="tones[toast.tone]" />
          <p class="flex-1 text-ink">{{ toast.message }}</p>
          <button
            type="button"
            class="rounded-md p-1 text-ink-subtle hover:bg-surface-sunken hover:text-ink"
            aria-label="Dispensar aviso"
            @click="dismiss(toast.id)"
          >
            <AppIcon name="close" :size="16" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
