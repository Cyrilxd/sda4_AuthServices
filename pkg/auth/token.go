package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Secret key for signing JWT tokens (should be stored securely, e.g., in environment variables)
var jwtSecret = []byte("topSecret")

// Claims struct to store custom and standard JWT claims
type Claims struct {
	Username           string `json:"username"` // Custom claim for the username
	jwt.StandardClaims        // Embeds standard JWT claims
}

// GenerateToken creates a JWT for a user
func GenerateToken(username string) (string, error) {
	// Define the JWT claims
	claims := jwt.MapClaims{
		"username": username,                              // Custom claim for the user's username
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Expiry time set to 24 hours from now
	}

	// Create a new JWT with the specified claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key and return it
	return token.SignedString(jwtSecret)
}

// ValidateToken verifies the JWT token and extracts the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{} // Struct to hold the parsed claims

	// Parse the token string and populate the claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for verifying the token
		return jwtSecret, nil
	})
	if err != nil {
		// Return an error if the token is invalid
		return nil, errors.New("invalid token")
	}

	// Ensure the token is valid and not expired
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Additional check: Ensure the username claim is present
	if claims.Username == "" {
		return nil, errors.New("invalid token: username missing")
	}

	// Return the claims if validation succeeds
	return claims, nil
}
