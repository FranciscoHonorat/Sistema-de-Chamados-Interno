<script setup lang="ts">
import { computed, ref } from 'vue'

import AppIcon from './ui/AppIcon.vue'

const props = defineProps<{
  id: string
  label: string
  type?: 'text' | 'password'
  autocomplete?: string
  error?: string
  hint?: string
  multiline?: boolean
  rows?: number
  placeholder?: string
  maxlength?: number
}>()

const value = defineModel<string>({ default: '' })

const revealed = ref(false)
const inputType = computed(() => (props.type === 'password' && revealed.value ? 'text' : (props.type ?? 'text')))
const describedBy = computed(() => [props.error ? `${props.id}-error` : '', props.hint ? `${props.id}-hint` : ''].filter(Boolean).join(' ') || undefined)
</script>

<template>
  <div>
    <label :for="id" class="field-label">{{ label }}</label>
    <textarea
      v-if="multiline"
      :id="id"
      v-model="value"
      :rows="rows ?? 4"
      :placeholder="placeholder"
      :maxlength="maxlength"
      :aria-invalid="error ? 'true' : undefined"
      :aria-describedby="describedBy"
      class="field-input resize-y"
    />
    <div v-else class="relative">
      <input
        :id="id"
        v-model="value"
        :type="inputType"
        :autocomplete="autocomplete"
        :placeholder="placeholder"
        :maxlength="maxlength"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="describedBy"
        class="field-input"
        :class="type === 'password' ? 'pr-11' : ''"
      />
      <button
        v-if="type === 'password'"
        type="button"
        class="absolute inset-y-0 right-0 grid w-10 place-items-center rounded-r-lg text-ink-subtle hover:text-ink"
        :aria-label="revealed ? 'Ocultar o que foi digitado' : 'Mostrar o que foi digitado'"
        :aria-pressed="revealed"
        tabindex="-1"
        @click="revealed = !revealed"
      >
        <AppIcon :name="revealed ? 'eyeSlash' : 'eye'" :size="18" />
      </button>
    </div>
    <p v-if="hint && !error" :id="`${id}-hint`" class="mt-1.5 text-xs text-ink-subtle">{{ hint }}</p>
    <p v-if="error" :id="`${id}-error`" class="mt-1.5 text-sm text-rose-600 dark:text-rose-400">{{ error }}</p>
  </div>
</template>
