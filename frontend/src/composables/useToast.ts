import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
}

const toasts = ref<Toast[]>([])
let nextId = 0

function addToast(message: string, type: ToastType = 'info') {
  const id = nextId++
  const toast: Toast = { id, message, type }

  if (toasts.value.length >= 3) {
    toasts.value.shift()
  }

  toasts.value.push(toast)

  setTimeout(() => {
    removeToast(id)
  }, 4000)
}

function removeToast(id: number) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

export function useToast() {
  return {
    toasts,
    addToast,
    removeToast
  }
}
