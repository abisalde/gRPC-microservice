package clients

import (
	"context"

	"github.com/abisalde/gprc-microservice/auth/pkg/ent/proto/entpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn       *grpc.ClientConn
	userClient entpb.UserServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &AuthClient{
		userClient: entpb.NewUserServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) CreateUser(ctx context.Context, req *entpb.CreateUserRequest) (*entpb.User, error) {
	return c.userClient.Create(ctx, req)
}

func (c *AuthClient) GetUser(ctx context.Context, req *entpb.GetUserRequest) (*entpb.User, error) {
	return c.userClient.Get(ctx, req)
}

func (c *AuthClient) GetUserByEmail(ctx context.Context, email string) (*entpb.User, error) {
	return nil, nil
}
