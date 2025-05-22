package e2e

import (
	"campaigntool/tests/suite"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"campaigntool/protos/campaignv1"

	"campaigntool/protos/ssov1"
	"github.com/brianvoe/gofakeit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	cc         *grpc.ClientConn
	authClient ssov1.AuthClient
)

func SetUp() {
	composeFilePath := "../../../sso-service/docker/docker-compose.yaml"
	composeCmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	output, err := composeCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("docker-compose up failed: %v\nOutput: %s", err, string(output))
	}
	// require.NoError(t, composeCmd.Run())

	fmt.Println("Waiting for SSO service to start...")
	waitForServer("localhost:50051")

	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	authClient = ssov1.NewAuthClient(cc)
	fmt.Println("Auth client created")
}

func TearDown() {
	if cc != nil {
		cc.Close()
	}
	composeFilePath := "../../../sso-service/docker/docker-compose.yaml"
	downCmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	output, err := downCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("docker-compose down failed: %v\nOutput: %s", err, string(output))
	}
	fmt.Println("Docker Compose Down Output:", string(output))
}

func TestMain(m *testing.M) {
	fmt.Println("Starting E2E tests...")

	SetUp()

	exitCode := m.Run()

	TearDown()

	os.Exit(exitCode)
}

func TestE2E_CreateDeleteGetCampaign_Success(t *testing.T) {
	ctx, suite := suite.New(t)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	campaignName := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	// Check created campaigns

	getCreatedResp, err := suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns := getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(0), campaigns[0].PlayersCount, "expected campaign players to be 0")
	assert.Empty(t, campaigns[0].PlayersId)

	// Delete campaigns

	delResp, err := suite.CampaignClient.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.True(t, delResp.GetSuccess())

	// Check created campaigns again

	getCreatedResp, err = suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	assert.Error(t, err)
	assert.Empty(t, getCreatedResp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaigns not found", st.Message(), "unexpected error message")
}

func TestE2E_Create_CampaignAlreadyExists(t *testing.T) {
	ctx, suite := suite.New(t)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	campaignName := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	createResp, err = suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})
	assert.Error(t, err)
	assert.Empty(t, createResp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code")
	require.Equal(t, "campaign with this name already exists", st.Message(), "unexpected error message")
}

func TestE2E_Create_EmptyCampaignName(t *testing.T) {
	ctx, suite := suite.New(t)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	campaignName := ""
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})

	assert.Error(t, err)
	assert.Empty(t, createResp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.InvalidArgument, st.Code(), "unexpected error code")
	require.Equal(t, "name is required", st.Message(), "unexpected error message")
}

func TestE2E_Create_Unauthenticated(t *testing.T) {
	ctx, suite := suite.New(t)

	campaignName := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: int64(0),
	})
	assert.Error(t, err)
	assert.Empty(t, createResp)

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.Unauthenticated, st.Code(), "unexpected error code")
	require.Equal(t, "unauthenticated", st.Message(), "unexpected error message")
}

func TestE2E_Delete_CampaignNotFound(t *testing.T) {
	ctx, suite := suite.New(t)

	logName := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     logName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := logResp.GetToken()
	require.NotEmpty(t, token, "expected non-empty token")

	delResp, err := suite.CampaignClient.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: int64(123),
		UserId: masterId,
	})
	assert.Error(t, err)
	require.False(t, delResp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.NotFound, st.Code(), "unexpected error code")
	require.Equal(t, "campaign not found", st.Message(), "unexpected error message")
}

func TestE2E_Delete_Unauthenticated(t *testing.T) {
	ctx, suite := suite.New(t)

	delResp, err := suite.CampaignClient.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: int64(123),
		UserId: int64(0),
	})
	assert.Error(t, err)
	assert.False(t, delResp.GetSuccess())

	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.Unauthenticated, st.Code(), "unexpected error code")
	require.Equal(t, "unauthenticated", st.Message(), "unexpected error message")
}

