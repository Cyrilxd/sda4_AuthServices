package main

import (
	"cloudservices/pkg/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "github.com/swaggo/http-swagger"
	_ "sda4_AuthServices/docs" // Import generated Swagger docs
	"sda4_AuthServices/pkg/handlers"
)

// @title Auth Services API
// @version 1.0
// @description This is a proof-of-concept for authentication services.
// @termsOfService tbd

// @contact.name API Support
// @contact.email cyril.heiniger@students.bfh.ch

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /auth
func main() {
	http.HandleFunc("/auth/register", handlers.RegisterUserHandler)
	http.HandleFunc("/auth/token", handlers.GenerateTokenHandler)
	http.HandleFunc("/auth/validate", handlers.ValidateTokenHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.HandleFunc("/auth/profile", handlers.ProfileHandler)

	// Swagger UI endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	log.Println("Authentication service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
