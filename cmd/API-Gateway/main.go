package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	authgrpc "github.com/SlashLight/todo-list/internal/clients/auth-service/grpc"
	taskgrpc "github.com/SlashLight/todo-list/internal/clients/task-service/grpc"
	"github.com/SlashLight/todo-list/internal/config"
	"github.com/SlashLight/todo-list/internal/http/api-gateway/handlers"
	"github.com/SlashLight/todo-list/internal/http/api-gateway/router"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.APIGatewayConfig.Env)

	log.Info("starting app")

	authService, err := authgrpc.New(fmt.Sprintf("localhost:%d", cfg.GRPCConfig.AuthConfig.Port), log, 3, cfg.GRPCConfig.AuthConfig.Timeout)
	if err != nil {
		log.Error("failed to create auth service client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	taskService, err := taskgrpc.New(fmt.Sprintf("localhost:%d", cfg.GRPCConfig.TaskConfig.Port), log, 3, cfg.GRPCConfig.TaskConfig.Timeout)
	if err != nil {
		log.Error("failed to create task service client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	gateway := handlers.New(authService, taskService, log)
	mux := router.New(gateway, cfg.AuthConfig.SecretKey)

	log.Info("initializing mux")
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Error("failed to start HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
