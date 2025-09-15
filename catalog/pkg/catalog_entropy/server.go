package catalog_entropy

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/abisalde/gprc-microservice/catalog/internal/model"
	"github.com/abisalde/gprc-microservice/catalog/internal/service"
	"github.com/abisalde/gprc-microservice/catalog/pkg/ent"
	"github.com/abisalde/gprc-microservice/catalog/pkg/ent/proto/entpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type grpcServer struct {
	entpb.UnimplementedCatalogServiceServer
	entpb.UnimplementedCatalogExtendedServiceServer
	service service.ProductCatalogService
}

var (
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = ""
)

func ListenGRPC(s service.ProductCatalogService) error {

	lis, err := net.Listen("tcp", fmt.Sprint(":%w", 50052))

	if err != nil {
		return err
	}

	serve := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(serve, healthcheck)

	entpb.RegisterCatalogServiceServer(serve, &grpcServer{service: s})

	entpb.RegisterCatalogExtendedServiceServer(serve, &grpcServer{service: s})

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

func (s *grpcServer) Create(ctx context.Context, r *entpb.CreateCatalogRequest) (*entpb.Catalog, error) {
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

func (s *grpcServer) Update(ctx context.Context, r *entpb.UpdateCatalogRequest) (*entpb.Catalog, error) {

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

func (s *grpcServer) GetProduct(ctx context.Context, r *entpb.GetCatalogRequest) (*entpb.Catalog, error) {
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

func (s *grpcServer) GetProducts(ctx context.Context, r *entpb.GetProductsRequest) (*entpb.GetProductsResponse, error) {

	skip := uint64(0)
	if r.GetPageToken() != "" {
		if parsedSkip, err := strconv.ParseUint(r.GetPageToken(), 10, 64); err == nil {
			skip = parsedSkip
		}
	}

	take := uint64(r.GetPageSize())
	if take == 0 {
		take = 20
	}

	var products []*ent.Catalog
	var err error

	switch {
	case len(r.GetIds()) > 0:
		uuids, parseErr := parseUUIDs(r.GetIds())
		if parseErr != nil {
			return nil, parseErr
		}
		products, err = s.service.GetProductsByIDs(ctx, uuids)
	case r.GetQuery() != nil && r.GetQuery().Value != "":
		products, err = s.service.SearchProducts(ctx, r.GetQuery().Value, skip, take)
	default:
		products, err = s.service.GetProducts(ctx, skip, take)
	}

	if err != nil {
		return nil, err
	}

	catalogList := make([]*entpb.Catalog, len(products))
	for i, product := range products {
		catalogList[i] = entCatalogToProto(product)
	}

	nextPageToken := ""
	if len(r.GetIds()) == 0 && len(products) == int(take) {
		nextPageToken = strconv.FormatUint(uint64(skip+take), 10)
	}

	return &entpb.GetProductsResponse{
		Products:      catalogList,
		NextPageToken: nextPageToken,
		TotalCount:    int64(len(products)),
	}, nil

}

func (s *grpcServer) BatchCreate(ctx context.Context, r *entpb.BatchCreateCatalogsRequest) (*entpb.BatchCreateCatalogsResponse, error) {
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

	createdCatalogs := make([]*entpb.Catalog, len(createdProducts))
	for i, product := range createdProducts {
		createdCatalogs[i] = entCatalogToProto(product)
	}

	return &entpb.BatchCreateCatalogsResponse{
		Catalogs: createdCatalogs,
	}, nil
}

func entCatalogToProto(catalog *ent.Catalog) *entpb.Catalog {
	if catalog == nil {
		return nil
	}

	protoCatalog := &entpb.Catalog{
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
