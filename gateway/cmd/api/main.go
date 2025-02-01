package main

import (
	"log"

	"github.com/PakornBank/go-grpc-example/gateway/internal/config"
	"github.com/PakornBank/go-grpc-example/gateway/internal/di"
	"github.com/PakornBank/go-grpc-example/gateway/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	container := di.NewContainer(cfg)
	defer container.Close()

	r := gin.Default()

	router.SetupRoutes(r, container)

	log.Printf("Starting gateway server on port %s\n", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
