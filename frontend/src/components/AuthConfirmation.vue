<script setup lang="ts">
import { paths } from '../router/paths'
import AppIcon from './ui/AppIcon.vue'
</script>

<template>
  <div class="space-y-6 text-center">
    <span class="mx-auto grid size-12 place-items-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-300">
      <AppIcon name="check" :size="26" />
    </span>
    <p role="status" class="text-sm text-ink"><slot /></p>
    <RouterLink :to="paths.login" class="inline-block text-sm font-semibold text-brand-600 hover:underline dark:text-brand-300"
      >Voltar para o login</RouterLink
    >
  </div>
</template>
