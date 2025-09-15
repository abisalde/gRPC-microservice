package main

import (
	"context"
	"log"
	"time"

	"github.com/abisalde/gprc-microservice/catalog/internal/es"
	"github.com/abisalde/gprc-microservice/catalog/internal/repository"
	"github.com/abisalde/gprc-microservice/catalog/internal/service"
	"github.com/abisalde/gprc-microservice/catalog/pkg/catalog_entropy"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

func setupElasticSearch() (*es.ElasticClient, error) {
	elastic, err := es.Connect()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if err := elastic.HealthCheck(ctx); err != nil {
		elastic.Close()
		return nil, err
	}

	return elastic, nil
}

func main() {

	elasticSearch, err := setupElasticSearch()

	if err != nil {
		log.Fatalf("❌ Failed to setup database: %v", err)
	}
	defer elasticSearch.Close()

	r := repository.NewCatalogRepository(elasticSearch)

	s := service.NewProductCatalogService(r)
	catalog_entropy.ListenGRPC(s)

}
