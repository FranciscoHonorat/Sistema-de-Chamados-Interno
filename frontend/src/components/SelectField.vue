<script setup lang="ts">
defineProps<{
  id: string
  label: string
  options: ReadonlyArray<{ value: string; label: string }>
  placeholder?: string
}>()

const value = defineModel<string>({ default: '' })
</script>

<template>
  <div>
    <label :for="id" class="field-label">{{ label }}</label>
    <select :id="id" v-model="value" class="field-input appearance-auto">
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
    </select>
  </div>
</template>
