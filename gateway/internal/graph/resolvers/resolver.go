package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/abisalde/gprc-microservice/auth/pkg/auth_entropy"
)

type Resolver struct {
	*Server
}

type Server struct {
	authClient auth_entropy.Client
	// catalogClient catalog_entropy.Client
}

func NewResolverGraphServer(authURL, catalogURL string) (*Resolver, error) {

	authClient, err := auth_entropy.NewClient(authURL)
	if err != nil {
		return nil, err
	}

	// catalogClient, err := catalog_entropy.NewClient(catalogURL)
	// if err != nil {
	// 	authClient.Close()
	// 	return nil, err
	// }

	server := &Server{
		authClient: *authClient,
	}

	return &Resolver{Server: server}, nil
}
