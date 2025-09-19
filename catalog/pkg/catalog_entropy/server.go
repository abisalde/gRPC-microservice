package catalog_entropy

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/abisalde/grpc-microservice/catalog/internal/model"
	"github.com/abisalde/grpc-microservice/catalog/internal/service"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent/proto/catalog_pbuf"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type grpcServer struct {
	catalog_pbuf.UnimplementedCatalogServiceServer
	catalog_pbuf.UnimplementedCatalogExtendedServiceServer
	service service.ProductCatalogService
}

var (
	sleep = flag.Duration("catalog-sleep", time.Second*5, "catalog duration between changes in health")

	system = "catalog.CatalogService"
)

func ListenGRPC(s service.ProductCatalogService) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))

	if err != nil {
		return err
	}

	serve := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(serve, healthcheck)

	catalog_pbuf.RegisterCatalogServiceServer(serve, &grpcServer{service: s})

	catalog_pbuf.RegisterCatalogExtendedServiceServer(serve, &grpcServer{service: s})

	reflection.Register(serve)

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
	return serve.Serve(lis)
}

func (s *grpcServer) Create(ctx context.Context, r *catalog_pbuf.CreateCatalogRequest) (*catalog_pbuf.Catalog, error) {
	catalog := &model.CreateCatalog{
		Name:        r.GetCatalog().GetName(),
		Description: r.GetCatalog().Description.GetValue(),
		Price:       r.GetCatalog().GetPrice(),
	}

	c, err := s.service.CreateProduct(ctx, *catalog)

	if err != nil {
		return nil, err
	}
	return entCatalogToProto(c), nil

}

func (s *grpcServer) Update(ctx context.Context, r *catalog_pbuf.UpdateCatalogRequest) (*catalog_pbuf.Catalog, error) {

	idUUID, err := uuid.Parse(string(r.GetCatalog().GetId()))
	if err != nil {
		return nil, err
	}

	catalog := &model.CreateCatalog{
		Name:        r.GetCatalog().GetName(),
		Description: r.GetCatalog().Description.GetValue(),
		Price:       r.GetCatalog().GetPrice(),
	}

	c, err := s.service.UpdateProduct(ctx, idUUID, catalog)

	if err != nil {
		return nil, err
	}

	return entCatalogToProto(c), nil
}

func (s *grpcServer) Get(ctx context.Context, r *catalog_pbuf.GetCatalogRequest) (*catalog_pbuf.Catalog, error) {
	idUUID, err := uuid.Parse(string(r.GetId()))
	if err != nil {
		return nil, err
	}

	p, err := s.service.GetProductByID(ctx, idUUID)
	if err != nil {
		return nil, err
	}

	return entCatalogToProto(p), nil
}

func (s *grpcServer) BatchCreate(ctx context.Context, r *catalog_pbuf.BatchCreateCatalogsRequest) (*catalog_pbuf.BatchCreateCatalogsResponse, error) {
	var createRequests []model.CreateCatalog
	for _, req := range r.GetRequests() {
		createCatalog := model.CreateCatalog{
			Name:        req.GetCatalog().GetName(),
			Description: req.GetCatalog().GetDescription().Value,
			Price:       req.GetCatalog().GetPrice(),
		}
		createRequests = append(createRequests, createCatalog)
	}

	var createdProducts []*ent.Catalog
	for _, createReq := range createRequests {
		product, err := s.service.CreateProduct(ctx, createReq)
		if err != nil {
			return nil, err
		}
		createdProducts = append(createdProducts, product)
	}

	createdCatalogs := make([]*catalog_pbuf.Catalog, len(createdProducts))
	for i, product := range createdProducts {
		createdCatalogs[i] = entCatalogToProto(product)
	}

	return &catalog_pbuf.BatchCreateCatalogsResponse{
		Catalogs: createdCatalogs,
	}, nil
}

func entCatalogToProto(catalog *ent.Catalog) *catalog_pbuf.Catalog {
	if catalog == nil {
		return nil
	}

	protoCatalog := &catalog_pbuf.Catalog{
		Id:        catalog.ID[:],
		Name:      catalog.Name,
		Price:     catalog.Price,
		CreatedAt: timestamppb.New(catalog.CreatedAt),
		UpdatedAt: timestamppb.New(catalog.UpdatedAt),
	}
	if catalog.Description != "" {
		protoCatalog.Description = wrapperspb.String(catalog.Description)
	}

	return protoCatalog
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(ids))
	for i, idStr := range ids {
		idUUID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID format for id %s: %v", idStr, err)
		}
		uuids[i] = idUUID
	}
	return uuids, nil
}
