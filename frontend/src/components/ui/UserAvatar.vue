<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{ name?: string; size?: 'sm' | 'md' }>(), { size: 'md' })

// The initials are drawn with CSS content, so the avatar adds no text to the
// accessible name or to the text of the element around it.
const initials = computed(() =>
  (props.name ?? '?')
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join(''),
)

const palette = ['bg-brand-600', 'bg-violet-600', 'bg-amber-600', 'bg-rose-600', 'bg-sky-600', 'bg-cyan-700']
const color = computed(() => {
  const name = props.name ?? ''
  let hash = 0
  for (const char of name) hash = (hash * 31 + char.charCodeAt(0)) >>> 0
  return palette[hash % palette.length]
})
</script>

<template>
  <span
    aria-hidden="true"
    :data-initials="initials"
    class="inline-grid shrink-0 place-items-center rounded-full font-semibold text-white before:content-[attr(data-initials)]"
    :class="[color, size === 'sm' ? 'size-7 text-[0.65rem]' : 'size-9 text-xs']"
  ></span>
</template>
