package main

import (
	"context"
	"log"
	"time"

	"github.com/abisalde/gprc-microservice/auth/internal/database"
	"github.com/abisalde/gprc-microservice/auth/internal/repository"
	"github.com/abisalde/gprc-microservice/auth/internal/service"
	"github.com/abisalde/gprc-microservice/auth/pkg/auth_entropy"
)

func setupDatabase() (*database.Database, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	if err := db.HealthCheck(ctx); err != nil {
		db.Close()
		return nil, err
	}

	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db, nil
}

func main() {

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("❌ Failed to setup database: %v", err)
	}
	defer db.Close()

	r := repository.NewUserRepository(db.Client)

	s := service.NewUserService(r)
	auth_entropy.ListenGRPC(s, db)
}
