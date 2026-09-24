import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'

import type { Role } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import { fakeTicketsApi } from '../test/fakeTicketsApi'
import { ticketsApiKey, type TicketsApi } from '../tickets/ticketsApi'
import NewTicketForm from './NewTicketForm.vue'

function renderForm(api: TicketsApi = fakeTicketsApi(), role: Role = 'user') {
  const session = createSession()
  session.start({ id: `${role}-1`, name: 'Alguém', role })
  const rendered = render(NewTicketForm, { global: { provide: { [ticketsApiKey]: api, [sessionKey]: session } } })
  return { api, ...rendered }
}

async function fillAndSubmit() {
  await userEvent.type(screen.getByLabelText('Título'), 'Impressora')
  await userEvent.type(screen.getByLabelText('Descrição'), 'Não imprime')
  await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))
}

describe('NewTicketForm', () => {
  it('opens a ticket with title, description and priority and reports its id', async () => {
    const { api, emitted } = renderForm(fakeTicketsApi({ open: vi.fn().mockResolvedValue('t-1') }))

    await userEvent.selectOptions(screen.getByLabelText('Prioridade'), 'Alta')
    await fillAndSubmit()

    expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'High', auto_assign: true })
    expect(emitted().opened).toEqual([['t-1']])
  })

  it('starts with medium priority and automatic assignment', () => {
    renderForm()

    expect(screen.getByLabelText('Prioridade')).toHaveValue('Medium')
    expect(screen.getByLabelText('Responsável')).toHaveValue('auto')
  })

  it('lets the assignment be decided later', async () => {
    const { api } = renderForm()

    await userEvent.selectOptions(screen.getByLabelText('Responsável'), 'Definir depois')
    await fillAndSubmit()

    expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'Medium' })
  })

  it('does not offer specific agents to a regular user', () => {
    const api = fakeTicketsApi()
    renderForm(api)

    expect(screen.getAllByRole('option', { name: /Automático|Definir depois|Baixa|Média|Alta/ })).toHaveLength(5)
    expect(api.responsibles).not.toHaveBeenCalled()
  })

  it('lets an administrator pick the agent', async () => {
    const api = fakeTicketsApi({ responsibles: vi.fn().mockResolvedValue([{ id: 'agent-2', name: 'Bruno Lima' }]) })
    renderForm(api, 'admin')

    await userEvent.selectOptions(screen.getByLabelText('Responsável'), await screen.findByRole('option', { name: 'Bruno Lima' }))
    await fillAndSubmit()

    expect(api.open).toHaveBeenCalledWith({ title: 'Impressora', description: 'Não imprime', priority: 'Medium', assignee_id: 'agent-2' })
  })

  it('requires a title and a description', async () => {
    const { api } = renderForm()

    await userEvent.click(screen.getByRole('button', { name: 'Abrir chamado' }))

    expect(screen.getByText('Informe o título')).toBeInTheDocument()
    expect(screen.getByText('Informe a descrição')).toBeInTheDocument()
    expect(api.open).not.toHaveBeenCalled()
  })

  it('explains when the ticket cannot be opened', async () => {
    const { emitted } = renderForm(fakeTicketsApi({ open: vi.fn().mockRejectedValue(new Error('internal error')) }))

    await fillAndSubmit()

    expect(await screen.findByRole('alert')).toHaveTextContent('Não foi possível abrir o chamado')
    expect(emitted().opened).toBeUndefined()
  })
})
