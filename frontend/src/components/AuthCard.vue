<script setup lang="ts">
import BrandMark from './BrandMark.vue'
import ThemeToggle from './ThemeToggle.vue'
import AppIcon from './ui/AppIcon.vue'

const highlights = [
  { icon: 'sparkles', text: 'Distribuição automática para o atendente menos ocupado' },
  { icon: 'chat', text: 'Conversa e histórico completo em cada chamado' },
  { icon: 'bell', text: 'Notificações de tudo que muda nos seus chamados' },
] as const
</script>

<template>
  <main class="grid min-h-dvh lg:grid-cols-[1.1fr_1fr]">
    <section
      class="relative hidden overflow-hidden bg-gradient-to-br from-brand-800 via-brand-600 to-teal-600 p-12 text-white lg:flex lg:flex-col lg:justify-between"
    >
      <div aria-hidden="true" class="absolute -top-24 -right-24 size-96 rounded-full bg-white/10 blur-3xl"></div>
      <div aria-hidden="true" class="absolute -bottom-32 -left-16 size-96 rounded-full bg-lime-300/20 blur-3xl"></div>
      <BrandMark inverted class="relative" />
      <div class="relative max-w-md">
        <p class="text-3xl leading-tight font-semibold tracking-tight">Chamados internos, do pedido à solução.</p>
        <ul class="mt-8 space-y-4">
          <li v-for="item in highlights" :key="item.text" class="flex items-center gap-3 text-brand-50">
            <span class="grid size-9 place-items-center rounded-lg bg-white/15"><AppIcon :name="item.icon" /></span>
            {{ item.text }}
          </li>
        </ul>
      </div>
      <p class="relative text-sm text-brand-100">Suporte interno com rastreabilidade de ponta a ponta.</p>
    </section>

    <section class="relative flex items-center justify-center px-4 py-12 sm:px-8">
      <div class="absolute top-4 right-4"><ThemeToggle /></div>
      <div class="w-full max-w-sm">
        <header class="mb-8 text-center lg:text-left">
          <svg aria-hidden="true" viewBox="0 0 32 32" class="mx-auto mb-4 size-11 lg:hidden">
            <rect width="32" height="32" rx="8" fill="#15803d" />
            <path d="M9 11.5A2.5 2.5 0 0 1 11.5 9h9a2.5 2.5 0 0 1 2.5 2.5v6a2.5 2.5 0 0 1-2.5 2.5H15l-4.2 3.2c-.5.4-1.3 0-1.3-.6V20A2.5 2.5 0 0 1 9 17.5z" fill="#fff" />
          </svg>
          <h1 class="text-2xl font-semibold tracking-tight text-ink">CodeTicket</h1>
          <p class="mt-1 text-sm text-ink-muted">Sistema de chamados internos</p>
        </header>
        <div class="card p-6 sm:p-8">
          <slot />
        </div>
      </div>
    </section>
  </main>
</template>
