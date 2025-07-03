package tests

import (
	"campaigntool/internal/domain/models"
	grpccampaign "campaigntool/internal/grpc/campaign"
	"context"
	"net"
	"os"
	"strconv"
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

func TestGRPC_GetCreatedCampaign_Success(t *testing.T) {
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

	userId := int64(123)
	expectedCampaigns := make([]*models.Campaign, 2)
	for i := 0; i < 2; i++ {
		expectedCampaigns[i] = &models.Campaign{
			Id:           int64(i),
			Name:         "valid-campaign-name" + strconv.Itoa(i),
			Description:  "valid-campaign-description",
			PlayersCount: 4,
			PlayersId:    &[]int64{1, 2, 3, 4},
			CreatedAt:    time.Now(),
		}
	}

	mockGameProvider.EXPECT().
		CreatedCampaigns(gomock.Any(), userId).
		Return(expectedCampaigns, nil)

	resp, err := client.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.IsType(t, []*campaignv1.Campaign{}, resp.Campaigns)
	assert.Equal(t, expectedCampaigns[0].Name, resp.Campaigns[0].Name)
	assert.Equal(t, expectedCampaigns[1].Name, resp.Campaigns[1].Name)
	assert.Equal(t, expectedCampaigns[0].Id, resp.Campaigns[0].CampaignId)
	assert.Equal(t, expectedCampaigns[1].Id, resp.Campaigns[1].CampaignId)
	assert.Equal(t, expectedCampaigns[0].PlayersCount, resp.Campaigns[0].PlayersCount)
	assert.Equal(t, expectedCampaigns[1].PlayersCount, resp.Campaigns[1].PlayersCount)
	assert.Equal(t, expectedCampaigns[0].Description, *resp.Campaigns[0].Description)
	assert.Equal(t, expectedCampaigns[1].Description, *resp.Campaigns[1].Description)

	for i := 0; i < 2; i++ {
		for j := 0; j < int(expectedCampaigns[i].PlayersCount); j++ {
			players := *(expectedCampaigns[i].PlayersId)
			assert.Equal(t, players[j], resp.Campaigns[i].PlayersId[j])
		}
	}
}

func TestGRPC_GetCreatedCampaign_NoCampaigns(t *testing.T) {
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

	userId := int64(123)

	mockGameProvider.EXPECT().
		CreatedCampaigns(gomock.Any(), userId).
		Return(nil, models.ErrNoCampaigns)

	resp, err := client.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: userId,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaigns not found", st.Message(), "unexpected error message")
}


func TestGRPC_GetCurrentCampaign_Success(t *testing.T) {
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

	userId := int64(123)
	expectedCampaigns := make([]*models.CampaignForPlayer, 2)
	for i := 0; i < 2; i++ {
		expectedCampaigns[i] = &models.CampaignForPlayer{
			Id:       int64(i),
			Name:     "valid-campaign-name" + strconv.Itoa(i),
			MasterId: int64(i + 1),
		}
	}

	mockGameProvider.EXPECT().
		CurrentCampaigns(gomock.Any(), userId).
		Return(expectedCampaigns, nil)

	resp, err := client.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{
		UserId: userId,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.IsType(t, []*campaignv1.CampaignForPlayer{}, resp.Campaigns)
	assert.Equal(t, expectedCampaigns[0].Name, resp.Campaigns[0].Name)
	assert.Equal(t, expectedCampaigns[1].Name, resp.Campaigns[1].Name)
	assert.Equal(t, expectedCampaigns[0].Id, resp.Campaigns[0].CampaignId)
	assert.Equal(t, expectedCampaigns[1].Id, resp.Campaigns[1].CampaignId)
	assert.Equal(t, expectedCampaigns[0].MasterId, resp.Campaigns[0].MasterId)
	assert.Equal(t, expectedCampaigns[1].MasterId, resp.Campaigns[1].MasterId)
}

func TestGRPC_GetCurrentCampaign_NoCampaigns(t *testing.T) {
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

	userId := int64(123)

	mockGameProvider.EXPECT().
		CurrentCampaigns(gomock.Any(), userId).
		Return(nil, models.ErrNoCampaigns)

	resp, err := client.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{
		UserId: userId,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaigns not found", st.Message(), "unexpected error message")
}
