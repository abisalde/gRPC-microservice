package main

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

func (c *AuthClient) GetUserByID(ctx context.Context, id int64) (*entpb.User, error) {
	return c.userClient.Get(ctx, &entpb.GetUserRequest{Id: id})
}

func (c *AuthClient) CreateUser(ctx context.Context, user *entpb.User) (*entpb.User, error) {
	return c.userClient.Create(ctx, &entpb.CreateUserRequest{User: user})
}

func (c *AuthClient) GetUserByEmail(ctx context.Context, email string) (*entpb.User, error) {
	// return c.userClient.Get(ctx, &entpb.GetUserRequest{})
	return nil, nil
}
