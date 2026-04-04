<template>
  <div class="config-view" style="padding: 0 0 100px 0;">
    <!-- Header -->
    <div class="sticky-header">
      <h1 class="font-display" style="font-size: 1.2rem; font-weight: 700;">Configurações</h1>
    </div>

    <!-- ===== APARÊNCIA ===== -->
    <section style="padding: 0 16px 8px;">
      <h2 class="section-title">Aparência</h2>
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p style="font-weight: 500; font-size: 0.95rem;">Tema</p>
            <p class="text-xs text-secondary">
              {{ theme === 'dark' ? 'Modo escuro ativado' : 'Modo claro ativado' }}
            </p>
          </div>
          <button
            class="theme-toggle-btn"
            @click="toggleTheme"
            :title="theme === 'dark' ? 'Mudar para tema claro' : 'Mudar para tema escuro'"
          >
            <!-- Sun icon (shown in dark mode) -->
            <svg v-if="theme === 'dark'" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/>
              <line x1="21" y1="12" x2="23" y2="12"/>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <!-- Moon icon (shown in light mode) -->
            <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </button>
        </div>
      </div>
    </section>

    <!-- ===== RENDA ===== -->
    <section style="padding: 0 16px 8px;">
      <h2 class="section-title">Renda</h2>
      <div class="card">
        <div v-if="!editandoRenda" class="flex items-center justify-between">
          <div>
            <p class="text-sm text-secondary" style="margin-bottom: 2px;">Lista ativa</p>
            <p class="font-display" style="font-size: 1.4rem; font-weight: 700; color: var(--color-primary);">
              R$ {{ formatCurrency(listaAtiva?.renda ?? 0) }}
            </p>
          </div>
          <button class="btn btn-outline btn-sm" @click="iniciarEdicaoRenda">
            Editar
          </button>
        </div>

        <div v-else class="flex items-center" style="gap: 8px;">
          <input
            v-model.number="novaRenda"
            type="number"
            step="0.01"
            min="0"
            placeholder="0,00"
            style="flex: 1;"
            @keyup.enter="salvarRenda"
            @keyup.escape="cancelarEdicaoRenda"
          />
          <button class="btn btn-ghost btn-sm" @click="cancelarEdicaoRenda">Cancelar</button>
          <button class="btn btn-primary btn-sm" :disabled="salvandoRenda" @click="salvarRenda">
            {{ salvandoRenda ? '...' : 'Salvar' }}
          </button>
        </div>

        <p v-if="!listaAtiva" class="text-sm text-secondary" style="margin-top: 4px;">
          Nenhuma lista ativa encontrada.
        </p>
      </div>
    </section>

    <!-- ===== CATEGORIAS ===== -->
    <section style="padding: 8px 16px;">
      <div class="flex items-center justify-between" style="margin-bottom: 8px; cursor: pointer;" @click="configSections.categorias = !configSections.categorias">
        <div class="flex items-center" style="gap: 8px;">
          <svg
            class="config-chevron"
            :class="{ collapsed: !configSections.categorias }"
            width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"
          ><polyline points="6 9 12 15 18 9"/></svg>
          <h2 class="section-title" style="margin-bottom: 0;">Categorias</h2>
          <span class="text-xs text-secondary">({{ categorias.length }})</span>
        </div>
        <button class="btn btn-primary btn-sm" @click.stop="abrirModalNovaCategoria">
          + Nova
        </button>
      </div>

      <div v-show="configSections.categorias">

      <div v-if="loadingCategorias" class="loading-spinner" style="margin: 24px auto;"></div>

      <div v-else-if="categorias.length === 0" class="card text-center text-secondary">
        <p class="text-sm">Nenhuma categoria cadastrada.</p>
      </div>

      <div v-else class="card" style="padding: 0; overflow: hidden;">
        <div
          v-for="(cat, idx) in categoriasOrdenadas"
          :key="cat.id"
          class="flex items-center"
          style="padding: 12px 16px; gap: 12px;"
          :style="idx < categoriasOrdenadas.length - 1 ? 'border-bottom: 1px solid var(--color-border);' : ''"
        >
          <!-- Indicador de cor -->
          <span
            style="width: 14px; height: 14px; border-radius: 50%; flex-shrink: 0; display: block;"
            :style="{ background: cat.cor }"
          ></span>

          <!-- Nome e ordem -->
          <div style="flex: 1; min-width: 0;">
            <p style="font-weight: 500; font-size: 0.95rem;" class="truncate">{{ cat.nome }}</p>
            <p class="text-xs text-secondary">Ordem: {{ cat.ordem }}</p>
          </div>

          <!-- Ações -->
          <div class="flex" style="gap: 4px;">
            <button
              class="btn btn-ghost btn-sm"
              style="padding: 6px 10px;"
              @click="abrirModalEditarCategoria(cat)"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
            </button>
            <button
              class="btn btn-ghost btn-sm"
              style="padding: 6px 10px; color: var(--color-danger);"
              @click="confirmarDeleteCategoria(cat)"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                <path d="M10 11v6M14 11v6"/>
                <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      </div>
    </section>

    <!-- ===== PRODUTOS ===== -->
    <section style="padding: 8px 16px;">
      <div class="flex items-center justify-between" style="margin-bottom: 8px; cursor: pointer;" @click="configSections.produtos = !configSections.produtos">
        <div class="flex items-center" style="gap: 8px;">
          <svg
            class="config-chevron"
            :class="{ collapsed: !configSections.produtos }"
            width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"
          ><polyline points="6 9 12 15 18 9"/></svg>
          <h2 class="section-title" style="margin-bottom: 0;">Produtos</h2>
          <span class="text-xs text-secondary">({{ produtos.length }})</span>
        </div>
        <button class="btn btn-primary btn-sm" @click.stop="abrirModalNovoProduto">
          + Novo
        </button>
      </div>

      <div v-show="configSections.produtos">

      <!-- Filtro por categoria -->
      <div style="margin-bottom: 8px;">
        <select v-model="filtroCategoriaId" style="width: 100%; font-size: 0.85rem;">
          <option :value="null">Todas as categorias</option>
          <option
            v-for="cat in categoriasOrdenadas"
            :key="cat.id"
            :value="cat.id"
          >{{ cat.nome }}</option>
        </select>
      </div>

      <div v-if="loadingProdutos" class="loading-spinner" style="margin: 24px auto;"></div>

      <div v-else-if="produtosFiltrados.length === 0" class="card text-center text-secondary">
        <p class="text-sm">Nenhum produto cadastrado.</p>
      </div>

      <div v-else class="card" style="padding: 0; overflow: hidden;">
        <div
          v-for="(prod, idx) in produtosFiltrados"
          :key="prod.id"
          class="flex items-center"
          style="padding: 12px 16px; gap: 12px;"
          :style="idx < produtosFiltrados.length - 1 ? 'border-bottom: 1px solid var(--color-border);' : ''"
        >
          <!-- Indicador de cor da categoria -->
          <span
            style="width: 14px; height: 14px; border-radius: 50%; flex-shrink: 0; display: block;"
            :style="{ background: corDaCategoria(prod.categoria_id) }"
          ></span>

          <!-- Nome e detalhes -->
          <div style="flex: 1; min-width: 0;">
            <p style="font-weight: 500; font-size: 0.95rem;" class="truncate">{{ prod.nome }}</p>
            <p class="text-xs text-secondary">
              {{ prod.unidade }} &middot; R$ {{ formatCurrency(prod.preco_ref) }}
            </p>
          </div>

          <!-- Status -->
          <span
            v-if="prod.ativo"
            class="badge badge-primary"
            style="flex-shrink: 0;"
          >ativo</span>
          <span
            v-else
            class="text-xs text-secondary"
            style="flex-shrink: 0;"
          >inativo</span>

          <!-- Ações -->
          <div class="flex" style="gap: 4px; flex-shrink: 0;">
            <button
              class="btn btn-ghost btn-sm"
              style="padding: 6px 10px;"
              @click="abrirModalEditarProduto(prod)"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
            </button>
            <button
              class="btn btn-ghost btn-sm"
              style="padding: 6px 10px; color: var(--color-danger);"
              @click="confirmarDeleteProduto(prod)"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                <path d="M10 11v6M14 11v6"/>
                <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      </div>
    </section>

    <!-- ===== SOBRE ===== -->
    <section style="padding: 8px 16px;">
      <h2 class="section-title">Sobre</h2>
      <div class="card text-center" style="padding: 24px 16px;">
        <p class="font-display" style="font-size: 1.3rem; font-weight: 700; color: var(--color-primary);">Gastou.app</p>
        <p class="text-sm text-secondary" style="margin-top: 4px;">v1.0</p>
        <p class="text-xs text-secondary" style="margin-top: 4px;">by 9LEVEL</p>
        <p class="text-sm text-secondary" style="margin-top: 8px;">Controle de compras familiar</p>
      </div>
    </section>

    <!-- ===== MODAL CATEGORIA ===== -->
    <Modal
      :show="modalCategoriaAberto"
      :title="categoriaEditando ? 'Editar Categoria' : 'Nova Categoria'"
      @close="fecharModalCategoria"
    >
      <form style="padding: 0 16px 16px;" @submit.prevent="salvarCategoria">
        <div class="form-group">
          <label>Nome <span style="color: var(--color-danger);">*</span></label>
          <input
            v-model="formCategoria.nome"
            type="text"
            placeholder="Ex: Hortifruti"
            required
          />
        </div>

        <div class="form-group">
          <label>Cor <span style="color: var(--color-danger);">*</span></label>
          <div class="flex items-center" style="gap: 10px;">
            <span
              style="width: 20px; height: 20px; border-radius: 50%; flex-shrink: 0; display: block; border: 2px solid var(--color-border);"
              :style="{ background: formCategoria.cor }"
            ></span>
            <select v-model="formCategoria.cor" style="flex: 1;">
              <option
                v-for="cor in coresPredefinidas"
                :key="cor.hex"
                :value="cor.hex"
              >{{ cor.nome }} ({{ cor.hex }})</option>
            </select>
          </div>
        </div>

        <div class="form-group">
          <label>Ordem</label>
          <input
            v-model.number="formCategoria.ordem"
            type="number"
            min="1"
            placeholder="1"
          />
        </div>

        <div class="flex" style="gap: 8px; margin-top: 8px;">
          <button type="button" class="btn btn-ghost" style="flex: 1;" @click="fecharModalCategoria">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary" style="flex: 2;" :disabled="salvandoCategoria">
            {{ salvandoCategoria ? 'Salvando...' : categoriaEditando ? 'Salvar' : 'Criar' }}
          </button>
        </div>
      </form>
    </Modal>

    <!-- ===== MODAL PRODUTO ===== -->
    <Modal
      :show="modalProdutoAberto"
      :title="produtoEditando ? 'Editar Produto' : 'Novo Produto'"
      @close="fecharModalProduto"
    >
      <form style="padding: 0 16px 16px;" @submit.prevent="salvarProduto">
        <div class="form-group">
          <label>Nome <span style="color: var(--color-danger);">*</span></label>
          <input
            v-model="formProduto.nome"
            type="text"
            placeholder="Ex: Arroz"
            required
          />
        </div>

        <div class="form-group">
          <label>Categoria <span style="color: var(--color-danger);">*</span></label>
          <select v-model.number="formProduto.categoria_id" required>
            <option
              v-for="cat in categoriasOrdenadas"
              :key="cat.id"
              :value="cat.id"
            >{{ cat.nome }}</option>
          </select>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Unidade</label>
            <input
              v-model="formProduto.unidade"
              type="text"
              placeholder="un"
            />
          </div>
          <div class="form-group">
            <label>Unidade preço</label>
            <input
              v-model="formProduto.unidade_preco"
              type="text"
              placeholder="un"
            />
          </div>
        </div>

        <div class="form-group">
          <label>Preço referência (R$)</label>
          <input
            v-model.number="formProduto.preco_ref"
            type="number"
            step="0.01"
            min="0"
            placeholder="0,00"
          />
        </div>

        <div class="form-group">
          <label class="flex items-center" style="gap: 8px; cursor: pointer;">
            <input
              v-model="formProduto.ativo"
              type="checkbox"
              style="width: auto; margin: 0;"
            />
            <span>Ativo</span>
          </label>
        </div>

        <div class="flex" style="gap: 8px; margin-top: 8px;">
          <button type="button" class="btn btn-ghost" style="flex: 1;" @click="fecharModalProduto">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary" style="flex: 2;" :disabled="salvandoProduto">
            {{ salvandoProduto ? 'Salvando...' : produtoEditando ? 'Salvar' : 'Criar' }}
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
import { useTheme } from '../composables/useTheme'
import { api } from '../api/client'
import Modal from '../components/Modal.vue'
import type { Categoria, CategoriaInput, Produto, ProdutoInput } from '../types'

