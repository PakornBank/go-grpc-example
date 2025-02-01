package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PakornBank/go-grpc-example/user/internal/model"
	"github.com/PakornBank/go-grpc-example/user/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrIDAlreadyExists    = errors.New("user ID already registered")
	ErrInvalidID          = errors.New("invalid user ID")
)

// Service defines the methods that a service must implement.
type Service interface {
	CreateUser(ctx context.Context, id, email, password string) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
}

// service is a struct that provides methods to interact with the user service.
type service struct {
	repository repository.Repository
}

// NewService creates a new instance of service with the provided repository and configuration.
func NewService(repository repository.Repository) Service {
	return &service{repository: repository}
}

// CreateUser handles the user registration process.
func (s *service) CreateUser(ctx context.Context, id, email, fullName string) (*model.User, error) {
	if existingUser, _ := s.repository.FindByID(ctx, id); existingUser != nil {
		return nil, ErrIDAlreadyExists
	}

	if existingUser, _ := s.repository.FindByEmail(ctx, email); existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	user := &model.User{
		ID:       parsedID,
		Email:    email,
		FullName: fullName,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *service) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repository.FindByID(ctx, id)
}
