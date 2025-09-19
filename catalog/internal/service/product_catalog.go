package service

import (
	"context"
	"errors"
	"time"

	"github.com/abisalde/grpc-microservice/catalog/internal/model"
	"github.com/abisalde/grpc-microservice/catalog/internal/repository"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent"
	"github.com/google/uuid"
)

type ProductCatalogService interface {
	CreateProduct(ctx context.Context, p model.CreateCatalog) (*ent.Catalog, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, p *model.CreateCatalog) (*ent.Catalog, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*ent.Catalog, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]*ent.Catalog, error)
	GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Catalog, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*ent.Catalog, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type catalogService struct {
	repository repository.CatalogRepository
}

func NewProductCatalogService(r repository.CatalogRepository) ProductCatalogService {
	return &catalogService{repository: r}
}

func (s *catalogService) CreateProduct(ctx context.Context, p model.CreateCatalog) (*ent.Catalog, error) {
	catalog := &ent.Catalog{
		ID:          uuid.New(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if err := s.repository.PostProductCatalog(ctx, "Create", *catalog); err != nil {
		return &ent.Catalog{}, err
	}

	return catalog, nil
}

func (s *catalogService) UpdateProduct(ctx context.Context, id uuid.UUID, p *model.CreateCatalog) (*ent.Catalog, error) {

	existing, err := s.repository.GetProductCatalogByID(ctx, id)

	if err != nil {
		return &ent.Catalog{}, err
	}

	catalog := &ent.Catalog{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	if p.Name != "" {
		catalog.Name = p.Name
	}
	if p.Description != "" {
		catalog.Description = p.Description
	}
	if p.Price != 0 {
		catalog.Price = p.Price
	}

	if err := s.repository.PostProductCatalog(ctx, "Update", *catalog); err != nil {
		return &ent.Catalog{}, err
	}

	return catalog, nil
}

func (s *catalogService) GetProductByID(ctx context.Context, id uuid.UUID) (*ent.Catalog, error) {
	if id == uuid.Nil {
		return nil, errors.New("product ID cannot be empty")
	}

	product, err := s.repository.GetProductCatalogByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]*ent.Catalog, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	products, err := s.repository.GetAllProductsCatalog(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Catalog, error) {
	if len(ids) == 0 {
		return []*ent.Catalog{}, nil
	}

	if len(ids) > 100 {
		return nil, errors.New("too many IDs requested (max 100)")
	}

	products, err := s.repository.GetProductCatalogWithIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*ent.Catalog, error) {
	if query == "" {
		return s.GetProducts(ctx, skip, take)
	}

	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	products, err := s.repository.SearchProducts(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *catalogService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}

	return s.repository.DeleteProductCatalog(ctx, id)
}

func ptrTime(t time.Time) *time.Time { return &t }
