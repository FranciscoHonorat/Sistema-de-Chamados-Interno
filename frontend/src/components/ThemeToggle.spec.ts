import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'

import ThemeToggle from './ThemeToggle.vue'

describe('ThemeToggle', () => {
  afterEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
  })

  it('goes from the system theme to light and dark, remembering the choice', async () => {
    render(ThemeToggle)

    await userEvent.click(screen.getByRole('button', { name: /Tema do sistema/ }))
    expect(screen.getByRole('button', { name: /Tema claro/ })).toBeInTheDocument()
    expect(document.documentElement).not.toHaveClass('dark')

    await userEvent.click(screen.getByRole('button', { name: /Tema claro/ }))
    expect(screen.getByRole('button', { name: /Tema escuro/ })).toBeInTheDocument()
    expect(document.documentElement).toHaveClass('dark')
    expect(localStorage.getItem('sys-called:theme')).toBe('dark')
  })
})
