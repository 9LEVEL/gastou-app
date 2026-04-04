package handlers

import (
	"net/http"

	"mercado-app/internal/models"
)

func (h *Handlers) ListCategorias(w http.ResponseWriter, r *http.Request) {
	cats, err := h.Repo.ListCategorias()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *Handlers) GetCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	cat, err := h.Repo.GetCategoria(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

func (h *Handlers) CreateCategoria(w http.ResponseWriter, r *http.Request) {
	var input models.CategoriaInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	cat, err := h.Repo.CreateCategoria(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, cat)
}

func (h *Handlers) UpdateCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.CategoriaInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	cat, err := h.Repo.UpdateCategoria(id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

func (h *Handlers) DeleteCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Repo.DeleteCategoria(id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Categoria removida"})
}
