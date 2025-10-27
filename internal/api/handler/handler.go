package handler

import (
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultRoute = "/"
	//SwaggerRoute = "/swagger/*any"
)

type Handler struct {
	envConf        *config.Config
	consulProvider *consul.ConsulProvider
}

func NewHandler(envConf *config.Config, cp *consul.ConsulProvider) *Handler {
	return &Handler{
		envConf:        envConf,
		consulProvider: cp,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.HeadersMiddleware())

	handler := r.Group("/api/v1")
	{
		handler.GET(DefaultRoute, h.Hello)
		//handler.GET(SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}
