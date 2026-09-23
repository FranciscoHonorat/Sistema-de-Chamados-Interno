<script setup lang="ts">
import { computed, inject, ref } from 'vue'

import { sessionKey } from '../auth/session'
import ClosingReport from '../components/ClosingReport.vue'
import CloseTicketModal from '../components/CloseTicketModal.vue'
import FormField from '../components/FormField.vue'
import PageLayout from '../components/PageLayout.vue'
import TicketActions from '../components/TicketActions.vue'
import TicketEditForm from '../components/TicketEditForm.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import PriorityBadge from '../components/ui/PriorityBadge.vue'
import SkeletonBlock from '../components/ui/SkeletonBlock.vue'
import StatusBadge from '../components/ui/StatusBadge.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'
import { paths } from '../router/paths'
import { formatDate, formatDateTime, formatRelative } from '../tickets/format'
import { allowedActions, type TicketAction } from '../tickets/permissions'
import { nameResolver } from '../tickets/responsibles'
import { useTicketDetail } from '../tickets/useTicketDetail'

const props = defineProps<{ id: string }>()

const session = inject(sessionKey)!
const { api, ticket, responsibles, notFound, loading, busy, actionError, perform } = useTicketDetail(props.id)

const reply = ref('')
const replyError = ref('')
const closing = ref(false)
const editing = ref(false)

const actions = computed(() =>
  ticket.value && session.user.value ? allowedActions(ticket.value, session.user.value) : [],
)
const can = (action: TicketAction) => actions.value.includes(action)
const panelActions = computed(() => actions.value.filter((action) => action !== 'respond'))
const responseCount = computed(() => ticket.value?.responses?.length ?? 0)
const responsibleOptions = computed(() => responsibles.value.map((r) => ({ value: r.id, label: r.name })))
// The tickets module only knows the support agents by name; the logged-in
// user is the other name the page can always show.
const nameOf = computed(() => {
  const resolve = nameResolver(responsibles.value)
  const me = session.user.value
  return (id: string | undefined) => (me && id === me.id ? me.name : resolve(id))
})
const isMine = (authorId: string) => authorId === session.user.value?.id

async function closeTicket(resolution: string) {
  if (await perform(() => api.close(props.id, resolution), 'Chamado fechado')) {
    closing.value = false
  }
}

async function saveEdit(changes: { title: string; description: string }) {
  if (await perform(() => api.edit(props.id, changes), 'Chamado atualizado')) {
    editing.value = false
  }
}

async function sendReply() {
  if (!reply.value.trim()) {
    replyError.value = 'Escreva uma resposta'
    return
  }
  replyError.value = ''
  if (await perform(() => api.respond(props.id, reply.value), 'Resposta enviada')) {
    reply.value = ''
  }
}
</script>

