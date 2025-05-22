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

func TestGRPC_DeleteCampaign_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockGameSaver, mockGameProvider := setupTest(t)
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

	mockGameProvider.EXPECT().
		GetCampaignPlayers(gomock.Any(), int(campaignId)).
		Return([]int{}, nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(gomock.Any(), campaignId, 1).
		Return(nil)

	resp, err := client.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_DeleteCampaign_NotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockGameSaver, mockGameProvider := setupTest(t)
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

	mockGameProvider.EXPECT().
		GetCampaignPlayers(gomock.Any(), int(campaignId)).
		Return([]int{}, nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(gomock.Any(), campaignId, 1).
		Return(models.ErrCampaignNotFound)

	resp, err := client.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "campaign not found")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_DeleteCampaign_InvalidToken(t *testing.T) {
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

	resp, err := client.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid token")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_DeleteCampaign_CampaignWithPlayers(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
	service, mockGameSaver, mockGameProvider := setupTest(t)
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

	mockGameProvider.EXPECT().
		GetCampaignPlayers(gomock.Any(), int(campaignId)).
		Return([]int{12, 13}, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, 12).
		Return(nil)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, 13).
		Return(nil)

	mockGameSaver.EXPECT().
		DeleteCampaign(gomock.Any(), campaignId, 1).
		Return(nil)

	resp, err := client.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}
