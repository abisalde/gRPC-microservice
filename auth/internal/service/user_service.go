package service

import (
	"context"
	"fmt"

	"github.com/abisalde/grpc-microservice/auth/internal/model"
	"github.com/abisalde/grpc-microservice/auth/internal/repository"
	"github.com/abisalde/grpc-microservice/auth/pkg/ent"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, input *model.RegisterUserInput) (*ent.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email %s is already registered", input.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.CreateNewUser(ctx, &model.RegisterUserInput{
		Email:           input.Email,
		Password:        string(hashedPassword),
		IsEmailVerified: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*ent.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) VerifyEmailExists(ctx context.Context, email string) (bool, error) {
	return s.userRepo.ExistsByEmail(ctx, email)
}

func (s *UserService) VerifyCredentials(ctx context.Context, email, password string) (*ent.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
