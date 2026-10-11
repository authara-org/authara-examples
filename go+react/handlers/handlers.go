package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/authara-org/authara-go/authara"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dist/index.html")
}

func Private(w http.ResponseWriter, r *http.Request) {
	userID, ok := authara.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "This response came from the Go backend.",
		"user_id": userID.String(),
	})
}
