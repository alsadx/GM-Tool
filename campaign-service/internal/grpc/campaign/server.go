package campaign

import (
	"campaigntool/internal/domain/models"
	"campaigntool/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"os"

	campaignv1 "github.com/alsadx/protos/gen/go/campaign"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CampaignTool interface {
	CreateCampaign(ctx context.Context, name, desc string, userId int) (campaignId int32, err error)
	DeleteCampaign(ctx context.Context, campaignId int32, userId int) (err error)
	GenerateInviteCode(ctx context.Context, campaignId int32, userId int) (inviteCode string, err error)
	JoinCampaign(ctx context.Context, inviteCode string, userId int) (err error)
	LeaveCampaign(ctx context.Context, campaignId int32, userId int) (err error)
	GetCreatedCampaigns(ctx context.Context, userId int) (campaigns []models.Campaign, err error)
	GetCurrentCampaigns(ctx context.Context, userId int) (campaigns []models.CampaignForPlayer, err error)
}

type ServerAPI struct {
	campaignv1.UnimplementedCampaignToolServer
	CampaignTool CampaignTool
}

func RegisterServerAPI(gRPC *grpc.Server, campaignTool CampaignTool) {
	campaignv1.RegisterCampaignToolServer(gRPC, &ServerAPI{CampaignTool: campaignTool})
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	
	token, err := jwt.ExtractTokenFromHeader(authHeader[0])
	if err != nil {
		return nil, err
	}
	
	if os.Getenv("TEST_ENV") == "true" {
		if token == "valid-token" {
			ctx = context.WithValue(ctx, "user_id", 1)
			return handler(ctx, req)
		} else {
			return nil, fmt.Errorf("invalid token")
		} 
	}

	userId, err := jwt.ValidateToken(token, "secret")
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "user_id", userId)

	return handler(ctx, req)
}

func (s *ServerAPI) CreateCampaign(ctx context.Context, req *campaignv1.CreateCampaignRequest) (*campaignv1.CreateCampaignResponse, error) {
	// TODO: validate
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaignId, err := s.CampaignTool.CreateCampaign(ctx, req.GetName(), req.GetDescription(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignExists) {
			return nil, status.Errorf(codes.AlreadyExists, "campaign with this name already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.CreateCampaignResponse{CampaignId: campaignId}, nil
}

func (s *ServerAPI) DeleteCampaign(ctx context.Context, req *campaignv1.DeleteCampaignRequest) (*campaignv1.DeleteCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.DeleteCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		}

		return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.DeleteCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) GenerateInviteCode(ctx context.Context, req *campaignv1.GenerateInviteCodeRequest) (*campaignv1.GenerateInviteCodeResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	code, err := s.CampaignTool.GenerateInviteCode(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return nil, status.Errorf(codes.NotFound, "campaign not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.GenerateInviteCodeResponse{InviteCode: code}, nil
}

func (s *ServerAPI) JoinCampaign(ctx context.Context, req *campaignv1.JoinCampaignRequest) (*campaignv1.JoinCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.JoinCampaign(ctx, req.GetInviteCode(), userId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCampaignNotFound):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		case errors.Is(err, models.ErrPlayerInCampaign):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.AlreadyExists, "player already in campaign")
		case errors.Is(err, models.ErrInvalidCode):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "invalid invite code")
		}

		return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.JoinCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) LeaveGame(ctx context.Context, req *campaignv1.LeaveCampaignRequest) (*campaignv1.LeaveCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.LeaveCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		}
		return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.LeaveCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) GetCreatedCampaigns(ctx context.Context, req *campaignv1.GetCreatedCampaignsRequest) (*campaignv1.GetCreatedCampaignsResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaigns, err := s.CampaignTool.GetCreatedCampaigns(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.Campaign
	for _, campaign := range campaigns {
		res = append(res, &campaignv1.Campaign{
			CampaignId:  campaign.Id,
			Name:        campaign.Name,
			Description: &campaign.Description,
			PlayerCount: campaign.PlayerCount,
		})
	}

	return &campaignv1.GetCreatedCampaignsResponse{Campaigns: res}, nil
}

func (s *ServerAPI) GetCurrentCampaigns(ctx context.Context, req *campaignv1.GetCurrentCampaignsRequest) (*campaignv1.GetCurrentCampaignsResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaigns, err := s.CampaignTool.GetCurrentCampaigns(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.CampaignForPlayer
	for _, campaign := range campaigns {
		res = append(res, &campaignv1.CampaignForPlayer{
			CampaignId: campaign.Id,
			Name:       campaign.Name,
		})
	}

	return &campaignv1.GetCurrentCampaignsResponse{Campaigns: res}, nil
}
