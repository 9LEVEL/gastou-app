<template>
  <div class="card compra-card" @click="$emit('select')">
    <div class="flex items-center justify-between">
      <div>
        <div style="font-weight: 600; font-size: 0.95rem;">{{ compra.local || 'Sem local' }}</div>
        <div class="text-sm text-secondary">{{ formatDate(compra.data) }}</div>
      </div>
      <div style="text-align: right;">
        <div class="num" style="font-size: 1.1rem; font-weight: 700;">
          R$ {{ formatCurrency(compra.total_calculado) }}
        </div>
        <div v-if="compra.itens" class="badge badge-primary" style="margin-top: 2px;">
          {{ compra.itens.length }} {{ compra.itens.length === 1 ? 'item' : 'itens' }}
        </div>
      </div>
    </div>
    <div v-if="compra.total_nfe && Math.abs(compra.total_nfe - compra.total_calculado) > 0.01" class="text-xs text-secondary" style="margin-top: 6px;">
      NF-e: R$ {{ formatCurrency(compra.total_nfe) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Compra } from '../types'

defineProps<{
  compra: Compra
}>()

defineEmits<{
  select: []
}>()

function formatCurrency(value: number): string {
  return value.toFixed(2).replace('.', ',')
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    return `${parts[2]}/${parts[1]}/${parts[0]}`
  }
  return dateStr
}
</script>

<style scoped>
.compra-card {
  cursor: pointer;
  transition: transform 0.1s ease;
  margin: 0 16px 8px;
}

.compra-card:active {
  transform: scale(0.98);
}
</style>