// ===== COMPOSABLES =====
const { listaAtiva, fetchListas, updateLista } = useLista()
const { addToast } = useToast()
const { theme, toggleTheme } = useTheme()

// ===== SEÇÕES COLAPSÁVEIS =====
const CONFIG_SECTIONS_KEY = 'gastou-config-sections'

function loadConfigSections(): Record<string, boolean> {
  try {
    const stored = localStorage.getItem(CONFIG_SECTIONS_KEY)
    if (stored) return JSON.parse(stored)
  } catch {
    // corrupted
  }
  return { categorias: false, produtos: false }
}

const configSections = reactive<Record<string, boolean>>(loadConfigSections())

watch(configSections, (val) => {
  try {
    localStorage.setItem(CONFIG_SECTIONS_KEY, JSON.stringify(val))
  } catch {
    // full
  }
}, { deep: true })

// ===== CATEGORIAS =====
const categorias = ref<Categoria[]>([])
const loadingCategorias = ref(false)

const categoriasOrdenadas = computed(() =>
  [...categorias.value].sort((a, b) => (a.ordem ?? 99) - (b.ordem ?? 99))
)

async function fetchCategorias() {
  try {
    loadingCategorias.value = true
    categorias.value = await api.get<Categoria[]>('/categorias')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    loadingCategorias.value = false
  }
}

