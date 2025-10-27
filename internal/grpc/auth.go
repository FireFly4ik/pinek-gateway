package grpc

import (
	"context"
	"fmt"
	authProto "gateway/internal/proto/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api authProto.AuthServiceClient
	cc  *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to connect to Auth gRPC service")
		return nil, fmt.Errorf("auth grpc client: dial: %w", err)
	}

	log.Debug().Msg("connected to Auth gRPC service with address: " + addr)

	return &AuthClient{
		api: authProto.NewAuthServiceClient(cc),
		cc:  cc,
	}, nil
}

func (c *AuthClient) RegisterRequest(ctx context.Context, login, password, username string) (string, string, string, error) {
	resp, err := c.api.Register(ctx, &authProto.RegisterRequest{
		Login:    login,
		Password: password,
		Username: username,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to register user via Auth gRPC service")
		return "", "", "", fmt.Errorf("auth grpc client: RegisterRequest: %w", err)
	}

	return resp.RefreshToken, resp.AccessToken, resp.Message, nil
}

func (c *AuthClient) Close() {
	log.Debug().Msg("closing Auth gRPC client connection with address: " + c.cc.Target())
	err := c.cc.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close Auth gRPC client connection")
	}
}
