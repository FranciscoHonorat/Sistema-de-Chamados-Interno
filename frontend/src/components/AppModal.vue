<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'

const props = defineProps<{ open: boolean; title: string; description?: string }>()
const emit = defineEmits<{ close: [] }>()

const titleId = useId()
const panel = ref<HTMLElement | null>(null)
let previouslyFocused: HTMLElement | null = null

watch(
  () => props.open,
  async (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
    if (open) {
      previouslyFocused = document.activeElement as HTMLElement | null
      await nextTick()
      panel.value?.focus()
    } else {
      previouslyFocused?.focus?.()
    }
  },
  { immediate: true },
)
onBeforeUnmount(() => (document.body.style.overflow = ''))

// Keeps Tab inside the dialog while it is open.
function trapFocus(event: KeyboardEvent) {
  if (!panel.value) return
  const focusable = panel.value.querySelectorAll<HTMLElement>(
    'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
  )
  if (focusable.length === 0) return
  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        data-testid="modal-backdrop"
        class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/50 backdrop-blur-[2px] sm:items-center sm:p-4"
        @click.self="emit('close')"
      >
        <section
          ref="panel"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="titleId"
          tabindex="-1"
          class="max-h-[92dvh] w-full animate-slide-up overflow-y-auto rounded-t-2xl bg-surface p-5 shadow-overlay focus:outline-none sm:max-w-lg sm:rounded-2xl sm:p-6"
          @keydown.esc="emit('close')"
          @keydown.tab="trapFocus"
        >
          <div class="mx-auto mb-4 h-1 w-10 rounded-full bg-line-strong sm:hidden" aria-hidden="true"></div>
          <header class="mb-5 flex items-start justify-between gap-4">
            <div>
              <h2 :id="titleId" class="text-lg font-semibold text-ink">{{ title }}</h2>
              <p v-if="description" class="mt-1 text-sm text-ink-muted">{{ description }}</p>
            </div>
            <button
              type="button"
              class="shrink-0 rounded-md px-2 py-1 text-sm font-medium text-ink-muted transition-colors hover:bg-surface-sunken hover:text-ink"
              @click="emit('close')"
            >
              ← Voltar
            </button>
          </header>
          <slot />
        </section>
      </div>
    </Transition>
  </Teleport>
</template>
