package auth_service

import (
	"context"

	"google.golang.org/grpc"

	authv1 "github.com/SlashLight/todo-list/api/gen/go/auth"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func RegisterServerAPI(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("implement me")
}
