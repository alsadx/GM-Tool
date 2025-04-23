package tests

import (
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"testing"
	"time"

	campaignv1 "github.com/alsadx/protos/gen/go/campaign"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGRPC_LeaveCampaign_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(grpccampaign.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int32(123)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, 1).
		Return(nil)

	resp, err := client.LeaveCampaign(ctx, &campaignv1.LeaveCampaignRequest{
		CampaignId: campaignId,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGRPC_LeaveCampaign_NotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(grpccampaign.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int32(123)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, 1).
		Return(models.ErrCampaignNotFound)

	resp, err := client.LeaveCampaign(ctx, &campaignv1.LeaveCampaignRequest{
		CampaignId:   campaignId,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "campaign not found")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_LeaveCampaign_InvalidToken(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(grpccampaign.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer invalid-token"))

	campaignId := int32(123)

	resp, err := client.LeaveCampaign(ctx, &campaignv1.LeaveCampaignRequest{
		CampaignId: campaignId,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid token")
	assert.Nil(t, resp)
}