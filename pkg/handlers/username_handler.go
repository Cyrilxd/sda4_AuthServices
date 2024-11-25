package handlers

import (
	"cloudservices/pkg/auth"
	"encoding/json"
	"net/http"
)

// LoginHandler authenticates users with username and password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := auth.Authenticate(req.Username, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// ProfileHandler returns user profile for authenticated users
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Simulated profile data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"username": "test_user", "role": "admin"})
}