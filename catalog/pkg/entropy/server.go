package entropy

import (
	"context"
	"fmt"
	"net"

	"github.com/abisalde/gprc-microservice/catalog/internal/model"
	"github.com/abisalde/gprc-microservice/catalog/internal/service"
	"github.com/abisalde/gprc-microservice/catalog/pkg/ent/proto/entpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	entpb.UnimplementedCatalogServiceServer
	service service.ProductCatalogService
}

func ListenGRPC(s service.ProductCatalogService) error {

	lis, err := net.Listen("tcp", fmt.Sprint(":%w", 50052))

	if err != nil {
		return err
	}

	serve := grpc.NewServer()
	entpb.RegisterCatalogServiceServer(serve, &grpcServer{service: s})
	reflection.Register(serve)
	return serve.Serve(lis)
}

func (s *grpcServer) CreateProduct(ctx context.Context, r *entpb.CreateCatalogRequest) (*entpb.Catalog, error) {
	catalog := &model.CreateCatalog{
		Name:        r.GetCatalog().GetName(),
		Description: r.GetCatalog().Description.GetValue(),
		Price:       r.GetCatalog().GetPrice(),
	}

	c, err := s.service.CreateProduct(ctx, *catalog)

	if err != nil {
		return nil, err
	}
	return &entpb.Catalog{
		Id:        c.ID[:],
		Name:      c.Name,
		Price:     c.Price,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
		DeletedAt: nil,
	}, nil

}

func (s *grpcServer) UpdateProduct(ctx context.Context, r *entpb.UpdateCatalogRequest) (*entpb.Catalog, error) {

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

	return &entpb.Catalog{
		Id:        c.ID[:],
		Name:      c.Name,
		Price:     c.Price,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
		DeletedAt: nil,
	}, nil
}
