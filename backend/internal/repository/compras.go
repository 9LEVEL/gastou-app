package repository

import (
	"database/sql"
	"mercado-app/internal/models"
)

func (r *Repository) ListCompras(listaID *int64) ([]models.Compra, error) {
	query := `
		SELECT c.id, c.lista_id, c.local, c.data, c.total_nfe, c.observacao, c.created_at,
			COALESCE(SUM(ci.preco_total), 0) AS total_calculado
		FROM compras c
		LEFT JOIN compra_itens ci ON ci.compra_id = c.id`

	var args []interface{}
	if listaID != nil {
		query += " WHERE c.lista_id = ?"
		args = append(args, *listaID)
	}
	query += " GROUP BY c.id ORDER BY c.data DESC"

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compras []models.Compra
	for rows.Next() {
		var c models.Compra
		if err := rows.Scan(&c.ID, &c.ListaID, &c.Local, &c.Data, &c.TotalNFe, &c.Observacao, &c.CreatedAt, &c.TotalCalculado); err != nil {
			return nil, err
		}
		compras = append(compras, c)
	}
	if compras == nil {
		compras = []models.Compra{}
	}
	return compras, rows.Err()
}

func (r *Repository) GetCompra(id int64) (*models.Compra, error) {
	var c models.Compra
	err := r.QueryRow(`
		SELECT c.id, c.lista_id, c.local, c.data, c.total_nfe, c.observacao, c.created_at,
			COALESCE((SELECT SUM(ci.preco_total) FROM compra_itens ci WHERE ci.compra_id = c.id), 0)
		FROM compras c WHERE c.id = ?`, id).
		Scan(&c.ID, &c.ListaID, &c.Local, &c.Data, &c.TotalNFe, &c.Observacao, &c.CreatedAt, &c.TotalCalculado)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Compra")
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.Query(`
		SELECT ci.id, ci.compra_id, ci.produto_id, ci.nome_nfe, ci.qtd, ci.unidade, ci.preco_unit, ci.preco_total, ci.lista_item_id,
			COALESCE(p.nome, '')
		FROM compra_itens ci
		LEFT JOIN produtos p ON p.id = ci.produto_id
		WHERE ci.compra_id = ?
		ORDER BY ci.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ci models.CompraItem
		if err := rows.Scan(&ci.ID, &ci.CompraID, &ci.ProdutoID, &ci.NomeNFe, &ci.Qtd, &ci.Unidade,
			&ci.PrecoUnit, &ci.PrecoTotal, &ci.ListaItemID, &ci.ProdutoNome); err != nil {
			return nil, err
		}
		c.Itens = append(c.Itens, ci)
	}
	if c.Itens == nil {
		c.Itens = []models.CompraItem{}
	}
	return &c, rows.Err()
}

func (r *Repository) CreateCompra(input models.CompraInput) (*models.Compra, error) {
	if input.Data == "" {
		return nil, models.ErrValidation("Data é obrigatória")
	}

	id, err := r.InsertReturningID(
		"INSERT INTO compras (lista_id, local, data, total_nfe, observacao) VALUES (?, ?, ?, ?, ?)",
		input.ListaID, input.Local, input.Data, input.TotalNFe, input.Observacao,
	)
	if err != nil {
		return nil, err
	}
	return r.GetCompra(id)
}

func (r *Repository) UpdateCompra(id int64, input models.CompraInput) (*models.Compra, error) {
	res, err := r.Exec(
		"UPDATE compras SET lista_id = ?, local = ?, data = ?, total_nfe = ?, observacao = ? WHERE id = ?",
		input.ListaID, input.Local, input.Data, input.TotalNFe, input.Observacao, id,
	)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, models.ErrNotFound("Compra")
	}
	return r.GetCompra(id)
}

func (r *Repository) DeleteCompra(id int64) error {
	res, err := r.Exec("DELETE FROM compras WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrNotFound("Compra")
	}
	return nil
}

func (r *Repository) AddCompraItem(compraID int64, input models.CompraItemInput) (*models.CompraItem, error) {
	if input.NomeNFe == "" {
		return nil, models.ErrValidation("Nome NFe é obrigatório")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	itemID, err := r.txInsertReturningID(tx,
		`INSERT INTO compra_itens (compra_id, produto_id, nome_nfe, qtd, unidade, preco_unit, preco_total, lista_item_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		compraID, input.ProdutoID, input.NomeNFe, input.Qtd, input.Unidade, input.PrecoUnit, input.PrecoTotal, input.ListaItemID,
	)
	if err != nil {
		return nil, err
	}

	// Update produto preco_ref
	if input.ProdutoID != nil {
		_, err = r.txExec(tx, "UPDATE produtos SET preco_ref = ? WHERE id = ?", input.PrecoUnit, *input.ProdutoID)
		if err != nil {
			return nil, err
		}

		// Mark matching lista_item as comprado
		var listaID *int64
		err = r.txQueryRow(tx, "SELECT lista_id FROM compras WHERE id = ?", compraID).Scan(&listaID)
		if err == nil && listaID != nil {
			_, err = r.txExec(tx, `
				UPDATE lista_itens SET comprado = 1
				WHERE id = (
					SELECT id FROM lista_itens
					WHERE lista_id = ? AND produto_id = ? AND comprado = 0
					LIMIT 1
				)`, *listaID, *input.ProdutoID)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.getCompraItem(itemID)
}

func (r *Repository) UpdateCompraItem(itemID int64, input models.CompraItemInput) (*models.CompraItem, error) {
	res, err := r.Exec(
		`UPDATE compra_itens SET produto_id = ?, nome_nfe = ?, qtd = ?, unidade = ?, preco_unit = ?, preco_total = ?, lista_item_id = ?
		 WHERE id = ?`,
		input.ProdutoID, input.NomeNFe, input.Qtd, input.Unidade, input.PrecoUnit, input.PrecoTotal, input.ListaItemID, itemID,
	)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, models.ErrNotFound("Item da compra")
	}
	return r.getCompraItem(itemID)
}

func (r *Repository) DeleteCompraItem(itemID int64) error {
	res, err := r.Exec("DELETE FROM compra_itens WHERE id = ?", itemID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrNotFound("Item da compra")
	}
	return nil
}

func (r *Repository) getCompraItem(id int64) (*models.CompraItem, error) {
	var ci models.CompraItem
	err := r.QueryRow(`
		SELECT ci.id, ci.compra_id, ci.produto_id, ci.nome_nfe, ci.qtd, ci.unidade, ci.preco_unit, ci.preco_total, ci.lista_item_id,
			COALESCE(p.nome, '')
		FROM compra_itens ci
		LEFT JOIN produtos p ON p.id = ci.produto_id
		WHERE ci.id = ?`, id).
		Scan(&ci.ID, &ci.CompraID, &ci.ProdutoID, &ci.NomeNFe, &ci.Qtd, &ci.Unidade,
			&ci.PrecoUnit, &ci.PrecoTotal, &ci.ListaItemID, &ci.ProdutoNome)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Item da compra")
	}
	if err != nil {
		return nil, err
	}
	return &ci, nil
}
