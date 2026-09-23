<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import BrandMark from '../components/BrandMark.vue'
import { homePathFor } from '../router/homePath'
import { paths } from '../router/paths'

const session = inject(sessionKey)!
const home = computed(() => (session.user.value ? homePathFor(session.user.value.role) : paths.login))
</script>

<template>
  <main class="grid min-h-dvh place-items-center px-6 py-16">
    <div class="text-center">
      <BrandMark class="mb-10" />
      <p class="text-sm font-semibold text-brand-600 dark:text-brand-300">Erro 404</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight text-ink sm:text-4xl">Página não encontrada</h1>
      <p class="mt-3 text-ink-muted">O endereço pode estar errado ou a página mudou de lugar.</p>
      <RouterLink
        :to="home"
        class="mt-8 inline-flex h-10 items-center rounded-lg bg-brand-600 px-4 text-sm font-medium text-white shadow-sm hover:bg-brand-700"
        >Voltar para o início</RouterLink
      >
    </div>
  </main>
</template>
