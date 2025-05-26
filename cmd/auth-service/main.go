package main

import (
	"log/slog"
	"os"

	"github.com/SlashLight/todo-list/internal/app"
	"github.com/SlashLight/todo-list/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting app")
	//TODO: [] init app
	application := app.New(log, cfg.GRPC.Port, cfg.AuthStoragePath, cfg.TokenTTL)
	application.GRPCSrv.MustRun()

	//TODO: [] start gRPC server
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
