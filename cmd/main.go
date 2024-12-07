package main

import (
	"cloudservices/pkg/handlers"
	"cloudservices/pkg/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/register", handlers.RegisterUserHandler)
	http.HandleFunc("/auth/validate", handlers.ValidateTokenHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)
	// http.HandleFunc("/list-users", handlers.ListUsersHandler)

	http.Handle("/auth/token/profile", middleware.TokenAuthMiddleware(http.HandlerFunc(handlers.ProfileHandler)))
	//http.Handle("/auth/basicauth/profile", handlers.ProfileHandler)

	log.Println("Authentication service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
