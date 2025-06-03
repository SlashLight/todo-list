package todo

import (
	"log/slog"

	grpcapp "github.com/SlashLight/todo-list/internal/app/todo/grpc"
	task_service "github.com/SlashLight/todo-list/internal/services/task-service"
	"github.com/SlashLight/todo-list/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	taskService := task_service.New(storage, log) // 7 days
	grpcApp := grpcapp.New(log, taskService, grpcPort)

	return &App{GRPCSrv: grpcApp}
}
