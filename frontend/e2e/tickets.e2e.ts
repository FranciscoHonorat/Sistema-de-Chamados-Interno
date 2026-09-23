import { expect, test } from '@playwright/test'

import { detail, loggedInPage, login, openTicket, unique } from './support.ts'

test('a ticket goes from opened to closed across the three roles', async ({ browser }) => {
  const title = unique('Impressora E2E')

  const requester = await loggedInPage(browser, 'usuario')
  const url = await openTicket(requester, title, 'A impressora do 2º andar mostra papel atolado.')
  await expect(detail(requester, 'Status')).toContainText('Aberto')

  const admin = await loggedInPage(browser, 'admin')
  await admin.goto(url)
  await admin.getByLabel('Atribuir a').selectOption({ label: 'Ana Souza' })
  await admin.getByRole('button', { name: 'Atribuir', exact: true }).click()
  await expect(detail(admin, 'Responsável')).toContainText('Ana Souza')

  const agent = await loggedInPage(browser, 'ana')
  await agent.goto(url)
  await agent.getByRole('button', { name: 'Iniciar atendimento' }).click()
  await expect(detail(agent, 'Status')).toContainText('Em andamento')
  await agent.getByLabel('Sua resposta').fill('Vou passar aí em 10 minutos.')
  await agent.getByRole('button', { name: 'Responder' }).click()
  await expect(agent.getByRole('list', { name: 'Respostas' })).toContainText('Vou passar aí em 10 minutos.')

  await requester.reload()
  await expect(detail(requester, 'Status')).toContainText('Em andamento')
  await expect(requester.getByRole('list', { name: 'Respostas' })).toContainText('Vou passar aí em 10 minutos.')
  await requester.getByLabel('Sua resposta').fill('Obrigado, estou na mesa ao lado da janela.')
  await requester.getByRole('button', { name: 'Responder' }).click()
  await expect(requester.getByRole('list', { name: 'Respostas' })).toContainText('Obrigado, estou na mesa ao lado da janela.')

  await agent.reload()
  await agent.getByRole('button', { name: 'Fechar chamado' }).click()
  const closing = agent.getByRole('dialog', { name: 'Fechar chamado' })
  await closing.getByLabel('O que foi feito').fill('Limpei o sensor de papel e testei a impressão.')
  await closing.getByRole('button', { name: 'Confirmar fechamento' }).click()
  await expect(detail(agent, 'Status')).toContainText('Fechado')

  await requester.reload()
  await expect(detail(requester, 'Status')).toContainText('Fechado')
  await expect(requester.getByRole('region', { name: 'Relatório de fechamento' })).toContainText('Limpei o sensor de papel')
  await expect(requester.getByLabel('Sua resposta')).toBeHidden()
})

test('the opening form asks for a title and a description', async ({ page }) => {
  await login(page, 'usuario')
  await page.goto('/chamados')
  await page.getByRole('button', { name: 'Abrir novo chamado' }).click()
  const dialog = page.getByRole('dialog', { name: 'Abrir novo chamado' })

  await dialog.getByRole('button', { name: 'Abrir chamado' }).click()

  await expect(dialog.getByText('Informe o título')).toBeVisible()
  await expect(dialog.getByText('Informe a descrição')).toBeVisible()
})

test('a support agent can hand a ticket to the least busy colleague', async ({ browser }) => {
  const title = unique('VPN E2E')
  const requester = await loggedInPage(browser, 'usuario')
  const url = await openTicket(requester, title, 'O cliente da VPN recusa minhas credenciais.')

  const agent = await loggedInPage(browser, 'bruno')
  await agent.goto(url)
  await expect(detail(agent, 'Responsável')).toContainText('—')
  await agent.getByRole('button', { name: 'Distribuir automaticamente' }).click()

  await expect(detail(agent, 'Responsável')).not.toContainText('—')
})

test('the list finds a ticket by its title', async ({ page }) => {
  const title = unique('Cadeira E2E')
  await login(page, 'usuario')
  await openTicket(page, title, 'O pistão da cadeira quebrou.')

  await page.goto('/chamados')
  await page.getByLabel('Buscar chamados pelo título').fill(title)

  await expect(page.getByRole('link', { name: title })).toBeVisible()
  await page.getByLabel('Buscar chamados pelo título').fill(unique('não existe'))
  await expect(page.getByText('Nenhum chamado encontrado')).toBeVisible()
})

test('a regular user does not see the support actions', async ({ page }) => {
  await login(page, 'usuario')
  await openTicket(page, unique('Monitor E2E'), 'O monitor pisca.')

  await expect(page.getByRole('button', { name: 'Iniciar atendimento' })).toBeHidden()
  await expect(page.getByLabel('Atribuir a')).toBeHidden()
})
