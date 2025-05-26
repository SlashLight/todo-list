package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/SlashLight/todo-list/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//TODO: init storage

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCSrv: grpcApp}
}