// ===== RENDA =====
const editandoRenda = ref(false)
const novaRenda = ref(0)
const salvandoRenda = ref(false)

function iniciarEdicaoRenda() {
  novaRenda.value = listaAtiva.value?.renda ?? 0
  editandoRenda.value = true
}

function cancelarEdicaoRenda() {
  editandoRenda.value = false
}

async function salvarRenda() {
  if (!listaAtiva.value) return
  salvandoRenda.value = true
  try {
    await updateLista(listaAtiva.value.id, { renda: novaRenda.value })
    addToast('Renda atualizada', 'success')
    editandoRenda.value = false
  } finally {
    salvandoRenda.value = false
  }
}

// ===== MODAL CATEGORIA =====
const modalCategoriaAberto = ref(false)
const categoriaEditando = ref<Categoria | null>(null)
const salvandoCategoria = ref(false)

const formCategoria = reactive<CategoriaInput>({
  nome: '',
  cor: '#2D7A5F',
  ordem: 1,
})

const coresPredefinidas = [
  { nome: 'Vermelho', hex: '#C0392B' },
  { nome: 'Âmbar', hex: '#C4850A' },
  { nome: 'Azul', hex: '#2980B9' },
  { nome: 'Verde', hex: '#2D7A5F' },
  { nome: 'Roxo', hex: '#6C5CE7' },
  { nome: 'Coral', hex: '#E17055' },
  { nome: 'Cinza', hex: '#636E72' },
]

