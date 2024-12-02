package handlers

import (
	"cloudservices/pkg/auth"
	"encoding/json"
	"net/http"
)

// TokenRequest represents the payload for token generation
type TokenRequest struct {
	Username string `json:"username"`
}

// GenerateTokenHandler handles token generation
// @Summary Generate a JWT token
// @Description Generates a JWT token for a given username
// @Tags Token
// @Accept json
// @Produce json
// @Param request body TokenRequest true "Username for token generation"
// @Success 200 {object} map[string]string "Generated token"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Failed to generate token"
// @Router /auth/token [post]
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ValidateTokenHandler validates a JWT token
// @Summary Validate a JWT token
// @Description Validates a JWT token and extracts claims
// @Tags Token
// @Accept json
// @Produce json
// @Param token query string true "JWT token to validate"
// @Success 200 {object} map[string]string "Extracted username from token"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Router /auth/validate [get]
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
