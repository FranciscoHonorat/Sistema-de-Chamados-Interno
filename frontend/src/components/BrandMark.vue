<script setup lang="ts">
withDefaults(defineProps<{ inverted?: boolean }>(), { inverted: false })
</script>

<template>
  <span class="inline-flex items-center gap-2.5 font-semibold tracking-tight">
    <svg aria-hidden="true" viewBox="0 0 32 32" class="size-8 shrink-0">
      <rect width="32" height="32" rx="8" :fill="inverted ? '#ffffff' : '#15803d'" />
      <path
        d="M9 11.5A2.5 2.5 0 0 1 11.5 9h9a2.5 2.5 0 0 1 2.5 2.5v6a2.5 2.5 0 0 1-2.5 2.5H15l-4.2 3.2c-.5.4-1.3 0-1.3-.6V20A2.5 2.5 0 0 1 9 17.5z"
        :fill="inverted ? '#15803d' : '#ffffff'"
      />
      <path
        d="m13 14.4 1.9 1.9 4.1-4.1"
        fill="none"
        :stroke="inverted ? '#ffffff' : '#15803d'"
        stroke-width="1.8"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
    <span class="text-lg" :class="inverted ? 'text-white' : 'text-ink'">CodeTicket</span>
  </span>
</template>
