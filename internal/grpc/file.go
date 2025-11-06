package grpc

import (
	"context"
	"fmt"
	fileProto "gateway/internal/proto/file"
	"gateway/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileClient struct {
	api fileProto.FileServiceClient
	cc  *grpc.ClientConn
}

func NewFileClient(addr string, metrics *middleware.Metrics) (*FileClient, error) {
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(metrics.GRPCClientMetricsInterceptor),
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to connect to File gRPC service")
		return nil, fmt.Errorf("file grpc client: dial: %w", err)
	}

	log.Debug().Msg("connected to File gRPC service with address: " + addr)

	return &FileClient{
		api: fileProto.NewFileServiceClient(cc),
		cc:  cc,
	}, nil
}

func (c *FileClient) UploadFile(ctx context.Context, fileName, fileOwner string, fileData []byte, chunkSize int) (string, string, string, error) {
	stream, err := c.api.UploadFile(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create upload stream: %w", err)
	}

	for offset := 0; offset < len(fileData); offset += chunkSize {
		end := offset + chunkSize
		if end > len(fileData) {
			end = len(fileData)
		}

		chunk := fileData[offset:end]
		req := &fileProto.UploadFileRequest{
			FileName:  fileName,
			FileOwner: fileOwner,
			Data:      chunk,
		}

		if err := stream.Send(req); err != nil {
			return "", "", "", fmt.Errorf("send chunk: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to receive response: %w", err)
	}

	return resp.FileId, resp.FileUrl, resp.Message, nil
}

func (c *FileClient) GetFile(ctx context.Context, fileID string) (string, error) {
	resp, err := c.api.GetFile(ctx, &fileProto.GetFileRequest{
		FileId: fileID,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to get file via File gRPC service")
		return "", fmt.Errorf("file grpc client: GetFile: %w", err)
	}

	return resp.FileUrl, nil
}

func (c *FileClient) GetFiles(ctx context.Context, fileIDs []string) ([]string, error) {
	resp, err := c.api.GetFiles(ctx, &fileProto.GetFilesRequest{
		FileIds: fileIDs,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to get files via File gRPC service")
		return nil, fmt.Errorf("file grpc client: GetFiles: %w", err)
	}

	return resp.FileUrls, nil
}

func (c *FileClient) DeleteFile(ctx context.Context, fileID string, fileOwner string) (string, error) {
	resp, err := c.api.DeleteFile(ctx, &fileProto.DeleteFileRequest{
		FileId:    fileID,
		FileOwner: fileOwner,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to delete file via File gRPC service")
		return "", fmt.Errorf("file grpc client: DeleteFile: %w", err)
	}

	return resp.Message, nil
}

func (c *FileClient) DeleteFiles(ctx context.Context, fileIDs []string, fileOwner string) (string, error) {
	resp, err := c.api.DeleteFiles(ctx, &fileProto.DeleteFilesRequest{
		FileIds:   fileIDs,
		FileOwner: fileOwner,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to delete files via File gRPC service")
		return "", fmt.Errorf("file grpc client: DeleteFiles: %w", err)
	}

	return resp.Message, nil
}

func (c *FileClient) Close() {
	log.Debug().Msg("closing File gRPC client connection with address: " + c.cc.Target())
	err := c.cc.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close File gRPC client connection")
	}
}
