package server

import (
	"context"
	"time"

	"github.com/abisalde/gprc-microservice/auth/internal/database"
)

func SetupDatabase() (*database.Database, error) {
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
