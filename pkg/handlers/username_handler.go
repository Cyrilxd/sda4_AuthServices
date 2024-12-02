package handlers

import (
	"cloudservices/pkg/auth"
	"cloudservices/pkg/middleware"
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
	username, err := middleware.RequireUsernameContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Simulated user profile data (replace with actual user data lookup if needed)
	profile := map[string]string{
		"username": username,
		"role":     "user", // Example: you could dynamically retrieve the user's role
	}

	// Respond with the user's profile
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Define a structure to parse the incoming JSON request
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if username and password are provided
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Register the user using the RegisterUser function
	err := auth.RegisterUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": req.Username,
	})
}
