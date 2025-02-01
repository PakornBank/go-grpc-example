package di

import (
	"log"

	"github.com/PakornBank/go-grpc-example/user/internal/config"
	"github.com/PakornBank/go-grpc-example/user/internal/database"
	"github.com/PakornBank/go-grpc-example/user/internal/repository"
	"github.com/PakornBank/go-grpc-example/user/internal/server"
	"github.com/PakornBank/go-grpc-example/user/internal/service"
)

type Container struct {
	Server *server.Server
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewDataBase(cfg)
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	r := repository.NewRepository(db)
	s := service.NewService(r)

	return &Container{
		Server: server.NewServer(s),
	}
}
