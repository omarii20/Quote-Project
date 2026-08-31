package main

import (
	"log"
	"net/http"

	"github.com/omarii20/Quote-Project/internal/business"
	"github.com/omarii20/Quote-Project/internal/config"
	"github.com/omarii20/Quote-Project/internal/customer"
	"github.com/omarii20/Quote-Project/internal/database"
	"github.com/omarii20/Quote-Project/internal/quote"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// Initialize business module
	businessRepo := business.NewRepository(db)
	businessService := business.NewService(businessRepo)
	businessHandler := business.NewHandler(businessService)

	business.RegisterRoutes(businessHandler)

	// Initialize customer module
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	customer.RegisterRoutes(customerHandler)

	// Initialize quote module
	quoteRepo := quote.NewRepository(db)
	quoteService := quote.NewService(quoteRepo)
	quoteHandler := quote.NewHandler(quoteService)

	quote.RegisterRoutes(quoteHandler)

	// Start the HTTP server
	log.Println("server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
