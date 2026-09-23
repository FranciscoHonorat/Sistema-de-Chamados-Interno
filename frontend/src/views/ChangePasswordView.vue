<script setup lang="ts">
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'

import { AccountRequestError, authServiceKey } from '../auth/authService'
import { sessionKey } from '../auth/session'
import { hasErrors } from '../auth/validateCredentials'
import { validateNewPassword } from '../auth/validateNewPassword'
import AuthCard from '../components/AuthCard.vue'
import BackLink from '../components/BackLink.vue'
import FormField from '../components/FormField.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import { homePathFor } from '../router/homePath'

const authService = inject(authServiceKey)!
const session = inject(sessionKey)!
const router = useRouter()

const forced = session.user.value?.mustChangePassword === true
const current = ref('')
const password = ref('')
const confirmation = ref('')
const errors = ref<Record<string, string | undefined>>({})
const failure = ref('')
const submitting = ref(false)

async function submit() {
  errors.value = {
    ...validateNewPassword(password.value, confirmation.value),
    ...(current.value ? {} : { current: 'Informe a senha atual' }),
  }
  if (hasErrors(errors.value)) {
    return
  }
  failure.value = ''
  submitting.value = true
  try {
    const user = await authService.changePassword(current.value, password.value)
    session.start(user)
    await router.push(homePathFor(user.role))
  } catch (error) {
    failure.value =
      error instanceof AccountRequestError && error.status === 401
        ? 'Senha atual incorreta'
        : 'Não foi possível trocar a senha'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthCard>
    <BackLink v-if="!forced" />
    <form class="space-y-5" novalidate @submit.prevent="submit">
      <div>
        <h2 class="text-lg font-semibold text-ink">Trocar senha</h2>
        <p class="text-sm text-ink-muted">A nova senha vale a partir do próximo acesso.</p>
      </div>
      <AlertMessage v-if="forced" tone="warning">Você entrou com uma senha temporária. Escolha uma nova senha para continuar.</AlertMessage>
      <AlertMessage v-if="failure" role="alert">{{ failure }}</AlertMessage>
      <FormField id="current" v-model="current" label="Senha atual" type="password" autocomplete="current-password" :error="errors.current" />
      <FormField
        id="password"
        v-model="password"
        label="Nova senha"
        type="password"
        autocomplete="new-password"
        hint="Pelo menos 8 caracteres."
        :error="errors.password"
      />
      <FormField
        id="confirmation"
        v-model="confirmation"
        label="Confirmar nova senha"
        type="password"
        autocomplete="new-password"
        :error="errors.confirmation"
      />
      <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Salvar nova senha</AppButton>
    </form>
  </AuthCard>
</template>
