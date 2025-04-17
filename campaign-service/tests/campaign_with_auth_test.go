package tests

import (
	"campaigntool/tests/suite"
	"testing"

	campaignv1 "github.com/alsadx/protos/gen/go/campaign"
	ssov1 "github.com/alsadx/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestCreateDeleteCampaign_WithAuth(t *testing.T) {
	ctx, suite := suite.New(t)

	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	authClient := ssov1.NewAuthClient(cc)
	defer cc.Close()
	require.NoError(t, err)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	md := metadata.Pairs("authorization", "Bearer " + token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	name := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        name,
		Description: &desc,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	delResp, err := suite.CampaignClient.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
	})
	require.NoError(t, err)
	assert.True(t, delResp.Success)
}

func TestCreatedCampaign_WithAuth(t *testing.T) {
	ctx, suite := suite.New(t)

	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	authClient := ssov1.NewAuthClient(cc)
	defer cc.Close()
	require.NoError(t, err)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	md := metadata.Pairs("authorization", "Bearer " + token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	names := make([]string, 10)
	descriptions := make([]string, 10)

	for i := 0; i < 10; i++ {
		name := gofakeit.Name()
		desc := gofakeit.Sentence(10)

		names[i] = name
		descriptions[i] = desc

		createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
			Name:        name,
			Description: &desc,
		})
		require.NoError(t, err)
		campaignId := createResp.GetCampaignId()
		assert.NotEmpty(t, campaignId, "campaign id should not be empty")
	}

	getResp, err := suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{})
	require.NoError(t, err)
	assert.Len(t, getResp.Campaigns, 10)

	for i := 0; i < 10; i++ {
		assert.Equal(t, names[9-i], getResp.Campaigns[i].Name)
		assert.Equal(t, descriptions[9-i], *getResp.Campaigns[i].Description)
	}
}

func TestJoinCampaign(t *testing.T) {
	// register master and create game
	ctx, suite := suite.New(t)

	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	authClient := ssov1.NewAuthClient(cc)
	defer cc.Close()
	require.NoError(t, err)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	masterToken := logResp.GetToken()
	require.NotEmpty(t, masterToken, "expected non-empty token")

	md := metadata.Pairs("authorization", "Bearer " + masterToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	name := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        name,
		Description: &desc,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	// Generate invite code
	inviteResp, err := suite.CampaignClient.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
	})
	require.NoError(t, err)
	inviteCode := inviteResp.GetInviteCode()
	assert.NotEmpty(t, inviteCode, "invite code should not be empty")

	// Register player and join campaign
	cc, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %w", err)
	}
	authClient = ssov1.NewAuthClient(cc)
	defer cc.Close()
	require.NoError(t, err)

	logName = gofakeit.Name()
	email = gofakeit.Email()
	password = gofakeit.Password(true, true, true, true, false, 6)

	regResp, err = authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	logResp, err = authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	playerToken := logResp.GetToken()
	require.NotEmpty(t, playerToken, "expected non-empty token")

	md = metadata.Pairs("authorization", "Bearer " + playerToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	joinResp, err := suite.CampaignClient.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
	})
	require.NoError(t, err)
	assert.Equal(t, true, joinResp.GetSuccess())
}