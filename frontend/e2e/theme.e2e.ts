import { expect, test } from '@playwright/test'

test.use({ colorScheme: 'light' })

test('the chosen theme survives a reload', async ({ page }) => {
  await page.goto('/')
  const html = page.locator('html')
  await expect(html).not.toHaveClass(/dark/)

  const toggle = page.getByRole('button', { name: /clique para trocar/ })
  while (!(await html.getAttribute('class'))?.includes('dark')) {
    await toggle.click()
  }
  await page.reload()

  await expect(html).toHaveClass(/dark/)
  await expect(page.getByRole('button', { name: 'Tema escuro (clique para trocar)' })).toBeVisible()
})
