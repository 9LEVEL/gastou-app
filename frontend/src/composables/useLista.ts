import { ref, computed, watch } from 'vue'
import { api } from '../api/client'
import { useToast } from './useToast'
import type {
  Lista,
  ListaItem,
  ListaItemInput,
  ListaItemUpdateInput,
  ListaInput,
  ListaUpdateInput
} from '../types'

const LISTA_ATIVA_KEY = 'gastou-lista-ativa-id'

const listas = ref<Lista[]>([])
const listaAtiva = ref<Lista | null>(null)
const itens = ref<ListaItem[]>([])
const loading = ref(false)

const { addToast } = useToast()

const totalMensalAjustado = computed(() => {
  return itens.value.reduce((sum, item) => {
    return sum + (item.qtd * item.preco_estimado) / (item.duracao_meses || 1)
  }, 0)
})

const totalComprado = computed(() => {
  return itens.value
    .filter(i => i.comprado)
    .reduce((sum, item) => {
      return sum + (item.qtd * item.preco_estimado) / (item.duracao_meses || 1)
    }, 0)
})

const itensFaltantes = computed(() => {
  return itens.value.filter(i => !i.comprado).length
})

const percentualRenda = computed(() => {
  if (!listaAtiva.value || listaAtiva.value.renda === 0) return 0
  return (totalMensalAjustado.value / listaAtiva.value.renda) * 100
})

async function fetchListas() {
  try {
    loading.value = true
    listas.value = await api.get<Lista[]>('/listas')
    if (listas.value.length > 0 && !listaAtiva.value) {
      // Try to restore from localStorage
      const storedId = localStorage.getItem(LISTA_ATIVA_KEY)
      if (storedId) {
        const stored = listas.value.find(l => l.id === Number(storedId))
        if (stored) {
          listaAtiva.value = stored
          return
        }
      }
      // Fallback: find the 'ativa' one or first
      const ativa = listas.value.find(l => l.status === 'ativa')
      listaAtiva.value = ativa || listas.value[0]
    }
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    loading.value = false
  }
}

watch(listaAtiva, (lista) => {
  if (lista) {
    try {
      localStorage.setItem(LISTA_ATIVA_KEY, String(lista.id))
    } catch {
      // localStorage full
    }
  }
})

async function fetchItens(listaId: number) {
  try {
    loading.value = true
    itens.value = await api.get<ListaItem[]>(`/listas/${listaId}/itens`)
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  } finally {
    loading.value = false
  }
}

async function toggleCheck(listaId: number, itemId: number) {
  try {
    const item = itens.value.find(i => i.id === itemId)
    if (item) {
      item.comprado = !item.comprado
    }
    await api.patch<ListaItem>(`/listas/${listaId}/itens/${itemId}/check`)
  } catch (e: unknown) {
    const item = itens.value.find(i => i.id === itemId)
    if (item) {
      item.comprado = !item.comprado
    }
    addToast((e as Error).message, 'error')
  }
}

async function addItem(listaId: number, input: ListaItemInput) {
  try {
    const newItem = await api.post<ListaItem>(`/listas/${listaId}/itens`, input)
    itens.value.push(newItem)
    addToast('Item adicionado', 'success')
    return newItem
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function updateItem(listaId: number, itemId: number, input: ListaItemUpdateInput) {
  try {
    const updated = await api.put<ListaItem>(`/listas/${listaId}/itens/${itemId}`, input)
    const index = itens.value.findIndex(i => i.id === itemId)
    if (index !== -1) {
      itens.value[index] = { ...itens.value[index], ...updated }
    }
    addToast('Item atualizado', 'success')
    return updated
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function deleteItem(listaId: number, itemId: number) {
  try {
    await api.del(`/listas/${listaId}/itens/${itemId}`)
    itens.value = itens.value.filter(i => i.id !== itemId)
    addToast('Item removido', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

async function copiarLista(fromId: number, input: ListaInput) {
  try {
    const nova = await api.post<Lista>(`/listas?copiar_de=${fromId}`, input)
    listas.value.push(nova)
    addToast('Lista copiada com sucesso', 'success')
    return nova
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function limparChecks(listaId: number) {
  try {
    for (const item of itens.value) {
      if (item.comprado) {
        item.comprado = false
        await api.patch(`/listas/${listaId}/itens/${item.id}/check`)
      }
    }
    addToast('Checks limpos', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

async function updateLista(id: number, input: ListaUpdateInput) {
  try {
    const updated = await api.put<Lista>(`/listas/${id}`, input)
    const index = listas.value.findIndex(l => l.id === id)
    if (index !== -1) {
      listas.value[index] = { ...listas.value[index], ...updated }
    }
    if (listaAtiva.value?.id === id) {
      listaAtiva.value = { ...listaAtiva.value, ...updated }
    }
    return updated
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function createLista(input: ListaInput) {
  try {
    const nova = await api.post<Lista>('/listas', input)
    listas.value.push(nova)
    addToast('Lista criada', 'success')
    return nova
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
    return null
  }
}

async function deleteLista(id: number) {
  try {
    await api.del(`/listas/${id}`)
    listas.value = listas.value.filter(l => l.id !== id)
    if (listaAtiva.value?.id === id) {
      listaAtiva.value = listas.value.length > 0 ? listas.value[0] : null
    }
    addToast('Lista removida', 'success')
  } catch (e: unknown) {
    addToast((e as Error).message, 'error')
  }
}

export function useLista() {
  return {
    listas,
    listaAtiva,
    itens,
    loading,
    totalMensalAjustado,
    totalComprado,
    itensFaltantes,
    percentualRenda,
    fetchListas,
    fetchItens,
    toggleCheck,
    addItem,
    updateItem,
    deleteItem,
    copiarLista,
    limparChecks,
    updateLista,
    createLista,
    deleteLista
  }
}
