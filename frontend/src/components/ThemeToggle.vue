<script setup lang="ts">
import { computed } from 'vue'

import { useTheme } from '../composables/useTheme'
import AppIcon from './ui/AppIcon.vue'

const { mode, cycle } = useTheme()

const labels = { system: 'Tema do sistema', light: 'Tema claro', dark: 'Tema escuro' } as const
const icons = { system: 'computer', light: 'sun', dark: 'moon' } as const
const label = computed(() => `${labels[mode.value]} (clique para trocar)`)
</script>

<template>
  <button
    type="button"
    class="grid size-10 place-items-center rounded-lg text-ink-muted transition-colors hover:bg-surface-sunken hover:text-ink"
    :aria-label="label"
    :title="label"
    @click="cycle"
  >
    <AppIcon :name="icons[mode]" />
  </button>
</template>
