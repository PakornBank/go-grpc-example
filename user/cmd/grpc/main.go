package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/PakornBank/go-grpc-example/user/internal/config"
	"github.com/PakornBank/go-grpc-example/user/internal/di"
	"github.com/PakornBank/go-grpc-example/user/internal/security"
	pb "github.com/PakornBank/go-grpc-example/user/proto/user/v1"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	creds := security.NewCredentials(cfg)

	container := di.NewContainer(cfg)

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterUserServiceServer(s, container.Server)

	// Handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	go func() {
		log.Printf("Starting User gRPC server on port %s\n", cfg.ServerPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve: ", err)
		}
	}()

	// Wait for termination signal
	<-quit
	fmt.Println("\nShutting down gRPC server...")

	// Gracefully stop the gRPC server
	s.GracefulStop()
	fmt.Println("gRPC server stopped")

	// Close the container
	fmt.Println("Closing container...")
	if err := container.Close(); err != nil {
		log.Printf("Error closing container: %v", err)
	}
	fmt.Println("Server exited gracefully.")
}
