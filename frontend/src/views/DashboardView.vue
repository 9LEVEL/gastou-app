<template>
  <div>
    <!-- Header -->
    <div class="sticky-header">
      <div class="flex items-center justify-between" style="padding: 0 16px;">
        <div>
          <h1 class="font-display" style="margin: 0; font-size: 1.25rem;">Dashboard</h1>
          <div v-if="listaAtiva" class="text-secondary text-sm">
            {{ listaAtiva.nome }}
          </div>
        </div>
      </div>
    </div>

    <!-- Empty state: sem lista ativa -->
    <div v-if="!listaAtiva" class="empty-state" style="padding: 48px 16px; text-align: center;">
      <div class="text-secondary">Nenhuma lista ativa encontrada.</div>
    </div>

    <template v-else>
      <!-- Loading -->
      <div v-if="loading" style="display: flex; justify-content: center; padding: 32px;">
        <div class="loading-spinner"></div>
      </div>

      <template v-else>
        <!-- 1. Cards de resumo 2x2 -->
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px; padding: 16px 16px 0;">
          <!-- Planejado -->
          <div class="card" style="padding: 12px;">
            <div class="text-xs text-secondary" style="margin-bottom: 4px;">Planejado</div>
            <div class="num font-display" style="font-size: 1.1rem; color: var(--color-primary);">
              R$ {{ fmt(resumo?.total_planejado ?? 0) }}
            </div>
          </div>

          <!-- Gasto -->
          <div class="card" style="padding: 12px;">
            <div class="text-xs text-secondary" style="margin-bottom: 4px;">Gasto</div>
            <div
              class="num font-display"
              style="font-size: 1.1rem;"
              :style="{ color: gastoExcedeu ? 'var(--color-danger)' : 'inherit' }"
            >
              R$ {{ fmt(resumo?.total_gasto ?? 0) }}
            </div>
          </div>

          <!-- Saldo -->
          <div class="card" style="padding: 12px;">
            <div class="text-xs text-secondary" style="margin-bottom: 4px;">Saldo</div>
            <div
              class="num font-display"
              style="font-size: 1.1rem;"
              :style="{ color: saldoPositivo ? 'var(--color-success)' : 'var(--color-danger)' }"
            >
              R$ {{ fmt(resumo?.saldo_restante ?? 0) }}
            </div>
          </div>

          <!-- Progresso -->
          <div class="card" style="padding: 12px;">
            <div class="text-xs text-secondary" style="margin-bottom: 4px;">Comprados</div>
            <div class="num font-display" style="font-size: 1rem;">
              {{ resumo?.itens_comprados ?? 0 }}/{{ resumo?.itens_total ?? 0 }}
            </div>
            <div style="margin-top: 4px; height: 4px; background: var(--color-border); border-radius: 2px; overflow: hidden;">
              <div style="height: 100%; border-radius: 2px; background: var(--color-primary); transition: width 0.5s ease;" :style="{ width: progressoPct + '%' }"></div>
            </div>
            <div style="margin-top: 4px;">
              <span
                class="chip"
                :class="chipRendaClass"
                style="font-size: 0.7rem;"
              >
                {{ fmt(resumo?.percentual_renda ?? 0) }}% renda
              </span>
            </div>
          </div>
        </div>

        <!-- 2. Estratégia Semanal -->
        <div v-if="resumo" style="padding: 8px 16px 0;">
          <div class="card" style="padding: 12px;">
            <div class="text-xs text-secondary" style="margin-bottom: 4px;">Estratégia Semanal</div>
            <div class="text-sm">
              Você pode gastar
              <span class="num" style="font-weight: 600; color: var(--color-primary);">
                ~R$ {{ fmt(gastoSemanal) }}
              </span>
              por semana
            </div>
          </div>
        </div>

        <!-- 3. Comparativo por Categoria -->
        <div style="padding: 16px 16px 0;">
          <div class="text-sm" style="font-weight: 600; margin-bottom: 8px;">Planejado vs Real</div>
          <div v-if="comparativo.length === 0" class="text-secondary text-sm">
            Sem dados de categoria.
          </div>
          <BarChart v-else :items="barChartItems" />
        </div>

        <!-- 4. Evolução Mensal -->
        <div style="padding: 16px 16px 24px;">
          <div class="text-sm" style="font-weight: 600; margin-bottom: 8px;">Evolução</div>
          <div v-if="evolucao.length === 0" class="text-secondary text-sm">
            Sem histórico mensal.
          </div>
          <div v-else>
            <div
              v-for="(item, index) in evolucao"
              :key="index"
              style="margin-bottom: 10px;"
            >
              <div class="flex items-center justify-between text-sm" style="margin-bottom: 3px;">
                <span>{{ formatMes(item.mes, item.ano) }}</span>
                <span class="num text-xs" style="font-weight: 600;">R$ {{ fmt(item.gasto) }}</span>
              </div>
              <div style="height: 6px; background: var(--color-border); border-radius: 3px; overflow: hidden;">
                <div
                  style="height: 100%; border-radius: 3px; background: var(--color-primary); transition: width 0.5s ease;"
                  :style="{ width: maxEvolucao > 0 ? (item.gasto / maxEvolucao * 100) + '%' : '0%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useLista } from '../composables/useLista'
