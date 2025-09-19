package repository

import (
	"context"

	"github.com/abisalde/grpc-microservice/auth/internal/model"
	"github.com/abisalde/grpc-microservice/auth/pkg/ent"
	"github.com/abisalde/grpc-microservice/auth/pkg/ent/user"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	GetByID(ctx context.Context, id int64) (*ent.User, error)
	CreateNewUser(ctx context.Context, input *model.RegisterUserInput) (*ent.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)
}

func (r *userRepository) CreateNewUser(ctx context.Context, input *model.RegisterUserInput) (*ent.User, error) {
	create := r.client.User.
		Create().
		SetEmail(input.Email).
		SetPasswordHash(input.Password).
		SetNillableIsEmailVerified(&input.IsEmailVerified)

	return create.Save(ctx)
}
