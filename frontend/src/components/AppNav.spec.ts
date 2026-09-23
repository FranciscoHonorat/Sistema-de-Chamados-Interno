import { render, screen, within } from '@testing-library/vue'
import { createMemoryHistory, createRouter } from 'vue-router'

import type { User } from '../auth/authService'
import { createSession, sessionKey } from '../auth/session'
import AppNav from './AppNav.vue'

async function renderNavFor(user: User, at: string) {
  const session = createSession()
  session.start(user)
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  })
  await router.push(at)
  render(AppNav, { global: { plugins: [router], provide: { [sessionKey]: session } } })
}

function menu() {
  return within(screen.getByRole('navigation', { name: 'Menu principal' }))
}

describe('AppNav', () => {
  it('lists the home, the menu of the role and the account settings', async () => {
    await renderNavFor({ id: 'admin-1', name: 'Administradora', role: 'admin' }, '/admin')

    expect(menu().getAllByRole('link').map((link) => link.textContent?.trim())).toEqual([
      'Início',
      'Todos os chamados',
      'Usuários',
      'Suportes',
      'Trocar senha',
    ])
  })

  it('marks the page being shown, telling filtered lists apart', async () => {
    await renderNavFor({ id: 'agent-1', name: 'Ana Souza', role: 'support' }, '/chamados?status=In+Progress')

    expect(menu().getByRole('link', { name: 'Em atendimento' })).toHaveAttribute('aria-current', 'page')
    expect(menu().getByRole('link', { name: 'Chamados abertos' })).not.toHaveAttribute('aria-current')
  })
})
