package handlers

import (
	"mercado-app/internal/models"
	"net/http"
	"strconv"
)

func (h *Handlers) ListCompras(w http.ResponseWriter, r *http.Request) {
	var listaID *int64

	if v := r.URL.Query().Get("lista_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			listaID = &id
		}
	}

	compras, err := h.Repo.ListCompras(listaID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, compras)
}

func (h *Handlers) GetCompra(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	compra, err := h.Repo.GetCompra(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, compra)
}

func (h *Handlers) CreateCompra(w http.ResponseWriter, r *http.Request) {
	var input models.CompraInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	compra, err := h.Repo.CreateCompra(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, compra)
}

func (h *Handlers) UpdateCompra(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.CompraInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	compra, err := h.Repo.UpdateCompra(id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, compra)
}

func (h *Handlers) DeleteCompra(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Repo.DeleteCompra(id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Compra removida"})
}

func (h *Handlers) AddCompraItem(w http.ResponseWriter, r *http.Request) {
	compraID, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.CompraItemInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Repo.AddCompraItem(compraID, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handlers) UpdateCompraItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseID(r, "itemId")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.CompraItemInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Repo.UpdateCompraItem(itemID, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handlers) DeleteCompraItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseID(r, "itemId")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Repo.DeleteCompraItem(itemID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Item da compra removido"})
}
