package tests

import (
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/alsadx/gm-protos/gen/go/campaignv1"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGRPC_GenerateInviteCode_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockGameSaver, _ := setupTest(t)
	srv := grpccampaign.ServerAPI{
		CampaignTool: service,
	}

	campaignv1.RegisterCampaignToolServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := campaignv1.NewCampaignToolClient(clientConn)

	campaignId := int64(123)
	userId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameSaver.EXPECT().
		SetInviteCode(gomock.Any(), campaignId, userId, gomock.Any()).
		Return(nil)

	resp, err := client.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
		UserId:     userId,
	})

	inviteCode := resp.GetInviteCode()

	require.NoError(t, err)
	assert.NotEmpty(t, inviteCode)
	require.Condition(t, func() bool { return len(inviteCode) == 6 })
}

func TestGRPC_GenerateInviteCode_NotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockGameSaver, _ := setupTest(t)
	srv := grpccampaign.ServerAPI{
		CampaignTool: service,
	}

	campaignv1.RegisterCampaignToolServer(server, &srv)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(listener)
	defer server.Stop()

	serverAddress := listener.Addr().String()

	clientConn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	require.NoError(t, err)
	defer clientConn.Close()

	client := campaignv1.NewCampaignToolClient(clientConn)

	campaignId := int64(123)
	userId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameSaver.EXPECT().
		SetInviteCode(gomock.Any(), campaignId, userId, gomock.Any()).
		Return(models.ErrCampaignNotFound)

	resp, err := client.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
		UserId:     userId,
	})

	inviteCode := resp.GetInviteCode()

	require.Error(t, err)
	assert.Empty(t, inviteCode)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}
