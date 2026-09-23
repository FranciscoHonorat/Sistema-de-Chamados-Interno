<script setup lang="ts">
import { ref } from 'vue'

import FormField from './FormField.vue'
import AppButton from './ui/AppButton.vue'

const props = defineProps<{ title: string; description: string; saving?: boolean }>()
const emit = defineEmits<{ save: [changes: { title: string; description: string }]; cancel: [] }>()

const draftTitle = ref(props.title)
const draftDescription = ref(props.description)
</script>

<template>
  <form
    class="space-y-4 rounded-xl border border-brand-200 bg-brand-50/50 p-4 dark:border-brand-500/30 dark:bg-brand-500/5"
    @submit.prevent="emit('save', { title: draftTitle, description: draftDescription })"
  >
    <FormField id="edit-title" v-model="draftTitle" label="Título" :maxlength="200" />
    <FormField id="edit-description" v-model="draftDescription" label="Descrição" multiline />
    <div class="flex flex-wrap justify-end gap-2">
      <AppButton variant="ghost" @click="emit('cancel')">Cancelar</AppButton>
      <AppButton type="submit" variant="primary" :loading="saving">Salvar</AppButton>
    </div>
  </form>
</template>
