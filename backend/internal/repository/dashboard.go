package repository

import (
	"mercado-app/internal/models"
)

func (r *Repository) GetResumo(listaID int64) (*models.DashboardResumo, error) {
	var d models.DashboardResumo

	err := r.QueryRow(`
		SELECT
			COALESCE(SUM(li.qtd * li.preco_estimado / li.duracao_meses), 0),
			COUNT(li.id),
			COALESCE(SUM(li.comprado), 0)
		FROM lista_itens li
		WHERE li.lista_id = ?`, listaID).
		Scan(&d.TotalPlanejado, &d.ItensTotal, &d.ItensComprados)
	if err != nil {
		return nil, err
	}

	err = r.QueryRow(`
		SELECT COALESCE(SUM(ci.preco_total), 0)
		FROM compra_itens ci
		JOIN compras c ON c.id = ci.compra_id
		WHERE c.lista_id = ?`, listaID).
		Scan(&d.TotalGasto)
	if err != nil {
		return nil, err
	}

	err = r.QueryRow("SELECT renda FROM listas WHERE id = ?", listaID).Scan(&d.Renda)
	if err != nil {
		return nil, err
	}

	d.SaldoRestante = d.TotalPlanejado - d.TotalGasto
	d.ItensFaltantes = d.ItensTotal - d.ItensComprados
	if d.Renda > 0 {
		d.PercentualRenda = (d.TotalPlanejado / d.Renda) * 100
	}

	return &d, nil
}

func (r *Repository) GetComparativo(listaID int64) ([]models.CategoriaComparativo, error) {
	rows, err := r.Query(`
		SELECT c.id, c.nome, c.cor,
			COALESCE(SUM(li.qtd * li.preco_estimado / li.duracao_meses), 0) AS planejado,
			COALESCE((
				SELECT SUM(ci.preco_total)
				FROM compra_itens ci
				JOIN compras co ON co.id = ci.compra_id
				WHERE co.lista_id = ? AND ci.produto_id IN (
					SELECT li2.produto_id FROM lista_itens li2
					JOIN produtos p2 ON p2.id = li2.produto_id
					WHERE li2.lista_id = ? AND p2.categoria_id = c.id
				)
			), 0) AS real_val
		FROM categorias c
		LEFT JOIN produtos p ON p.categoria_id = c.id
		LEFT JOIN lista_itens li ON li.produto_id = p.id AND li.lista_id = ?
		GROUP BY c.id
		HAVING planejado > 0 OR real_val > 0
		ORDER BY c.ordem`, listaID, listaID, listaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CategoriaComparativo
	for rows.Next() {
		var cc models.CategoriaComparativo
		if err := rows.Scan(&cc.CategoriaID, &cc.CategoriaNome, &cc.CategoriaCor, &cc.Planejado, &cc.Real); err != nil {
			return nil, err
		}
		result = append(result, cc)
	}
	if result == nil {
		result = []models.CategoriaComparativo{}
	}
	return result, rows.Err()
}

func (r *Repository) GetEvolucao() ([]models.EvolucaoMensal, error) {
	rows, err := r.Query(`
		SELECT l.mes, l.ano, l.nome, COALESCE(SUM(ci.preco_total), 0) AS gasto
		FROM listas l
		LEFT JOIN compras c ON c.lista_id = l.id
		LEFT JOIN compra_itens ci ON ci.compra_id = c.id
		GROUP BY l.id
		ORDER BY l.ano, l.mes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.EvolucaoMensal
	for rows.Next() {
		var e models.EvolucaoMensal
		if err := rows.Scan(&e.Mes, &e.Ano, &e.Nome, &e.Gasto); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	if result == nil {
		result = []models.EvolucaoMensal{}
	}
	return result, rows.Err()
}
