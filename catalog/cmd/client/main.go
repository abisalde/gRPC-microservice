package main

import (
	"context"

	"github.com/abisalde/gprc-microservice/catalog/pkg/ent/proto/entpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CatalogClient struct {
	conn          *grpc.ClientConn
	catalogClient entpb.CatalogServiceClient
}

func NewCatalogClient(addr string) (*CatalogClient, error) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &CatalogClient{
		catalogClient: entpb.NewCatalogServiceClient(conn),
		conn:          conn,
	}, nil
}

func (c *CatalogClient) Close() error {
	return c.conn.Close()
}

func (c *CatalogClient) GetCatalogByID(ctx context.Context, id int64) (*entpb.Catalog, error) {
	return c.catalogClient.Get(ctx, &entpb.GetCatalogRequest{Id: id})
}

func (c *CatalogClient) CreateCatalog(ctx context.Context, catalog *entpb.Catalog) (*entpb.Catalog, error) {
	return c.catalogClient.Create(ctx, &entpb.CreateCatalogRequest{Catalog: catalog})
}
