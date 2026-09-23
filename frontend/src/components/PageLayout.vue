<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import AppHeader from './AppHeader.vue'
import AppNav from './AppNav.vue'
import BackLink from './BackLink.vue'
import MobileDrawer from './MobileDrawer.vue'
import ToastHost from './ToastHost.vue'

defineProps<{ back?: boolean; backTo?: string }>()

const menuOpen = ref(false)
const route = useRoute()
watch(() => route.fullPath, () => (menuOpen.value = false))
</script>

<template>
  <div class="min-h-dvh">
    <a
      href="#conteudo"
      class="sr-only focus:not-sr-only focus:fixed focus:top-3 focus:left-3 focus:z-[70] focus:rounded-lg focus:bg-brand-600 focus:px-4 focus:py-2 focus:text-white"
      >Pular para o conteúdo</a
    >

    <aside class="fixed inset-y-0 left-0 z-40 hidden w-64 border-r border-line bg-surface lg:block">
      <AppNav />
    </aside>

    <div class="lg:pl-64">
      <AppHeader @open-menu="menuOpen = true" />
      <main id="conteudo" tabindex="-1" class="mx-auto w-full max-w-6xl px-4 pt-6 pb-16 focus:outline-none sm:px-6 lg:px-8 lg:pt-8">
        <BackLink v-if="back" :to="backTo" />
        <slot />
      </main>
    </div>

    <MobileDrawer :open="menuOpen" @close="menuOpen = false">
      <AppNav @navigate="menuOpen = false" />
    </MobileDrawer>
    <ToastHost />
  </div>
</template>
