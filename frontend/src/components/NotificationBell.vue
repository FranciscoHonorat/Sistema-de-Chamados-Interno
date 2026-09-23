<script setup lang="ts">
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'

import { usePolling } from '../composables/usePolling'
import { paths, ticketDetailPath } from '../router/paths'
import { formatDate } from '../tickets/format'
import { ticketsApiKey, type Notification, type NotificationFeed } from '../tickets/ticketsApi'
import AppIcon from './ui/AppIcon.vue'

const props = withDefaults(defineProps<{ pollEveryMs?: number }>(), { pollEveryMs: 30_000 })

const api = inject(ticketsApiKey)!

const feed = ref<NotificationFeed>({ unread: 0, items: [] })
const open = ref(false)
const root = ref<HTMLElement | null>(null)

const label = computed(() =>
  feed.value.unread > 0 ? `Notificações (${feed.value.unread} não lidas)` : 'Notificações',
)
const badge = computed(() => (feed.value.unread > 99 ? '99+' : String(feed.value.unread)))

async function refresh() {
  feed.value = await api.notifications().catch(() => feed.value)
}

async function toggle() {
  open.value = !open.value
  if (open.value && feed.value.unread > 0) {
    await api.markNotificationsRead()
    feed.value = { unread: 0, items: feed.value.items.map((item) => ({ ...item, unread: false })) }
  }
}

function destination(item: Notification): string {
  return item.ticket_id ? ticketDetailPath(item.ticket_id) : paths.users
}

function closeOnOutsideClick(event: MouseEvent) {
  if (root.value && !root.value.contains(event.target as Node)) {
    open.value = false
  }
}

watch(open, (isOpen) => {
  if (isOpen) {
    document.addEventListener('click', closeOnOutsideClick)
  } else {
    document.removeEventListener('click', closeOnOutsideClick)
  }
})
onBeforeUnmount(() => document.removeEventListener('click', closeOnOutsideClick))

usePolling(refresh, props.pollEveryMs)
</script>

<template>
  <div ref="root" class="relative" @keydown.esc="open = false">
    <button
      type="button"
      :aria-label="label"
      :aria-expanded="open"
      class="relative grid size-10 place-items-center rounded-lg text-ink-muted transition-colors hover:bg-surface-sunken hover:text-ink"
      :class="open ? 'bg-surface-sunken text-ink' : ''"
      @click="toggle"
    >
      <AppIcon name="bell" :size="22" />
      <span
        v-if="feed.unread > 0"
        class="absolute top-1 right-1 grid h-4.5 min-w-4.5 place-items-center rounded-full bg-rose-600 px-1 text-[0.65rem] leading-none font-semibold text-white ring-2 ring-surface"
        >{{ badge }}</span
      >
    </button>

    <section
      v-if="open"
      aria-label="Notificações"
      class="card fixed inset-x-3 top-[4.25rem] z-40 animate-slide-up overflow-hidden shadow-overlay sm:absolute sm:inset-x-auto sm:top-auto sm:right-0 sm:mt-2 sm:w-96"
    >
      <header class="flex items-center justify-between border-b border-line px-4 py-3">
        <h2 class="text-sm font-semibold text-ink">Notificações</h2>
        <span class="text-xs text-ink-subtle">{{ feed.items.length }} no total</span>
      </header>
      <div v-if="feed.items.length === 0" class="flex flex-col items-center gap-2 px-4 py-10 text-center">
        <AppIcon name="bell" :size="28" class="text-ink-subtle" />
        <p class="text-sm text-ink-muted">Nenhuma notificação</p>
      </div>
      <ul v-else class="max-h-[60vh] divide-y divide-line overflow-y-auto sm:max-h-96">
        <li v-for="item in feed.items" :key="item.id">
          <RouterLink
            :to="destination(item)"
            class="flex gap-3 px-4 py-3 text-sm transition-colors hover:bg-surface-muted"
            :class="item.unread ? 'font-medium text-ink' : 'text-ink-muted'"
            @click="open = false"
          >
            <span
              aria-hidden="true"
              class="mt-1.5 size-2 shrink-0 rounded-full"
              :class="item.unread ? 'bg-brand-500' : 'bg-transparent'"
            ></span>
            <span class="min-w-0 flex-1">
              {{ item.message }}
              <span class="mt-0.5 block text-xs font-normal text-ink-subtle">{{ formatDate(item.created_at) }}</span>
            </span>
          </RouterLink>
        </li>
      </ul>
    </section>
  </div>
</template>
