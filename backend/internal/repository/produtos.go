package repository

import (
	"database/sql"
	"fmt"
	"mercado-app/internal/models"
	"strings"
)

func (r *Repository) ListProdutos(categoriaID *int64, ativo *bool) ([]models.Produto, error) {
	query := `SELECT p.id, p.nome, p.categoria_id, p.unidade, p.unidade_preco, p.preco_ref, p.ativo, p.created_at,
		c.nome, c.cor
		FROM produtos p JOIN categorias c ON c.id = p.categoria_id`

	var conditions []string
	var args []interface{}

	if categoriaID != nil {
		conditions = append(conditions, "p.categoria_id = ?")
		args = append(args, *categoriaID)
	}
	if ativo != nil {
		val := 0
		if *ativo {
			val = 1
		}
		conditions = append(conditions, "p.ativo = ?")
		args = append(args, val)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY c.ordem, p.nome"

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prods []models.Produto
	for rows.Next() {
		var p models.Produto
		var ativoInt int
		if err := rows.Scan(&p.ID, &p.Nome, &p.CategoriaID, &p.Unidade, &p.UnidadePreco,
			&p.PrecoRef, &ativoInt, &p.CreatedAt, &p.CategoriaNome, &p.CategoriaCor); err != nil {
			return nil, err
		}
		p.Ativo = ativoInt == 1
		prods = append(prods, p)
	}
	if prods == nil {
		prods = []models.Produto{}
	}
	return prods, rows.Err()
}

func (r *Repository) GetProduto(id int64) (*models.Produto, error) {
	var p models.Produto
	var ativoInt int
	err := r.QueryRow(`SELECT p.id, p.nome, p.categoria_id, p.unidade, p.unidade_preco, p.preco_ref, p.ativo, p.created_at,
		c.nome, c.cor
		FROM produtos p JOIN categorias c ON c.id = p.categoria_id WHERE p.id = ?`, id).
		Scan(&p.ID, &p.Nome, &p.CategoriaID, &p.Unidade, &p.UnidadePreco,
			&p.PrecoRef, &ativoInt, &p.CreatedAt, &p.CategoriaNome, &p.CategoriaCor)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Produto")
	}
	if err != nil {
		return nil, err
	}
	p.Ativo = ativoInt == 1
	return &p, nil
}

func (r *Repository) CreateProduto(input models.ProdutoInput) (*models.Produto, error) {
	if input.Nome == "" {
		return nil, models.ErrValidation("Nome é obrigatório")
	}
	if input.CategoriaID == 0 {
		return nil, models.ErrValidation("Categoria é obrigatória")
	}
	if input.Unidade == "" {
		input.Unidade = "un"
	}
	if input.UnidadePreco == "" {
		input.UnidadePreco = "un"
	}

	ativo := 1
	if input.Ativo != nil && !*input.Ativo {
		ativo = 0
	}

	id, err := r.InsertReturningID(
		"INSERT INTO produtos (nome, categoria_id, unidade, unidade_preco, preco_ref, ativo) VALUES (?, ?, ?, ?, ?, ?)",
		input.Nome, input.CategoriaID, input.Unidade, input.UnidadePreco, input.PrecoRef, ativo,
	)
	if err != nil {
		return nil, err
	}
	return r.GetProduto(id)
}

func (r *Repository) UpdateProduto(id int64, input models.ProdutoInput) (*models.Produto, error) {
	if input.Nome == "" {
		return nil, models.ErrValidation("Nome é obrigatório")
	}

	ativo := 1
	if input.Ativo != nil && !*input.Ativo {
		ativo = 0
	}

	res, err := r.Exec(
		`UPDATE produtos SET nome = ?, categoria_id = ?, unidade = ?, unidade_preco = ?, preco_ref = ?, ativo = ? WHERE id = ?`,
		input.Nome, input.CategoriaID, input.Unidade, input.UnidadePreco, input.PrecoRef, ativo, id,
	)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, models.ErrNotFound("Produto")
	}
	return r.GetProduto(id)
}

func (r *Repository) GetHistoricoPrecos(produtoID int64) ([]models.PrecoHistorico, error) {
	rows, err := r.Query(`
		SELECT c.data, c.local, ci.preco_unit, ci.qtd, ci.preco_total
		FROM compra_itens ci
		JOIN compras c ON c.id = ci.compra_id
		WHERE ci.produto_id = ?
		ORDER BY c.data DESC`,
		produtoID)
	if err != nil {
		return nil, fmt.Errorf("querying historico: %w", err)
	}
	defer rows.Close()

	var hist []models.PrecoHistorico
	for rows.Next() {
		var h models.PrecoHistorico
		if err := rows.Scan(&h.Data, &h.Local, &h.PrecoUnit, &h.Qtd, &h.PrecoTotal); err != nil {
			return nil, err
		}
		hist = append(hist, h)
	}
	if hist == nil {
		hist = []models.PrecoHistorico{}
	}
	return hist, rows.Err()
}
