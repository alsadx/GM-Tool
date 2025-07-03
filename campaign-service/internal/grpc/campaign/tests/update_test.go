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

func TestGRPC_UpdateCampaign_Success(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "New Test Campaign"
	newDescription := "This is a new test campaign"

	updatedCampaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        newName,
		Description: newDescription,
		MasterId:    userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	mockGameSaver.EXPECT().
		UpdateCampaign(gomock.Any(), updatedCampaign).
		Return(nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_UpdateCampaign_Success_EmptyName(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := ""
	newDescription := "This is a new test campaign"

	updatedCampaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: newDescription,
		MasterId:    userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	mockGameSaver.EXPECT().
		UpdateCampaign(gomock.Any(), updatedCampaign).
		Return(nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_UpdateCampaign_Success_EmptyDescription(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "New Test Campaign"
	newDescription := ""

	updatedCampaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        newName,
		Description: description,
		MasterId:    userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	mockGameSaver.EXPECT().
		UpdateCampaign(gomock.Any(), updatedCampaign).
		Return(nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestGRPC_UpdateCampaign_NotMaster(t *testing.T) {
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

	newName := "New Test Campaign"
	newDescription := "This is a new test campaign"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(false, nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.PermissionDenied, st.Code(), "unexpected error code")
	require.Equal(t, "user is not master", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_NoChanges(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "Test Campaign"
	newDescription := "This is a test campaign"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name: &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "no changes", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_SameName_EmptyDescription(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "Test Campaign"
	newDescription := ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name: &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "no changes", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_EmptyName_SameDescription(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := ""
	newDescription := "This is a test campaign"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name: &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "no changes", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_EmptyArguments(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := ""
	newDescription := ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name: &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "no changes", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_CampaignNotFound(t *testing.T) {
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

	name := "Test Campaign"
	description := "This is a test campaign"
	campaignId := int64(123)
	userId := int64(123)

	campaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        name,
		Description: description,
		MasterId:    userId,
	}

	newName := "New Test Campaign"
	newDescription := "This is a new test campaign"

	updatedCampaign := &models.CampaignInfo{
		Id:          campaignId,
		Name:        newName,
		Description: newDescription,
		MasterId:    userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(campaign, nil)

	mockGameSaver.EXPECT().
		UpdateCampaign(gomock.Any(), updatedCampaign).
		Return(models.ErrCampaignNotFound)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}

func TestGRPC_UpdateCampaign_CampaignNotFound_InProvider(t *testing.T) {
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

	newName := "New Test Campaign"
	newDescription := "This is a new test campaign"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockGameProvider.EXPECT().
		IsMaster(gomock.Any(), campaignId, userId).
		Return(true, nil)

	mockGameProvider.EXPECT().
		GetCampaign(gomock.Any(), campaignId).
		Return(nil, models.ErrCampaignNotFound)

	resp, err := client.UpdateCampaign(ctx, &campaignv1.UpdateCampaignRequest{
		CampaignId:  campaignId,
		Name:        &newName,
		Description: &newDescription,
		UserId:      userId,
	})

	require.Error(t, err)
	assert.False(t, resp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	assert.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}