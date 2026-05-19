package http

import (
	"encoding/json"
	"main/internal/ports"
	"net/http"
)

type Handler struct {
	uc ports.UserUseCase
}

func NewHandler(uc ports.UserUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	user, err := h.uc.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
