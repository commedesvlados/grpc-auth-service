package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/commedesvlados/grpc_service/internal/app"
	cfg "github.com/commedesvlados/grpc_service/internal/config"
)

func main() {
	cfg.MustLoadVariables()

	log := setupLogger(cfg.C.Env)

	log.Info("[App Main] starting application")

	application := app.New(log, cfg.E.GRPC.Port, cfg.E.Database.Path, cfg.E.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("[App Main] stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("[App Main] application stopped")
}

func setupLogger(env string) (log *slog.Logger) {
	switch env {
	case cfg.EnvironmentProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case cfg.EnvironmentDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case cfg.EnvironmentLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
