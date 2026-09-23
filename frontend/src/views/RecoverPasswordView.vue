<script setup lang="ts">
import { inject, ref } from 'vue'

import { authServiceKey } from '../auth/authService'
import AuthCard from '../components/AuthCard.vue'
import AuthConfirmation from '../components/AuthConfirmation.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import { paths } from '../router/paths'

const authService = inject(authServiceKey)!

const username = ref('')
const error = ref('')
const failure = ref('')
const sent = ref(false)
const submitting = ref(false)

async function submit() {
  error.value = username.value.trim() ? '' : 'Informe o username'
  if (error.value) {
    return
  }
  failure.value = ''
  submitting.value = true
  try {
    await authService.requestPasswordReset(username.value.trim())
    sent.value = true
  } catch {
    failure.value = 'Não foi possível enviar o pedido'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthCard>
    <BackLink :to="paths.login" />
    <AuthConfirmation v-if="sent">Pedido enviado! O administrador vai gerar uma senha temporária para você.</AuthConfirmation>
    <form v-else class="space-y-5" novalidate @submit.prevent="submit">
      <div>
        <h2 class="text-lg font-semibold text-ink">Recuperar senha</h2>
        <p class="text-sm text-ink-muted">Informe seu username e o administrador vai gerar uma senha temporária.</p>
      </div>
      <AlertMessage v-if="failure" role="alert">{{ failure }}</AlertMessage>
      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="error" />
      <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Pedir nova senha</AppButton>
    </form>
  </AuthCard>
</template>