<template>
  <PageLayout back :back-to="paths.tickets">
    <div v-if="loading" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]" aria-hidden="true">
      <div class="card space-y-4 p-6">
        <div class="h-7 w-2/3 animate-pulse rounded bg-surface-sunken"></div>
        <SkeletonBlock :lines="4" />
      </div>
      <div class="card p-6"><SkeletonBlock :lines="5" /></div>
    </div>

    <div v-else-if="notFound" class="card">
      <EmptyState icon="inbox" title="Chamado não encontrado" description="Ele não existe ou você não tem acesso a ele." />
    </div>

    <!-- One article for the whole ticket. The aside comes before the
         conversation in the DOM (and in the reading order on phones), and
         moves to the right column on large screens. -->
    <article v-if="ticket" class="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
      <div class="space-y-4 lg:col-start-1 lg:row-start-1">
        <AlertMessage v-if="actionError" role="alert">{{ actionError }}</AlertMessage>

        <section class="card p-5 sm:p-6">
          <TicketEditForm
            v-if="editing"
            class="mb-6"
            :title="ticket.title"
            :description="ticket.description"
            :saving="busy"
            @save="saveEdit"
            @cancel="editing = false"
          />
          <header>
            <div class="mb-3 flex flex-wrap items-center gap-2 text-xs text-ink-subtle">
              <StatusBadge :status="ticket.status" />
              <span>Aberto {{ formatRelative(ticket.created_at) }}</span>
            </div>
            <h1 class="text-xl font-semibold tracking-tight text-ink break-words sm:text-2xl">{{ ticket.title }}</h1>
            <p v-if="!editing" class="mt-3 text-sm leading-relaxed whitespace-pre-line text-ink-muted break-words sm:text-base">
              {{ ticket.description }}
            </p>
          </header>
        </section>
      </div>

      <aside class="space-y-4 lg:sticky lg:top-24 lg:col-start-2 lg:row-span-2 lg:row-start-1">
        <section class="card p-5">
          <h2 class="mb-3 text-sm font-semibold text-ink">Detalhes</h2>
          <ul aria-label="Detalhes" class="space-y-3 text-sm">
            <li class="flex items-center justify-between gap-3">
              <span class="text-ink-subtle">Status</span><StatusBadge :status="ticket.status" />
            </li>
            <li class="flex items-center justify-between gap-3">
              <span class="text-ink-subtle">Prioridade</span><PriorityBadge :priority="ticket.priority" />
            </li>
            <li class="flex items-center justify-between gap-3">
              <span class="text-ink-subtle">Responsável</span>
              <span class="inline-flex items-center gap-2 font-medium text-ink"
                ><UserAvatar v-if="ticket.assignee_id" :name="nameOf(ticket.assignee_id)" size="sm" /><span>{{
                  nameOf(ticket.assignee_id)
                }}</span></span
              >
            </li>
            <li class="flex items-center justify-between gap-3">
              <span class="text-ink-subtle">Aberto em</span
              ><span class="font-medium text-ink tabular-nums" :title="formatDateTime(ticket.created_at)">{{ formatDate(ticket.created_at) }}</span>
            </li>
            <li v-if="ticket.closed_at" class="flex items-center justify-between gap-3">
              <span class="text-ink-subtle">Fechado em</span
              ><span class="font-medium text-ink tabular-nums" :title="formatDateTime(ticket.closed_at)">{{ formatDate(ticket.closed_at) }}</span>
            </li>
          </ul>
        </section>

        <TicketActions
          v-if="panelActions.length > 0"
          :actions="actions"
          :responsibles="responsibleOptions"
          :current-priority="ticket.priority"
          :busy="busy"
          @assign="(assigneeId) => perform(() => api.assign(id, assigneeId), 'Chamado atribuído')"
          @auto-assign="perform(() => api.autoAssign(id), 'Chamado distribuído automaticamente')"
          @change-priority="(priority) => perform(() => api.changePriority(id, priority), 'Prioridade alterada')"
          @start="perform(() => api.start(id), 'Atendimento iniciado')"
          @close="closing = true"
          @edit="editing = true"
        />
      </aside>

      <div class="space-y-4 lg:col-start-1 lg:row-start-2">
        <ClosingReport
          v-if="ticket.closed_at"
          :resolution="ticket.resolution"
          :opened-at="ticket.created_at"
          :closed-at="ticket.closed_at"
        />

        <section class="card p-5 sm:p-6">
          <h2 class="mb-4 flex items-center gap-2 text-sm font-semibold text-ink">
            Respostas
            <span class="rounded-full bg-surface-sunken px-2 py-0.5 text-xs font-medium text-ink-muted tabular-nums">{{ responseCount }}</span>
          </h2>

          <p v-if="responseCount === 0" class="rounded-lg border border-dashed border-line-strong px-4 py-6 text-center text-sm text-ink-muted">
            Ninguém respondeu ainda.
          </p>
          <ul v-else aria-label="Respostas" class="space-y-4">
            <li
              v-for="response in ticket.responses ?? []"
              :key="response.response_id"
              class="flex gap-3"
              :class="isMine(response.author_id) ? 'flex-row-reverse' : ''"
            >
              <UserAvatar :name="nameOf(response.author_id)" size="sm" />
              <div
                class="max-w-[85%] rounded-2xl px-4 py-2.5 sm:max-w-[75%]"
                :class="
                  isMine(response.author_id)
                    ? 'rounded-tr-sm bg-brand-600 text-white'
                    : 'rounded-tl-sm bg-surface-sunken text-ink'
                "
              >
                <p class="text-xs" :class="isMine(response.author_id) ? 'text-brand-100' : 'text-ink-subtle'">
                  {{ nameOf(response.author_id) }} · {{ formatDateTime(response.created_at) }}
                </p>
                <p class="mt-0.5 text-sm whitespace-pre-line break-words">{{ response.content }}</p>
              </div>
            </li>
          </ul>

          <form v-if="can('respond')" class="mt-6 space-y-3 border-t border-line pt-5" novalidate @submit.prevent="sendReply">
            <FormField
              id="reply"
              v-model="reply"
              label="Sua resposta"
              multiline
              :rows="3"
              placeholder="Escreva uma mensagem… (Ctrl + Enter envia)"
              :error="replyError"
              @keydown.ctrl.enter.prevent="sendReply"
              @keydown.meta.enter.prevent="sendReply"
            />
            <div class="flex justify-end">
              <AppButton type="submit" variant="primary" icon="send" :loading="busy">Responder</AppButton>
            </div>
          </form>
        </section>
      </div>
    </article>

    <CloseTicketModal
      v-if="ticket"
      :open="closing"
      :opened-at="ticket.created_at"
      :saving="busy"
      @close="closing = false"
      @confirm="closeTicket"
    />
  </PageLayout>
</template>
