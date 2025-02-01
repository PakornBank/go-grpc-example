package server

import (
	"context"
	"errors"
	"log"

	"github.com/PakornBank/go-grpc-example/user/internal/service"
	pb "github.com/PakornBank/go-grpc-example/user/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server handles authentication gRPC requests.
type Server struct {
	pb.UnimplementedUserServiceServer
	service service.Service
}

// NewServer creates a new Server instance.
func NewServer(service service.Service) *Server {
	return &Server{
		service: service,
	}
}

// CreateUser handles user registration.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := s.service.CreateUser(ctx, req.UserId, req.Email, req.FullName)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		if errors.Is(err, service.ErrIDAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user ID already exists")
		}
		if errors.Is(err, service.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, "invalid user ID")
		}
		log.Printf("create user error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.CreateUserResponse{User: &pb.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{User: &pb.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}}, nil
}
