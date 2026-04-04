import { ref, watch } from 'vue'

type Theme = 'light' | 'dark'

const STORAGE_KEY = 'gastou-theme'

function getInitialTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'dark' || stored === 'light') return stored
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  return 'light'
}

const theme = ref<Theme>(getInitialTheme())

function applyTheme(t: Theme) {
  document.documentElement.setAttribute('data-theme', t)
}

// Apply immediately on module load
applyTheme(theme.value)

watch(theme, (newTheme) => {
  applyTheme(newTheme)
  localStorage.setItem(STORAGE_KEY, newTheme)
})

function toggleTheme() {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}

export function useTheme() {
  return {
    theme,
    toggleTheme
  }
}
