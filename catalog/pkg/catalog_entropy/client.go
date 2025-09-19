package catalog_entropy

import (
	"context"
	"time"

	"github.com/abisalde/grpc-microservice/catalog/internal/model"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent/proto/catalog_pbuf"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Client struct {
	conn    *grpc.ClientConn
	service catalog_pbuf.CatalogServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{
		service: catalog_pbuf.NewCatalogServiceClient(conn),
		conn:    conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateProduct(ctx context.Context, p *model.CreateCatalog) (*ent.Catalog, error) {
	r, err := c.service.Create(ctx, &catalog_pbuf.CreateCatalogRequest{
		Catalog: &catalog_pbuf.Catalog{
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

	r, err := c.service.Update(ctx, &catalog_pbuf.UpdateCatalogRequest{
		Catalog: &catalog_pbuf.Catalog{
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

	r, err := c.service.Get(ctx, &catalog_pbuf.GetCatalogRequest{
		Id: []byte(id),
	})

	if err != nil {
		return nil, err
	}

	return c.protoToEntCatalog(r), nil
}

func (c *Client) protoToEntCatalog(protoCatalog *catalog_pbuf.Catalog) *ent.Catalog {
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
