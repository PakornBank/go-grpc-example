package di

import (
	"log"

	"github.com/PakornBank/go-grpc-example/auth/internal/config"
	"github.com/PakornBank/go-grpc-example/auth/internal/database"
	"github.com/PakornBank/go-grpc-example/auth/internal/repository"
	"github.com/PakornBank/go-grpc-example/auth/internal/server"
	"github.com/PakornBank/go-grpc-example/auth/internal/service"
	"gorm.io/gorm"
)

type Container struct {
	Server *server.Server
	DB     *gorm.DB
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewDataBase(cfg)
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	r := repository.NewRepository(db)
	s := service.NewService(r, cfg)

	return &Container{
		Server: server.NewServer(s),
		DB:     db,
	}
}

func (c *Container) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
