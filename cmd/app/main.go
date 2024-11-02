package main

import (
	"context"
	"fmt"
	taskpb "github.com/makebeepboop/protos/gen/go/task"
	taskgrpc "github.com/makebeepboop/workflow-app/internal/clients/task/grpc"
	"github.com/makebeepboop/workflow-app/internal/config"
	"github.com/makebeepboop/workflow-app/internal/lib/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal       = "local"
	envDevelopment = "development"
	envProduction  = "production"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting application",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.Graphql.Port),
	)

	taskClient, err := taskgrpc.New(
		context.Background(),
		log,
		cfg.Clients.Task.Address,
		cfg.Clients.Task.Timeout,
		cfg.Clients.Task.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init task client", sl.Err(err))
		os.Exit(1)
	}

	// tests
	res, err := taskClient.Api.Status(context.Background(), &taskpb.StatusRequest{TaskId: 1})
	if err != nil {
		log.Error("failed to call api", sl.Err(err))
	}

	fmt.Println(res.GetStatus())
	//application := app.New(log, cfg.Graphql.Port)
	//go application.Graphql.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signalStop := <-stop
	log.Info("stopping application", slog.String("signal", signalStop.String()))
	//application.Graphql.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDevelopment:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProduction:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log
}
