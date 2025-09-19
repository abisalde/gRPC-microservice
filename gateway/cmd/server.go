package server

import (
	"context"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/abisalde/grpc-microservice/gateway/internal/graph"
	"github.com/abisalde/grpc-microservice/gateway/internal/graph/resolvers"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func SetupGraphQLServer(r *resolvers.Resolver) (server *handler.Server, port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: r}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		op := graphql.GetOperationContext(ctx)
		log.Println("Operation Name:::::", op.OperationName)
		log.Println("GraphQL operation received")
		return next(ctx)
	})

	return srv, port
}
