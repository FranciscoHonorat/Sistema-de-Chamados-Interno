<script setup lang="ts">
import { inject, ref } from 'vue'

import { AccountRequestError, authServiceKey } from '../auth/authService'
import { hasErrors } from '../auth/validateCredentials'
import { validateNewPassword } from '../auth/validateNewPassword'
import AuthCard from '../components/AuthCard.vue'
import AuthConfirmation from '../components/AuthConfirmation.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import { paths } from '../router/paths'

const authService = inject(authServiceKey)!

const name = ref('')
const username = ref('')
const password = ref('')
const confirmation = ref('')
const errors = ref<Record<string, string | undefined>>({})
const failure = ref('')
const created = ref(false)
const submitting = ref(false)

function validate() {
  const found: Record<string, string | undefined> = { ...validateNewPassword(password.value, confirmation.value) }
  if (!name.value.trim()) {
    found.name = 'Informe o nome'
  }
  if (!username.value.trim()) {
    found.username = 'Informe o username'
  }
  return found
}

async function submit() {
  errors.value = validate()
  if (hasErrors(errors.value)) {
    return
  }
  failure.value = ''
  submitting.value = true
  try {
    await authService.signUp({ name: name.value.trim(), username: username.value.trim(), password: password.value })
    created.value = true
  } catch (error) {
    failure.value =
      error instanceof AccountRequestError && error.status === 409
        ? 'Este username já está em uso'
        : 'Não foi possível criar a conta'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthCard>
    <BackLink :to="paths.login" />
    <AuthConfirmation v-if="created">Conta criada! Aguarde a aprovação do administrador para entrar.</AuthConfirmation>
    <form v-else class="space-y-5" novalidate @submit.prevent="submit">
      <div>
        <h2 class="text-lg font-semibold text-ink">Criar conta</h2>
        <p class="text-sm text-ink-muted">Depois de criada, a conta precisa ser aprovada pelo administrador.</p>
      </div>
      <AlertMessage v-if="failure" role="alert">{{ failure }}</AlertMessage>
      <FormField id="name" v-model="name" label="Nome" autocomplete="name" :error="errors.name" />
      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="errors.username" />
      <FormField
        id="password"
        v-model="password"
        label="Senha"
        type="password"
        autocomplete="new-password"
        hint="Pelo menos 8 caracteres."
        :error="errors.password"
      />
      <FormField
        id="confirmation"
        v-model="confirmation"
        label="Confirmar senha"
        type="password"
        autocomplete="new-password"
        :error="errors.confirmation"
      />
      <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Criar conta</AppButton>
    </form>
  </AuthCard>
</template>
