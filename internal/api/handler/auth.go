package handler

import (
	"context"
	"gateway/internal/api/models"
	"gateway/internal/grpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	authServiceAddress, err := h.consulProvider.GetService("auth-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	grpcConn, err := grpc.NewAuthClient(authServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	accessToken, refreshToken, message, err := grpcConn.RegisterRequest(context.Background(), req.Login, req.Password, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	grpcConn.Close()

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	authServiceAddress, err := h.consulProvider.GetService("auth-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	grpcConn, err := grpc.NewAuthClient(authServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	accessToken, refreshToken, message, err := grpcConn.LoginRequest(context.Background(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	grpcConn.Close()

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	authServiceAddress, err := h.consulProvider.GetService("auth-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	grpcConn, err := grpc.NewAuthClient(authServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	accessToken, refreshToken, message, err := grpcConn.RefreshRequest(context.Background(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refreshToken failed"})
		return
	}

	grpcConn.Close()

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	authServiceAddress, err := h.consulProvider.GetService("auth-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	grpcConn, err := grpc.NewAuthClient(authServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to auth service"})
		return
	}

	message, err := grpcConn.LogoutRequest(context.Background(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	grpcConn.Close()

	resp := models.LogoutResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}
