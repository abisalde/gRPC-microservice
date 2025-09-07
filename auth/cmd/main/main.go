package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
		Service:   "authentication-service",
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
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:         ":50051",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Authentication Service starting on :50051")
	log.Println("Endpoints available:")
	log.Println("  - http://localhost:50051/")
	log.Println("  - http://localhost:50051/health")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
