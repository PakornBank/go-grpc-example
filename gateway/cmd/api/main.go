package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	go func() {
		log.Printf("Starting gateway server on port %s\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	// Wait for termination signal
	<-quit
	log.Println("\nShutting down API Gateway...")

	// Create a timeout context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API Gateway forced to shutdown: %s", err)
	}

	log.Println("API Gateway exited gracefully")
}
