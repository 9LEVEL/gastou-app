<template>
  <div class="app-container">
    <div v-if="isOffline" class="offline-banner">
      Sem conexão — mostrando dados em cache
    </div>
    <router-view />
    <Toast />
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { useOffline } from './composables/useOffline'
import { useTheme } from './composables/useTheme'
import Toast from './components/Toast.vue'
import BottomNav from './components/BottomNav.vue'

const { isOffline } = useOffline()
useTheme()
</script>
