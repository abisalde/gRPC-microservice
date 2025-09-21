package main

import (
	"context"
	"log"
	"time"

	"github.com/abisalde/grpc-microservice/auth/internal/database"
	"github.com/abisalde/grpc-microservice/auth/internal/repository"
	"github.com/abisalde/grpc-microservice/auth/internal/service"
	"github.com/abisalde/grpc-microservice/auth/pkg/auth_entropy"
)

func setupDatabase() (*database.Database, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.HealthCheck(ctx); err != nil {
		db.Close()
		return nil, err
	}

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
	if err := auth_entropy.ListenGRPC(s, db); err != nil {
		log.Fatalf("❌ Failed to start AUTH gRPC server: %v", err)
	}
}
