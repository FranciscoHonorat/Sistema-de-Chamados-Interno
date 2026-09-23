<script setup lang="ts">
import { computed, ref } from 'vue'

import { formatDuration } from '../tickets/format'
import AppModal from './AppModal.vue'
import FormField from './FormField.vue'
import AppButton from './ui/AppButton.vue'
import AppIcon from './ui/AppIcon.vue'

const props = defineProps<{ open: boolean; openedAt: string; saving?: boolean }>()
const emit = defineEmits<{ close: []; confirm: [resolution: string] }>()

const resolution = ref('')
const error = ref('')
const elapsed = computed(() => formatDuration(props.openedAt, new Date().toISOString()))

function confirm() {
  if (!resolution.value.trim()) {
    error.value = 'Descreva o que foi feito'
    return
  }
  error.value = ''
  emit('confirm', resolution.value.trim())
}
</script>

<template>
  <AppModal :open="open" title="Fechar chamado" description="O relatório fica registrado no chamado para quem o abriu." @close="emit('close')">
    <form class="space-y-4" novalidate @submit.prevent="confirm">
      <p class="flex items-center gap-2 rounded-lg bg-surface-sunken px-3 py-2.5 text-sm text-ink-muted">
        <AppIcon name="clock" :size="18" />Tempo de atendimento até agora: {{ elapsed }}
      </p>
      <FormField id="resolution" v-model="resolution" label="O que foi feito" multiline :error="error" placeholder="Ex.: troquei o cabo de rede e testei com o usuário." />
      <AppButton type="submit" variant="primary" block size="lg" :loading="saving">Confirmar fechamento</AppButton>
    </form>
  </AppModal>
</template>
