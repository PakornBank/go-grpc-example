package main

import (
	"log"
	"net"

	"github.com/PakornBank/go-grpc-example/auth/internal/config"
	"github.com/PakornBank/go-grpc-example/auth/internal/di"
	pb "github.com/PakornBank/go-grpc-example/auth/proto/auth/v1"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	container := di.NewContainer(cfg)

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, container.Server)

	log.Printf("Starting Auth gRPC server on port %s\n", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
