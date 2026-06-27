package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/samkhedikar/tou-api/internal/db"
	"github.com/samkhedikar/tou-api/internal/handler"
	"github.com/samkhedikar/tou-api/internal/repository"
	"github.com/samkhedikar/tou-api/internal/service"
)

func main() {
	database, err := db.New()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	// Repositories — DB layer
	chargerRepo := repository.NewChargerRepo(database)
	pricingRepo := repository.NewPricingRepo(database)

	// Services — business logic layer
	chargerSvc := service.NewChargerService(chargerRepo)
	pricingSvc := service.NewPricingService(chargerRepo, pricingRepo)

	// Handlers — HTTP layer
	h := handler.New(chargerSvc, pricingSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      h.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("tou-api listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
