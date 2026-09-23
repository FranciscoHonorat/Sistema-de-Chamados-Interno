<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import { sessionKey } from '../auth/session'
import { useSessionExit } from '../auth/useSessionExit'
import NewTicketModal from '../components/NewTicketModal.vue'
import PageLayout from '../components/PageLayout.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppIcon from '../components/ui/AppIcon.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import PriorityBadge from '../components/ui/PriorityBadge.vue'
import StatusBadge from '../components/ui/StatusBadge.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'
import { useToast } from '../composables/useToast'
import { paths, ticketDetailPath } from '../router/paths'
import { formatDate } from '../tickets/format'
import { canOpenTickets } from '../tickets/permissions'
import { nameResolver } from '../tickets/responsibles'
import { TicketStatus } from '../tickets/status'
import { SessionExpiredError, ticketsApiKey, type Responsible, type Ticket } from '../tickets/ticketsApi'

const api = inject(ticketsApiKey)!
const session = inject(sessionKey)!
const { expire } = useSessionExit()
const route = useRoute()
const toast = useToast()

const tickets = ref<Ticket[]>([])
const responsibles = ref<Responsible[]>([])
const loaded = ref(false)
const failed = ref(false)
const search = ref('')

const filters: Array<{ label: string; status?: TicketStatus }> = [
  { label: 'Todos' },
  { label: 'Abertos', status: TicketStatus.Open },
  { label: 'Em atendimento', status: TicketStatus.InProgress },
  { label: 'Fechados', status: TicketStatus.Closed },
]

const nameOf = computed(() => nameResolver(responsibles.value))

const currentStatus = computed(() => route.query.status as TicketStatus | undefined)
const countFor = (status?: TicketStatus) =>
  status ? tickets.value.filter((t) => t.status === status).length : tickets.value.length

const visible = computed(() => {
  const term = search.value.trim().toLocaleLowerCase('pt-BR')
  return tickets.value
    .filter((t) => !currentStatus.value || t.status === currentStatus.value)
    .filter((t) => !term || t.title.toLocaleLowerCase('pt-BR').includes(term))
    .sort((a, b) => b.created_at.localeCompare(a.created_at))
})

const descriptions: Record<string, string> = {
  user: 'Os chamados que você abriu.',
  support: 'Os chamados abertos e os que estão com você.',
  admin: 'Todos os chamados da empresa.',
}
const description = computed(() => descriptions[session.user.value?.role ?? ''] ?? '')

const creating = ref(false)
const canOpen = computed(() => (session.user.value ? canOpenTickets(session.user.value.role) : false))

async function load() {
  failed.value = false
  try {
    const [list, people] = await Promise.all([api.list(), api.responsibles()])
    tickets.value = list
    responsibles.value = people
    loaded.value = true
  } catch (error) {
    if (error instanceof SessionExpiredError) {
      await expire()
      return
    }
    failed.value = true
  }
}

async function onOpened() {
  creating.value = false
  toast.show('Chamado aberto')
  tickets.value = await api.list()
}

onMounted(load)
</script>

<template>
  <PageLayout back>
    <PageHeader title="Chamados" :description="description">
      <template #actions>
        <AppButton v-if="canOpen" variant="primary" icon="plus" @click="creating = true">Abrir novo chamado</AppButton>
      </template>
    </PageHeader>

    <div class="mb-4 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
      <nav aria-label="Filtrar por situação" class="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0">
        <div class="inline-flex gap-1 rounded-xl border border-line bg-surface p-1 shadow-card">
          <RouterLink
            v-for="filter in filters"
            :key="filter.label"
            v-slot="{ href, navigate }"
            :to="{ path: paths.tickets, query: filter.status ? { status: filter.status } : {} }"
            custom
          >
            <a
              :href="href"
              :aria-current="currentStatus === filter.status ? 'page' : undefined"
              class="inline-flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium whitespace-nowrap transition-colors"
              :class="
                currentStatus === filter.status
                  ? 'bg-brand-600 text-white shadow-sm'
                  : 'text-ink-muted hover:bg-surface-sunken hover:text-ink'
              "
              @click="navigate"
              >{{ filter.label
              }}<span
                v-if="loaded"
                aria-hidden="true"
                class="rounded-full px-1.5 text-xs tabular-nums"
                :class="currentStatus === filter.status ? 'bg-white/20' : 'bg-surface-sunken'"
                >{{ countFor(filter.status) }}</span
              ></a
            >
          </RouterLink>
        </div>
      </nav>

      <div class="relative w-full lg:w-72">
        <AppIcon name="search" :size="18" class="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-ink-subtle" />
        <input v-model="search" type="search" aria-label="Buscar chamados pelo título" placeholder="Buscar pelo título…" class="field-input pl-10" />
      </div>
    </div>

    <AlertMessage v-if="failed" role="alert" class="mb-4">
      <p>Não foi possível carregar os chamados.</p>
      <button type="button" class="mt-1 font-semibold underline underline-offset-2" @click="load">Tentar novamente</button>
    </AlertMessage>

    <div v-else-if="!loaded" class="card space-y-4 p-5" aria-hidden="true">
      <div v-for="n in 5" :key="n" class="flex animate-pulse items-center gap-4">
        <div class="h-4 flex-1 rounded bg-surface-sunken"></div>
        <div class="h-5 w-24 rounded-full bg-surface-sunken"></div>
        <div class="hidden h-4 w-20 rounded bg-surface-sunken sm:block"></div>
      </div>
    </div>

    <div v-else-if="visible.length === 0" class="card">
      <EmptyState
        title="Nenhum chamado encontrado"
        :description="search ? 'Nenhum título combina com a busca.' : 'Quando houver chamados nesta situação, eles aparecem aqui.'"
      />
    </div>

    <div v-else class="md:overflow-hidden md:rounded-xl md:border md:border-line md:bg-surface md:shadow-card">
      <table class="responsive-table">
        <thead class="bg-surface-muted">
          <tr>
            <th scope="col">Título</th>
            <th scope="col">Status</th>
            <th scope="col">Prioridade</th>
            <th scope="col">Responsável</th>
            <th scope="col">Aberto em</th>
            <th scope="col">Fechado em</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ticket in visible" :key="ticket.ticket_id" class="relative">
            <td class="block pb-2 md:table-cell md:max-w-xs">
              <RouterLink
                :to="ticketDetailPath(ticket.ticket_id)"
                class="block truncate font-semibold text-ink after:absolute after:inset-0 hover:text-brand-600 md:font-medium dark:hover:text-brand-300"
                >{{ ticket.title }}</RouterLink
              >
            </td>
            <td data-label="Status"><StatusBadge :status="ticket.status" /></td>
            <td data-label="Prioridade"><PriorityBadge :priority="ticket.priority" /></td>
            <td data-label="Responsável">
              <span class="inline-flex items-center gap-2 text-ink-muted"
                ><UserAvatar v-if="ticket.assignee_id" :name="nameOf(ticket.assignee_id)" size="sm" />{{ nameOf(ticket.assignee_id) }}</span
              >
            </td>
            <td data-label="Aberto em" class="text-ink-muted tabular-nums">{{ formatDate(ticket.created_at) }}</td>
            <td data-label="Fechado em" class="text-ink-muted tabular-nums">{{ formatDate(ticket.closed_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <NewTicketModal :open="creating" @close="creating = false" @opened="onOpened" />
  </PageLayout>
</template>
