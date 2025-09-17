package auth_entropy

import (
	"context"

	"github.com/abisalde/gprc-microservice/auth/pkg/ent/proto/entpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service entpb.UserServiceClient
}

func NewClient(addr string) (*Client, error) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{
		service: entpb.NewUserServiceClient(conn),
		conn:    conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUserByID(ctx context.Context, id int64) (*entpb.User, error) {
	return c.service.Get(ctx, &entpb.GetUserRequest{Id: id})
}

func (c *Client) CreateUser(ctx context.Context, user *entpb.User) (*entpb.User, error) {
	return c.service.Create(ctx, &entpb.CreateUserRequest{User: user})
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*entpb.User, error) {
	// return c.userClient.Get(ctx, &entpb.GetUserRequest{})
	return nil, nil
}
