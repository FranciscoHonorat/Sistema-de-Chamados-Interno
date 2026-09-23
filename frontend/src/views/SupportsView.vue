<script setup lang="ts">
import { computed, inject } from 'vue'

import PageLayout from '../components/PageLayout.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'
import { useLoad } from '../composables/useLoad'
import { ticketsApiKey, type AgentWorkload } from '../tickets/ticketsApi'

const api = inject(ticketsApiKey)!

const { data: workloads, failed, loading } = useLoad<AgentWorkload[]>(() => api.supportWorkload(), [])

const total = (agent: AgentWorkload) => agent.open + agent.in_progress + agent.closed
const busiest = computed(() => Math.max(1, ...workloads.value.map((agent) => agent.open + agent.in_progress)))
const totals = computed(() => ({
  open: workloads.value.reduce((sum, agent) => sum + agent.open, 0),
  inProgress: workloads.value.reduce((sum, agent) => sum + agent.in_progress, 0),
  closed: workloads.value.reduce((sum, agent) => sum + agent.closed, 0),
}))
const share = (value: number) => `${Math.round((value / busiest.value) * 100)}%`
</script>

<template>
  <PageLayout back>
    <PageHeader title="Suportes" description="Carga de trabalho de cada atendente. A distribuição automática escolhe quem tem menos chamados em aberto." />

    <AlertMessage v-if="failed" role="alert">Não foi possível carregar os suportes</AlertMessage>
    <template v-else>
      <dl class="mb-6 grid grid-cols-3 gap-3 sm:gap-4">
        <div class="card p-4">
          <dt class="text-xs font-medium text-ink-subtle sm:text-sm">Abertos</dt>
          <dd class="mt-1 text-2xl font-semibold tabular-nums text-ink">{{ totals.open }}</dd>
        </div>
        <div class="card p-4">
          <dt class="text-xs font-medium text-ink-subtle sm:text-sm">Em andamento</dt>
          <dd class="mt-1 text-2xl font-semibold tabular-nums text-ink">{{ totals.inProgress }}</dd>
        </div>
        <div class="card p-4">
          <dt class="text-xs font-medium text-ink-subtle sm:text-sm">Fechados</dt>
          <dd class="mt-1 text-2xl font-semibold tabular-nums text-ink">{{ totals.closed }}</dd>
        </div>
      </dl>

      <div v-if="loading" class="card space-y-4 p-5" aria-hidden="true">
        <div v-for="n in 3" :key="n" class="h-10 animate-pulse rounded-lg bg-surface-sunken"></div>
      </div>

      <div v-else-if="workloads.length === 0" class="card">
        <EmptyState icon="users" title="Nenhum atendente cadastrado" description="Os funcionários com perfil de suporte aparecem aqui." />
      </div>

      <div v-else class="md:overflow-hidden md:rounded-xl md:border md:border-line md:bg-surface md:shadow-card">
        <table class="responsive-table">
          <thead class="bg-surface-muted">
            <tr>
              <th scope="col">Nome</th>
              <th scope="col">Abertos</th>
              <th scope="col">Em andamento</th>
              <th scope="col">Fechados</th>
              <th scope="col">Total</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="agent in workloads" :key="agent.id">
              <td class="block pb-2 md:table-cell">
                <span class="flex items-center gap-3">
                  <UserAvatar :name="agent.name" />
                  <span class="min-w-0 flex-1">
                    <span class="block font-semibold text-ink md:font-medium">{{ agent.name }}</span>
                    <span aria-hidden="true" class="mt-1.5 flex h-1.5 w-full max-w-48 overflow-hidden rounded-full bg-surface-sunken">
                      <span class="bg-sky-500" :style="{ width: share(agent.open) }"></span>
                      <span class="bg-amber-500" :style="{ width: share(agent.in_progress) }"></span>
                    </span>
                  </span>
                </span>
              </td>
              <td data-label="Abertos" class="tabular-nums">{{ agent.open }}</td>
              <td data-label="Em andamento" class="tabular-nums">{{ agent.in_progress }}</td>
              <td data-label="Fechados" class="tabular-nums">{{ agent.closed }}</td>
              <td data-label="Total" class="font-semibold tabular-nums">{{ total(agent) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </PageLayout>
</template>
