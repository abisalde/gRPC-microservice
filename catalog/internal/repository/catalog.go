package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/abisalde/grpc-microservice/catalog/internal/es"
	"github.com/abisalde/grpc-microservice/catalog/pkg/ent"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/google/uuid"
)

type CatalogRepository interface {
	Close()
	PostProductCatalog(ctx context.Context, method string, p ent.Catalog) error
	GetProductCatalogByID(ctx context.Context, id uuid.UUID) (*ent.Catalog, error)
	GetAllProductsCatalog(ctx context.Context, skip uint64, take uint64) ([]*ent.Catalog, error)
	GetProductCatalogWithIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Catalog, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*ent.Catalog, error)
	DeleteProductCatalog(ctx context.Context, id uuid.UUID) error
}

type elasticRepository struct {
	client *es.ElasticClient
}

func NewCatalogRepository(esClient *es.ElasticClient) CatalogRepository {
	return &elasticRepository{client: esClient}
}

func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PostProductCatalog(ctx context.Context, method string, p ent.Catalog) error {
	doc := map[string]interface{}{
		"id":          p.ID.String(),
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
	}
	_, err := r.client.Client.Index(r.client.Index).
		Id(p.ID.String()).
		Document(doc).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to create catalog item: %w", err)
	}

	log.Printf("✅ %sd catalog item %s in Elasticsearch", method, p.ID)
	return nil
}

func (r *elasticRepository) GetProductCatalogByID(ctx context.Context, id uuid.UUID) (*ent.Catalog, error) {
	resp, err := r.client.Client.Get(r.client.Index, id.String()).Do(ctx)

	if err != nil {
		return &ent.Catalog{}, fmt.Errorf("failed to get catalog item: %w", err)
	}

	if !resp.Found {
		return nil, &ent.NotFoundError{}
	}

	p := ent.Catalog{}
	if err := json.Unmarshal(resp.Source_, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal catalog item: %w", err)
	}

	return &ent.Catalog{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}, nil
}

func (r *elasticRepository) GetAllProductsCatalog(ctx context.Context, skip uint64, take uint64) ([]*ent.Catalog, error) {

	resp, err := r.client.Client.Search().
		Index(r.client.Index).
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		From(int(skip)).
		Size(int(take)).Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*ent.Catalog
	for _, hit := range resp.Hits.Hits {
		var source map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &source); err != nil {
			log.Printf("Warning: Failed to unmarshal catalog item %v: %v", hit.Id_, err)
			continue
		}
		products = append(products, r.unmarshalCatalog(source))
	}

	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*ent.Catalog, error) {

	multiMatchQuery := &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query:  query,
			Fields: []string{"name", "description"},
		},
	}

	resp, err := r.client.Client.Search().
		Index(r.client.Index).
		Query(multiMatchQuery).
		From(int(skip)).
		Size(int(take)).Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*ent.Catalog
	for _, hit := range resp.Hits.Hits {
		var source map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &source); err != nil {
			log.Printf("Warning: Failed to unmarshal catalog item %d: %v", hit.Id_, err)
			continue
		}
		products = append(products, r.unmarshalCatalog(source))
	}

	return products, nil
}

func (r *elasticRepository) GetProductCatalogWithIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Catalog, error) {
	if len(ids) == 0 {
		return []*ent.Catalog{}, nil
	}

	docs := make([]types.MgetOperationVariant, 0, len(ids))
	for _, id := range ids {
		docs = append(docs, &types.MgetOperation{
			Id_:    id.String(),
			Index_: &r.client.Index,
		})
	}

	resp, err := r.client.Client.Mget().
		Index(r.client.Index).
		Docs(docs...).
		Do(ctx)

	if err != nil {
		log.Println("MultiGet error:", err)
		return nil, err
	}

	var catalogs []*ent.Catalog
	for i, doc := range resp.Docs {
		switch item := doc.(type) {
		case *types.GetResult:
			if item.Found && item.Source_ != nil {
				var source map[string]interface{}
				if err := json.Unmarshal(item.Source_, &source); err != nil {
					log.Printf("Warning: failed to unmarshal document %s: %v", item.Id_, err)
					continue
				}

				catalog := r.unmarshalCatalog(source)
				catalogs = append(catalogs, catalog)
			} else if !item.Found {
				log.Printf("Warning: document %s not found", ids[i])
			}
		case *types.MultiGetError:
			log.Printf("Warning: error getting document %s: %v", ids[i], item.Error.Reason)
		default:
			log.Printf("Warning: unknown response type for document %s: %T", ids[i], doc)
		}
	}

	return catalogs, nil
}

func (r *elasticRepository) DeleteProductCatalog(ctx context.Context, id uuid.UUID) error {
	_, err := r.client.Client.Delete(r.client.Index, id.String()).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete catalog item: %w", err)
	}

	log.Printf("✅ Deleted catalog item %s from Elasticsearch", id)
	return nil
}

func (r *elasticRepository) unmarshalCatalog(source map[string]interface{}) *ent.Catalog {
	catalog := &ent.Catalog{}
	if id, ok := source["id"].(uuid.UUID); ok {
		catalog.ID = id
	}

	if name, ok := source["name"].(string); ok {
		catalog.Name = name
	}

	if description, ok := source["description"].(string); ok {
		catalog.Description = description
	}

	if price, ok := source["price"].(float64); ok {
		catalog.Price = price
	}

	if createdAt, ok := source["created_at"].(time.Time); ok {
		catalog.CreatedAt = createdAt
	}

	if updatedAt, ok := source["updated_at"].(time.Time); ok {
		catalog.UpdatedAt = updatedAt
	}

	return catalog
}

// func (r *elasticRepository) unmarshalCatalog(source map[string]interface{}) *ent.Catalog {
// 	return &ent.Catalog{
// 		ID:          source["_id"].(uuid.UUID),
// 		Name:        source["name"].(string),
// 		Description: source["description"].(string),
// 		Price:       source["price"].(float64),
// 		CreatedAt:   source["created_at"].(time.Time),
// 		UpdatedAt:   source["updated_at"].(time.Time),
// 		DeletedAt:   source["deleted_at"].(*time.Time),
// 	}
// }
