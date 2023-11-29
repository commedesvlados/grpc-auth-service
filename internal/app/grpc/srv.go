package grpcsrv

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	authgrpc "github.com/commedesvlados/grpc_service/internal/grpc/auth"
)

type Srv struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app.
func New(log *slog.Logger, authService authgrpc.Auth, port int) *Srv {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &Srv{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server.
func (a *Srv) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *Srv) Run() error {
	const fn = "grpcsrv.Run"

	log := a.log.With(
		slog.String("fn", fn),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *Srv) Stop() {
	const fn = "grpcsrv.Stop"

	a.log.With(slog.String("fn", fn)).
		Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
