package main

import (
	"context"
	"log"
	"time"

	"github.com/abisalde/grpc-microservice/catalog/internal/es"
	"github.com/abisalde/grpc-microservice/catalog/internal/repository"
	"github.com/abisalde/grpc-microservice/catalog/internal/service"
	"github.com/abisalde/grpc-microservice/catalog/pkg/catalog_entropy"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
