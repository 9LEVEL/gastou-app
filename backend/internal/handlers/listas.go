package handlers

import (
	"mercado-app/internal/models"
	"net/http"
	"strconv"
)

func (h *Handlers) ListListas(w http.ResponseWriter, r *http.Request) {
	listas, err := h.Repo.ListListas()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, listas)
}

func (h *Handlers) GetLista(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	lista, err := h.Repo.GetLista(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, lista)
}

func (h *Handlers) CreateLista(w http.ResponseWriter, r *http.Request) {
	var input models.ListaInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}

	var lista *models.Lista
	var err error

	if copiarDe := r.URL.Query().Get("copiar_de"); copiarDe != "" {
		srcID, parseErr := strconv.ParseInt(copiarDe, 10, 64)
		if parseErr != nil {
			writeError(w, models.ErrValidation("copiar_de inválido"))
			return
		}
		lista, err = h.Repo.CopiarLista(input, srcID)
	} else {
		lista, err = h.Repo.CreateLista(input)
	}

	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, lista)
}

func (h *Handlers) UpdateLista(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.ListaUpdateInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	lista, err := h.Repo.UpdateLista(id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, lista)
}

func (h *Handlers) DeleteLista(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Repo.DeleteLista(id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Lista removida"})
}

// ---------- Itens ----------

func (h *Handlers) ListItens(w http.ResponseWriter, r *http.Request) {
	listaID, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	itens, err := h.Repo.ListItens(listaID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, itens)
}

func (h *Handlers) AddItem(w http.ResponseWriter, r *http.Request) {
	listaID, err := parseID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.ListaItemInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Repo.AddItem(listaID, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseID(r, "itemId")
	if err != nil {
		writeError(w, err)
		return
	}
	var input models.ListaItemUpdateInput
	if err := decodeBody(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Repo.UpdateItem(itemID, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handlers) ToggleCheck(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseID(r, "itemId")
	if err != nil {
		writeError(w, err)
		return
	}
	item, err := h.Repo.ToggleCheck(itemID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseID(r, "itemId")
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.Repo.DeleteItem(itemID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Item removido"})
}
