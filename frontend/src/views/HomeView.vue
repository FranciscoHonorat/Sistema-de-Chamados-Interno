<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'

import { sessionKey } from '../auth/session'
import NewTicketModal from '../components/NewTicketModal.vue'
import PageLayout from '../components/PageLayout.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppIcon from '../components/ui/AppIcon.vue'
import StatusBadge from '../components/ui/StatusBadge.vue'
import { menuFor } from '../router/menu'
import { paths, ticketDetailPath } from '../router/paths'
import { formatRelative } from '../tickets/format'
import { canOpenTickets } from '../tickets/permissions'
import { TicketStatus } from '../tickets/status'
import { ticketsApiKey, type Ticket } from '../tickets/ticketsApi'

const session = inject(sessionKey)!
const api = inject(ticketsApiKey)!

const user = computed(() => session.user.value)
const role = computed(() => user.value?.role)
const menu = computed(() => (role.value ? menuFor(role.value) : []))

const subtitles = {
  user: 'Abra um chamado quando precisar de ajuda e acompanhe o andamento por aqui.',
  support: 'Veja o que está esperando atendimento e o que está com você.',
  admin: 'Acompanhe os chamados, as contas e a carga de trabalho do suporte.',
} as const

const tickets = ref<Ticket[]>([])
const loaded = ref(false)

const stats = computed(() => [
  { label: 'Abertos', value: tickets.value.filter((t) => t.status === TicketStatus.Open).length, tone: 'bg-sky-500' },
  { label: 'Em andamento', value: tickets.value.filter((t) => t.status === TicketStatus.InProgress).length, tone: 'bg-amber-500' },
  { label: 'Fechados', value: tickets.value.filter((t) => t.status === TicketStatus.Closed).length, tone: 'bg-emerald-500' },
])
const recent = computed(() =>
  [...tickets.value].sort((a, b) => b.created_at.localeCompare(a.created_at)).slice(0, 5),
)

async function loadTickets() {
  try {
    tickets.value = await api.list()
  } catch {
    tickets.value = []
  } finally {
    loaded.value = true
  }
}

const creating = ref(false)
const openedTicketId = ref<string | null>(null)

function onOpened(ticketId: string) {
  creating.value = false
  openedTicketId.value = ticketId
  loadTickets()
}

onMounted(loadTickets)
</script>

<template>
  <PageLayout>
    <section class="mb-8">
      <h1 class="text-2xl font-semibold tracking-tight text-ink sm:text-3xl">Olá, {{ user?.name }}</h1>
      <p v-if="role" class="mt-1 text-sm text-ink-muted">{{ subtitles[role] }}</p>
    </section>

    <AlertMessage v-if="openedTicketId" role="status" tone="success" class="mb-6">
      Chamado aberto com sucesso.
      <RouterLink :to="ticketDetailPath(openedTicketId)" class="font-semibold underline underline-offset-2">Ver chamado</RouterLink>
    </AlertMessage>

    <nav aria-label="Atalhos" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <button
        v-if="role && canOpenTickets(role)"
        type="button"
        aria-label="Abrir novo chamado"
        aria-describedby="new-ticket-hint"
        class="group flex flex-col items-start gap-3 rounded-xl bg-brand-600 p-5 text-left text-white shadow-card transition hover:bg-brand-700 focus-visible:outline-offset-4"
        @click="creating = true"
      >
        <span class="grid size-10 place-items-center rounded-lg bg-white/15"><AppIcon name="plus" :size="22" /></span>
        <span>
          <span class="block text-base font-semibold">Abrir novo chamado</span>
          <span id="new-ticket-hint" class="mt-0.5 block text-sm text-brand-100">Descreva o problema e o suporte assume</span>
        </span>
      </button>

      <div
        v-for="item in menu"
        :key="item.label"
        class="group relative flex flex-col items-start gap-3 rounded-xl border border-line bg-surface p-5 shadow-card transition hover:border-brand-300 hover:shadow-md dark:hover:border-brand-500/50"
      >
        <span
          class="grid size-10 place-items-center rounded-lg bg-brand-50 text-brand-600 transition group-hover:bg-brand-600 group-hover:text-white dark:bg-brand-500/10 dark:text-brand-300"
          ><AppIcon :name="item.icon" :size="22"
        /></span>
        <span>
          <RouterLink :to="item.to" class="text-base font-semibold text-ink after:absolute after:inset-0 after:rounded-xl focus-visible:outline-none">{{
            item.label
          }}</RouterLink>
          <span class="mt-0.5 block text-sm text-ink-muted">{{ item.description }}</span>
        </span>
        <AppIcon name="chevronRight" class="absolute top-5 right-5 text-ink-subtle transition group-hover:translate-x-0.5 group-hover:text-brand-600" />
      </div>
    </nav>

    <div class="mt-8 grid gap-6 lg:grid-cols-3">
      <section aria-labelledby="summary-title" class="card p-5">
        <h2 id="summary-title" class="text-sm font-semibold text-ink">Resumo</h2>
        <p class="mt-0.5 text-xs text-ink-subtle">Chamados que você pode ver</p>
        <dl class="mt-4 space-y-3">
          <div v-for="stat in stats" :key="stat.label" class="flex items-center justify-between gap-3">
            <dt class="flex items-center gap-2 text-sm text-ink-muted">
              <span aria-hidden="true" class="size-2 rounded-full" :class="stat.tone"></span>{{ stat.label }}
            </dt>
            <dd class="text-lg font-semibold tabular-nums text-ink">{{ loaded ? stat.value : '–' }}</dd>
          </div>
        </dl>
      </section>

      <section aria-labelledby="recent-title" class="card overflow-hidden lg:col-span-2">
        <header class="flex items-center justify-between border-b border-line px-5 py-4">
          <h2 id="recent-title" class="text-sm font-semibold text-ink">Chamados recentes</h2>
          <RouterLink :to="paths.tickets" class="text-sm font-medium text-brand-600 hover:underline dark:text-brand-300">Ver todos</RouterLink>
        </header>
        <div v-if="!loaded" class="animate-pulse space-y-4 p-5" aria-hidden="true">
          <div v-for="n in 3" :key="n" class="h-10 rounded-lg bg-surface-sunken"></div>
        </div>
        <p v-else-if="recent.length === 0" class="px-5 py-10 text-center text-sm text-ink-muted">Nenhum chamado por enquanto.</p>
        <ul v-else class="divide-y divide-line">
          <li v-for="ticket in recent" :key="ticket.ticket_id" class="relative flex items-center gap-4 px-5 py-3.5 hover:bg-surface-muted">
            <div class="min-w-0 flex-1">
              <RouterLink :to="ticketDetailPath(ticket.ticket_id)" class="block truncate text-sm font-medium text-ink after:absolute after:inset-0">{{
                ticket.title
              }}</RouterLink>
              <p class="text-xs text-ink-subtle">{{ formatRelative(ticket.created_at) }}</p>
            </div>
            <StatusBadge :status="ticket.status" />
          </li>
        </ul>
      </section>
    </div>

    <NewTicketModal :open="creating" @close="creating = false" @opened="onOpened" />
  </PageLayout>
</template>
