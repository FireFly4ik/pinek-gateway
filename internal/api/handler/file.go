package handler

import (
	"context"
	"gateway/internal/api/models"
	"gateway/internal/grpc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

const (
	MaxFileSize = 20 * 1024 * 1024 // 20 MB
	ChunkSize   = 1024 * 1024      // 1 MB
)

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("unable to get file from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get file from request"})
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds maximum allowed size"})
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("unable to read file data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to read file data"})
		return
	}

	fileServiceAddress, err := h.consulProvider.GetService("file-service")
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}

	grpcConn, err := grpc.NewFileClient(fileServiceAddress, h.metrics)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}
	defer grpcConn.Close()

	fileID, fileURL, message, err := grpcConn.UploadFile(context.Background(), header.Filename, fileData, ChunkSize)
	if err != nil {
		log.Error().Err(err).Msg("file upload failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":  fileID,
		"file_url": fileURL,
		"message":  message,
	})
}

func (h *Handler) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_id is required"})
		return
	}

	fileServiceAddress, err := h.consulProvider.GetService("file-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}

	grpcConn, err := grpc.NewFileClient(fileServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}
	defer grpcConn.Close()

	fileURL, err := grpcConn.GetFile(context.Background(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_url": fileURL,
	})
}

func (h *Handler) GetFiles(c *gin.Context) {
	var req models.UploadDeleteFilesResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	fileServiceAddress, err := h.consulProvider.GetService("file-service")
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}

	grpcConn, err := grpc.NewFileClient(fileServiceAddress, h.metrics)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}
	defer grpcConn.Close()

	fileURLs, err := grpcConn.GetFiles(context.Background(), req.FileIDs)
	if err != nil {
		log.Error().Err(err).Msg("failed to get files")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_urls": fileURLs,
	})
}

func (h *Handler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_id is required"})
		return
	}

	fileServiceAddress, err := h.consulProvider.GetService("file-service")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}

	grpcConn, err := grpc.NewFileClient(fileServiceAddress, h.metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}
	defer grpcConn.Close()

	message, err := grpcConn.DeleteFile(context.Background(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *Handler) DeleteFiles(c *gin.Context) {
	var req models.UploadDeleteFilesResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}

	fileServiceAddress, err := h.consulProvider.GetService("file-service")
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}

	grpcConn, err := grpc.NewFileClient(fileServiceAddress, h.metrics)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to file service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to file service"})
		return
	}
	defer grpcConn.Close()

	message, err := grpcConn.DeleteFiles(context.Background(), req.FileIDs)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete files")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
