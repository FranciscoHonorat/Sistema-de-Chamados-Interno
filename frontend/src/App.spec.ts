import { render, screen, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import { createMemoryHistory } from 'vue-router'

import App from './App.vue'
import { authServiceKey, type User } from './auth/authService'
import { createSession, sessionKey } from './auth/session'
import { createAppRouter } from './router'
import { paths } from './router/paths'
import { fakeAuthService } from './test/fakeAuthService'
import { fakeTicketsApi } from './test/fakeTicketsApi'
import { ticketsApiKey } from './tickets/ticketsApi'

async function renderAppAt(path: string, loggedInAs: User) {
  const session = createSession()
  const router = createAppRouter(session, createMemoryHistory())
  const authService = fakeAuthService({ login: vi.fn().mockResolvedValue(loggedInAs) })
  await router.push(path)
  render(App, {
    global: {
      plugins: [router],
      provide: { [sessionKey]: session, [authServiceKey]: authService, [ticketsApiKey]: fakeTicketsApi() },
    },
  })
  return router
}

describe('App', () => {
  const ana: User = { id: 'agent-1', name: 'Ana Souza', role: 'support' }

  it('takes a support agent from the login to the support home', async () => {
    const router = await renderAppAt(paths.login, ana)

    await userEvent.type(screen.getByLabelText('Username'), 'ana')
    await userEvent.type(screen.getByLabelText('Senha'), 'senha123')
    await userEvent.click(screen.getByRole('button', { name: 'Entrar' }))

    expect(await screen.findByText('Olá, Ana Souza')).toBeInTheDocument()
    expect(within(screen.getByRole('navigation', { name: 'Atalhos' })).getByRole('link', { name: 'Chamados abertos' })).toBeInTheDocument()
    expect(router.currentRoute.value.path).toBe(paths.supportHome)
  })

  it('shows the login to a visitor opening a restricted area', async () => {
    const router = await renderAppAt(paths.adminHome, ana)

    expect(await screen.findByRole('heading', { name: 'CodeTicket' })).toBeInTheDocument()
    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(router.currentRoute.value.path).toBe(paths.login)
  })

  it('shows a page not found for unknown addresses', async () => {
    await renderAppAt('/nada-aqui', ana)

    expect(await screen.findByRole('heading', { name: 'Página não encontrada' })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Voltar para o início' })).toHaveAttribute('href', paths.login)
  })

  it('names the browser tab after the page', async () => {
    await renderAppAt(paths.recoverPassword, ana)

    expect(await screen.findByRole('button', { name: 'Pedir nova senha' })).toBeInTheDocument()
    expect(document.title).toBe('Recuperar senha · CodeTicket')
  })

  it('lets a visitor create an account or ask for a new password', async () => {
    const router = await renderAppAt(paths.login, ana)

    await userEvent.click(screen.getByRole('link', { name: 'Criar conta' }))
    expect(await screen.findByLabelText('Confirmar senha')).toBeInTheDocument()

    await router.push(paths.recoverPassword)
    expect(await screen.findByRole('button', { name: 'Pedir nova senha' })).toBeInTheDocument()
  })

  it('makes whoever logs in with a temporary password choose a new one first', async () => {
    const router = await renderAppAt(paths.login, { ...ana, mustChangePassword: true })

    await userEvent.type(screen.getByLabelText('Username'), 'ana')
    await userEvent.type(screen.getByLabelText('Senha'), 'Temp-1234')
    await userEvent.click(screen.getByRole('button', { name: 'Entrar' }))

    expect(await screen.findByRole('button', { name: 'Salvar nova senha' })).toBeInTheDocument()
    expect(router.currentRoute.value.path).toBe(paths.changePassword)
  })
})
