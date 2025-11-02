package grpc

import (
	"context"
	"fmt"
	authProto "gateway/internal/proto/auth"
	"gateway/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api authProto.AuthServiceClient
	cc  *grpc.ClientConn
}

func NewAuthClient(addr string, metrics *middleware.Metrics) (*AuthClient, error) {
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(metrics.GRPCClientMetricsInterceptor),
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

	return resp.AccessToken, resp.RefreshToken, resp.Message, nil
}

func (c *AuthClient) LoginRequest(ctx context.Context, login, password string) (string, string, string, error) {
	resp, err := c.api.Login(ctx, &authProto.LoginRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to login user via Auth gRPC service")
		return "", "", "", fmt.Errorf("auth grpc client: LoginRequest: %w", err)
	}

	return resp.AccessToken, resp.RefreshToken, resp.Message, nil
}

func (c *AuthClient) RefreshRequest(ctx context.Context, refreshToken string) (string, string, string, error) {
	resp, err := c.api.Refresh(ctx, &authProto.RefreshRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to refresh tokens via Auth gRPC service")
		return "", "", "", fmt.Errorf("auth grpc client: RefreshRequest: %w", err)
	}

	return resp.AccessToken, resp.RefreshToken, resp.Message, nil
}

func (c *AuthClient) LogoutRequest(ctx context.Context, refreshToken string) (string, error) {
	resp, err := c.api.Logout(ctx, &authProto.LogoutRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to logout user via Auth gRPC service")
		return "", fmt.Errorf("auth grpc client: LogoutRequest: %w", err)
	}

	return resp.Message, nil
}

func (c *AuthClient) Close() {
	log.Debug().Msg("closing Auth gRPC client connection with address: " + c.cc.Target())
	err := c.cc.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close Auth gRPC client connection")
	}
}
