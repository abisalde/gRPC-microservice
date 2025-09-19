package auth_entropy

import (
	"context"
	"log"

	"github.com/abisalde/grpc-microservice/auth/pkg/ent/proto/auth_pbuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service auth_pbuf.UserServiceClient
}

func NewClient(addr string) (*Client, error) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{
		service: auth_pbuf.NewUserServiceClient(conn),
		conn:    conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUserByID(ctx context.Context, id int64) (*auth_pbuf.User, error) {
	return c.service.Get(ctx, &auth_pbuf.GetUserRequest{Id: id})
}

func (c *Client) CreateUser(ctx context.Context, user *auth_pbuf.User) (*auth_pbuf.User, error) {
	log.Print("✅ I made it here")
	return c.service.Create(ctx, &auth_pbuf.CreateUserRequest{User: user})
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*auth_pbuf.User, error) {
	// return c.userClient.Get(ctx, &entpb.GetUserRequest{})
	return nil, nil
}