import { api } from '../api/client'
import BarChart from '../components/BarChart.vue'

interface DashboardResumo {
  total_planejado: number
  total_gasto: number
  saldo_restante: number
  percentual_renda: number
  renda: number
  itens_faltantes: number
  itens_comprados: number
  itens_total: number
}

interface CategoriaComparativo {
  categoria_id: number
  categoria_nome: string
  categoria_cor: string
  planejado: number
  real: number
}

interface EvolucaoMensal {
  mes: number
  ano: number
  nome: string
  gasto: number
}

const { listaAtiva, fetchListas } = useLista()

const loading = ref(false)
const resumo = ref<DashboardResumo | null>(null)
const comparativo = ref<CategoriaComparativo[]>([])
const evolucao = ref<EvolucaoMensal[]>([])

function fmt(v: number): string {
  return v.toFixed(2).replace('.', ',')
}

function formatMes(mes: number, ano: number): string {
  const nomes = ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez']
  return `${nomes[mes - 1] ?? mes}/${ano}`
}

const gastoExcedeu = computed(() => {
  if (!resumo.value) return false
  return resumo.value.total_gasto > resumo.value.total_planejado
})

const saldoPositivo = computed(() => {
  if (!resumo.value) return true
  return resumo.value.saldo_restante >= 0
})

const chipRendaClass = computed(() => {
  const pct = resumo.value?.percentual_renda ?? 0
  if (pct >= 80) return 'chip-danger'
  if (pct >= 60) return 'chip-warning'
  return 'chip-primary'
})

const barChartItems = computed(() =>
  comparativo.value.map(c => ({
    label: c.categoria_nome,
    planejado: c.planejado,
    real: c.real,
    cor: c.categoria_cor
  }))
)

const maxEvolucao = computed(() => {
  return evolucao.value.reduce((max, item) => Math.max(max, item.gasto), 0)
})

const progressoPct = computed(() => {
  if (!resumo.value || resumo.value.itens_total === 0) return 0
  return (resumo.value.itens_comprados / resumo.value.itens_total) * 100
})

const gastoSemanal = computed(() => {
  if (!resumo.value) return 0
  const saldo = resumo.value.saldo_restante
  if (saldo <= 0) return 0
  const hoje = new Date()
  const diasNoMes = new Date(hoje.getFullYear(), hoje.getMonth() + 1, 0).getDate()
  const diasRestantes = diasNoMes - hoje.getDate()
  const semanasRestantes = diasRestantes / 7
  if (semanasRestantes <= 0) return saldo
  return saldo / semanasRestantes
})

async function fetchDashboard(id: number) {
  loading.value = true
  try {
    const [r, c, e] = await Promise.allSettled([
      api.get<DashboardResumo>(`/listas/${id}/dashboard/resumo`),
      api.get<CategoriaComparativo[]>(`/listas/${id}/dashboard/comparativo`),
      api.get<EvolucaoMensal[]>('/dashboard/evolucao')
    ])

    resumo.value = r.status === 'fulfilled' ? r.value : null
    comparativo.value = c.status === 'fulfilled' ? c.value : []
    evolucao.value = e.status === 'fulfilled' ? e.value : []
  } finally {
    loading.value = false
  }
}

watch(listaAtiva, (lista) => {
  if (lista) {
    fetchDashboard(lista.id)
  }
})

onMounted(async () => {
  await fetchListas()
  if (listaAtiva.value) {
    fetchDashboard(listaAtiva.value.id)
  }
})
</script>
