package tests

import (
	"auth"
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"testing"
	"time"

	"protos/gen/go/campaignv1"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGRPC_JoinCampaign_Success(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int32(123)
	inviteCode := "ABCDEF"

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, 1).
		Return(nil)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGRPC_JoinCampaign_InvalidInviteCode(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
	service, _, mockGameProvider := setupTest(t)
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

	inviteCode := "ABCDEF"

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(int32(0), models.ErrCampaignNotFound)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid invite code")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_JoinCampaign_NotFound(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int32(123)
	inviteCode := "ABCDEF"

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, 1).
		Return(models.ErrCampaignNotFound)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "campaign not found")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_JoinCampaign_AlreadyJoined(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int32(123)
	inviteCode := "ABCDEF"

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, 1).
		Return(models.ErrPlayerInCampaign)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "player is already in campaign")
	assert.False(t, resp.GetSuccess())
}

func TestGRPC_JoinCampaign_InvalidToken(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))
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

	inviteCode := "ABCDEF"

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid token")
	assert.Nil(t, resp)
}
