package handler

import (
	"gateway/internal/api/models"
	"gateway/internal/grpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreatePost создание поста
// @Summary Создание поста
// @Description Создание нового поста
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post body models.CreatePostRequest true "Создание поста"
// @Success 200 {object} models.CreatePostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	postId, err := grpcConn.CreatePost(c.Request.Context(), req.PostId, userId, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	grpcConn.Close()

	resp := models.CreatePostResponse{
		PostId: postId,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePost обновление поста
// @Summary Обновление поста
// @Description Обновление существующего поста
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID поста"
// @Param post body models.UpdatePostRequest true "Обновление поста"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post/{id} [post]
func (h *Handler) UpdatePost(c *gin.Context) {
	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	err = grpcConn.UpdatePost(c.Request.Context(), postId, userId, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	grpcConn.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// GetPost получение поста по ID
// @Summary Получение поста по ID
// @Description Получение поста по его ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {object} models.GetPostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post/{id} [get]
func (h *Handler) GetPost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	postId, userId, title, description, extension, tags, err := grpcConn.GetPost(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
		return
	}

	grpcConn.Close()

	tagsResp := make([]models.Tag, len(tags))
	for i, tag := range tags {
		tagsResp[i] = models.Tag{
			TagId: tag[0],
			Name:  tag[1],
		}
	}

	resp := models.GetPostResponse{
		Post: models.Post{
			PostId:      postId,
			UserId:      userId,
			Title:       title,
			Description: description,
			Extension:   extension,
			Tags:        tagsResp,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetPosts получение нескольких постов по ID
// @Summary Получение нескольких постов по ID
// @Description Получение нескольких постов по их ID
// @Tags post
// @Accept json
// @Produce json
// @Param posts_ids query string true "IDs постов, через запятую (например: id1,id2,id3)"
// @Success 200 {object} models.GetPostsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post [get]
func (h *Handler) GetPosts(c *gin.Context) {
	postsIds := c.QueryArray("posts_ids")
	if len(postsIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "posts_ids are required"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	postsIds, userIds, titles, descriptions, extensions, tags, err := grpcConn.GetPosts(c.Request.Context(), postsIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	grpcConn.Close()

	postsResp := make([]models.Post, len(postsIds))
	for i := range postsIds {
		tagsResp := make([]models.Tag, len(tags[i]))
		for j, tag := range tags[i] {
			tagsResp[j] = models.Tag{
				TagId: tag[0],
				Name:  tag[1],
			}
		}

		postsResp[i] = models.Post{
			PostId:      postsIds[i],
			UserId:      userIds[i],
			Title:       titles[i],
			Description: descriptions[i],
			Extension:   extensions[i],
			Tags:        tagsResp,
		}
	}

	resp := models.GetPostsResponse{
		Posts: postsResp,
	}

	c.JSON(http.StatusOK, resp)
}

// SearchPosts поиск постов по параметрам
// @Summary Поиск постов по параметрам
// @Description Поиск постов по заданным параметрам
// @Tags post
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param user_ids query string false "IDs пользователей, через запятую (например: id1,id2,id3)"
// @Param tag_ids query string false "IDs тегов, через запятую (например: id1,id2,id3)"
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Success 200 {object} models.SearchPostsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post/search [get]
func (h *Handler) SearchPosts(c *gin.Context) {
	query := c.Query("query")
	userIds := c.QueryArray("user_ids")
	tagIds := c.QueryArray("tag_ids")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	postsIds, userIds, titles, descriptions, extensions, tags, err := grpcConn.SearchPosts(c.Request.Context(), query, userIds, tagIds, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	grpcConn.Close()

	postsResp := make([]models.Post, len(postsIds))
	for i := range postsIds {
		tagsResp := make([]models.Tag, len(tags[i]))
		for j, tag := range tags[i] {
			tagsResp[j] = models.Tag{
				TagId: tag[0],
				Name:  tag[1],
			}
		}

		postsResp[i] = models.Post{
			PostId:      postsIds[i],
			UserId:      userIds[i],
			Title:       titles[i],
			Description: descriptions[i],
			Extension:   extensions[i],
			Tags:        tagsResp,
		}
	}

	resp := models.SearchPostsResponse{
		Posts: postsResp,
	}

	c.JSON(http.StatusOK, resp)
}

// DeletePost удаление поста по ID
// @Summary Удаление поста по ID
// @Description Удаление поста по его ID
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID поста"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/post/{id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	err = grpcConn.DeletePost(c.Request.Context(), postId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	grpcConn.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// UpdateBoard обновление доски
// @Summary Обновление доски
// @Description Обновление существующей доски
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param board body models.UpdateBoardRequest true "Обновление доски"
// @Success 200 {object} models.UpdateBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board [post]
func (h *Handler) CreateBoard(c *gin.Context) {
	var req models.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	boardId, err := grpcConn.CreateBoard(c.Request.Context(), userId, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board"})
		return
	}

	grpcConn.Close()

	resp := models.CreateBoardResponse{
		BoardId: boardId,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateBoard обновление доски
// @Summary Обновление доски
// @Description Обновление существующей доски
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID доски"
// @Param board body models.UpdateBoardRequest true "Обновление доски"
// @Success 200 {object} models.UpdateBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board/{id} [post]
func (h *Handler) UpdateBoard(c *gin.Context) {
	var req models.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	boardId := c.Param("id")
	if boardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.UpdateBoard(c.Request.Context(), boardId, userId, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update board"})
		return
	}

	grpcConn.Close()

	resp := models.UpdateBoardResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

// GetBoard получение доски по ID
// @Summary Получение доски по ID
// @Description Получение доски по ее ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "ID доски"
// @Success 200 {object} models.GetBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board/{id} [get]
func (h *Handler) GetBoard(c *gin.Context) {
	boardId := c.Param("id")
	if boardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	boardId, userId, name, description, posts, err := grpcConn.GetBoard(c.Request.Context(), boardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get board"})
		return
	}

	grpcConn.Close()

	postStructs := make([]models.Post, len(posts))
	for i, post := range posts {
		postStructs[i] = models.Post{
			PostId:      post[0],
			UserId:      post[1],
			Title:       post[2],
			Description: post[3],
			Extension:   post[4],
		}
	}

	resp := models.GetBoardResponse{
		Board: models.Board{
			BoardId:     boardId,
			UserId:      userId,
			Name:        name,
			Description: description,
			Posts:       postStructs,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetBoards получение нескольких досок по ID
// @Summary Получение нескольких досок по ID
// @Description Получение нескольких досок по их ID
// @Tags post
// @Accept json
// @Produce json
// @Param boards_ids query string true "IDs досок, через запятую (например: id1,id2,id3)"
// @Success 200 {object} models.GetBoardsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board [get]
func (h *Handler) GetBoards(c *gin.Context) {
	boardsIds := c.QueryArray("boards_ids")
	if len(boardsIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "boards_ids are required"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	boardsIds, userIds, names, descriptions, postsList, err := grpcConn.GetBoards(c.Request.Context(), boardsIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get boards"})
		return
	}

	grpcConn.Close()

	boardsResp := make([]models.Board, len(boardsIds))
	for i := range boardsIds {
		postStructs := make([]models.Post, len(postsList[i]))
		for j, post := range postsList[i] {
			postStructs[j] = models.Post{
				PostId:      post[0],
				UserId:      post[1],
				Title:       post[2],
				Description: post[3],
				Extension:   post[4],
			}
		}

		boardsResp[i] = models.Board{
			BoardId:     boardsIds[i],
			UserId:      userIds[i],
			Name:        names[i],
			Description: descriptions[i],
			Posts:       postStructs,
		}
	}

	resp := models.GetBoardsResponse{
		Boards: boardsResp,
	}

	c.JSON(http.StatusOK, resp)
}

// SearchBoards поиск досок по параметрам
// @Summary Поиск досок по параметрам
// @Description Поиск досок по заданным параметрам
// @Tags post
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param user_ids query string false "IDs пользователей, через запятую (например: id1,id2,id3)"
// @Param post_ids query string false "IDs постов, через запятую (например: id1,id2,id3)"
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Success 200 {object} models.SearchBoardsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board/search [get]
func (h *Handler) SearchBoards(c *gin.Context) {
	query := c.Query("query")
	userIds := c.QueryArray("user_ids")
	postIds := c.QueryArray("post_ids")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	boardsIds, userIds, names, descriptions, postsList, err := grpcConn.SearchBoards(c.Request.Context(), query, userIds, postIds, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search boards"})
		return
	}

	grpcConn.Close()

	boardsResp := make([]models.Board, len(boardsIds))
	for i := range boardsIds {
		postStructs := make([]models.Post, len(postsList[i]))
		for j, post := range postsList[i] {
			postStructs[j] = models.Post{
				PostId:      post[0],
				UserId:      post[1],
				Title:       post[2],
				Description: post[3],
				Extension:   post[4],
			}
		}

		boardsResp[i] = models.Board{
			BoardId:     boardsIds[i],
			UserId:      userIds[i],
			Name:        names[i],
			Description: descriptions[i],
			Posts:       postStructs,
		}
	}

	resp := models.SearchBoardsResponse{
		Boards: boardsResp,
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteBoard удаление доски по ID
// @Summary Удаление доски по ID
// @Description Удаление доски по ее ID
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID доски"
// @Success 200 {object} models.DeleteBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/board/{id} [delete]
func (h *Handler) DeleteBoard(c *gin.Context) {
	boardId := c.Param("id")
	if boardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.DeleteBoard(c.Request.Context(), boardId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete board"})
		return
	}

	grpcConn.Close()

	resp := models.DeleteBoardResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

// GetBoards получение нескольких досок по ID
// @Summary Получение нескольких досок по ID
// @Description Получение нескольких досок по их ID
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param tag body models.CreateTagRequest true "Создание тега"
// @Success 200 {object} models.CreateTagResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/tag [post]
func (h *Handler) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	tagId, err := grpcConn.CreateTag(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	grpcConn.Close()

	resp := models.CreateTagResponse{
		TagId: tagId,
	}

	c.JSON(http.StatusOK, resp)
}

// GetTag получение тега по ID
// @Summary Получение тега по ID
// @Description Получение тега по его ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "ID тега"
// @Success 200 {object} models.GetTagResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/tag/{id} [get]
func (h *Handler) GetTag(c *gin.Context) {
	tagId := c.Param("id")
	if tagId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID is required"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	tagId, name, err := grpcConn.GetTag(c.Request.Context(), tagId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tag"})
		return
	}

	grpcConn.Close()

	resp := models.GetTagResponse{
		Tag: models.Tag{
			TagId: tagId,
			Name:  name,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// SearchTags поиск тегов по параметрам
// @Summary Поиск тегов по параметрам
// @Description Поиск тегов по заданным параметрам
// @Tags post
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Success 200 {object} models.SearchTagsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/tag/search [get]
func (h *Handler) SearchTags(c *gin.Context) {
	query := c.Query("query")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	tagsIds, names, err := grpcConn.SearchTags(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tags"})
		return
	}

	grpcConn.Close()

	tagsResp := make([]models.Tag, len(tagsIds))
	for i := range tagsIds {
		tagsResp[i] = models.Tag{
			TagId: tagsIds[i],
			Name:  names[i],
		}
	}

	resp := models.SearchTagsResponse{
		Tags: tagsResp,
	}

	c.JSON(http.StatusOK, resp)
}

// PinPostToBoard прикрепление поста к доске
// @Summary Прикрепление поста к доске
// @Description Прикрепление поста к доске
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path string true "ID поста"
// @Param board_id path string true "ID доски"
// @Success 200 {object} models.PinPostToBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pin/{post_id}/{board_id} [post]
func (h *Handler) PinPostToBoard(c *gin.Context) {
	postId := c.Param("post_id")
	boardId := c.Param("board_id")
	if postId == "" || boardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID and Board ID are required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.PinPostToBoard(c.Request.Context(), postId, boardId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pin post to board"})
		return
	}

	grpcConn.Close()

	resp := models.PinPostToBoardResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

// UnpinPostFromBoard открепление поста от доски
// @Summary Открепление поста от доски
// @Description Открепление поста от доски
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path string true "ID поста"
// @Param board_id path string true "ID доски"
// @Success 200 {object} models.UnpinPostFromBoardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /unpin/{post_id}/{board_id} [post]
func (h *Handler) UnpinPostFromBoard(c *gin.Context) {
	postId := c.Param("post_id")
	boardId := c.Param("board_id")

	if postId == "" || boardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID and Board ID are required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.UnpinPostFromBoard(c.Request.Context(), postId, boardId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unpin post from board"})
		return
	}

	grpcConn.Close()

	resp := models.UnpinPostFromBoardResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

// AddTagToPost добавление тега к посту
// @Summary Добавление тега к посту
// @Description Добавление тега к посту
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path string true "ID поста"
// @Param tag_id path string true "ID тега"
// @Success 200 {object} models.AddTagToPostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tag/{post_id}/{tag_id} [post]
func (h *Handler) AddTagToPost(c *gin.Context) {
	postId := c.Param("post_id")
	tagId := c.Param("tag_id")
	if postId == "" || tagId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID and Tag ID are required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.AddTagToPost(c.Request.Context(), postId, tagId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to post"})
		return
	}

	grpcConn.Close()

	resp := models.AddTagToPostResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

// AddTagToPost добавление тега к посту
// @Summary Добавление тега к посту
// @Description Добавление тега к посту
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path string true "ID поста"
// @Param tag_id path string true "ID тега"
// @Success 200 {object} models.RemoveTagFromPostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /untag/{post_id}/{tag_id} [post]
func (h *Handler) RemoveTagFromPost(c *gin.Context) {
	postId := c.Param("post_id")
	tagId := c.Param("tag_id")

	if postId == "" || tagId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID and Tag ID are required"})
		return
	}

	userId := c.GetString("user_id")

	postServiceAddress, err := h.consulProvider.GetService("post-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	grpcConn, err := grpc.NewPostClient(postServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to post service"})
		return
	}

	message, err := grpcConn.RemoveTagFromPost(c.Request.Context(), postId, tagId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from post"})
		return
	}

	grpcConn.Close()

	resp := models.RemoveTagFromPostResponse{
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}
