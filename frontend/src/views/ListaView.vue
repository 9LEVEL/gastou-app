<template>
  <div class="lista-view">
    <!-- Sticky Header -->
    <div class="sticky-header">
      <div v-if="listaAtiva" class="lista-header-content">
        <div class="flex items-center justify-between" style="margin-bottom: 8px;">
          <div class="flex items-center" style="gap: 6px; flex: 1; min-width: 0;">
            <select
              :value="listaAtiva.id"
              @change="trocarLista(Number(($event.target as HTMLSelectElement).value))"
              style="font-family: var(--font-display); font-size: 1.1rem; font-weight: 700; border: none; background: transparent; color: var(--color-text); padding: 0; height: auto; min-height: auto; cursor: pointer; max-width: 100%;"
            >
              <option
                v-for="l in listas"
                :key="l.id"
                :value="l.id"
              >{{ l.nome }} ({{ meses[l.mes - 1] }}/{{ l.ano }})</option>
            </select>
          </div>
        </div>
        <div class="flex" style="gap: 8px; flex-wrap: wrap;">
          <span class="chip chip-primary">
            R$ {{ formatCurrency(totalMensalAjustado) }} planejado
          </span>
          <span
            class="chip"
            :class="percentualRenda > 30 ? 'chip-danger' : 'chip-primary'"
          >
            {{ percentualRenda.toFixed(0) }}% da renda
          </span>
        </div>
        <div class="flex" style="gap: 6px; margin-top: 8px;">
          <button
            class="filter-chip"
            :class="{ active: filtroLista === 'todos' }"
            @click="filtroLista = 'todos'"
          >Todos ({{ itens.length }})</button>
          <button
            class="filter-chip"
            :class="{ active: filtroLista === 'faltam' }"
            @click="filtroLista = 'faltam'"
          >Faltam ({{ itensFaltantes }})</button>
          <button
            class="filter-chip"
            :class="{ active: filtroLista === 'comprados' }"
            @click="filtroLista = 'comprados'"
          >Comprados ({{ itensComprados.length }})</button>
        </div>
      </div>
      <div v-else class="lista-header-content">
        <h1 class="font-display" style="font-size: 1.1rem; font-weight: 700;">Lista de Compras</h1>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-spinner" style="margin: 48px auto;"></div>

    <!-- Empty State -->
    <div v-else-if="!listaAtiva" class="empty-state">
      <p>Nenhuma lista encontrada.</p>
      <p class="text-sm text-secondary">Crie uma lista para começar.</p>
    </div>

    <div v-else-if="grouped.length === 0 && !loading" class="empty-state">
      <p>Nenhum item na lista.</p>
      <p class="text-sm text-secondary">Toque no + para adicionar itens.</p>
    </div>

    <!-- Category Groups -->
    <div v-else class="lista-grupos" :style="{ paddingBottom: itensComprados.length > 0 ? '140px' : '80px' }">
      <CategoriaGroup
        v-for="grupo in grouped"
        :key="grupo.nome"
        :nome="grupo.nome"
        :cor="grupo.cor"
        :count="grupo.items.length"
        v-model="expandedMap[grupo.nome]"
      >
        <ListaItem
          v-for="item in grupo.items"
          :key="item.id"
          :item="item"
          @check="handleCheck(item)"
          @edit="handleEdit(item)"
          @delete="handleDelete(item)"
        />
      </CategoriaGroup>
    </div>

    <!-- FAB -->
    <button class="fab" :style="{ bottom: itensComprados.length > 0 ? '150px' : '90px' }" @click="abrirModalAdicionar" title="Adicionar item">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="5" x2="12" y2="19" />
        <line x1="5" y1="12" x2="19" y2="12" />
      </svg>
    </button>

    <!-- Cart Summary Bar -->
    <Transition name="slide-up">
      <div v-if="itensComprados.length > 0" class="cart-bar">
        <div class="cart-bar-info">
          <span class="cart-bar-count">{{ itensComprados.length }} {{ itensComprados.length === 1 ? 'item' : 'itens' }}</span>
          <span class="cart-bar-total num">R$ {{ formatCurrency(totalCarrinho) }}</span>
        </div>
        <button
          class="btn btn-primary btn-sm"
          style="white-space: nowrap;"
          :disabled="finalizando"
          @click="finalizarCompra"
        >
          {{ finalizando ? 'Enviando...' : 'Finalizar Compra' }}
        </button>
      </div>
    </Transition>

    <!-- Modal Adicionar/Editar Item -->
    <Modal
      :show="modalAberto"
      :title="itemEditando ? 'Editar Item' : 'Adicionar Item'"
      @close="fecharModal"
    >
      <form class="modal-form" @submit.prevent="salvarItem" style="padding: 16px;">
        <!-- Produto (busca) -->
        <div class="form-group" style="position: relative;">
          <label class="text-sm" style="font-weight: 500; display: block; margin-bottom: 4px;">
            Produto <span style="color: var(--color-danger);">*</span>
          </label>
          <div v-if="itemEditando" style="padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; background: var(--color-surface); color: var(--color-text-secondary); font-size: 0.95rem;">
            {{ itemEditando.produto_nome }}
          </div>
          <template v-else>
            <input
              ref="buscaProdutoRef"
              v-model="buscaProduto"
              type="text"
              placeholder="Digite para buscar..."
              autocomplete="off"
              style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 0.95rem; background: var(--color-surface);"
              @focus="showSugestoes = true"
              @input="showSugestoes = true"
            />
            <div v-if="form.produto_id && !showSugestoes" class="text-xs" style="margin-top: 4px; color: var(--color-primary); font-weight: 500;">
              {{ produtoSelecionadoNome }}
            </div>
            <div
              v-if="showSugestoes && buscaProduto.length >= 1"
              class="search-dropdown"
            >
              <div
                v-for="prod in produtosFiltrados"
                :key="prod.id"
                class="search-dropdown-item"
                @mousedown.prevent="selecionarProduto(prod)"
              >
                <div class="truncate" style="font-size: 0.9rem;">{{ prod.nome }}</div>
                <div class="text-xs text-secondary">{{ prod.unidade }} · R$ {{ formatCurrency(prod.preco_ref) }}</div>
              </div>
              <div v-if="produtosFiltrados.length === 0" class="search-dropdown-item text-secondary text-sm" style="cursor: default;">
                Nenhum produto encontrado
              </div>
            </div>
          </template>
        </div>

        <!-- Quantidade e Duração -->
        <div class="form-row" style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
          <div class="form-group">
            <label class="text-sm" style="font-weight: 500; display: block; margin-bottom: 4px;">
              Quantidade <span style="color: var(--color-danger);">*</span>
            </label>
            <input
              v-model.number="form.qtd"
              type="number"
              step="0.1"
              min="0.1"
              required
              style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 0.95rem; background: var(--color-surface);"
            />
          </div>
          <div class="form-group">
            <label class="text-sm" style="font-weight: 500; display: block; margin-bottom: 4px;">
              Duração (meses)
            </label>
            <input
              v-model.number="form.duracao_meses"
              type="number"
              step="1"
              min="1"
              style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 0.95rem; background: var(--color-surface);"
            />
          </div>
        </div>

        <!-- Preço Estimado -->
        <div class="form-group">
          <label class="text-sm" style="font-weight: 500; display: block; margin-bottom: 4px;">
            Preço estimado (R$) <span style="color: var(--color-danger);">*</span>
          </label>
          <input
            v-model.number="form.preco_estimado"
            type="number"
            step="0.01"
            min="0"
            required
            style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 0.95rem; background: var(--color-surface);"
          />
        </div>

        <!-- Observação -->
        <div class="form-group">
          <label class="text-sm" style="font-weight: 500; display: block; margin-bottom: 4px;">
            Observação
          </label>
          <textarea
            v-model="form.observacao"
            rows="2"
            placeholder="Opcional"
            style="width: 100%; padding: 10px 12px; border: 1px solid var(--color-border); border-radius: 8px; font-size: 0.95rem; background: var(--color-surface); resize: none;"
          ></textarea>
        </div>

        <!-- Ações -->
        <div class="flex" style="gap: 8px; margin-top: 8px;">
          <button
            v-if="itemEditando"
            type="button"
            class="btn btn-outline"
            style="color: var(--color-danger); border-color: var(--color-danger);"
            @click="confirmarDelete"
          >
            Remover
          </button>
          <button type="button" class="btn btn-ghost" style="flex: 1;" @click="fecharModal">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary" style="flex: 2;" :disabled="salvando">
            {{ salvando ? 'Salvando...' : itemEditando ? 'Salvar' : 'Adicionar' }}
          </button>
        </div>
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted } from 'vue'
import { useLista } from '../composables/useLista'
import { useToast } from '../composables/useToast'
import { api } from '../api/client'
import CategoriaGroup from '../components/CategoriaGroup.vue'
import ListaItem from '../components/ListaItem.vue'
import Modal from '../components/Modal.vue'
import type { ListaItem as ListaItemType, Produto } from '../types'

