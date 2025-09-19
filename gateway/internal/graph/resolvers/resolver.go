package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
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

func NewResolverGraphServer(authURL, catalogURL string) (*Resolver, error) {

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
