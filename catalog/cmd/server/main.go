package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abisalde/gprc-microservice/catalog/internal/es"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "Hello, Catalog Service")
	log.Printf("Received request on / from %s", r.RemoteAddr)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "catalog-service",
		Version:   "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Health check requested from %s", r.RemoteAddr)
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

	// svc := entpb.NewCatalogService(elasticSearch.Client)

}