function abrirModalNovaCategoria() {
  categoriaEditando.value = null
  formCategoria.nome = ''
  formCategoria.cor = '#2D7A5F'
  formCategoria.ordem = (categorias.value.length + 1)
  modalCategoriaAberto.value = true
}

function abrirModalEditarCategoria(cat: Categoria) {
  categoriaEditando.value = cat
  formCategoria.nome = cat.nome
  formCategoria.cor = cat.cor
  formCategoria.ordem = cat.ordem ?? 1
  modalCategoriaAberto.value = true
}

function fecharModalCategoria() {
  modalCategoriaAberto.value = false
  categoriaEditando.value = null
}

async function salvarCategoria() {
  salvandoCategoria.value = true
  try {
    if (categoriaEditando.value) {
      const updated = await api.put<Categoria>(`/categorias/${categoriaEditando.value.id}`, {
        nome: formCategoria.nome,
        cor: formCategoria.cor,
        ordem: formCategoria.ordem,
      })
      const idx = categorias.value.findIndex(c => c.id === categoriaEditando.value!.id)
      if (idx !== -1) categorias.value[idx] = updated
      addToast('Categoria atualizada', 'success')
    } else {
      const nova = await api.post<Categoria>('/categorias', {
        nome: formCategoria.nome,
        cor: formCategoria.cor,
        ordem: formCategoria.ordem,
      })
      categorias.value.push(nova)
      addToast('Categoria criada', 'success')
    }
    fecharModalCategoria()
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    salvandoCategoria.value = false
  }
}

