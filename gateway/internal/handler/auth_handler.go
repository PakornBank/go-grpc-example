package handler

import (
	"net/http"

	authPB "github.com/PakornBank/go-grpc-example/auth/proto/auth/v1"
	userPB "github.com/PakornBank/go-grpc-example/user/proto/user/v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterInput is a struct that contains the input fields for the Register method.
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginInput is a struct that contains the input fields for the Login method.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	authClient authPB.AuthServiceClient
	userClient userPB.UserServiceClient
}

func NewAuthHandler(authClient authPB.AuthServiceClient, userClient userPB.UserServiceClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		userClient: userClient,
	}
}

// Register handles the user registration process.
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authClient.Register(c.Request.Context(), &authPB.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if _, err := h.userClient.CreateUser(c.Request.Context(), &userPB.CreateUserRequest{
		UserId:   res.UserId,
		Email:    input.Email,
		FullName: input.FullName,
	}); err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusCreated)
}

// Login handles the user login process.
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authClient.Login(c.Request.Context(), &authPB.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token.Token})
}