const {
  listas,
  listaAtiva,
  itens,
  loading,
  totalMensalAjustado,
  itensFaltantes,
  percentualRenda,
  fetchListas,
  fetchItens,
  toggleCheck,
  addItem,
  updateItem,
  deleteItem,
} = useLista()

function trocarLista(id: number) {
  const lista = listas.value.find(l => l.id === id)
  if (lista) {
    listaAtiva.value = lista
  }
}

const { addToast } = useToast()

// Produtos para o modal
const produtos = ref<Produto[]>([])

// Modal state
const modalAberto = ref(false)
const itemEditando = ref<ListaItemType | null>(null)
const salvando = ref(false)

// Form state
const form = reactive({
  produto_id: 0 as number | '',
  qtd: 1,
  preco_estimado: 0,
  duracao_meses: 1,
  observacao: '',
})

// Expanded state por categoria (usando nome como chave)
const EXPANDED_KEY = 'gastou-expanded-categories'

function loadExpandedMap(): Record<string, boolean> {
  try {
    const stored = localStorage.getItem(EXPANDED_KEY)
    if (stored) return JSON.parse(stored)
  } catch {
    // corrupted data, ignore
  }
  return {}
}

const expandedMap = reactive<Record<string, boolean>>(loadExpandedMap())

// Meses em português
const meses = [
  'Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun',
  'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez',
]

