<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import { homePathFor } from '../router/homePath'
import { paths } from '../router/paths'

const props = defineProps<{ to?: string }>()

const session = inject(sessionKey)!
const destination = computed(
  () => props.to ?? (session.user.value ? homePathFor(session.user.value.role) : paths.login),
)
</script>

<template>
  <RouterLink
    :to="destination"
    class="mb-4 inline-flex items-center gap-1 rounded-md text-sm font-medium text-ink-muted transition-colors hover:text-brand-600 dark:hover:text-brand-300"
    >← Voltar</RouterLink
  >
</template>
