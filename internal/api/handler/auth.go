package handler

import (
	"context"
	"gateway/internal/api/models"
	"gateway/internal/grpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register регистрация пользователя
// @Summary Регистрация пользователя
// @Description Регистрация пользователя с логином, паролем и именем пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "Регистрационные данные"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
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

	if message == "User already exists" {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

// Login вход пользователя
// @Summary Вход пользователя
// @Description Аутентификация пользователя с логином и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
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

	if message == "Invalid credentials" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

// Login вход пользователя
// @Summary Вход пользователя
// @Description Аутентификация пользователя с логином и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.RefreshRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
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

	if message == "Invalid credentials" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	resp := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      message,
	}

	c.JSON(http.StatusOK, resp)
}

// Logout выход пользователя
// @Summary Выход пользователя
// @Description Выход пользователя и аннулирование refresh токена
// @Tags auth
// @Accept json
// @Produce json
// @Param logout body models.LogoutRequest true "Данные для выхода"
// @Success 200 {object} models.LogoutResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
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
