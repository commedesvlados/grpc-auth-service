package app

import (
	"log/slog"
	"time"

	grpcsrv "github.com/commedesvlados/grpc_service/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcsrv.Srv
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcSrv := grpcsrv.New(log, grpcPort)

	return &App{GRPCSrv: grpcSrv}
}
