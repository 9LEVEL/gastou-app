<template>
  <div>
    <!-- Header -->
    <div class="sticky-header">
      <div class="flex items-center justify-between">
        <div class="flex items-center" style="gap: 8px;">
          <h1 style="font-size: 1.3rem; font-weight: 700;">Compras</h1>
          <span class="badge badge-primary">{{ compras.length }}</span>
        </div>
        <div style="text-align: right;">
          <div class="text-xs text-secondary">gasto no mês</div>
          <div class="num" style="font-size: 1.1rem; font-weight: 700; color: var(--color-primary);">
            R$ {{ formatCurrency(totalMes) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-spinner"></div>

    <!-- Empty State -->
    <div v-else-if="compras.length === 0" class="empty-state">
      <div style="font-size: 2.5rem;">🛒</div>
      <p>Nenhuma compra registrada</p>
      <p class="text-xs" style="margin-top: 4px;">Toque em + para registrar uma compra</p>
    </div>

    <!-- Lista de Compras -->
    <div v-else style="padding-top: 8px; padding-bottom: 8px;">
      <CompraCard
        v-for="compra in compras"
        :key="compra.id"
        :compra="compra"
        @select="abrirDetalhes(compra.id)"
      />
    </div>

    <!-- FAB Nova Compra -->
    <button class="fab" @click="showNovaCompra = true" aria-label="Nova compra">
      +
    </button>

    <!-- Modal: Nova Compra -->
    <Modal :show="showNovaCompra" title="Nova Compra" @close="fecharNovaCompra">
      <form @submit.prevent="salvarNovaCompra">
        <div class="form-group">
          <label>Local *</label>
          <input
            v-model="formCompra.local"
            type="text"
            placeholder="Ex: Supermercado BIG"
            required
          />
        </div>

        <div class="form-group">
          <label>Data *</label>
          <input
            v-model="formCompra.data"
            type="date"
            required
          />
        </div>

        <div class="form-group">
          <label>Lista vinculada</label>
          <select v-model="formCompra.lista_id">
            <option :value="null">Nenhuma</option>
            <option v-for="lista in listas" :key="lista.id" :value="lista.id">
              {{ lista.nome }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>Total NF-e (R$)</label>
          <input
            v-model.number="formCompra.total_nfe"
            type="number"
            step="0.01"
            min="0"
            placeholder="0,00"
          />
        </div>

        <div class="form-group">
          <label>Observação</label>
          <textarea v-model="formCompra.observacao" placeholder="Observações opcionais..."></textarea>
        </div>

        <button type="submit" class="btn btn-primary" style="width: 100%;">
          Registrar Compra
        </button>
      </form>
    </Modal>

    <!-- Modal: Detalhes da Compra -->
    <Modal
      :show="showDetalhes"
      :title="showAddItem ? 'Adicionar Item' : (compraAtiva?.local || 'Detalhes da Compra')"
      @close="fecharDetalhes"
    >
      <!-- Formulário de Adicionar Item -->
      <div v-if="showAddItem">
        <form @submit.prevent="salvarItem">
          <div class="form-group">
            <label>Produto (opcional)</label>
            <select v-model="formItem.produto_id" @change="preencherProduto">
              <option :value="null">Selecionar produto...</option>
              <option v-for="produto in produtos" :key="produto.id" :value="produto.id">
                {{ produto.nome }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Nome na NF-e *</label>
            <input
              v-model="formItem.nome_nfe"
              type="text"
              placeholder="Nome do produto"
              required
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Qtd *</label>
              <input
                v-model.number="formItem.qtd"
                type="number"
                step="0.001"
                min="0.001"
                required
              />
            </div>
            <div class="form-group">
              <label>Unidade *</label>
              <input
                v-model="formItem.unidade"
                type="text"
                placeholder="un"
                required
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Preço unit. (R$) *</label>
              <input
                v-model.number="formItem.preco_unit"
                type="number"
                step="0.01"
                min="0"
                required
              />
            </div>
            <div class="form-group">
              <label>Total (R$)</label>
              <input
                v-model.number="formItem.preco_total"
                type="number"
                step="0.01"
                min="0"
              />
            </div>
          </div>

          <div class="flex" style="gap: 8px;">
            <button
              type="button"
              class="btn btn-outline"
              style="flex: 1;"
              @click="showAddItem = false"
            >
              Cancelar
            </button>
            <button type="submit" class="btn btn-primary" style="flex: 1;">
              Adicionar
            </button>
          </div>
        </form>
      </div>

      <!-- Detalhe da Compra -->
      <div v-else>
        <!-- Info da Compra -->
        <div class="card" style="margin-bottom: 16px;">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-xs text-secondary">Data</div>
              <div style="font-weight: 600;">{{ formatDate(compraAtiva?.data || '') }}</div>
            </div>
            <div style="text-align: right;">
              <div class="text-xs text-secondary">Total</div>
              <div class="num" style="font-size: 1.2rem; font-weight: 700; color: var(--color-primary);">
                R$ {{ formatCurrency(compraAtiva?.total_calculado || 0) }}
              </div>
            </div>
          </div>
          <div v-if="compraAtiva?.total_nfe && Math.abs(compraAtiva.total_nfe - (compraAtiva?.total_calculado || 0)) > 0.01"
            class="text-xs text-secondary"
            style="margin-top: 6px;">
            NF-e: R$ {{ formatCurrency(compraAtiva.total_nfe) }}
          </div>
          <div v-if="compraAtiva?.observacao" class="text-sm text-secondary" style="margin-top: 6px;">
            {{ compraAtiva.observacao }}
          </div>
          <div v-if="compraAtiva?.lista_id" class="mt-8">
            <span class="badge badge-primary">
              Lista: {{ nomeLista(compraAtiva.lista_id) }}
            </span>
          </div>
        </div>

        <!-- Lista de Itens -->
        <div style="margin-bottom: 12px;">
          <div class="flex items-center justify-between" style="margin-bottom: 8px;">
            <div style="font-weight: 600; font-size: 0.9rem;">
              Itens
              <span class="badge badge-primary" style="margin-left: 4px;">
                {{ compraAtiva?.itens?.length || 0 }}
              </span>
            </div>
            <button class="btn btn-sm btn-outline" @click="abrirAddItem">
              + Adicionar
            </button>
          </div>

          <!-- Loading itens -->
          <div v-if="loading" class="loading-spinner" style="margin: 20px auto;"></div>

          <!-- Empty itens -->
          <div v-else-if="!compraAtiva?.itens?.length" class="empty-state" style="padding: 20px;">
            <p>Nenhum item registrado</p>
          </div>

          <!-- Tabela de Itens -->
          <div v-else>
            <div
              v-for="item in compraAtiva.itens"
              :key="item.id"
              class="flex items-center justify-between"
              style="padding: 10px 0; border-bottom: 1px solid var(--color-border);"
            >
              <div style="flex: 1; min-width: 0;">
                <div class="truncate" style="font-weight: 500; font-size: 0.9rem;">
                  {{ item.nome_nfe }}
                </div>
                <div class="text-xs text-secondary">
                  {{ formatQtd(item.qtd) }} {{ item.unidade }} × R$ {{ formatCurrency(item.preco_unit) }}
                </div>
              </div>
              <div style="display: flex; align-items: center; gap: 8px; flex-shrink: 0; margin-left: 8px;">
                <div class="num" style="font-weight: 700; font-size: 0.95rem; white-space: nowrap;">
                  R$ {{ formatCurrency(item.preco_total) }}
                </div>
                <button
                  class="btn-ghost"
                  style="padding: 4px 6px; color: var(--color-danger); font-size: 0.85rem;"
                  @click="removerItem(item.id)"
                  aria-label="Remover item"
                >
                  ✕
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Botão Excluir Compra -->
        <button
          class="btn btn-danger"
          style="width: 100%; margin-top: 8px;"
          @click="confirmarExclusao"
        >
          Excluir Compra
        </button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useCompras } from '../composables/useCompras'
import { useLista } from '../composables/useLista'
import { api } from '../api/client'
import CompraCard from '../components/CompraCard.vue'
import Modal from '../components/Modal.vue'
import type { CompraInput, CompraItemInput, Produto } from '../types'

// Composables
const { compras, compraAtiva, loading, fetchCompras, fetchCompra, createCompra, deleteCompra, addCompraItem, deleteCompraItem } = useCompras()
const { listas, listaAtiva, fetchListas } = useLista()

// Modal state
const showNovaCompra = ref(false)
const showDetalhes = ref(false)
const showAddItem = ref(false)

// Produtos para o select de item
const produtos = ref<Produto[]>([])

// Formulário nova compra
const formCompraInicial = (): CompraInput & { total_nfe: number | null; observacao: string; lista_id: number | null } => ({
  local: '',
  data: new Date().toISOString().split('T')[0],
  lista_id: listaAtiva.value?.id ?? null,
  total_nfe: null,
  observacao: ''
})

const formCompra = ref(formCompraInicial())

// Formulário novo item
const formItemInicial = (): CompraItemInput & { produto_id: number | null } => ({
  produto_id: null,
  nome_nfe: '',
  qtd: 1,
  unidade: 'un',
  preco_unit: 0,
  preco_total: 0
})

const formItem = ref(formItemInicial())

// Total gasto no mês atual
const totalMes = computed(() => {
  const now = new Date()
  const mes = now.getMonth() + 1
  const ano = now.getFullYear()
  return compras.value
    .filter(c => {
      if (!c.data) return false
      const parts = c.data.split('-')
      return parseInt(parts[0]) === ano && parseInt(parts[1]) === mes
    })
    .reduce((sum, c) => sum + (c.total_calculado || 0), 0)
})

// Auto-calcular preco_total ao mudar qtd ou preco_unit
watch(
  () => [formItem.value.qtd, formItem.value.preco_unit],
  ([qtd, preco_unit]) => {
    if (qtd && preco_unit) {
      formItem.value.preco_total = parseFloat(((qtd as number) * (preco_unit as number)).toFixed(2))
    }
  }
)

// Funções de formatação
function formatCurrency(value: number): string {
  return (value || 0).toFixed(2).replace('.', ',')
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    return `${parts[2]}/${parts[1]}/${parts[0]}`
  }
  return dateStr
}

function formatQtd(qtd: number): string {
  return Number.isInteger(qtd) ? String(qtd) : qtd.toFixed(3).replace(/\.?0+$/, '')
}

function nomeLista(listaId: number): string {
  return listas.value.find(l => l.id === listaId)?.nome || `Lista #${listaId}`
}

// Ações nova compra
function fecharNovaCompra() {
  showNovaCompra.value = false
  formCompra.value = formCompraInicial()
}

async function salvarNovaCompra() {
  const input: CompraInput = {
    local: formCompra.value.local,
    data: formCompra.value.data,
    lista_id: formCompra.value.lista_id || undefined,
    total_nfe: formCompra.value.total_nfe || undefined,
    observacao: formCompra.value.observacao || undefined
  }
  const nova = await createCompra(input)
  if (nova) {
    fecharNovaCompra()
  }
}

// Ações detalhes
async function abrirDetalhes(id: number) {
  showDetalhes.value = true
  showAddItem.value = false
  await fetchCompra(id)
}

function fecharDetalhes() {
  showDetalhes.value = false
  showAddItem.value = false
  formItem.value = formItemInicial()
}

async function confirmarExclusao() {
  if (!compraAtiva.value) return
  if (!window.confirm(`Excluir compra em "${compraAtiva.value.local}"?`)) return
  await deleteCompra(compraAtiva.value.id)
  fecharDetalhes()
}

// Ações item
function abrirAddItem() {
  formItem.value = formItemInicial()
  showAddItem.value = true
}

function preencherProduto() {
  const produtoId = formItem.value.produto_id
  if (!produtoId) return
  const produto = produtos.value.find(p => p.id === produtoId)
  if (produto) {
    formItem.value.nome_nfe = produto.nome
    formItem.value.unidade = produto.unidade
    if (produto.preco_ref > 0) {
      formItem.value.preco_unit = produto.preco_ref
      formItem.value.preco_total = parseFloat((formItem.value.qtd * produto.preco_ref).toFixed(2))
    }
  }
}

async function salvarItem() {
  if (!compraAtiva.value) return
  const input: CompraItemInput = {
    produto_id: formItem.value.produto_id || undefined,
    nome_nfe: formItem.value.nome_nfe,
    qtd: formItem.value.qtd,
    unidade: formItem.value.unidade,
    preco_unit: formItem.value.preco_unit,
    preco_total: formItem.value.preco_total
  }
  const item = await addCompraItem(compraAtiva.value.id, input)
  if (item) {
    showAddItem.value = false
    formItem.value = formItemInicial()
  }
}

async function removerItem(itemId: number) {
  if (!compraAtiva.value) return
  if (!window.confirm('Remover este item?')) return
  await deleteCompraItem(compraAtiva.value.id, itemId)
}

// Inicialização
onMounted(async () => {
  await fetchListas()
  const listaId = listaAtiva.value?.id
  await fetchCompras(listaId)

  // Pré-popular lista no formulário
  if (listaAtiva.value) {
    formCompra.value.lista_id = listaAtiva.value.id
  }

  // Buscar produtos para o formulário de item
  try {
    produtos.value = await api.get<Produto[]>('/produtos?ativo=true')
  } catch {
    // silencioso — select fica vazio
  }
})
</script>
