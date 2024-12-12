package auth

import (
	"errors"
	"log"
	"sync"
)

var (
	// ErrUserAlreadyExists is returned when attempting to add a duplicate user
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"`
}

// Global in-memory user store and mutex
var (
	users = make(map[string]*User)
	mu    sync.Mutex
)

// GetUser retrieves a user by username
func GetUser(username string) (*User, bool) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[username]
	return user, exists
}

// AddUser adds a new user to the global store
func AddUser(user *User) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	users[user.Username] = user
	log.Printf("User added: %+v\n", user)

	return nil
}

// ListUsers returns a copy of all users (for debugging or inspection)
func ListUsers() []*User {
	mu.Lock()
	defer mu.Unlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList
}
