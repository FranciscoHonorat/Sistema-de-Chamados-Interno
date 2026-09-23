<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'

import AppIcon from './ui/AppIcon.vue'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const panel = ref<HTMLElement | null>(null)

watch(
  () => props.open,
  async (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
    if (open) {
      await nextTick()
      panel.value?.focus()
    }
  },
)
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 lg:hidden" @keydown.esc="emit('close')">
      <div class="absolute inset-0 animate-fade-in bg-slate-950/50" aria-hidden="true" @click="emit('close')"></div>
      <div
        ref="panel"
        role="dialog"
        aria-modal="true"
        aria-label="Menu"
        tabindex="-1"
        class="absolute inset-y-0 left-0 flex w-72 max-w-[85vw] animate-slide-in-left flex-col bg-surface shadow-overlay focus:outline-none"
      >
        <button
          type="button"
          class="absolute top-3 right-3 grid size-10 place-items-center rounded-lg text-ink-muted hover:bg-surface-sunken hover:text-ink"
          aria-label="Fechar menu"
          @click="emit('close')"
        >
          <AppIcon name="close" />
        </button>
        <slot />
      </div>
    </div>
  </Teleport>
</template>
