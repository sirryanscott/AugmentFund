package main

import (
	"log"
	"net/http"

	"github.com/AugmentFund/internal/handlers"
	"github.com/AugmentFund/internal/services"
	"github.com/AugmentFund/internal/stores"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// TODO use router
// TODO errors
func main() {
	newDataStore := stores.NewDataStore()
	userHandler := handlers.NewUserHandler(services.NewUserService(newDataStore))
	fundHandler := handlers.NewFundHandler(services.NewFundService(newDataStore))
	r := chi.NewRouter()

	// CORS middleware configuration
	cors := cors.New(cors.Options{
		// Allow any origin to access the resource (you can also specify a specific origin here)
		AllowedOrigins: []string{"*"},
		// Allow methods
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Allow headers
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// Allow credentials
		AllowCredentials: true,
		// Preflight request cache duration
		MaxAge: 300, // MaxAge is 5 minutes
	})

	// Apply the CORS middleware
	r.Use(cors.Handler)

	// User routes
	r.Get("/users", userHandler.GetUsers)
	r.Get("/users/{id}", userHandler.GetUser)
	r.Post("/users", userHandler.CreateUser)

	// Cap table routes
	r.Get("/cap-tables", fundHandler.GetCapTables)
	r.Get("/cap-table/{id}", fundHandler.GetCapTableByID)

	// Fund routes
	r.Post("/fund", fundHandler.CreateFund)
	r.Post("/transfer", fundHandler.CreateTransfer)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
