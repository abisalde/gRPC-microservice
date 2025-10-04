package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"log"
	"net"
	"time"

	"github.com/abisalde/grpc-microservice/auth/pkg/auth_entropy"
	"github.com/abisalde/grpc-microservice/catalog/pkg/catalog_entropy"
)

type Resolver struct {
	*Server
}

type Server struct {
	authClient    *auth_entropy.Client
	catalogClient *catalog_entropy.Client
}

func debugServiceResolution() {
	log.Println("=== Debugging Service Resolution ===")

	services := map[string]string{
		"auth-service":    "50051",
		"catalog-service": "50052",
	}

	// Test DNS resolution

	for host, port := range services {
		addrs, err := net.LookupHost(host)
		if err != nil {
			log.Printf("❌ DNS lookup failed for %s: %v", host, err)
		} else {
			log.Printf("✅ %s resolves to: %v", host, addrs)
		}

		// Test port connectivity with correct port
		address := net.JoinHostPort(host, port)
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			log.Printf("❌ TCP connection failed to %s:%s: %v", host, port, err)
		} else {
			conn.Close()
			log.Printf("✅ TCP connection successful to %s:%s", host, port)
		}
	}
}

func NewResolverGraphServer(authURL, catalogURL string) (*Resolver, error) {

	debugServiceResolution()

	authClient, err := auth_entropy.NewClient(authURL)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog_entropy.NewClient(catalogURL)
	if err != nil {
		authClient.Close()
		return nil, err
	}

	server := &Server{
		authClient:    authClient,
		catalogClient: catalogClient,
	}

	return &Resolver{Server: server}, nil
}

func (r *Resolver) GetAuthClient() *auth_entropy.Client {
	if r.Server != nil && r.Server.authClient != nil {
		return r.Server.authClient
	}
	return nil
}

func (r *Resolver) GetCatalogClient() *catalog_entropy.Client {
	if r.Server != nil && r.Server.catalogClient != nil {
		return r.Server.catalogClient
	}
	return nil
}
