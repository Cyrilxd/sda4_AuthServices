package handlers

import (
	"cloudservices/pkg/auth"
	"cloudservices/pkg/middleware"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// User struct to store user information
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"` // Optional field for JWT token
}

// LoginRequest represents the structure of login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the structure of login responses
type LoginResponse struct {
	Token string `json:"token"` // JWT token returned upon successful login
}

// LoginHandler authenticates users with username and password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	// Decode the JSON payload into the LoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve user information from storage
	user, exists := auth.GetUser(req.Username)
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the stored password hash with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	log.Println("bcrypt comparison error:", err)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	user.Token = token // Store the token in the user object

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// ProfileHandler returns user profile for authenticated users
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the context (middleware validation)
	username, err := middleware.RequireUsernameContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Simulated profile response for demonstration
	profile := map[string]string{
		"username": username,
		"role":     "user", // Replace with actual role retrieval if needed
	}

	// Respond with the user's profile
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser auth.User
	// Decode the JSON payload into the User struct
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure both username and password are provided
	if newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword) // Update the password to the hashed value

	// Add the user to storage
	if err := auth.AddUser(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User " + newUser.Username + " registered successfully"))
}

// ListUsersHandler retrieves and returns a list of all registered users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all users from storage
	users := auth.ListUsers()

	// Encode the user list into JSON and send it as a response
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
