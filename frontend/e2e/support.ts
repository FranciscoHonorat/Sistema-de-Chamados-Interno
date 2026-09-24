import { expect, type Browser, type Page } from '@playwright/test'

export const demoPassword = 'senha123'

export function unique(prefix: string): string {
  return `${prefix} ${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`
}

export async function login(page: Page, username: string, password = demoPassword): Promise<void> {
  await page.goto('/')
  await page.getByLabel('Username').fill(username)
  await page.getByLabel('Senha', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page.getByRole('button', { name: 'Sair' })).toBeVisible()
}

export async function loggedInPage(browser: Browser, username: string): Promise<Page> {
  const page = await (await browser.newContext()).newPage()
  await login(page, username)
  return page
}

export async function openTicket(page: Page, title: string, description: string, assignee?: string): Promise<string> {
  await page.goto('/chamados')
  await page.getByRole('button', { name: 'Abrir novo chamado' }).click()
  const dialog = page.getByRole('dialog', { name: 'Abrir novo chamado' })
  await dialog.getByLabel('Título').fill(title)
  await dialog.getByLabel('Descrição').fill(description)
  if (assignee) {
    await dialog.getByLabel('Responsável').selectOption({ label: assignee })
  }
  await dialog.getByRole('button', { name: 'Abrir chamado' }).click()
  await expect(dialog).toBeHidden()
  await page.getByRole('link', { name: title }).click()
  await expect(page.getByRole('heading', { name: title })).toBeVisible()
  return page.url()
}

export function detail(page: Page, field: string) {
  return page.getByRole('list', { name: 'Detalhes' }).getByRole('listitem').filter({ hasText: field })
}
