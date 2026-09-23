import { inject, onMounted, ref } from 'vue'

import { useSessionExit } from '../auth/useSessionExit'
import { useToast } from '../composables/useToast'
import { SessionExpiredError, ticketsApiKey, type Responsible, type TicketDetail } from './ticketsApi'

export function useTicketDetail(id: string) {
  const api = inject(ticketsApiKey)!
  const { expire } = useSessionExit()
  const toast = useToast()

  const ticket = ref<TicketDetail | null>(null)
  const responsibles = ref<Responsible[]>([])
  const notFound = ref(false)
  const loading = ref(true)
  const busy = ref(false)
  const actionError = ref('')

  async function handleFailure(error: unknown, explain: (message: string) => void) {
    if (error instanceof SessionExpiredError) {
      await expire()
      return
    }
    explain(error instanceof Error ? error.message : String(error))
  }

  /** Runs an action, reloads the ticket and says whether it worked. */
  async function perform(action: () => Promise<unknown>, success?: string): Promise<boolean> {
    busy.value = true
    try {
      await action()
      actionError.value = ''
      ticket.value = await api.get(id)
      if (success) {
        toast.show(success)
      }
      return true
    } catch (error) {
      await handleFailure(error, (message) => {
        actionError.value = `Não foi possível concluir a ação: ${message}`
      })
      return false
    } finally {
      busy.value = false
    }
  }

  onMounted(async () => {
    try {
      const [loaded, list] = await Promise.all([api.get(id), api.responsibles()])
      ticket.value = loaded
      responsibles.value = list
    } catch (error) {
      await handleFailure(error, () => {
        notFound.value = true
      })
    } finally {
      loading.value = false
    }
  })

  return { api, ticket, responsibles, notFound, loading, busy, actionError, perform }
}
