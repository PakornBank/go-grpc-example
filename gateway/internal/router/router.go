package router

import (
	"github.com/PakornBank/go-grpc-example/gateway/internal/di"
	"github.com/PakornBank/go-grpc-example/gateway/internal/routes"
	"github.com/gin-gonic/gin"
)

// SetupRoutes call functions to register routes on gin router.
func SetupRoutes(router *gin.Engine, container *di.Container) {
	group := router.Group("/api")
	routes.RegisterAuthRoutes(group, container.AuthHandler)
}
