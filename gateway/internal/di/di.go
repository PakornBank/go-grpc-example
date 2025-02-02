package di

import (
	"context"
	"log"
	"time"

	authPB "github.com/PakornBank/go-grpc-example/auth/proto/auth/v1"
	"github.com/PakornBank/go-grpc-example/gateway/internal/config"
	"github.com/PakornBank/go-grpc-example/gateway/internal/handler"
	"github.com/PakornBank/go-grpc-example/gateway/internal/security"
	userPB "github.com/PakornBank/go-grpc-example/user/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	AuthConn    *grpc.ClientConn
	UserConn    *grpc.ClientConn
}

// NewGRPCConnection establishes a gRPC connection with a timeout
func NewGRPCConnection(address string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewContainer(cfg *config.Config) *Container {
	creds := security.NewCredentials(cfg)

	authConn, err := NewGRPCConnection(cfg.AuthServiceAddr, creds)
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}

	userConn, err := NewGRPCConnection(cfg.UserServiceAddr, creds)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	authClient := authPB.NewAuthServiceClient(authConn)
	userClient := userPB.NewUserServiceClient(userConn)

	authHandler := handler.NewAuthHandler(authClient, userClient)

	return &Container{
		AuthHandler: authHandler,
		AuthConn:    authConn,
		UserConn:    userConn,
	}
}

// Close ensures that both gRPC connections are properly closed
func (c *Container) Close() {
	if c.AuthConn != nil {
		c.AuthConn.Close()
	}
	if c.UserConn != nil {
		c.UserConn.Close()
	}
}
