<template>
  <div class="lista-item card" :class="{ 'item-checked': item.comprado }">
    <div class="lista-item-row">
      <div class="checkbox-touch" @click.stop="$emit('check')">
        <div class="check-box" :class="{ checked: item.comprado }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
      </div>

      <div class="lista-item-info" @click="$emit('edit')">
        <div class="lista-item-name truncate">{{ item.produto_nome }}</div>
        <div class="lista-item-formula text-xs text-secondary">
          {{ formatQtd(item.qtd) }} &times; R$ {{ formatCurrency(item.preco_estimado) }}
          <span v-if="item.duracao_meses > 1" class="badge badge-info" style="margin-left: 4px;">
            rende {{ item.duracao_meses }} meses
          </span>
        </div>
      </div>

      <div class="lista-item-total">
        <div class="num" style="font-size: 0.95rem; font-weight: 600;">
          R$ {{ formatCurrency(subtotal) }}
        </div>
        <div v-if="item.duracao_meses > 1" class="text-xs text-secondary num">
          R$ {{ formatCurrency(monthly) }}/m
        </div>
      </div>

      <button class="lista-item-delete" @click.stop="$emit('delete')" title="Remover">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="3 6 5 6 21 6" />
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
        </svg>
      </button>
    </div>
    <div v-if="item.observacao" class="lista-item-obs text-xs text-secondary" style="margin-top: 4px; padding-left: 44px;">
      {{ item.observacao }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ListaItem } from '../types'

const props = defineProps<{
  item: ListaItem
}>()

defineEmits<{
  check: []
  edit: []
  delete: []
}>()

const subtotal = computed(() => props.item.qtd * props.item.preco_estimado)
const monthly = computed(() => subtotal.value / (props.item.duracao_meses || 1))

function formatCurrency(value: number): string {
  return value.toFixed(2).replace('.', ',')
}

function formatQtd(value: number): string {
  if (Number.isInteger(value)) return value.toString()
  return value.toFixed(2).replace('.', ',')
}
</script>

<style scoped>
.lista-item {
  padding: 8px;
  margin: 0 16px 4px;
}

.lista-item-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.lista-item-info {
  flex: 1;
  min-width: 0;
  cursor: pointer;
}

.lista-item-name {
  font-weight: 500;
  font-size: 0.9rem;
}

.lista-item-total {
  text-align: right;
  flex-shrink: 0;
  padding-right: 4px;
}

.lista-item-delete {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--color-text-secondary);
  border-radius: 50%;
  transition: background 0.15s ease, color 0.15s ease;
}

.lista-item-delete:hover {
  background: var(--color-danger-light);
  color: var(--color-danger);
}
</style>
