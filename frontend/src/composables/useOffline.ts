import { ref, onMounted, onUnmounted } from 'vue'

const isOffline = ref(!navigator.onLine)

let listenersAttached = false

function handleOnline() {
  isOffline.value = false
}

function handleOffline() {
  isOffline.value = true
}

export function useOffline() {
  onMounted(() => {
    if (!listenersAttached) {
      window.addEventListener('online', handleOnline)
      window.addEventListener('offline', handleOffline)
      listenersAttached = true
    }
  })

  onUnmounted(() => {
    // Keep listeners alive as the composable is shared
  })

  return {
    isOffline
  }
}
