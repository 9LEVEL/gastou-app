export interface Categoria {
  id: number
  nome: string
  cor: string
  icone: string | null
  ordem: number
}

export interface CategoriaInput {
  nome: string
  cor: string
  icone?: string | null
  ordem?: number
}

export interface Produto {
  id: number
  nome: string
  categoria_id: number
  unidade: string
  unidade_preco: string
  preco_ref: number
  ativo: boolean
  created_at: string
  categoria_nome?: string
  categoria_cor?: string
}

export interface ProdutoInput {
  nome: string
  categoria_id: number
  unidade: string
  unidade_preco: string
  preco_ref: number
  ativo?: boolean
}

export interface Lista {
  id: number
  nome: string
  mes: number
  ano: number
  renda: number
  status: string
  created_at: string
  total_planejado: number
  total_comprado: number
  itens_total: number
  itens_comprados: number
}

export interface ListaInput {
  nome: string
  mes: number
  ano: number
  renda: number
}

export interface ListaUpdateInput {
  nome?: string
  renda?: number
  status?: string
}

export interface ListaItem {
  id: number
  lista_id: number
  produto_id: number
  qtd: number
  preco_estimado: number
  duracao_meses: number
  comprado: boolean
  observacao: string | null
  produto_nome?: string
  produto_unidade?: string
  categoria_id?: number
  categoria_nome?: string
  categoria_cor?: string
  categoria_ordem?: number
}

export interface ListaItemInput {
  produto_id: number
  qtd: number
  preco_estimado: number
  duracao_meses: number
  observacao?: string | null
}

export interface ListaItemUpdateInput {
  qtd?: number
  preco_estimado?: number
  duracao_meses?: number
  observacao?: string | null
}

export interface Compra {
  id: number
  lista_id: number | null
  local: string
  data: string
  total_nfe: number | null
  observacao: string | null
  created_at: string
  total_calculado: number
  itens?: CompraItem[]
}

export interface CompraInput {
  lista_id?: number | null
  local: string
  data: string
  total_nfe?: number | null
  observacao?: string | null
}

export interface CompraItem {
  id: number
  compra_id: number
  produto_id: number | null
  nome_nfe: string
  qtd: number
  unidade: string
  preco_unit: number
  preco_total: number
  lista_item_id: number | null
  produto_nome?: string
}

export interface CompraItemInput {
  produto_id?: number | null
  nome_nfe: string
  qtd: number
  unidade: string
  preco_unit: number
  preco_total: number
  lista_item_id?: number | null
}

export interface PrecoHistorico {
  data: string
  local: string
  preco_unit: number
  qtd: number
  preco_total: number
}

export interface DashboardResumo {
  total_planejado: number
  total_gasto: number
  saldo_restante: number
  percentual_renda: number
  renda: number
  itens_faltantes: number
  itens_comprados: number
  itens_total: number
}

export interface CategoriaComparativo {
  categoria_id: number
  categoria_nome: string
  categoria_cor: string
  planejado: number
  real: number
}

export interface EvolucaoMensal {
  mes: number
  ano: number
  nome: string
  gasto: number
}

export interface ApiError {
  message: string
  code: string
}
