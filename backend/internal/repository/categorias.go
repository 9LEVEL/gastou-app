package repository

import (
	"database/sql"
	"mercado-app/internal/models"
)

func (r *Repository) ListCategorias() ([]models.Categoria, error) {
	rows, err := r.Query("SELECT id, nome, cor, icone, ordem FROM categorias ORDER BY ordem, nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Categoria
	for rows.Next() {
		var c models.Categoria
		if err := rows.Scan(&c.ID, &c.Nome, &c.Cor, &c.Icone, &c.Ordem); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []models.Categoria{}
	}
	return cats, rows.Err()
}

func (r *Repository) GetCategoria(id int64) (*models.Categoria, error) {
	var c models.Categoria
	err := r.QueryRow("SELECT id, nome, cor, icone, ordem FROM categorias WHERE id = ?", id).
		Scan(&c.ID, &c.Nome, &c.Cor, &c.Icone, &c.Ordem)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound("Categoria")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) CreateCategoria(input models.CategoriaInput) (*models.Categoria, error) {
	if input.Nome == "" {
		return nil, models.ErrValidation("Nome é obrigatório")
	}
	if input.Cor == "" {
		input.Cor = "#636E72"
	}

	id, err := r.InsertReturningID(
		"INSERT INTO categorias (nome, cor, icone, ordem) VALUES (?, ?, ?, ?)",
		input.Nome, input.Cor, input.Icone, input.Ordem,
	)
	if err != nil {
		return nil, err
	}
	return r.GetCategoria(id)
}

func (r *Repository) UpdateCategoria(id int64, input models.CategoriaInput) (*models.Categoria, error) {
	if input.Nome == "" {
		return nil, models.ErrValidation("Nome é obrigatório")
	}
	if input.Cor == "" {
		input.Cor = "#636E72"
	}

	res, err := r.Exec(
		"UPDATE categorias SET nome = ?, cor = ?, icone = ?, ordem = ? WHERE id = ?",
		input.Nome, input.Cor, input.Icone, input.Ordem, id,
	)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, models.ErrNotFound("Categoria")
	}
	return r.GetCategoria(id)
}

func (r *Repository) DeleteCategoria(id int64) error {
	var count int
	err := r.QueryRow("SELECT COUNT(*) FROM produtos WHERE categoria_id = ?", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return models.ErrConflict("Categoria possui produtos vinculados")
	}

	res, err := r.Exec("DELETE FROM categorias WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.ErrNotFound("Categoria")
	}
	return nil
}
