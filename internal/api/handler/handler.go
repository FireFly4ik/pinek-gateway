package handler

import (
	"crypto/rsa"
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	DefaultRoute = "/"
	SwaggerRoute = "/swagger/*any"

	RegisterRoute = "/register"
	LoginRoute    = "/login"
	RefreshRoute  = "/refresh"
	LogoutRoute   = "/logout"

	UploadFileRoute  = "/upload"
	GetFileRoute     = "/:id"
	GetFilesRoute    = "/files"
	DeleteFileRoute  = "/:id"
	DeleteFilesRoute = "/files"

	CreatePostRoute    = "/post"
	UpdatePostRoute    = "/post/:id"
	GetPostRoute       = "/post/:id"
	GetPostsRoute      = "/post"
	SearchPostsRoute   = "/post/search"
	DeletePostRoute    = "/post/:id"
	CreateBoardRoute   = "/board"
	UpdateBoardRoute   = "/board/:id"
	GetBoardRoute      = "/board/:id"
	GetBoardsRoute     = "/board"
	SearchBoardsRoute  = "/board/search"
	DeleteBoardRoute   = "/board/:id"
	CreateTagRoute     = "/tag"
	GetTagRoute        = "/tag/:id"
	SearchTagsRoute    = "/tag/search"
	PinPostToBoard     = "/pin/:post_id/:board_id"
	UnpinPostFromBoard = "/unpin/:post_id/:board_id"
	AddTagToPost       = "/tag/:post_id/:tag_id"
	RemoveTagFromPost  = "/untag/:post_id/:tag_id"
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
		handler.GET(SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

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

		postGroup := handler.Group("/post")
		{
			postGroup.POST(CreatePostRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.CreatePost)
			postGroup.POST(UpdatePostRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.UpdatePost)
			postGroup.GET(GetPostRoute, h.GetPost)
			postGroup.GET(GetPostsRoute, h.GetPosts)
			postGroup.GET(SearchPostsRoute, h.SearchPosts)
			postGroup.DELETE(DeletePostRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.DeletePost)
			postGroup.POST(CreateBoardRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.CreateBoard)
			postGroup.POST(UpdateBoardRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.UpdateBoard)
			postGroup.GET(GetBoardRoute, h.GetBoard)
			postGroup.GET(GetBoardsRoute, h.GetBoards)
			postGroup.GET(SearchBoardsRoute, h.SearchBoards)
			postGroup.DELETE(DeleteBoardRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.DeleteBoard)
			postGroup.POST(CreateTagRoute, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.CreateTag)
			postGroup.GET(GetTagRoute, h.GetTag)
			postGroup.GET(SearchTagsRoute, h.SearchTags)
			postGroup.POST(PinPostToBoard, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.PinPostToBoard)
			postGroup.POST(UnpinPostFromBoard, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.UnpinPostFromBoard)
			postGroup.POST(AddTagToPost, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.AddTagToPost)
			postGroup.POST(RemoveTagFromPost, middleware.JWTAccessMiddleware(h.rsaPubKey, h.metrics), h.RemoveTagFromPost)
		}
	}

	return r
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}
