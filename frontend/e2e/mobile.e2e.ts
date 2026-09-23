import { expect, test } from '@playwright/test'

import { login } from './support.ts'

test('the menu opens as a drawer on a phone', async ({ page }) => {
  await login(page, 'usuario')

  await page.getByRole('button', { name: 'Abrir menu' }).click()
  const menu = page.getByRole('dialog', { name: 'Menu' })
  await menu.getByRole('link', { name: 'Meus chamados' }).click()

  await expect(page).toHaveURL('/chamados')
  await expect(menu).toBeHidden()
  await expect(page.getByRole('heading', { name: 'Chamados', level: 1 })).toBeVisible()
})
