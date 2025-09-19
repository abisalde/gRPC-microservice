package clients

import (
	"context"

	"github.com/abisalde/grpc-microservice/auth/pkg/ent/proto/auth_pbuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn       *grpc.ClientConn
	userClient auth_pbuf.UserServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &AuthClient{
		userClient: auth_pbuf.NewUserServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) CreateUser(ctx context.Context, req *auth_pbuf.CreateUserRequest) (*auth_pbuf.User, error) {
	return c.userClient.Create(ctx, req)
}

func (c *AuthClient) GetUser(ctx context.Context, req *auth_pbuf.GetUserRequest) (*auth_pbuf.User, error) {
	return c.userClient.Get(ctx, req)
}

func (c *AuthClient) GetUserByEmail(ctx context.Context, email string) (*auth_pbuf.User, error) {
	return nil, nil
}