// Filtro de itens
const filtroLista = ref<'todos' | 'faltam' | 'comprados'>('todos')

// Busca de produto
const buscaProduto = ref('')
const showSugestoes = ref(false)
const buscaProdutoRef = ref<HTMLInputElement | null>(null)

function semAcento(s: string): string {
  return s.normalize('NFD').replace(/[\u0300-\u036f]/g, '')
}

const produtosFiltrados = computed(() => {
  const termo = semAcento(buscaProduto.value.toLowerCase().trim())
  if (!termo) return produtos.value.slice(0, 10)
  return produtos.value.filter(p =>
    semAcento(p.nome.toLowerCase()).includes(termo)
  ).slice(0, 8)
})

const produtoSelecionadoNome = computed(() => {
  if (!form.produto_id) return ''
  const p = produtos.value.find(p => p.id === form.produto_id)
  return p?.nome ?? ''
})

function selecionarProduto(prod: Produto) {
  form.produto_id = prod.id
  form.preco_estimado = prod.preco_ref
  buscaProduto.value = prod.nome
  showSugestoes.value = false
}

// Agrupamento por categoria
const grouped = computed(() => {
  const filtrados = filtroLista.value === 'todos'
    ? itens.value
    : filtroLista.value === 'faltam'
      ? itens.value.filter(i => !i.comprado)
      : itens.value.filter(i => i.comprado)

  const map = new Map<number, { nome: string; cor: string; ordem: number; items: ListaItemType[] }>()
  for (const item of filtrados) {
    const key = item.categoria_id || 0
    if (!map.has(key)) {
      map.set(key, {
        nome: item.categoria_nome || 'Sem categoria',
        cor: item.categoria_cor || '#636E72',
        ordem: item.categoria_ordem || 99,
        items: [],
      })
    }
    map.get(key)!.items.push(item)
  }
  return [...map.values()].sort((a, b) => a.ordem - b.ordem)
})

// Inicializar expandedMap quando grouped muda
watch(grouped, (grupos) => {
  for (const grupo of grupos) {
    if (!(grupo.nome in expandedMap)) {
      expandedMap[grupo.nome] = true
    }
  }
})

// Persist expanded state
watch(expandedMap, (val) => {
  try {
    localStorage.setItem(EXPANDED_KEY, JSON.stringify(val))
  } catch {
    // localStorage full, ignore
  }
}, { deep: true })

function formatCurrency(v: number): string {
  return v.toFixed(2).replace('.', ',')
}

// Carrinho — itens marcados como comprado
const itensComprados = computed(() => itens.value.filter(i => i.comprado))

