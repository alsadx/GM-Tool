package suite

import (
	"campaigntool/internal/config"
	"context"
	"net"
	"strconv"
	"testing"

	"protos/gen/go/campaignv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	CampaignClient campaignv1.CampaignToolClient
}

const grpcHost = "localhost"

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../../config/local_test.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		CampaignClient: campaignv1.NewCampaignToolClient(cc),
	}
}
