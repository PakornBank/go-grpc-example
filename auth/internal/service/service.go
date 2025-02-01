package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PakornBank/go-grpc-example/auth/internal/config"
	"github.com/PakornBank/go-grpc-example/auth/internal/model"
	"github.com/PakornBank/go-grpc-example/auth/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnexpectedSigning  = errors.New("unexpected signing method")
)

// Service defines the methods that a service must implement.
type Service interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	VerifyToken(token string) (string, string, bool, error)
}

// service is a struct that provides methods to interact with the authentication service.
type service struct {
	repository  repository.Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewService creates a new instance of service with the provided repository and configuration.
func NewService(repository repository.Repository, config *config.Config) Service {
	return &service{
		repository:  repository,
		jwtSecret:   []byte(config.JWTSecret),
		tokenExpiry: config.TokenExpiry,
	}
}

// Register handles the user registration process.
func (s *service) Register(ctx context.Context, email, password string) (string, error) {
	existingCredential, _ := s.repository.FindByEmail(ctx, email)
	if existingCredential != nil {
		return "", ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	credential := &model.Credential{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repository.CreateUser(ctx, credential); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return credential.ID.String(), nil
}

// Login handles the user login process.
func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token : %w", err)
	}

	return token, nil
}

func (s *service) VerifyToken(token string) (string, string, bool, error) {
	if token == "" {
		return "", "", false, errors.New("empty token provided")
	}

	parsedToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigning
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", "", false, fmt.Errorf("token parse error: %w", err)
	}

	if !parsedToken.Valid {
		return "", "", false, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", false, errors.New("failed to parse claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", "", false, errors.New("missing or invalid exp claim in token")
	}
	if time.Now().Unix() > int64(exp) {
		return "", "", false, ErrTokenExpired
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", "", false, errors.New("missing or invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", "", false, errors.New("missing or invalid email in token")
	}

	return userID, email, true, nil
}

// generateToken generates a JWT token for the given user.
func (s *service) generateToken(user *model.Credential) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
