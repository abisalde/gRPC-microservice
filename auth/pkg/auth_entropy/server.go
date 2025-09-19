package auth_entropy

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/abisalde/grpc-microservice/auth/internal/database"
	"github.com/abisalde/grpc-microservice/auth/internal/service"
	"github.com/abisalde/grpc-microservice/auth/pkg/ent/proto/auth_pbuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	sleep = flag.Duration("auth-sleep", time.Second*5, "auth duration between changes in health")

	system = "auth.Service"
)

func ListenGRPC(s *service.UserService, db *database.Database) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))

	if err != nil {
		return err
	}

	svc := auth_pbuf.NewUserService(db.Client)

	server := grpc.NewServer()

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthcheck)

	auth_pbuf.RegisterUserServiceServer(server, svc)

	reflection.Register(server)

	go func() {
		next := healthgrpc.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus(system, next)

			if next == healthgrpc.HealthCheckResponse_SERVING {
				next = healthgrpc.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthgrpc.HealthCheckResponse_SERVING
			}

			time.Sleep(*sleep)
		}
	}()
	return server.Serve(lis)
}
