package handlers

import (
	"encoding/json"
	"mercado-app/internal/models"
	"mercado-app/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Repo *repository.Repository
}

func NewHandlers(repo *repository.Repository) *Handlers {
	return &Handlers{Repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*models.AppError); ok {
		writeJSON(w, appErr.Status, appErr)
		return
	}
	writeJSON(w, 500, models.ErrInternal())
}

func parseID(r *http.Request, param string) (int64, error) {
	s := chi.URLParam(r, param)
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, models.ErrValidation("ID inválido")
	}
	return id, nil
}

func decodeBody(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return models.ErrValidation("JSON inválido: " + err.Error())
	}
	return nil
}
