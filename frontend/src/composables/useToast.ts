import { readonly, ref } from 'vue'

export interface Toast {
  id: number
  message: string
  tone: 'success' | 'error' | 'info'
}

const toasts = ref<Toast[]>([])
let nextId = 1

function dismiss(id: number) {
  toasts.value = toasts.value.filter((toast) => toast.id !== id)
}

function show(message: string, tone: Toast['tone'] = 'success', durationMs = 4000) {
  const id = nextId++
  toasts.value = [...toasts.value.slice(-2), { id, message, tone }]
  setTimeout(() => dismiss(id), durationMs)
}

/** Short-lived feedback after an action, shown by ToastHost. */
export function useToast() {
  return { toasts: readonly(toasts), show, dismiss }
}
