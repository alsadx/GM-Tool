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

func TestGRPC_AddCharacter_Success(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsPlayer(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		AddCharacter(gomock.Any(), campaignId, userId, charId).
		Return(nil)

	resp, err := client.AddCharacter(ctx, &campaignv1.AddCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_AddCharacter_NotPlayer(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsPlayer(gomock.Any(), campaignId, userId).
		Return(false, nil)

	resp, err := client.AddCharacter(ctx, &campaignv1.AddCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.PermissionDenied, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "user is not player", st.Message(), "unexpected error message")
}

func TestGRPC_AddCharacter_CharacterInCampaign(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		IsPlayer(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameSaver.EXPECT().
		AddCharacter(gomock.Any(), campaignId, userId, charId).
		Return(models.ErrCharacterInCampaign)

	resp, err := client.AddCharacter(ctx, &campaignv1.AddCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "character already in campaign", st.Message(), "unexpected error message")
}

func TestGRPC_RemoveCharacter_Success(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(gomock.Any(), campaignId, userId, charId).
		Return(true, nil)

	mockGameSaver.EXPECT().
		RemoveCharacter(gomock.Any(), campaignId, userId, charId).
		Return(nil)

	resp, err := client.RemoveCharacter(ctx, &campaignv1.RemoveCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_RemoveCharacter_NotCharacterOwner(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(gomock.Any(), campaignId, userId, charId).
		Return(false, nil)

	resp, err := client.RemoveCharacter(ctx, &campaignv1.RemoveCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.PermissionDenied, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "user is not character owner", st.Message(), "unexpected error message")
}

func TestGRPC_RemoveCharacter_NotFound(t *testing.T) {
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
	charId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		IsCharacterOwner(gomock.Any(), campaignId, userId, charId).
		Return(false, nil)

	mockGameSaver.EXPECT().
		RemoveCharacter(gomock.Any(), campaignId, userId, charId).
		Return(models.ErrCharacterNotFound)

	resp, err := client.RemoveCharacter(ctx, &campaignv1.RemoveCharacterRequest{
		CampaignId:  campaignId,
		UserId:      userId,
		CharacterId: charId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "character not found", st.Message(), "unexpected error message")
}

func TestGRPC_GetCampaignCharacters_Success(t *testing.T) {
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

	characters := []*models.CampaignCharacter{
		{
			CharIds: []int64{1, 2, 3},
			UserId:        123,
		},
		{
			CharIds: []int64{6, 4},
			UserId:        64,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignCharacters(gomock.Any(), campaignId).
		Return(characters, nil)

	resp, err := client.GetCampaignCharacters(ctx, &campaignv1.GetCampaignCharactersRequest{
		CampaignId:  campaignId,
		UserId:      userId,
	})

	require.NoError(t, err)
	chars := resp.GetCharacters()
	assert.NotEmpty(t, chars)

	for i, char := range chars {
		assert.Equal(t, characters[i].CharIds, char.CharacterIds)
		assert.Equal(t, characters[i].UserId, char.UserId)
	}
}

func TestGRPC_GetCampaignCharacters_NotMaster(t *testing.T) {
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
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)


	resp, err := client.GetCampaignCharacters(ctx, &campaignv1.GetCampaignCharactersRequest{
		CampaignId:  campaignId,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.PermissionDenied, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "user is not master", st.Message(), "unexpected error message")
}

func TestGRPC_GetCampaignCharacters_CampaignNotFound(t *testing.T) {
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
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaignCharacters(gomock.Any(), campaignId).
		Return(nil, models.ErrCampaignNotFound)

	resp, err := client.GetCampaignCharacters(ctx, &campaignv1.GetCampaignCharactersRequest{
		CampaignId:  campaignId,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.Empty(t, resp)
	
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code: %v", st.Code())
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}

func TestGRPC_GetPlayerCharacters_Success_Master(t *testing.T) {
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
	playerId := int64(111)

	characters := []int64{1, 2, 3}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(gomock.Any(), campaignId, playerId).
		Return(characters, nil)

	resp, err := client.GetPlayerCharacters(ctx, &campaignv1.GetPlayerCharactersRequest{
		CampaignId:  campaignId,
		PlayerId:    playerId,
		UserId:      userId,
	})

	require.NoError(t, err)
	chars := resp.GetCharactersId()
	assert.NotEmpty(t, chars)

	for i := range chars {
		assert.Equal(t, characters[i], chars[i])
		assert.Equal(t, characters[i], chars[i])
	}
}

func TestGRPC_GetPlayerCharacters_Success_NotMaster(t *testing.T) {
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

	characters := []int64{1, 2, 3}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(gomock.Any(), campaignId, playerId).
		Return(characters, nil)

	resp, err := client.GetPlayerCharacters(ctx, &campaignv1.GetPlayerCharactersRequest{
		CampaignId:  campaignId,
		PlayerId:    playerId,
		UserId:      userId,
	})

	require.NoError(t, err)
	chars := resp.GetCharactersId()
	assert.NotEmpty(t, chars)

	for i := range chars {
		assert.Equal(t, characters[i], chars[i])
		assert.Equal(t, characters[i], chars[i])
	}
}

func TestGRPC_GetPlayerCharacters_NotCharactersOwner(t *testing.T) {
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
	playerId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	resp, err := client.GetPlayerCharacters(ctx, &campaignv1.GetPlayerCharactersRequest{
		CampaignId:  campaignId,
		PlayerId:    playerId,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.PermissionDenied, st.Code(), "unexpected error code")
	require.Equal(t, "user is not character owner", st.Message(), "unexpected error message")
}

func TestGRPC_GetPlayerCharacters_NotFound(t *testing.T) {
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
	playerId := int64(111)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetPlayerCharacters(gomock.Any(), campaignId, playerId).
		Return(nil, models.ErrCampaignNotFound)

	resp, err := client.GetPlayerCharacters(ctx, &campaignv1.GetPlayerCharactersRequest{
		CampaignId:  campaignId,
		PlayerId:    playerId,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.Empty(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}