package handler

import (
	"gateway/internal/api/models"
	"gateway/internal/grpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
			Name:  tag[0],
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

func (h *Handler) GetPosts(c *gin.Context) {
	var req models.GetPostsRequest
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

	postsIds, userIds, titles, descriptions, extensions, tags, err := grpcConn.GetPosts(c.Request.Context(), req.PostsIds)
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

func (h *Handler) SearchPosts(c *gin.Context) {
	var req models.SearchPostsRequest
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

	postsIds, userIds, titles, descriptions, extensions, tags, err := grpcConn.SearchPosts(c.Request.Context(), req.Query, req.UserIds, req.TagIds, req.Limit, req.Offset)
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

func (h *Handler) GetBoards(c *gin.Context) {
	var req models.GetBoardsRequest
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

	boardsIds, userIds, names, descriptions, postsList, err := grpcConn.GetBoards(c.Request.Context(), req.BoardsIds)
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

func (h *Handler) SearchBoards(c *gin.Context) {
	var req models.SearchBoardsRequest
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

	boardsIds, userIds, names, descriptions, postsList, err := grpcConn.SearchBoards(c.Request.Context(), req.Query, req.UserIds, req.PostIds, req.Limit, req.Offset)
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

			boardsResp[i] = models.Board{
				BoardId:     boardsIds[i],
				UserId:      userIds[i],
				Name:        names[i],
				Description: descriptions[i],
				Posts:       postStructs,
			}
		}
	}

	resp := models.SearchBoardsResponse{
		Boards: boardsResp,
	}

	c.JSON(http.StatusOK, resp)
}

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
