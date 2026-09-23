<script setup lang="ts">
import { inject } from 'vue'

import { authServiceKey, type User } from '../auth/authService'
import { useLoginForm } from '../auth/useLoginForm'
import AuthCard from '../components/AuthCard.vue'
import FormField from '../components/FormField.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import { paths } from '../router/paths'

const emit = defineEmits<{ authenticated: [user: User] }>()

const { username, password, errors, loginError, submitting, submit } = useLoginForm(inject(authServiceKey)!, (user) =>
  emit('authenticated', user),
)
</script>

<template>
  <AuthCard>
    <form class="space-y-5" novalidate @submit.prevent="submit">
      <div>
        <h2 class="text-lg font-semibold text-ink">Entrar</h2>
        <p class="text-sm text-ink-muted">Use o username e a senha da sua conta.</p>
      </div>

      <AlertMessage v-if="loginError" role="alert">{{ loginError }}</AlertMessage>

      <FormField id="username" v-model="username" label="Username" autocomplete="username" :error="errors.username" />
      <FormField
        id="password"
        v-model="password"
        label="Senha"
        type="password"
        autocomplete="current-password"
        :error="errors.password"
      />

      <AppButton type="submit" variant="primary" block size="lg" :loading="submitting">Entrar</AppButton>
    </form>

    <footer class="mt-6 flex items-center justify-between border-t border-line pt-5 text-sm">
      <RouterLink :to="paths.register" class="font-medium text-brand-600 hover:underline dark:text-brand-300">Criar conta</RouterLink>
      <RouterLink :to="paths.recoverPassword" class="font-medium text-ink-muted hover:text-ink hover:underline">Recuperar senha</RouterLink>
    </footer>
  </AuthCard>
</template>
