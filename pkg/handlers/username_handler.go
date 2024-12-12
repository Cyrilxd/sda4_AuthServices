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
	Token    string `json:"token,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler authenticates users with username and password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Println("username:", req.Username)
	log.Println("password:", req.Password)
	// Retrieve user from storage
	user, exists := auth.GetUser(req.Username)
	log.Println("user:", user)
	log.Println("exists:", exists)
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Println("Password in request:", req.Password)
	log.Println("Stored hashed password:", user.Password)

	log.Println("Password in request:", req.Password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	log.Println("Hashed password during registration:", string(hashedPassword))

	// Validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	log.Println("bcrypt comparison error:", err)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	user.Token = token // Store the token in the user object

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser auth.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Be sure that username and password are set
	if newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Hashed password during registration:", string(hashedPassword))

	// save hash in user struct
	newUser.Password = string(hashedPassword)

	// Add the user to storage
	if err := auth.AddUser(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	users := auth.ListUsers()

	// Konvertiere die Liste der Benutzer in JSON
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	log.Println("Stored User Password: ", string(response))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User " + newUser.Username + " registered successfully"))
}

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Hole die Liste der Benutzer
	users := auth.ListUsers()

	// Konvertiere die Liste der Benutzer in JSON
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}

	// Setze Header und sende die Antwort
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
