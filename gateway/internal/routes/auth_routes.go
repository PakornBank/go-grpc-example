package routes

import (
	"github.com/PakornBank/go-grpc-example/gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers the auth routes with the provided gin router group and handler.
func RegisterAuthRoutes(group *gin.RouterGroup, h *handler.AuthHandler) {
	auth := group.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}
