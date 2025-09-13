package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/abisalde/gprc-microservice/auth/internal/database"
	"github.com/abisalde/gprc-microservice/auth/internal/service"
	"github.com/abisalde/gprc-microservice/auth/pkg/ent/proto/entpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	entpb.UnimplementedUserServiceServer
	service *service.UserService
}

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

	svc := entpb.NewUserService(db.Client)

	server := grpc.NewServer()

	entpb.RegisterUserServiceServer(server, svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}
