package tests

import (
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
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

func TestGRPC_RemovePlayer_Success(t *testing.T) {
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
	userId := int64(123)
	playerId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, userId).
		Return(nil)

	resp, err := client.RemovePlayer(ctx, &campaignv1.RemovePlayerRequest{
		CampaignId: campaignId,
		UserId:     userId,
		PlayerId:   playerId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_RemovePlayer_NotMaster(t *testing.T) {
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

	campaignId := int64(123)
	userId := int64(123)
	playerId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	resp, err := client.RemovePlayer(ctx, &campaignv1.RemovePlayerRequest{
		CampaignId: campaignId,
		UserId:     userId,
		PlayerId:   playerId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.PermissionDenied, st.Code())
	require.Equal(t, "user is not master", st.Message())
}

func TestGRPC_RemovePlayer_NotFound(t *testing.T) {
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
	userId := int64(123)
	playerId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemovePlayer(gomock.Any(), campaignId, userId).
		Return(models.ErrCampaignNotFound)

	resp, err := client.RemovePlayer(ctx, &campaignv1.RemovePlayerRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		PlayerId: playerId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "campaign not found", st.Message())
}

func TestGRPC_GetCampaignPlayers_Success(t *testing.T) {
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

	campaignId := int64(123)
	userId := int64(123)
	players := []int64{123, 456, 789}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		GetCampaignPlayers(gomock.Any(), campaignId).
		Return(players, nil)

	resp, err := client.GetCampaignPlayers(ctx, &campaignv1.GetCampaignPlayersRequest{
		CampaignId:  campaignId,
		UserId:      userId,
	})

	require.NoError(t, err)
	respPlayers := resp.GetPlayersId()
	assert.NotEmpty(t, respPlayers)

	for i, playerId := range players {
		assert.Equal(t, playerId, respPlayers[i])
	}
}

func TestGRPC_GetCampaignPlayers_CampaignNotFound(t *testing.T) {
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

	campaignId := int64(123)
	userId := int64(123)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		GetCampaignPlayers(gomock.Any(), campaignId).
		Return(nil, models.ErrCampaignNotFound)

	resp, err := client.GetCampaignPlayers(ctx, &campaignv1.GetCampaignPlayersRequest{
		CampaignId:  campaignId,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "campaign not found", st.Message())
}