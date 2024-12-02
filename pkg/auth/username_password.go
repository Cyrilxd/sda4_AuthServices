package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]string{} // Simulated database

// RegisterUser stores a hashed password for a username
func RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users[username] = string(hashedPassword)

	return nil
}

// Authenticate checks the username and password
func Authenticate(username, password string) error {
	storedPassword, exists := users[username]
	if !exists {
		return errors.New("user not found")
	}
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
}
