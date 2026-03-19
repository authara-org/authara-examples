package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/authara-org/authara-go/authara"
)

type Handler struct {
	client *authara.Client
}

func New(autharaBaseURL string) *Handler {
	return &Handler{
		client: authara.NewClient(autharaBaseURL),
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dist/index.html")
}

func (h *Handler) Private(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dist/index.html")
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.GetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"email":         user.Email,
		"username":      user.Username,
	})
}
