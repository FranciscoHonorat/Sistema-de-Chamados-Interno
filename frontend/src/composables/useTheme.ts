import { ref } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

const storageKey = 'sys-called:theme'
const order: ThemeMode[] = ['system', 'light', 'dark']

function read(): ThemeMode {
  try {
    const saved = localStorage.getItem(storageKey)
    return saved === 'light' || saved === 'dark' ? saved : 'system'
  } catch {
    return 'system'
  }
}

function systemPrefersDark(): boolean {
  return typeof window.matchMedia === 'function' && window.matchMedia('(prefers-color-scheme: dark)').matches
}

const mode = ref<ThemeMode>(read())

function apply() {
  const dark = mode.value === 'dark' || (mode.value === 'system' && systemPrefersDark())
  document.documentElement.classList.toggle('dark', dark)
}

if (typeof window.matchMedia === 'function') {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener?.('change', () => {
    if (mode.value === 'system') apply()
  })
}

/** The colour theme, shared by the whole app and remembered per browser. */
export function useTheme() {
  function set(next: ThemeMode) {
    mode.value = next
    try {
      localStorage.setItem(storageKey, next)
    } catch {
      /* storage unavailable: the choice lasts for this visit */
    }
    apply()
  }

  function cycle() {
    set(order[(order.indexOf(mode.value) + 1) % order.length])
  }

  return { mode, set, cycle }
}
