package main

import (
	"cloudservices/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/token", handlers.GenerateTokenHandler)
	http.HandleFunc("/auth/validate", handlers.ValidateTokenHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.HandleFunc("/auth/profile", handlers.ProfileHandler)

	log.Println("Authentication service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
