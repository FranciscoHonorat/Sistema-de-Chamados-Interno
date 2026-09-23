<script setup lang="ts">
import { inject, ref } from 'vue'

import { priorityOptions } from '../tickets/priority'
import { ticketsApiKey } from '../tickets/ticketsApi'
import FormField from './FormField.vue'
import SelectField from './SelectField.vue'
import AlertMessage from './ui/AlertMessage.vue'
import AppButton from './ui/AppButton.vue'

const emit = defineEmits<{ opened: [ticketId: string] }>()

const api = inject(ticketsApiKey)!

const title = ref('')
const description = ref('')
const priority = ref('Medium')
const titleError = ref('')
const descriptionError = ref('')
const submitError = ref('')
const submitting = ref(false)

async function submit() {
  titleError.value = title.value.trim() ? '' : 'Informe o título'
  descriptionError.value = description.value.trim() ? '' : 'Informe a descrição'
  if (titleError.value || descriptionError.value) {
    return
  }

  submitting.value = true
  submitError.value = ''
  try {
    emit('opened', await api.open({ title: title.value, description: description.value, priority: priority.value }))
  } catch {
    submitError.value = 'Não foi possível abrir o chamado'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <form class="space-y-4" novalidate @submit.prevent="submit">
    <AlertMessage v-if="submitError" role="alert">{{ submitError }}</AlertMessage>

    <FormField id="title" v-model="title" label="Título" :error="titleError" :maxlength="200" placeholder="Ex.: Impressora do 2º andar não imprime" />
    <FormField
      id="description"
      v-model="description"
      label="Descrição"
      multiline
      :error="descriptionError"
      hint="Conte o que aconteceu, onde e desde quando."
    />
    <SelectField id="priority" v-model="priority" label="Prioridade" :options="priorityOptions" />

    <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Abrir chamado</AppButton>
  </form>
</template>