const totalCarrinho = computed(() => {
  return itensComprados.value.reduce((sum, item) => {
    return sum + item.qtd * item.preco_estimado
  }, 0)
})

// Handlers de item
function handleCheck(item: ListaItemType) {
  if (!listaAtiva.value) return
  toggleCheck(listaAtiva.value.id, item.id)
}

function handleEdit(item: ListaItemType) {
  itemEditando.value = item
  form.produto_id = item.produto_id
  form.qtd = item.qtd
  form.preco_estimado = item.preco_estimado
  form.duracao_meses = item.duracao_meses
  form.observacao = item.observacao || ''
  modalAberto.value = true
}

async function handleDelete(item: ListaItemType) {
  if (!listaAtiva.value) return
  await deleteItem(listaAtiva.value.id, item.id)
}

// Modal
function abrirModalAdicionar() {
  itemEditando.value = null
  resetForm()
  modalAberto.value = true
}

function fecharModal() {
  modalAberto.value = false
  itemEditando.value = null
  resetForm()
}

function resetForm() {
  form.produto_id = ''
  form.qtd = 1
  form.preco_estimado = 0
  form.duracao_meses = 1
  form.observacao = ''
  buscaProduto.value = ''
  showSugestoes.value = false
}

async function salvarItem() {
  if (!listaAtiva.value) return
  if (!form.produto_id) {
    addToast('Selecione um produto', 'error')
    return
  }

  salvando.value = true
  try {
    let result
    if (itemEditando.value) {
      result = await updateItem(listaAtiva.value.id, itemEditando.value.id, {
        qtd: form.qtd,
        preco_estimado: form.preco_estimado,
        duracao_meses: form.duracao_meses,
        observacao: form.observacao || null,
      })
    } else {
      result = await addItem(listaAtiva.value.id, {
        produto_id: form.produto_id as number,
        qtd: form.qtd,
        preco_estimado: form.preco_estimado,
        duracao_meses: form.duracao_meses,
        observacao: form.observacao || null,
      })
    }
    if (result) fecharModal()
  } finally {
    salvando.value = false
  }
}

async function confirmarDelete() {
  if (!listaAtiva.value || !itemEditando.value) return
  await deleteItem(listaAtiva.value.id, itemEditando.value.id)
  fecharModal()
}

// Finalizar compra
const finalizando = ref(false)

async function finalizarCompra() {
  if (!listaAtiva.value || itensComprados.value.length === 0) return

  const confirmMsg = `Finalizar compra com ${itensComprados.value.length} ${itensComprados.value.length === 1 ? 'item' : 'itens'} (R$ ${formatCurrency(totalCarrinho.value)})?`
  if (!confirm(confirmMsg)) return

  finalizando.value = true
  try {
    // 1. Create the compra
    const hoje = new Date().toISOString().split('T')[0]
    const compraRes = await api.post<{ id: number }>('/compras', {
      lista_id: listaAtiva.value.id,
      local: '',
      data: hoje,
    })

    if (!compraRes?.id) {
      addToast('Erro ao criar compra', 'error')
      return
    }

    // 2. Add each checked item to the compra
    let erros = 0
    for (const item of itensComprados.value) {
      try {
        await api.post(`/compras/${compraRes.id}/itens`, {
          produto_id: item.produto_id,
          nome_nfe: item.produto_nome || 'Item',
          qtd: item.qtd,
          unidade: item.produto_unidade || 'un',
          preco_unit: item.preco_estimado,
          preco_total: item.qtd * item.preco_estimado,
          lista_item_id: item.id,
        })
      } catch {
        erros++
      }
    }

    if (erros > 0) {
      addToast(`Compra criada com ${erros} ${erros === 1 ? 'erro' : 'erros'}`, 'warning')
    } else {
      addToast(`Compra registrada! ${itensComprados.value.length} itens`, 'success')
    }

    // 3. Refresh itens to get updated state
    await fetchItens(listaAtiva.value.id)

  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    finalizando.value = false
  }
}

// Watch listaAtiva para buscar itens
watch(listaAtiva, (lista) => {
  if (lista) {
    fetchItens(lista.id)
  }
})

// Mount
onMounted(async () => {
  await fetchListas()
  if (listaAtiva.value) {
    await fetchItens(listaAtiva.value.id)
  }
  try {
    produtos.value = await api.get<Produto[]>('/produtos?ativo=true')
  } catch {
    // silencioso — lista funciona sem produtos
  }
})
</script>