func TestE2E_JoinWithInviteCode_LeaveCampaign_Success(t *testing.T) {
	ctx, suite := suite.New(t)

	// Register master and create campaign

	masterName := gofakeit.Name()
	masterEmail := gofakeit.Email()
	masterPassword := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    masterEmail,
		Password: masterPassword,
		Name:     masterName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    masterEmail,
		Password: masterPassword,
	})
	require.NoError(t, err)
	masterToken := logResp.GetToken()
	require.NotEmpty(t, masterToken, "expected non-empty token")

	campaignName := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	genResp, err := suite.CampaignClient.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
		UserId: masterId,
	})
	require.NoError(t, err)
	inviteCode := genResp.GetInviteCode()
	assert.NotEmpty(t, inviteCode, "invite code should not be empty")

	// Check created campaigns

	getCreatedResp, err := suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns := getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(0), campaigns[0].PlayersCount, "expected campaign players to be 0")
	assert.Empty(t, campaigns[0].PlayersId)

	// Register player and join with invite code

	playerName := gofakeit.Name()
	playerEmail := gofakeit.Email()
	playerPassword := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err = authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    playerEmail,
		Password: playerPassword,
		Name:     playerName,
	})
	require.NoError(t, err)
	playerId := regResp.GetUserId()
	assert.NotEmpty(t, playerId)

	logResp, err = authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    playerEmail,
		Password: playerPassword,
	})
	require.NoError(t, err)
	playerToken := logResp.GetToken()
	require.NotEmpty(t, playerToken, "expected non-empty token")

	joinResp, err := suite.CampaignClient.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
		UserId: playerId,
	})
	require.NoError(t, err)
	assert.True(t, joinResp.GetSuccess())

	// Get current campaigns

	getCurrentResp, err := suite.CampaignClient.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{
		UserId: playerId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCurrentResp, "no current campaigns")
	currentCampaigns := getCurrentResp.GetCampaigns()
	assert.Equal(t, campaignName, currentCampaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, campaignId, currentCampaigns[0].CampaignId, "expected campaign id to match")
	assert.Equal(t, masterId, currentCampaigns[0].MasterId, "expectd master id to match: ", masterId, currentCampaigns[0].MasterId)

	// Check created campaigns again

	getCreatedResp, err = suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns = getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(1), campaigns[0].PlayersCount, "expected campaign players to be 1")
	assert.Equal(t, playerId, campaigns[0].PlayersId[0])

	// Try to join again

	joinResp, err = suite.CampaignClient.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
		UserId: playerId,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "error is not a gRPC status error")

	require.Equal(t, codes.AlreadyExists, st.Code(), "unexpected error code")
	require.Equal(t, "player is already in campaign", st.Message(), "unexpected error message")
	assert.False(t, joinResp.GetSuccess(), "expected join to fail")

	// Leave campaign

	leaveResp, err := suite.CampaignClient.LeaveCampaign(ctx, &campaignv1.LeaveCampaignRequest{
		CampaignId: campaignId,
		UserId: playerId,
	})
	require.NoError(t, err)
	assert.True(t, leaveResp.GetSuccess())

	// Check current campaigns again

	getCurrentResp, err = suite.CampaignClient.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{
		UserId: playerId,
	})
	require.Error(t, err)
	assert.Empty(t, getCurrentResp, "should be no current campaigns")

	// Check created campaigns again

	getCreatedResp, err = suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns = getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(0), campaigns[0].PlayersCount, "expected campaign players to be 0")
	assert.Empty(t, campaigns[0].PlayersId)
}

