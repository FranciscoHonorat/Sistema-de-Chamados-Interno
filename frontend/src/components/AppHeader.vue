<script setup lang="ts">
import { computed, inject } from 'vue'

import { sessionKey } from '../auth/session'
import { useSessionExit } from '../auth/useSessionExit'
import { roleLabel } from '../employees/roles'
import NotificationBell from './NotificationBell.vue'
import ThemeToggle from './ThemeToggle.vue'
import AppIcon from './ui/AppIcon.vue'
import UserAvatar from './ui/UserAvatar.vue'

const emit = defineEmits<{ openMenu: [] }>()

const session = inject(sessionKey)!
const user = computed(() => session.user.value)
const { logout } = useSessionExit()
</script>

<template>
  <header
    class="sticky top-0 z-30 flex h-16 items-center gap-2 border-b border-line bg-surface/85 px-3 backdrop-blur supports-[backdrop-filter]:bg-surface/70 sm:px-6 lg:px-8"
  >
    <button
      type="button"
      class="grid size-10 place-items-center rounded-lg text-ink-muted hover:bg-surface-sunken hover:text-ink lg:hidden"
      aria-label="Abrir menu"
      @click="emit('openMenu')"
    >
      <AppIcon name="menu" :size="22" />
    </button>
    <span class="text-base font-semibold text-ink lg:hidden">CodeTicket</span>

    <div class="ml-auto flex items-center gap-1">
      <ThemeToggle />
      <NotificationBell />

      <div v-if="user" class="ml-2 hidden items-center gap-2.5 border-l border-line pl-3 sm:flex">
        <UserAvatar :name="user.name" />
        <div class="leading-tight">
          <p class="max-w-40 truncate text-sm font-medium text-ink">{{ user.name }}</p>
          <p class="text-xs text-ink-subtle">{{ roleLabel(user.role) }}</p>
        </div>
      </div>

      <button
        type="button"
        class="ml-1 inline-flex h-10 items-center gap-2 rounded-lg px-2.5 text-sm font-medium text-ink-muted transition-colors hover:bg-surface-sunken hover:text-ink"
        title="Sair"
        @click="logout"
      >
        <AppIcon name="logout" />
        <span class="sr-only sm:not-sr-only">Sair</span>
      </button>
    </div>
  </header>
</template>
