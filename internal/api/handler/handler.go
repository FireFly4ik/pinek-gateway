package handler

import (
	"crypto/rsa"
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	metrics        *middleware.Metrics
}

func NewHandler(
	envConf *config.Config,
	cp *consul.ConsulProvider,
	rsaPubKey *rsa.PublicKey,
	metrics *middleware.Metrics,
) *Handler {
	return &Handler{
		envConf:        envConf,
		consulProvider: cp,
		rsaPubKey:      rsaPubKey,
		metrics:        metrics,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.MetricsMiddleware(h.metrics))
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(h.metrics.Reg, promhttp.HandlerOpts{})))

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
			fileGroup.POST(UploadFileRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.UploadFile)
			fileGroup.GET(GetFileRoute, h.GetFile)
			fileGroup.GET(GetFilesRoute, h.GetFiles)
			fileGroup.DELETE(DeleteFileRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.DeleteFile)
			fileGroup.DELETE(DeleteFilesRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.DeleteFiles)
		}
	}

	return r
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}
