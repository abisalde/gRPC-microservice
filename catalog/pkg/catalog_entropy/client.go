package catalog_entropy

import (
	"context"
	"time"

	"github.com/abisalde/gprc-microservice/catalog/internal/model"
	"github.com/abisalde/gprc-microservice/catalog/pkg/ent"
	"github.com/abisalde/gprc-microservice/catalog/pkg/ent/proto/entpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Client struct {
	conn    *grpc.ClientConn
	service entpb.CatalogServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{
		service: entpb.NewCatalogServiceClient(conn),
		conn:    conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateProduct(ctx context.Context, p *model.CreateCatalog) (*ent.Catalog, error) {
	r, err := c.service.Create(ctx, &entpb.CreateCatalogRequest{
		Catalog: &entpb.Catalog{
			Name:        p.Name,
			Description: wrapperspb.String(p.Description),
			Price:       p.Price,
		},
	})

	if err != nil {
		return nil, err
	}

	return c.protoToEntCatalog(r), nil
}

func (c *Client) UpdateProduct(ctx context.Context, id string, p *model.CreateCatalog) (*ent.Catalog, error) {

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	r, err := c.service.Update(ctx, &entpb.UpdateCatalogRequest{
		Catalog: &entpb.Catalog{
			Id:          uid[:],
			Name:        p.Name,
			Description: wrapperspb.String(p.Description),
			Price:       p.Price,
		},
	})

	if err != nil {
		return nil, err
	}

	return c.protoToEntCatalog(r), nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*ent.Catalog, error) {

	r, err := c.service.Get(ctx, &entpb.GetCatalogRequest{
		Id: []byte(id),
	})

	if err != nil {
		return nil, err
	}

	return c.protoToEntCatalog(r), nil
}

func (c *Client) protoToEntCatalog(protoCatalog *entpb.Catalog) *ent.Catalog {
	if protoCatalog == nil {
		return &ent.Catalog{}
	}

	return &ent.Catalog{
		ID:   uuid.UUID(protoCatalog.Id),
		Name: protoCatalog.Name,
		Description: func() string {
			if protoCatalog.Description != nil {
				return protoCatalog.Description.Value
			}
			return ""
		}(),
		Price: protoCatalog.Price,
		DeletedAt: func() *time.Time {
			if protoCatalog.DeletedAt != nil {
				t := protoCatalog.DeletedAt.AsTime()
				return &t
			}
			return nil
		}(),
		CreatedAt: protoCatalog.CreatedAt.AsTime(),
		UpdatedAt: protoCatalog.UpdatedAt.AsTime(),
	}
}
