<script setup lang="ts">
import { ref } from 'vue'

import type { TicketAction } from '../tickets/permissions'
import { priorityOptions } from '../tickets/priority'
import SelectField from './SelectField.vue'
import AppButton from './ui/AppButton.vue'

const props = defineProps<{
  actions: TicketAction[]
  responsibles: ReadonlyArray<{ value: string; label: string }>
  busy?: boolean
  currentPriority?: string
}>()

const emit = defineEmits<{
  edit: []
  assign: [assigneeId: string]
  autoAssign: []
  changePriority: [priority: string]
  start: []
  close: []
}>()

const selectedAssignee = ref('')
const selectedPriority = ref(props.currentPriority ?? 'Medium')

const can = (action: TicketAction) => props.actions.includes(action)
</script>

<template>
  <section aria-label="Ações" class="card divide-y divide-line">
    <h2 class="px-5 py-4 text-sm font-semibold text-ink">Ações</h2>

    <div v-if="can('edit')" class="px-5 py-4">
      <AppButton size="sm" variant="ghost" icon="pencil" :disabled="busy" @click="emit('edit')">Editar</AppButton>
    </div>
    <div v-if="can('manage')" class="space-y-4 px-5 py-4">
      <div class="space-y-2">
        <SelectField id="assignee" v-model="selectedAssignee" label="Atribuir a" placeholder="Escolha um atendente" :options="responsibles" />
        <div class="flex flex-wrap gap-2">
          <AppButton size="sm" :disabled="busy || !selectedAssignee" @click="emit('assign', selectedAssignee)">Atribuir</AppButton>
          <AppButton size="sm" variant="ghost" icon="sparkles" :disabled="busy" @click="emit('autoAssign')">Distribuir automaticamente</AppButton>
        </div>
      </div>
      <div class="space-y-2">
        <SelectField id="new-priority" v-model="selectedPriority" label="Nova prioridade" :options="priorityOptions" />
        <AppButton size="sm" icon="flag" :disabled="busy" @click="emit('changePriority', selectedPriority)">Alterar prioridade</AppButton>
      </div>
    </div>

    <div v-if="can('start') || can('close')" class="grid gap-2 px-5 py-4">
      <AppButton v-if="can('start')" variant="primary" icon="play" block :disabled="busy" @click="emit('start')">Iniciar atendimento</AppButton>
      <AppButton v-if="can('close')" icon="check" block :disabled="busy" @click="emit('close')">Fechar chamado</AppButton>
    </div>

  </section>
</template>
