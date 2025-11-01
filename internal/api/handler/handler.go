package handler

import (
	"crypto/rsa"
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultRoute = "/"
	//SwaggerRoute = "/swagger/*any"

	RegisterRoute = "/register"
	LoginRoute    = "/login"
	RefreshRoute  = "/refresh"
	LogoutRoute   = "/logout"

	UploadFileRoute  = "/upload"
	GetFileRoute     = "/:id"
	GetFilesRoute    = "/files"
	DeleteFileRoute  = "/:id"
	DeleteFilesRoute = "/files"
)

type Handler struct {
	envConf        *config.Config
	consulProvider *consul.ConsulProvider
	rsaPubKey      *rsa.PublicKey
}

func NewHandler(envConf *config.Config, cp *consul.ConsulProvider, rsaPubKey *rsa.PublicKey) *Handler {
	return &Handler{
		envConf:        envConf,
		consulProvider: cp,
		rsaPubKey:      rsaPubKey,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	handler := r.Group("/api/v1")
	{
		handler.GET(DefaultRoute, h.Hello)
		//handler.GET(SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

		authGroup := handler.Group("/auth")
		{
			authGroup.POST(RegisterRoute, h.Register)
			authGroup.POST(LoginRoute, h.Login)
			authGroup.POST(RefreshRoute, h.Refresh)
			authGroup.POST(LogoutRoute, h.Logout)
		}

		fileGroup := handler.Group("/file")
		{
			fileGroup.POST(UploadFileRoute, h.UploadFile).Use(middleware.JWTAccessMiddleware(h.rsaPubKey))
			fileGroup.GET(GetFileRoute, h.GetFile)
			fileGroup.GET(GetFilesRoute, h.GetFiles)
			fileGroup.DELETE(DeleteFileRoute, h.DeleteFile).Use(middleware.JWTAccessMiddleware(h.rsaPubKey))
			fileGroup.DELETE(DeleteFilesRoute, h.DeleteFiles).Use(middleware.JWTAccessMiddleware(h.rsaPubKey))
		}
	}

	return r
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}
