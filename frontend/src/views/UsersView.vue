<script setup lang="ts">
import { computed, inject, ref } from 'vue'

import AppModal from '../components/AppModal.vue'
import PageLayout from '../components/PageLayout.vue'
import AlertMessage from '../components/ui/AlertMessage.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppIcon from '../components/ui/AppIcon.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'
import { useLoad } from '../composables/useLoad'
import { useToast } from '../composables/useToast'
import { employeesApiKey, type Employee } from '../employees/employeesApi'
import { needsAttention, roleLabel, situationLabel } from '../employees/roles'

const api = inject(employeesApiKey)!
const toast = useToast()

const { data: employees, failed, loading, reload } = useLoad<Employee[]>(() => api.list(), [])
const issued = ref<{ name: string; password: string } | null>(null)
const actionFailed = ref(false)
const working = ref<string | null>(null)
const copied = ref(false)
const search = ref('')

const pending = computed(() => employees.value.filter(needsAttention).length)
const visible = computed(() => {
  const term = search.value.trim().toLocaleLowerCase('pt-BR')
  return term
    ? employees.value.filter((e) => `${e.name} ${e.username}`.toLocaleLowerCase('pt-BR').includes(term))
    : employees.value
})

async function run(employee: Employee, action: () => Promise<void>) {
  actionFailed.value = false
  working.value = employee.id
  try {
    await action()
  } catch {
    actionFailed.value = true
  } finally {
    working.value = null
  }
}

function approve(employee: Employee) {
  return run(employee, async () => {
    await api.approve(employee.id)
    toast.show(`Conta de ${employee.name} aprovada`)
    await reload()
  })
}

function issueTemporaryPassword(employee: Employee) {
  return run(employee, async () => {
    copied.value = false
    issued.value = { name: employee.name, password: await api.issueTemporaryPassword(employee.id) }
    await reload()
  })
}

async function copyPassword() {
  if (!issued.value) return
  try {
    await navigator.clipboard.writeText(issued.value.password)
    copied.value = true
  } catch {
    copied.value = false
  }
}

const roleTones: Record<string, string> = {
  admin: 'bg-violet-50 text-violet-700 ring-violet-600/20 dark:bg-violet-400/10 dark:text-violet-300 dark:ring-violet-400/30',
  support: 'bg-sky-50 text-sky-700 ring-sky-600/20 dark:bg-sky-400/10 dark:text-sky-300 dark:ring-sky-400/30',
  user: 'bg-surface-sunken text-ink-muted ring-line-strong',
}
</script>

<template>
  <PageLayout back>
    <PageHeader
      title="Usuários"
      :description="pending > 0 ? `${pending} ${pending === 1 ? 'conta pede' : 'contas pedem'} sua atenção.` : 'Contas, perfis e situação de cada funcionário.'"
    />

    <AlertMessage v-if="failed" role="alert">Não foi possível carregar os usuários</AlertMessage>
    <template v-else>
      <AlertMessage v-if="actionFailed" role="alert" class="mb-4">Não foi possível concluir a ação</AlertMessage>

      <div class="relative mb-4 w-full sm:w-72">
        <AppIcon name="search" :size="18" class="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-ink-subtle" />
        <input v-model="search" type="search" aria-label="Buscar por nome ou username" placeholder="Buscar por nome ou username…" class="field-input pl-10" />
      </div>

      <div v-if="loading && employees.length === 0" class="card space-y-4 p-5" aria-hidden="true">
        <div v-for="n in 4" :key="n" class="flex animate-pulse items-center gap-3">
          <div class="size-9 rounded-full bg-surface-sunken"></div>
          <div class="h-4 flex-1 rounded bg-surface-sunken"></div>
        </div>
      </div>

      <div v-else class="md:overflow-hidden md:rounded-xl md:border md:border-line md:bg-surface md:shadow-card">
        <table class="responsive-table">
          <thead class="bg-surface-muted">
            <tr>
              <th scope="col">Nome</th>
              <th scope="col">Username</th>
              <th scope="col">Perfil</th>
              <th scope="col">Situação</th>
              <th scope="col"><span class="sr-only">Ações</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="employee in visible" :key="employee.id">
              <td class="block pb-2 md:table-cell">
                <span class="inline-flex items-center gap-3 font-semibold text-ink md:font-medium"
                  ><UserAvatar :name="employee.name" />{{ employee.name }}</span
                >
              </td>
              <td data-label="Username" class="text-ink-muted">{{ employee.username }}</td>
              <td data-label="Perfil">
                <span class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset" :class="roleTones[employee.role]">{{
                  roleLabel(employee.role)
                }}</span>
              </td>
              <td data-label="Situação">
                <span
                  class="inline-flex items-center gap-1.5 text-sm"
                  :class="needsAttention(employee) ? 'font-medium text-amber-700 dark:text-amber-300' : 'text-ink-muted'"
                  ><span
                    aria-hidden="true"
                    class="size-1.5 rounded-full"
                    :class="needsAttention(employee) ? 'bg-amber-500' : 'bg-emerald-500'"
                  ></span
                  >{{ situationLabel(employee) }}</span
                >
              </td>
              <td class="pt-3 md:pt-3 md:text-right">
                <AppButton
                  v-if="employee.status === 'pending'"
                  size="sm"
                  variant="primary"
                  icon="check"
                  :aria-label="`Aprovar ${employee.name}`"
                  :loading="working === employee.id"
                  class="w-full md:w-auto"
                  @click="approve(employee)"
                  >Aprovar</AppButton
                >
                <AppButton
                  v-else
                  size="sm"
                  icon="key"
                  :aria-label="`Gerar senha temporária para ${employee.name}`"
                  :loading="working === employee.id"
                  class="w-full md:w-auto"
                  @click="issueTemporaryPassword(employee)"
                  >Gerar senha temporária</AppButton
                >
              </td>
            </tr>
          </tbody>
        </table>
        <p v-if="visible.length === 0" class="px-5 py-10 text-center text-sm text-ink-muted">Ninguém combina com a busca.</p>
      </div>
    </template>

    <AppModal :open="issued !== null" title="Senha temporária" @close="issued = null">
      <div v-if="issued" class="space-y-4">
        <div class="flex items-center gap-2 rounded-xl border border-dashed border-line-strong bg-surface-sunken p-3">
          <p class="flex-1 text-center font-mono text-xl tracking-widest text-ink select-all">{{ issued.password }}</p>
          <AppButton size="sm" variant="ghost" :aria-label="copied ? 'Senha copiada' : 'Copiar a senha'" @click="copyPassword">{{
            copied ? 'Copiada' : 'Copiar'
          }}</AppButton>
        </div>
        <AlertMessage tone="info">
          Entregue esta senha pessoalmente: ela só aparece agora. {{ issued.name }} vai precisar trocá-la no próximo acesso.
        </AlertMessage>
        <AppButton variant="primary" block @click="issued = null">Fechar</AppButton>
      </div>
    </AppModal>
  </PageLayout>
</template>
