package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/SlashLight/todo-list/api/gen/go/auth"
)

type Client struct {
	api authv1.AuthClient
	log *slog.Logger
}

func New(addr string, log *slog.Logger, retries int, timeout time.Duration) (*Client, error) {
	const op = "auth.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.DeadlineExceeded, codes.Aborted),
		grpcretry.WithMax(uint(retries)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadSent, grpclog.PayloadReceived),
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{api: authv1.NewAuthClient(conn)}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) Register(ctx context.Context, email, password string) (string, error) {
	const op = "auth.grpc.Register"

	resp, err := c.api.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.UserId, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.grpc.Login"

	resp, err := c.api.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}
