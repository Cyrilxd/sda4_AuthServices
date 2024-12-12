package handlers

import (
	"cloudservices/pkg/auth"
	"encoding/json"
	"net/http"
)

// ValidateTokenHandler validates a JWT token
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"username": claims.Username})
}
