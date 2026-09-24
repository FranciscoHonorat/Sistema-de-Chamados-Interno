<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'

import { sessionKey } from '../auth/session'
import { priorityOptions } from '../tickets/priority'
import { ticketsApiKey, type NewTicket, type Responsible } from '../tickets/ticketsApi'
import FormField from './FormField.vue'
import SelectField from './SelectField.vue'
import AlertMessage from './ui/AlertMessage.vue'
import AppButton from './ui/AppButton.vue'

const emit = defineEmits<{ opened: [ticketId: string] }>()

const api = inject(ticketsApiKey)!
const session = inject(sessionKey)!

const AUTO = 'auto'
const LATER = 'later'

const title = ref('')
const description = ref('')
const priority = ref('Medium')
const assignee = ref(AUTO)
const responsibles = ref<Responsible[]>([])
const titleError = ref('')
const descriptionError = ref('')
const submitError = ref('')
const submitting = ref(false)

const canChooseAssignee = computed(() => session.user.value?.role === 'admin')

const assigneeOptions = computed(() => [
  { value: AUTO, label: 'Automático (quem tem menos chamados em aberto)' },
  { value: LATER, label: 'Definir depois' },
  ...responsibles.value.map((r) => ({ value: r.id, label: r.name })),
])

onMounted(async () => {
  if (!canChooseAssignee.value) {
    return
  }
  try {
    responsibles.value = await api.responsibles()
  } catch {
    responsibles.value = []
  }
})

function assignment(): Pick<NewTicket, 'assignee_id' | 'auto_assign'> {
  if (assignee.value === AUTO) {
    return { auto_assign: true }
  }
  if (assignee.value === LATER) {
    return {}
  }
  return { assignee_id: assignee.value }
}

async function submit() {
  titleError.value = title.value.trim() ? '' : 'Informe o título'
  descriptionError.value = description.value.trim() ? '' : 'Informe a descrição'
  if (titleError.value || descriptionError.value) {
    return
  }

  submitting.value = true
  submitError.value = ''
  try {
    emit('opened', await api.open({ title: title.value, description: description.value, priority: priority.value, ...assignment() }))
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
    <SelectField
      id="new-assignee"
      v-model="assignee"
      label="Responsável"
      :options="assigneeOptions"
      :hint="canChooseAssignee ? undefined : 'Só o administrador escolhe um atendente específico.'"
    />

    <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Abrir chamado</AppButton>
  </form>
</template>
