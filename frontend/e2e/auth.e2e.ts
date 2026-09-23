import { expect, test } from '@playwright/test'

import { demoPassword, login, unique } from './support.ts'

test.describe('login', () => {
  test('refuses a wrong password', async ({ page }) => {
    await page.goto('/')
    await page.getByLabel('Username').fill('usuario')
    await page.getByLabel('Senha', { exact: true }).fill('senha-errada')
    await page.getByRole('button', { name: 'Entrar' }).click()

    await expect(page.getByRole('alert')).toContainText('Username ou senha inválidos')
    await expect(page).toHaveURL('/')
  })

  for (const { username, home, name } of [
    { username: 'usuario', home: '/usuario', name: 'Usuário Padrão' },
    { username: 'ana', home: '/suporte', name: 'Ana Souza' },
    { username: 'admin', home: '/admin', name: 'Administradora' },
  ]) {
    test(`takes ${username} to their home`, async ({ page }) => {
      await login(page, username)

      await expect(page).toHaveURL(home)
      await expect(page.getByRole('heading', { level: 1 })).toHaveText(`Olá, ${name}`)
    })
  }

  test('keeps a user out of another role area', async ({ page }) => {
    await login(page, 'usuario')

    await page.goto('/admin/usuarios')

    await expect(page).toHaveURL('/usuario')
  })

  test('logs out and protects the pages again', async ({ page }) => {
    await login(page, 'usuario')

    await page.getByRole('button', { name: 'Sair' }).click()
    await expect(page.getByRole('button', { name: 'Entrar' })).toBeVisible()

    await page.goto('/chamados')
    await expect(page).toHaveURL('/')
  })

  test('keeps the session after a reload', async ({ page }) => {
    await login(page, 'ana')

    await page.reload()

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Olá, Ana Souza')
  })
})

test('a new account can log in only after the administrator approves it', async ({ page, browser }) => {
  const name = unique('Pessoa E2E')
  const username = name.toLowerCase().replaceAll(' ', '.')

  await page.goto('/criar-conta')
  await page.getByLabel('Nome').fill(name)
  await page.getByLabel('Username').fill(username)
  await page.getByLabel('Senha', { exact: true }).fill(demoPassword)
  await page.getByLabel('Confirmar senha').fill(demoPassword)
  await page.getByRole('button', { name: 'Criar conta' }).click()
  await expect(page.getByText('Conta criada! Aguarde a aprovação do administrador para entrar.')).toBeVisible()

  await page.goto('/')
  await page.getByLabel('Username').fill(username)
  await page.getByLabel('Senha', { exact: true }).fill(demoPassword)
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page.getByRole('alert')).toContainText('Sua conta ainda aguarda a aprovação do administrador')

  const admin = await (await browser.newContext()).newPage()
  await login(admin, 'admin')
  await admin.getByRole('button', { name: /Notificações/ }).first().click()
  await expect(admin.getByText(`Nova conta aguardando aprovação: ${name}`)).toBeVisible()
  await admin.goto('/admin/usuarios')
  await admin.getByLabel('Buscar por nome ou username').fill(username)
  await admin.getByRole('button', { name: `Aprovar ${name}` }).click()
  await expect(admin.getByRole('button', { name: `Aprovar ${name}` })).toBeHidden()

  await login(page, username)
  await expect(page).toHaveURL('/usuario')
})