async function confirmarDeleteCategoria(cat: Categoria) {
  if (!confirm(`Remover a categoria "${cat.nome}"?`)) return
  try {
    await api.del(`/categorias/${cat.id}`)
    categorias.value = categorias.value.filter(c => c.id !== cat.id)
    addToast('Categoria removida', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

// ===== PRODUTOS =====
const produtos = ref<Produto[]>([])
const loadingProdutos = ref(false)
const filtroCategoriaId = ref<number | null>(null)

const produtosFiltrados = computed(() => {
  if (!filtroCategoriaId.value) return produtos.value
  return produtos.value.filter(p => p.categoria_id === filtroCategoriaId.value)
})

async function fetchProdutos() {
  try {
    loadingProdutos.value = true
    produtos.value = await api.get<Produto[]>('/produtos')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    loadingProdutos.value = false
  }
}

// ===== MODAL PRODUTO =====
const modalProdutoAberto = ref(false)
const produtoEditando = ref<Produto | null>(null)
const salvandoProduto = ref(false)

const formProduto = reactive<ProdutoInput>({
  nome: '',
  categoria_id: 0,
  unidade: 'un',
  unidade_preco: 'un',
  preco_ref: 0,
  ativo: true,
})

function abrirModalNovoProduto() {
  produtoEditando.value = null
  formProduto.nome = ''
  formProduto.categoria_id = categorias.value.length > 0 ? categorias.value[0].id : 0
  formProduto.unidade = 'un'
  formProduto.unidade_preco = 'un'
  formProduto.preco_ref = 0
  formProduto.ativo = true
  modalProdutoAberto.value = true
}

function abrirModalEditarProduto(prod: Produto) {
  produtoEditando.value = prod
  formProduto.nome = prod.nome
  formProduto.categoria_id = prod.categoria_id
  formProduto.unidade = prod.unidade
  formProduto.unidade_preco = prod.unidade_preco
  formProduto.preco_ref = prod.preco_ref
  formProduto.ativo = prod.ativo
  modalProdutoAberto.value = true
}

function fecharModalProduto() {
  modalProdutoAberto.value = false
  produtoEditando.value = null
}

async function salvarProduto() {
  salvandoProduto.value = true
  try {
    if (produtoEditando.value) {
      const updated = await api.put<Produto>(`/produtos/${produtoEditando.value.id}`, formProduto)
      const idx = produtos.value.findIndex(p => p.id === produtoEditando.value!.id)
      if (idx !== -1) produtos.value[idx] = updated
      addToast('Produto atualizado', 'success')
    } else {
      const novo = await api.post<Produto>('/produtos', formProduto)
      produtos.value.push(novo)
      addToast('Produto criado', 'success')
    }
    fecharModalProduto()
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    salvandoProduto.value = false
  }
}

async function confirmarDeleteProduto(prod: Produto) {
  if (!confirm(`Remover o produto "${prod.nome}"?`)) return
  try {
    const updated = await api.put<Produto>(`/produtos/${prod.id}`, { ...prod, ativo: false })
    const idx = produtos.value.findIndex(p => p.id === prod.id)
    if (idx !== -1) produtos.value[idx] = updated
    addToast('Produto desativado', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

function corDaCategoria(categoriaId: number): string {
  const cat = categorias.value.find(c => c.id === categoriaId)
  return cat?.cor ?? '#636E72'
}

// ===== UTILITÁRIOS =====

function formatCurrency(v: number): string {
  return v.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// ===== MOUNT =====
onMounted(async () => {
  await fetchListas()
  await fetchCategorias()
  await fetchProdutos()
})
</script>

<style scoped>
.section-title {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.theme-toggle-btn {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--color-primary-light);
  color: var(--color-primary);
  transition: transform 0.15s ease, background 0.15s ease;
}

.theme-toggle-btn:active {
  transform: scale(0.9);
}

.config-chevron {
  color: var(--color-text-secondary);
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.config-chevron.collapsed {
  transform: rotate(-90deg);
}
</style>
