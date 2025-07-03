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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGRPC_JoinCampaign_Success(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := int64(123)
	campaignId := int64(123)
	charId := int64(111)
	inviteCode := "ABCDEF"

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, userId, charId).
		Return(nil)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode:  inviteCode,
		UserId:      userId,
		CharacterId: charId,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGRPC_JoinCampaign_InvalidInviteCode(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	defer os.Setenv("TEST_ENV", "")

	server := grpc.NewServer()
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
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(int64(0), models.ErrCampaignNotFound)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode:  inviteCode,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	assert.Equal(t, "invalid invite code", st.Message(), "unexpected error message")
}

func TestGRPC_JoinCampaign_NotFound(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer valid-token"))

	campaignId := int64(123)
	inviteCode := "ABCDEF"
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, userId, charId).
		Return(models.ErrCampaignNotFound)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode:  inviteCode,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}

func TestGRPC_JoinCampaign_AlreadyJoined(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	campaignId := int64(123)
	inviteCode := "ABCDEF"
	userId := int64(123)
	charId := int64(111)

	mockGameProvider.EXPECT().
		CheckInviteCode(gomock.Any(), inviteCode).
		Return(campaignId, nil)

	mockGameSaver.EXPECT().
		AddPlayer(gomock.Any(), campaignId, userId, charId).
		Return(models.ErrPlayerInCampaign)

	resp, err := client.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode:  inviteCode,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code")
	require.Equal(t, "player is already in campaign", st.Message(), "unexpected error message")
}

