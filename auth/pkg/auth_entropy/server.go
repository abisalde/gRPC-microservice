package auth_entropy

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/abisalde/gprc-microservice/auth/internal/database"
	"github.com/abisalde/gprc-microservice/auth/internal/service"
	"github.com/abisalde/gprc-microservice/auth/pkg/ent/proto/entpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = ""
)

func ListenGRPC(s *service.UserService, db *database.Database) error {

	lis, err := net.Listen("tcp", fmt.Sprint(":%w", 50051))

	if err != nil {
		return err
	}

	svc := entpb.NewUserService(db.Client)

	server := grpc.NewServer()

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthcheck)

	entpb.RegisterUserServiceServer(server, svc)

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
