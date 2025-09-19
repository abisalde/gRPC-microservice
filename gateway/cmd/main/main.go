package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	server "github.com/abisalde/grpc-microservice/gateway/cmd"
	"github.com/abisalde/grpc-microservice/gateway/internal/graph/resolvers"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
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
		Service:   "graphql-gateway-service",
		Version:   "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Health check requested from %s", r.RemoteAddr)
}

func main() {

	authURL := "grpc-microservice-auth:50051"
	catalogURL := "grpc-microservice-catalog:50052"

	resolvers, err := resolvers.NewResolverGraphServer(authURL, catalogURL)

	if err != nil {
		log.Fatalf("❌ Failed to setup all client 🎛️: %v", err)
	}

	http.HandleFunc("/health", healthHandler)

	gqlSrv, port := server.SetupGraphQLServer(resolvers)

	http.Handle("/", playground.ApolloSandboxHandler("Microservice GraphQL playground", "/graphql"))
	http.Handle("/graphql", gqlSrv)

	log.Printf("🚀 Gateway server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
