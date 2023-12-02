package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	gspv1 "github.com/commedesvlados/grpc-service-protos/gen/go/grpc_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/commedesvlados/grpc_service/internal/config"
)

type Suite struct {
	*testing.T
	AuthClient gspv1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	config.MustLoadVariablesByPath("../.env.local", "../config/config.local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), config.E.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		AuthClient: gspv1.NewAuthClient(cc),
	}
}

func grpcAddress() string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(config.E.GRPC.Port))
}
