package app

import (
	"log/slog"
	"time"

	grpcsrv "github.com/commedesvlados/grpc_service/internal/app/grpc"
	"github.com/commedesvlados/grpc_service/internal/servises/auth"
	"github.com/commedesvlados/grpc_service/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcsrv.Srv
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcSrv := grpcsrv.New(log, authService, grpcPort)

	return &App{GRPCSrv: grpcSrv}
}
