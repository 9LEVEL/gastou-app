package handlers

import (
	"mercado-app/internal/models"
	"net/http"
	"strconv"
)

func (h *Handlers) ListProdutos(w http.ResponseWriter, r *http.Request) {
	var categoriaID *int64
	var ativo *bool

	if v := r.URL.Query().Get("categoria_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			categoriaID = &id
		}
	}
	if v := r.URL.Query().Get("ativo"); v != "" {
		b := v == "true" || v == "1"
		ativo = &b
	}

	prods, err := h.Repo.ListProdutos(categoriaID, ativo)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, prods)
}

func (h *Handlers) GetProduto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	prod, err := h.Repo.GetProduto(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, prod)
}

func (h *Handlers) CreateProduto(w http.ResponseWriter, r *http.Request) {
	var input models.ProdutoInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	prod, err := h.Repo.CreateProduto(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, prod)
}

func (h *Handlers) UpdateProduto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.ProdutoInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	prod, err := h.Repo.UpdateProduto(id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, prod)
}

func (h *Handlers) GetHistoricoPrecos(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	hist, err := h.Repo.GetHistoricoPrecos(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, hist)
}
