package auth

import (
	"log/slog"
	"time"

	grpcapp "github.com/SlashLight/todo-list/internal/app/auth/grpc"
	auth_service "github.com/SlashLight/todo-list/internal/services/auth-service"
	"github.com/SlashLight/todo-list/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth_service.New(storage, storage, log, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCSrv: grpcApp}
}
