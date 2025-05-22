package tests

import (
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"testing"
	"time"

	"campaigntool/protos/campaignv1"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockGameSaver.EXPECT().
		SetInviteCode(gomock.Any(), campaignId, gomock.Any()).
		Return(nil)

	resp, err := client.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	mockGameSaver.EXPECT().
		SetInviteCode(gomock.Any(), campaignId, gomock.Any()).
		Return(models.ErrCampaignNotFound)

	resp, err := client.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
	})

	inviteCode := resp.GetInviteCode()

	require.Error(t, err)
	assert.ErrorContains(t, err, "campaign not found")
	assert.Empty(t, inviteCode)
}

func TestGRPC_GenerateInviteCode_InvalidToken(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, _, _ := setupTest(t)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer invalid-token"))

	resp, err := client.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid token")
	assert.Nil(t, resp)
}
