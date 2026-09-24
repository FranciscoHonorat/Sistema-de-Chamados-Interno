<script setup lang="ts">
defineProps<{
  id: string
  label: string
  options: ReadonlyArray<{ value: string; label: string }>
  placeholder?: string
  hint?: string
}>()

const value = defineModel<string>({ default: '' })
</script>

<template>
  <div>
    <label :for="id" class="field-label">{{ label }}</label>
    <select :id="id" v-model="value" class="field-input appearance-auto" :aria-describedby="hint ? `${id}-hint` : undefined">
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
    </select>
    <p v-if="hint" :id="`${id}-hint`" class="mt-1.5 text-xs text-ink-subtle">{{ hint }}</p>
  </div>
</template>
