package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret = []byte("topSecret") // Use environment variables for production

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken erstellt ein JWT für den Benutzer
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token läuft nach 24 Stunden ab
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Ensure the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Additional checks: Ensure required claims are present
	if claims.Username == "" {
		return nil, errors.New("invalid token: username missing")
	}

	return claims, nil
}
