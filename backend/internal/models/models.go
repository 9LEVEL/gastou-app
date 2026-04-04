package models

import "fmt"

// ---------- AppError ----------

type AppError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string { return e.Message }

func ErrValidation(msg string) *AppError {
	return &AppError{Message: msg, Code: "VALIDATION_ERROR", Status: 400}
}

func ErrNotFound(entity string) *AppError {
	return &AppError{Message: fmt.Sprintf("%s não encontrado(a)", entity), Code: "NOT_FOUND", Status: 404}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Message: msg, Code: "CONFLICT", Status: 409}
}

func ErrInternal() *AppError {
	return &AppError{Message: "Erro interno do servidor", Code: "INTERNAL_ERROR", Status: 500}
}

// ---------- Categoria ----------

type Categoria struct {
	ID    int64   `json:"id"`
	Nome  string  `json:"nome"`
	Cor   string  `json:"cor"`
	Icone *string `json:"icone"`
	Ordem int     `json:"ordem"`
}

type CategoriaInput struct {
	Nome  string  `json:"nome"`
	Cor   string  `json:"cor"`
	Icone *string `json:"icone"`
	Ordem int     `json:"ordem"`
}

// ---------- Produto ----------

type Produto struct {
	ID            int64   `json:"id"`
	Nome          string  `json:"nome"`
	CategoriaID   int64   `json:"categoria_id"`
	Unidade       string  `json:"unidade"`
	UnidadePreco  string  `json:"unidade_preco"`
	PrecoRef      float64 `json:"preco_ref"`
	Ativo         bool    `json:"ativo"`
	CreatedAt     string  `json:"created_at"`
	CategoriaNome string  `json:"categoria_nome,omitempty"`
	CategoriaCor  string  `json:"categoria_cor,omitempty"`
}

type ProdutoInput struct {
	Nome         string  `json:"nome"`
	CategoriaID  int64   `json:"categoria_id"`
	Unidade      string  `json:"unidade"`
	UnidadePreco string  `json:"unidade_preco"`
	PrecoRef     float64 `json:"preco_ref"`
	Ativo        *bool   `json:"ativo"`
}

// ---------- Lista ----------

type Lista struct {
	ID              int64   `json:"id"`
	Nome            string  `json:"nome"`
	Mes             int     `json:"mes"`
	Ano             int     `json:"ano"`
	Renda           float64 `json:"renda"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	TotalPlanejado  float64 `json:"total_planejado"`
	TotalComprado   float64 `json:"total_comprado"`
	ItensTotal      int     `json:"itens_total"`
	ItensComprados  int     `json:"itens_comprados"`
}

type ListaInput struct {
	Nome  string  `json:"nome"`
	Mes   int     `json:"mes"`
	Ano   int     `json:"ano"`
	Renda float64 `json:"renda"`
}

type ListaUpdateInput struct {
	Nome   *string  `json:"nome"`
	Renda  *float64 `json:"renda"`
	Status *string  `json:"status"`
}

// ---------- ListaItem ----------

type ListaItem struct {
	ID             int64   `json:"id"`
	ListaID        int64   `json:"lista_id"`
	ProdutoID      int64   `json:"produto_id"`
	Qtd            float64 `json:"qtd"`
	PrecoEstimado  float64 `json:"preco_estimado"`
	DuracaoMeses   int     `json:"duracao_meses"`
	Comprado       bool    `json:"comprado"`
	Observacao     *string `json:"observacao"`
	ProdutoNome    string  `json:"produto_nome,omitempty"`
	ProdutoUnidade string  `json:"produto_unidade,omitempty"`
	CategoriaID    int64   `json:"categoria_id,omitempty"`
	CategoriaNome  string  `json:"categoria_nome,omitempty"`
	CategoriaCor   string  `json:"categoria_cor,omitempty"`
	CategoriaOrdem int     `json:"categoria_ordem,omitempty"`
}

type ListaItemInput struct {
	ProdutoID    int64   `json:"produto_id"`
	Qtd          float64 `json:"qtd"`
	PrecoEstimado float64 `json:"preco_estimado"`
	DuracaoMeses int     `json:"duracao_meses"`
	Observacao   *string `json:"observacao"`
}

type ListaItemUpdateInput struct {
	Qtd          *float64 `json:"qtd"`
	PrecoEstimado *float64 `json:"preco_estimado"`
	DuracaoMeses *int     `json:"duracao_meses"`
	Observacao   *string  `json:"observacao"`
}

// ---------- Compra ----------

type Compra struct {
	ID             int64       `json:"id"`
	ListaID        *int64      `json:"lista_id"`
	Local          string      `json:"local"`
	Data           string      `json:"data"`
	TotalNFe       *float64    `json:"total_nfe"`
	Observacao     *string     `json:"observacao"`
	CreatedAt      string      `json:"created_at"`
	TotalCalculado float64     `json:"total_calculado"`
	Itens          []CompraItem `json:"itens,omitempty"`
}

type CompraInput struct {
	ListaID    *int64   `json:"lista_id"`
	Local      string   `json:"local"`
	Data       string   `json:"data"`
	TotalNFe   *float64 `json:"total_nfe"`
	Observacao *string  `json:"observacao"`
}

// ---------- CompraItem ----------

type CompraItem struct {
	ID          int64   `json:"id"`
	CompraID    int64   `json:"compra_id"`
	ProdutoID   *int64  `json:"produto_id"`
	NomeNFe     string  `json:"nome_nfe"`
	Qtd         float64 `json:"qtd"`
	Unidade     string  `json:"unidade"`
	PrecoUnit   float64 `json:"preco_unit"`
	PrecoTotal  float64 `json:"preco_total"`
	ListaItemID *int64  `json:"lista_item_id"`
	ProdutoNome string  `json:"produto_nome,omitempty"`
}

type CompraItemInput struct {
	ProdutoID   *int64  `json:"produto_id"`
	NomeNFe     string  `json:"nome_nfe"`
	Qtd         float64 `json:"qtd"`
	Unidade     string  `json:"unidade"`
	PrecoUnit   float64 `json:"preco_unit"`
	PrecoTotal  float64 `json:"preco_total"`
	ListaItemID *int64  `json:"lista_item_id"`
}

// ---------- PrecoHistorico ----------

type PrecoHistorico struct {
	Data      string  `json:"data"`
	Local     string  `json:"local"`
	PrecoUnit float64 `json:"preco_unit"`
	Qtd       float64 `json:"qtd"`
	PrecoTotal float64 `json:"preco_total"`
}

// ---------- Dashboard ----------

type DashboardResumo struct {
	TotalPlanejado float64 `json:"total_planejado"`
	TotalGasto     float64 `json:"total_gasto"`
	SaldoRestante  float64 `json:"saldo_restante"`
	PercentualRenda float64 `json:"percentual_renda"`
	Renda          float64 `json:"renda"`
	ItensFaltantes int     `json:"itens_faltantes"`
	ItensComprados int     `json:"itens_comprados"`
	ItensTotal     int     `json:"itens_total"`
}

type CategoriaComparativo struct {
	CategoriaID   int64   `json:"categoria_id"`
	CategoriaNome string  `json:"categoria_nome"`
	CategoriaCor  string  `json:"categoria_cor"`
	Planejado     float64 `json:"planejado"`
	Real          float64 `json:"real"`
}

type EvolucaoMensal struct {
	Mes   int     `json:"mes"`
	Ano   int     `json:"ano"`
	Nome  string  `json:"nome"`
	Gasto float64 `json:"gasto"`
}
