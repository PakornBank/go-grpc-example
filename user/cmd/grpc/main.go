package main

import (
	"log"
	"net"

	"github.com/PakornBank/go-grpc-example/user/internal/config"
	"github.com/PakornBank/go-grpc-example/user/internal/di"
	pb "github.com/PakornBank/go-grpc-example/user/proto/user/v1"
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
	pb.RegisterUserServiceServer(s, container.Server)

	log.Printf("Starting User gRPC server on port %s\n", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
