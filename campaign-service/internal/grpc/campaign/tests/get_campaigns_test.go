package tests

import (
	"auth"
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"strconv"
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

func TestGRPC_GetCreatedCampaign_Success(t *testing.T) {
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

	expectedCampaigns := make([]*models.Campaign, 2)
	for i := 0; i < 2; i++ {
		expectedCampaigns[i] = &models.Campaign{
			Id:          int32(i),
			Name:        "valid-campaign-name" + strconv.Itoa(i),
			Description: "valid-campaign-description",
			PlayerCount: 4,
			CreatedAt:   time.Now(),
		}
	}

	mockGameProvider.EXPECT().
		CreatedCampaigns(gomock.Any(), 1).
		Return(expectedCampaigns, nil)

	resp, err := client.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{})

	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.IsType(t, []*campaignv1.Campaign{}, resp.Campaigns)
	assert.Equal(t, expectedCampaigns[0].Name, resp.Campaigns[0].Name)
	assert.Equal(t, expectedCampaigns[1].Name, resp.Campaigns[1].Name)
	assert.Equal(t, expectedCampaigns[0].Id, resp.Campaigns[0].CampaignId)
	assert.Equal(t, expectedCampaigns[1].Id, resp.Campaigns[1].CampaignId)
	assert.Equal(t, expectedCampaigns[0].PlayerCount, resp.Campaigns[0].PlayerCount)
	assert.Equal(t, expectedCampaigns[1].PlayerCount, resp.Campaigns[1].PlayerCount)
	assert.Equal(t, expectedCampaigns[0].Description, *resp.Campaigns[0].Description)
	assert.Equal(t, expectedCampaigns[1].Description, *resp.Campaigns[1].Description)
}

func TestGRPC_GetCreatedCampaign_NoCampaigns(t *testing.T) {
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

	mockGameProvider.EXPECT().
		CreatedCampaigns(gomock.Any(), 1).
		Return(nil, models.ErrNoCampaigns)

	resp, err := client.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaigns not found")
	assert.Nil(t, resp)
}

func TestGRPC_GetCreatedCampaign_InvalidToken(t *testing.T) {
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

	resp, err := client.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	assert.Nil(t, resp)
}

func TestGRPC_GetCurrentCampaign_Success(t *testing.T) {
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

	expectedCampaigns := make([]*models.CampaignForPlayer, 2)
	for i := 0; i < 2; i++ {
		expectedCampaigns[i] = &models.CampaignForPlayer{
			Id:   int32(i),
			Name: "valid-campaign-name" + strconv.Itoa(i),
		}
	}

	mockGameProvider.EXPECT().
		CurrentCampaigns(gomock.Any(), 1).
		Return(expectedCampaigns, nil)

	resp, err := client.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{})

	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.IsType(t, []*campaignv1.CampaignForPlayer{}, resp.Campaigns)
	assert.Equal(t, expectedCampaigns[0].Name, resp.Campaigns[0].Name)
	assert.Equal(t, expectedCampaigns[1].Name, resp.Campaigns[1].Name)
	assert.Equal(t, expectedCampaigns[0].Id, resp.Campaigns[0].CampaignId)
	assert.Equal(t, expectedCampaigns[1].Id, resp.Campaigns[1].CampaignId)
}

func TestGRPC_GetCurrentCampaign_NoCampaigns(t *testing.T) {
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

	mockGameProvider.EXPECT().
		CurrentCampaigns(gomock.Any(), 1).
		Return(nil, models.ErrNoCampaigns)

	resp, err := client.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaigns not found")
	assert.Nil(t, resp)
}

func TestGRPC_GetCurrentCampaign_InvalidToken(t *testing.T) {
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

	resp, err := client.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	assert.Nil(t, resp)
}
