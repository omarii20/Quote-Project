package main

import (
	"log"

	"github.com/omarii20/Quote-Project/internal/config"
	"github.com/omarii20/Quote-Project/internal/database"
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

	log.Println("connected to PostgreSQL successfully")
}
