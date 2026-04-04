package repository

import (
	"database/sql"
	"mercado-app/internal/models"
)

func (r *Repository) ListListas() ([]models.Lista, error) {
	rows, err := r.Query(`
		SELECT l.id, l.nome, l.mes, l.ano, l.renda, l.status, l.created_at,
			COALESCE(SUM(li.qtd * li.preco_estimado / li.duracao_meses), 0) AS total_planejado,
			COALESCE((SELECT SUM(ci.preco_total) FROM compra_itens ci JOIN compras c ON c.id = ci.compra_id WHERE c.lista_id = l.id), 0) AS total_comprado,
			COUNT(li.id) AS itens_total,
			COALESCE(SUM(li.comprado), 0) AS itens_comprados
		FROM listas l
		LEFT JOIN lista_itens li ON li.lista_id = l.id
		GROUP BY l.id
		ORDER BY l.ano DESC, l.mes DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listas []models.Lista
	for rows.Next() {
		var l models.Lista
		if err := rows.Scan(&l.ID, &l.Nome, &l.Mes, &l.Ano, &l.Renda, &l.Status, &l.CreatedAt,
			&l.TotalPlanejado, &l.TotalComprado, &l.ItensTotal, &l.ItensComprados); err != nil {
			return nil, err
		}
		listas = append(listas, l)
	}
	if listas == nil {
		listas = []models.Lista{}
	}
	return listas, rows.Err()
}

func (r *Repository) GetLista(id int64) (*models.Lista, error) {
	var l models.Lista
	err := r.QueryRow(`
		SELECT l.id, l.nome, l.mes, l.ano, l.renda, l.status, l.created_at,
			COALESCE(SUM(li.qtd * li.preco_estimado / li.duracao_meses), 0) AS total_planejado,
			COALESCE((SELECT SUM(ci.preco_total) FROM compra_itens ci JOIN compras c ON c.id = ci.compra_id WHERE c.lista_id = l.id), 0) AS total_comprado,
			COUNT(li.id) AS itens_total,
			COALESCE(SUM(li.comprado), 0) AS itens_comprados
		FROM listas l
		LEFT JOIN lista_itens li ON li.lista_id = l.id
		WHERE l.id = ?
		GROUP BY l.id`, id).
		Scan(&l.ID, &l.Nome, &l.Mes, &l.Ano, &l.Renda, &l.Status, &l.CreatedAt,
			&l.TotalPlanejado, &l.TotalComprado, &l.ItensTotal, &l.ItensComprados)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Lista")
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repository) CreateLista(input models.ListaInput) (*models.Lista, error) {
	if input.Nome == "" {
		return nil, models.ErrValidation("Nome é obrigatório")
	}
	if input.Mes < 1 || input.Mes > 12 {
		return nil, models.ErrValidation("Mês inválido")
	}
	if input.Ano < 2020 {
		return nil, models.ErrValidation("Ano inválido")
	}
	if input.Renda <= 0 {
		input.Renda = 6000
	}

	id, err := r.InsertReturningID(
		"INSERT INTO listas (nome, mes, ano, renda) VALUES (?, ?, ?, ?)",
		input.Nome, input.Mes, input.Ano, input.Renda,
	)
	if err != nil {
		return nil, err
	}
	return r.GetLista(id)
}

func (r *Repository) CopiarLista(input models.ListaInput, copiarDe int64) (*models.Lista, error) {
	lista, err := r.CreateLista(input)
	if err != nil {
		return nil, err
	}

	_, err = r.Exec(`
		INSERT INTO lista_itens (lista_id, produto_id, qtd, preco_estimado, duracao_meses, observacao)
		SELECT ?, li.produto_id, li.qtd, p.preco_ref, li.duracao_meses, li.observacao
		FROM lista_itens li
		JOIN produtos p ON p.id = li.produto_id
		WHERE li.lista_id = ?`,
		lista.ID, copiarDe)
	if err != nil {
		return nil, err
	}

	return r.GetLista(lista.ID)
}

func (r *Repository) UpdateLista(id int64, input models.ListaUpdateInput) (*models.Lista, error) {
	existing, err := r.GetLista(id)
	if err != nil {
		return nil, err
	}

	nome := existing.Nome
	renda := existing.Renda
	status := existing.Status

	if input.Nome != nil {
		nome = *input.Nome
	}
	if input.Renda != nil {
		renda = *input.Renda
	}
	if input.Status != nil {
		status = *input.Status
	}

	_, err = r.Exec("UPDATE listas SET nome = ?, renda = ?, status = ? WHERE id = ?",
		nome, renda, status, id)
	if err != nil {
		return nil, err
	}
	return r.GetLista(id)
}

func (r *Repository) DeleteLista(id int64) error {
	res, err := r.Exec("DELETE FROM listas WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrNotFound("Lista")
	}
	return nil
}

// ---------- Lista Itens ----------

func (r *Repository) ListItens(listaID int64) ([]models.ListaItem, error) {
	rows, err := r.Query(`
		SELECT li.id, li.lista_id, li.produto_id, li.qtd, li.preco_estimado, li.duracao_meses, li.comprado, li.observacao,
			p.nome, p.unidade, c.id, c.nome, c.cor, c.ordem
		FROM lista_itens li
		JOIN produtos p ON p.id = li.produto_id
		JOIN categorias c ON c.id = p.categoria_id
		WHERE li.lista_id = ?
		ORDER BY c.ordem, p.nome`, listaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []models.ListaItem
	for rows.Next() {
		var i models.ListaItem
		var comprado int
		if err := rows.Scan(&i.ID, &i.ListaID, &i.ProdutoID, &i.Qtd, &i.PrecoEstimado,
			&i.DuracaoMeses, &comprado, &i.Observacao,
			&i.ProdutoNome, &i.ProdutoUnidade, &i.CategoriaID, &i.CategoriaNome, &i.CategoriaCor, &i.CategoriaOrdem); err != nil {
			return nil, err
		}
		i.Comprado = comprado == 1
		itens = append(itens, i)
	}
	if itens == nil {
		itens = []models.ListaItem{}
	}
	return itens, rows.Err()
}

func (r *Repository) AddItem(listaID int64, input models.ListaItemInput) (*models.ListaItem, error) {
	if input.ProdutoID == 0 {
		return nil, models.ErrValidation("Produto é obrigatório")
	}
	if input.DuracaoMeses < 1 {
		input.DuracaoMeses = 1
	}
	if input.Qtd <= 0 {
		input.Qtd = 1
	}

	id, err := r.InsertReturningID(
		"INSERT INTO lista_itens (lista_id, produto_id, qtd, preco_estimado, duracao_meses, observacao) VALUES (?, ?, ?, ?, ?, ?)",
		listaID, input.ProdutoID, input.Qtd, input.PrecoEstimado, input.DuracaoMeses, input.Observacao,
	)
	if err != nil {
		return nil, err
	}
	return r.getItem(id)
}

func (r *Repository) UpdateItem(itemID int64, input models.ListaItemUpdateInput) (*models.ListaItem, error) {
	existing, err := r.getItem(itemID)
	if err != nil {
		return nil, err
	}

	qtd := existing.Qtd
	preco := existing.PrecoEstimado
	duracao := existing.DuracaoMeses
	obs := existing.Observacao

	if input.Qtd != nil {
		qtd = *input.Qtd
	}
	if input.PrecoEstimado != nil {
		preco = *input.PrecoEstimado
	}
	if input.DuracaoMeses != nil {
		duracao = *input.DuracaoMeses
	}
	if input.Observacao != nil {
		obs = input.Observacao
	}

	_, err = r.Exec(
		"UPDATE lista_itens SET qtd = ?, preco_estimado = ?, duracao_meses = ?, observacao = ? WHERE id = ?",
		qtd, preco, duracao, obs, itemID,
	)
	if err != nil {
		return nil, err
	}
	return r.getItem(itemID)
}

func (r *Repository) ToggleCheck(itemID int64) (*models.ListaItem, error) {
	res, err := r.Exec("UPDATE lista_itens SET comprado = 1 - comprado WHERE id = ?", itemID)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, models.ErrNotFound("Item")
	}
	return r.getItem(itemID)
}

func (r *Repository) DeleteItem(itemID int64) error {
	res, err := r.Exec("DELETE FROM lista_itens WHERE id = ?", itemID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrNotFound("Item")
	}
	return nil
}

func (r *Repository) getItem(id int64) (*models.ListaItem, error) {
	var i models.ListaItem
	var comprado int
	err := r.QueryRow(`
		SELECT li.id, li.lista_id, li.produto_id, li.qtd, li.preco_estimado, li.duracao_meses, li.comprado, li.observacao,
			p.nome, p.unidade, c.id, c.nome, c.cor, c.ordem
		FROM lista_itens li
		JOIN produtos p ON p.id = li.produto_id
		JOIN categorias c ON c.id = p.categoria_id
		WHERE li.id = ?`, id).
		Scan(&i.ID, &i.ListaID, &i.ProdutoID, &i.Qtd, &i.PrecoEstimado,
			&i.DuracaoMeses, &comprado, &i.Observacao,
			&i.ProdutoNome, &i.ProdutoUnidade, &i.CategoriaID, &i.CategoriaNome, &i.CategoriaCor, &i.CategoriaOrdem)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Item")
	}
	if err != nil {
		return nil, err
	}
	i.Comprado = comprado == 1
	return &i, nil
}
