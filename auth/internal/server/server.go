package server

import (
	"context"
	"errors"
	"log"

	"github.com/PakornBank/go-grpc-example/auth/internal/service"
	pb "github.com/PakornBank/go-grpc-example/auth/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server handles authentication gRPC requests.
type Server struct {
	pb.UnimplementedAuthServiceServer
	service service.Service
}

// NewServer creates a new Server instance.
func NewServer(service service.Service) *Server {
	return &Server{
		service: service,
	}
}

// Register handles user registration.
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userID, err := s.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		log.Printf("register error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.RegisterResponse{UserId: userID}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		log.Printf("login error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteUser(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("delete user error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	userID, email, valid, err := s.service.VerifyToken(req.Token)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		log.Printf("JWT verification failed: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &pb.VerifyTokenResponse{
		UserId: userID,
		Email:  email,
		Valid:  valid,
	}, nil
}
