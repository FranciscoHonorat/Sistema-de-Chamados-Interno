<script setup lang="ts">
import { computed, inject } from 'vue'
import { useRoute, useRouter, type RouteLocationRaw } from 'vue-router'

import { sessionKey } from '../auth/session'
import { roleLabel } from '../employees/roles'
import { homePathFor } from '../router/homePath'
import { menuFor, type MenuItem } from '../router/menu'
import { paths } from '../router/paths'
import BrandMark from './BrandMark.vue'
import AppIcon from './ui/AppIcon.vue'

const emit = defineEmits<{ navigate: [] }>()

const session = inject(sessionKey)!
const route = useRoute()
const router = useRouter()

const user = computed(() => session.user.value)
const items = computed<MenuItem[]>(() => {
  if (!user.value) return []
  const home: MenuItem = { label: 'Início', to: homePathFor(user.value.role), icon: 'home', description: '' }
  return [home, ...menuFor(user.value.role)]
})

function isActive(to: RouteLocationRaw): boolean {
  return router.resolve(to).fullPath === route.fullPath
}
</script>

<template>
  <div class="flex h-full flex-col">
    <div class="flex h-16 shrink-0 items-center px-5">
      <RouterLink v-if="user" :to="homePathFor(user.role)" aria-label="CodeTicket, início" @click="emit('navigate')">
        <BrandMark />
      </RouterLink>
    </div>

    <nav aria-label="Menu principal" class="flex-1 space-y-6 overflow-y-auto px-3 py-4">
      <div>
        <p class="px-3 pb-2 text-xs font-semibold tracking-wide text-ink-subtle uppercase">
          {{ user ? roleLabel(user.role) : '' }}
        </p>
        <ul class="space-y-1">
          <li v-for="item in items" :key="item.label">
            <RouterLink
              :to="item.to"
              :aria-current="isActive(item.to) ? 'page' : undefined"
              class="group flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
              :class="
                isActive(item.to)
                  ? 'bg-brand-50 text-brand-700 dark:bg-brand-500/15 dark:text-brand-200'
                  : 'text-ink-muted hover:bg-surface-sunken hover:text-ink'
              "
              @click="emit('navigate')"
            >
              <AppIcon :name="item.icon" />
              {{ item.label }}
            </RouterLink>
          </li>
        </ul>
      </div>

      <div>
        <p class="px-3 pb-2 text-xs font-semibold tracking-wide text-ink-subtle uppercase">Conta</p>
        <ul class="space-y-1">
          <li>
            <RouterLink
              :to="paths.changePassword"
              class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-ink-muted transition-colors hover:bg-surface-sunken hover:text-ink"
              @click="emit('navigate')"
            >
              <AppIcon name="key" />
              Trocar senha
            </RouterLink>
          </li>
        </ul>
      </div>
    </nav>
  </div>
</template>
