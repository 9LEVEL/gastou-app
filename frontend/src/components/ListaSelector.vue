<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
        <div class="lista-selector-panel">
          <!-- Handle -->
          <div class="lista-selector-handle"></div>

          <div class="lista-selector-header">
            <h2>Suas Listas</h2>
            <button class="modal-close" @click="$emit('close')">&times;</button>
          </div>

          <!-- Lista rows -->
          <div class="lista-selector-body">
            <div
              v-for="lista in listas"
              :key="lista.id"
              class="lista-selector-row"
              :class="{ active: lista.id === ativaId }"
              @click="$emit('select', lista)"
            >
              <div style="flex: 1; min-width: 0;">
                <div class="flex items-center" style="gap: 6px; margin-bottom: 2px;">
                  <span class="truncate" style="font-weight: 500; font-size: 0.95rem;">{{ lista.nome }}</span>
                  <span v-if="lista.id === ativaId" class="badge badge-primary">ativa</span>
                </div>
                <p class="text-xs text-secondary">
                  {{ meses[lista.mes - 1] }}/{{ lista.ano }}
                  <template v-if="lista.itens_total"> &middot; {{ lista.itens_total }} itens</template>
                  &middot; R$ {{ formatCurrency(lista.renda) }}
                </p>
              </div>

              <!-- Menu contextual -->
              <div class="lista-selector-actions" @click.stop>
                <button
                  class="lista-selector-menu-btn"
                  @click.stop="toggleMenu(lista.id)"
                  title="Opções"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <circle cx="12" cy="5" r="2"/>
                    <circle cx="12" cy="12" r="2"/>
                    <circle cx="12" cy="19" r="2"/>
                  </svg>
                </button>

                <Transition name="fade">
                  <div v-if="menuAberto === lista.id" class="lista-selector-dropdown">
                    <button @click="iniciarEdicao(lista)">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                      Editar
                    </button>
                    <button
                      v-if="listas.length > 1"
                      class="danger"
                      @click="$emit('delete', lista); menuAberto = null"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                        <path d="M10 11v6M14 11v6"/>
                      </svg>
                      Excluir
                    </button>
                  </div>
                </Transition>
              </div>
            </div>
          </div>

          <!-- Nova lista -->
          <button class="lista-selector-new" @click="$emit('create')">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            Nova Lista
          </button>

          <!-- Inline edit -->
          <Transition name="fade">
            <div v-if="editando" class="lista-selector-edit">
              <h3 style="font-size: 0.95rem; font-weight: 600; margin-bottom: 12px;">Editar Lista</h3>
              <div class="form-group">
                <label class="text-sm" style="font-weight: 500;">Nome</label>
                <input v-model="formEdit.nome" type="text" required />
              </div>
              <div class="form-group">
                <label class="text-sm" style="font-weight: 500;">Renda (R$)</label>
                <input v-model.number="formEdit.renda" type="number" step="0.01" min="0" />
              </div>
              <div class="flex" style="gap: 8px; margin-top: 8px;">
                <button class="btn btn-ghost btn-sm" style="flex: 1;" @click="cancelarEdicao">Cancelar</button>
                <button class="btn btn-primary btn-sm" style="flex: 2;" :disabled="salvandoEdit" @click="salvarEdicao">
                  {{ salvandoEdit ? 'Salvando...' : 'Salvar' }}
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { Lista } from '../types'

const props = defineProps<{
  show: boolean
  listas: Lista[]
  ativaId: number | null
}>()

const emit = defineEmits<{
  close: []
  select: [lista: Lista]
  create: []
  delete: [lista: Lista]
  update: [id: number, data: { nome: string; renda: number }]
}>()

const meses = [
  'Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun',
  'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez',
]

function formatCurrency(v: number): string {
  return v.toFixed(2).replace('.', ',')
}

// Menu contextual
const menuAberto = ref<number | null>(null)

function toggleMenu(id: number) {
  menuAberto.value = menuAberto.value === id ? null : id
}

// Edição inline
const editando = ref<Lista | null>(null)
const salvandoEdit = ref(false)
const formEdit = reactive({ nome: '', renda: 0 })

function iniciarEdicao(lista: Lista) {
  menuAberto.value = null
  editando.value = lista
  formEdit.nome = lista.nome
  formEdit.renda = lista.renda
}

function cancelarEdicao() {
  editando.value = null
}

async function salvarEdicao() {
  if (!editando.value) return
  salvandoEdit.value = true
  emit('update', editando.value.id, { nome: formEdit.nome, renda: formEdit.renda })
  editando.value = null
  salvandoEdit.value = false
}

// Reset state when closing
watch(() => props.show, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
    menuAberto.value = null
    editando.value = null
  }
})
</script>

<style scoped>
.lista-selector-panel {
  background: var(--color-card);
  border-radius: 16px 16px 0 0;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  padding: 8px 0 0;
  padding-bottom: env(safe-area-inset-bottom, 0px);
  animation: slideUp 0.25s ease;
}

.lista-selector-handle {
  width: 36px;
  height: 4px;
  border-radius: 2px;
  background: var(--color-border);
  margin: 0 auto 12px;
}

.lista-selector-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px 12px;
}

.lista-selector-header h2 {
  font-size: 1.1rem;
  font-weight: 600;
}

.lista-selector-body {
  border-top: 1px solid var(--color-border);
}

.lista-selector-row {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--color-border);
  transition: background 0.1s ease;
  position: relative;
}

.lista-selector-row:active {
  background: var(--color-primary-light);
}

.lista-selector-row.active {
  background: var(--color-primary-light);
}

.lista-selector-actions {
  position: relative;
  flex-shrink: 0;
}

.lista-selector-menu-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--color-text-secondary);
  transition: background 0.15s ease;
}

.lista-selector-menu-btn:active {
  background: var(--color-border);
}

.lista-selector-dropdown {
  position: absolute;
  right: 0;
  top: 100%;
  background: var(--color-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  z-index: 10;
  min-width: 140px;
  overflow: hidden;
}

.lista-selector-dropdown button {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 14px;
  font-size: 0.9rem;
  text-align: left;
  color: var(--color-text);
  background: transparent;
  border: none;
  cursor: pointer;
}

.lista-selector-dropdown button:active {
  background: var(--color-primary-light);
}

.lista-selector-dropdown button.danger {
  color: var(--color-danger);
}

.lista-selector-new {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 16px;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--color-primary);
  cursor: pointer;
  border: none;
  background: transparent;
  transition: background 0.15s ease;
}

.lista-selector-new:active {
  background: var(--color-primary-light);
}

.lista-selector-edit {
  padding: 16px;
  border-top: 1px solid var(--color-border);
  background: var(--color-surface);
}
</style>