func TestE2E_DeleteCampaignWithPlayers_Success(t *testing.T) {
	ctx, suite := suite.New(t)

	// Register master and create campaign

	masterName := gofakeit.Name()
	masterEmail := gofakeit.Email()
	masterPassword := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    masterEmail,
		Password: masterPassword,
		Name:     masterName,
	})
	require.NoError(t, err)
	masterId := regResp.GetUserId()
	assert.NotEmpty(t, masterId)

	logResp, err := authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    masterEmail,
		Password: masterPassword,
	})
	require.NoError(t, err)
	masterToken := logResp.GetToken()
	require.NotEmpty(t, masterToken, "expected non-empty token")

	campaignName := gofakeit.Name()
	desc := gofakeit.Sentence(10)

	createResp, err := suite.CampaignClient.CreateCampaign(ctx, &campaignv1.CreateCampaignRequest{
		Name:        campaignName,
		Description: &desc,
		UserId: masterId,
	})
	require.NoError(t, err)
	campaignId := createResp.GetCampaignId()
	assert.NotEmpty(t, campaignId, "campaign id should not be empty")

	genResp, err := suite.CampaignClient.GenerateInviteCode(ctx, &campaignv1.GenerateInviteCodeRequest{
		CampaignId: campaignId,
		UserId: masterId,
	})
	require.NoError(t, err)
	inviteCode := genResp.GetInviteCode()
	assert.NotEmpty(t, inviteCode, "invite code should not be empty")

	// Check created campaigns

	getCreatedResp, err := suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns := getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(0), campaigns[0].PlayersCount, "expected campaign players to be 0")
	assert.Empty(t, campaigns[0].PlayersId)

	// Register player and join with invite code

	playerName := gofakeit.Name()
	playerEmail := gofakeit.Email()
	playerPassword := gofakeit.Password(true, true, true, true, false, 6)

	regResp, err = authClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    playerEmail,
		Password: playerPassword,
		Name:     playerName,
	})
	require.NoError(t, err)
	playerId := regResp.GetUserId()
	assert.NotEmpty(t, playerId)

	logResp, err = authClient.Login(ctx, &ssov1.LoginRequest{
		Email:    playerEmail,
		Password: playerPassword,
	})
	require.NoError(t, err)
	playerToken := logResp.GetToken()
	require.NotEmpty(t, playerToken, "expected non-empty token")

	joinResp, err := suite.CampaignClient.JoinCampaign(ctx, &campaignv1.JoinCampaignRequest{
		InviteCode: inviteCode,
		UserId: playerId,
	})
	require.NoError(t, err)
	assert.True(t, joinResp.GetSuccess())

	// Get current campaigns

	getCurrentResp, err := suite.CampaignClient.GetCurrentCampaigns(ctx, &campaignv1.GetCurrentCampaignsRequest{
		UserId: playerId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCurrentResp, "no current campaigns")
	currentCampaigns := getCurrentResp.GetCampaigns()
	assert.Equal(t, campaignName, currentCampaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, campaignId, currentCampaigns[0].CampaignId, "expected campaign id to match")
	assert.Equal(t, masterId, currentCampaigns[0].MasterId, "expected master id match")

	// Check created campaigns again

	getCreatedResp, err = suite.CampaignClient.GetCreatedCampaigns(ctx, &campaignv1.GetCreatedCampaignsRequest{
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, getCreatedResp, "no created campaigns")
	campaigns = getCreatedResp.GetCampaigns()
	assert.NotEmpty(t, campaigns, "no created campaigns")
	assert.Equal(t, 1, len(campaigns), "expected one created campaign")
	assert.Equal(t, campaignId, int64(campaigns[0].CampaignId), "expected campaign id to match")
	assert.Equal(t, campaignName, campaigns[0].Name, "expected campaign name to match")
	assert.Equal(t, desc, *campaigns[0].Description, "expected campaign description to match")
	assert.Equal(t, int64(1), campaigns[0].PlayersCount, "expected campaign players to be 1")
	assert.Equal(t, playerId, campaigns[0].PlayersId[0])

	// Delete campaigns

	delResp, err := suite.CampaignClient.DeleteCampaign(ctx, &campaignv1.DeleteCampaignRequest{
		CampaignId: campaignId,
		UserId: masterId,
	})
	require.NoError(t, err)
	assert.True(t, delResp.GetSuccess())
}

func waitForServer(addr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Fatalf("server at %s is not ready after 2 minutes", addr)
		default:
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err == nil {
				client := ssov1.NewAuthClient(conn)
				resp, err := client.HealthCheck(ctx, &ssov1.HealthCheckRequest{})
				fmt.Println("trying to connect to server...")
				if err == nil && resp.Status == ssov1.HealthCheckResponse_SERVING {
					conn.Close()
					fmt.Println("Server is fully ready!")
					return
				}
				conn.Close()
			}
			time.Sleep(5 * time.Second)
		}
	}
}
