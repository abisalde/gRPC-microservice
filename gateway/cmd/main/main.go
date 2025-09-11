package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
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

	fmt.Fprintln(w, "Hello, Authentication Service")
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

	// server := &http.Server{
	// 	Addr:         ":8080",
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	IdleTimeout:  60 * time.Second,
	// }

	log.Println("Authentication Service starting on :8080")
	log.Println("Endpoints available:")
	log.Println("  - http://localhost:8080/")
	log.Println("  - http://localhost:8080/health")

	// s, err := server.NewGraphQLServer()

	// 	authService := auth.NewService(authClient)

	// // Configure GraphQL server
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
	// 	Resolvers: &graph.Resolver{
	// 		AuthService: authService,
	// 	},
	// }))

	http.HandleFunc("/health", healthHandler)

	http.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Gateway server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
