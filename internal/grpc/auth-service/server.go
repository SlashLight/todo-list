package auth_service

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/SlashLight/todo-list/api/gen/go/auth"
	auth_service "github.com/SlashLight/todo-list/internal/services/auth-service"
	"github.com/SlashLight/todo-list/pkg/my_err"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) (string, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	service Service
}

func RegisterServerAPI(gRPC *grpc.Server, service Service) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{service: service})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	err := validateLogin(req)
	if err != nil {
		return nil, err
	}

	token, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth_service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	userID, err := s.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, my_err.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}
