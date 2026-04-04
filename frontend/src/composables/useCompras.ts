import { ref } from 'vue'
import { api } from '../api/client'
import { useToast } from './useToast'
import type { Compra, CompraInput, CompraItem, CompraItemInput } from '../types'

const compras = ref<Compra[]>([])
const compraAtiva = ref<Compra | null>(null)
const loading = ref(false)

const { addToast } = useToast()

async function fetchCompras(listaId?: number) {
  try {
    loading.value = true
    const query = listaId ? `?lista_id=${listaId}` : ''
    compras.value = await api.get<Compra[]>(`/compras${query}`)
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    loading.value = false
  }
}

async function fetchCompra(id: number) {
  try {
    loading.value = true
    compraAtiva.value = await api.get<Compra>(`/compras/${id}`)
    return compraAtiva.value
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  } finally {
    loading.value = false
  }
}

async function createCompra(input: CompraInput) {
  try {
    const nova = await api.post<Compra>('/compras', input)
    compras.value.unshift(nova)
    compraAtiva.value = nova
    addToast('Compra registrada', 'success')
    return nova
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function updateCompra(id: number, input: CompraInput) {
  try {
    const updated = await api.put<Compra>(`/compras/${id}`, input)
    const index = compras.value.findIndex(c => c.id === id)
    if (index !== -1) {
      compras.value[index] = { ...compras.value[index], ...updated }
    }
    if (compraAtiva.value?.id === id) {
      compraAtiva.value = { ...compraAtiva.value, ...updated }
    }
    addToast('Compra atualizada', 'success')
    return updated
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function deleteCompra(id: number) {
  try {
    await api.del(`/compras/${id}`)
    compras.value = compras.value.filter(c => c.id !== id)
    if (compraAtiva.value?.id === id) {
      compraAtiva.value = null
    }
    addToast('Compra removida', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

async function addCompraItem(compraId: number, input: CompraItemInput) {
  try {
    const item = await api.post<CompraItem>(`/compras/${compraId}/itens`, input)
    if (compraAtiva.value?.id === compraId && compraAtiva.value.itens) {
      compraAtiva.value.itens.push(item)
      compraAtiva.value.total_calculado += item.preco_total
    }
    addToast('Item adicionado', 'success')
    return item
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function updateCompraItem(compraId: number, itemId: number, input: CompraItemInput) {
  try {
    const updated = await api.put<CompraItem>(`/compras/${compraId}/itens/${itemId}`, input)
    if (compraAtiva.value?.id === compraId && compraAtiva.value.itens) {
      const index = compraAtiva.value.itens.findIndex(i => i.id === itemId)
      if (index !== -1) {
        compraAtiva.value.itens[index] = { ...compraAtiva.value.itens[index], ...updated }
      }
    }
    addToast('Item atualizado', 'success')
    return updated
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function deleteCompraItem(compraId: number, itemId: number) {
  try {
    await api.del(`/compras/${compraId}/itens/${itemId}`)
    if (compraAtiva.value?.id === compraId && compraAtiva.value.itens) {
      const item = compraAtiva.value.itens.find(i => i.id === itemId)
      if (item) {
        compraAtiva.value.total_calculado -= item.preco_total
      }
      compraAtiva.value.itens = compraAtiva.value.itens.filter(i => i.id !== itemId)
    }
    addToast('Item removido', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

export function useCompras() {
  return {
    compras,
    compraAtiva,
    loading,
    fetchCompras,
    fetchCompra,
    createCompra,
    updateCompra,
    deleteCompra,
    addCompraItem,
    updateCompraItem,
    deleteCompraItem
  }
}
