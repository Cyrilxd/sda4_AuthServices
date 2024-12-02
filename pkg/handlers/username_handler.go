package handlers

import (
	"cloudservices/pkg/auth"
	"encoding/json"
	"net/http"
)

// LoginHandler authenticates users with username and password
// @Summary Authenticate a user
// @Description Authenticates a user using their username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body struct{Username string `json:"username"`; Password string `json:"password"`} true "User credentials"
// @Success 200 {object} map[string]string "Login successful message"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
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
// @Summary Get user profile
// @Description Retrieves profile information for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "User profile data"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Router /user/profile [get]
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Simulated profile data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"username": "test_user", "role": "admin"})
}

// RegisterUserHandler handles user registration requests
// @Summary Register a new user
// @Description Creates a new user in the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body struct{Username string `json:"username"`; Password string `json:"password"`} true "User credentials"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Failed to register user"
// @Router /auth/register [post]
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
